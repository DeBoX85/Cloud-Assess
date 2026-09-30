#!/usr/bin/env bash
# Publish a scoped proposal only. The protected branch is never a push target.
set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo 'Usage: publish-maintenance.sh tidy|rule-import' >&2
  exit 1
fi
case "$1" in
  tidy)
    paths=(go.mod go.sum)
    message='build: refresh Go dependency lock'
    ;;
  rule-import)
    paths=(.gitmodules internal/rules/custom internal/rules/upstream)
    message='data: import pinned assessment rule baseline'
    ;;
  *) echo 'Unknown maintenance kind' >&2; exit 1 ;;
esac
if [[ ! ${GITHUB_RUN_ID:-} =~ ^[0-9]+$ ]] || [ -z "${GITHUB_REPOSITORY:-}" ] || [ -z "${GITHUB_STEP_SUMMARY:-}" ]; then
  echo 'GitHub run ID, repository and step summary are required' >&2
  exit 1
fi
if [ "$(git branch --show-current)" != 'bootstrap/core-v1' ]; then
  echo 'Maintenance publication must start on bootstrap/core-v1' >&2
  exit 1
fi

# Do not accidentally commit an unrelated change already staged by another step.
while IFS= read -r -d '' path; do
  allowed=false
  for root in "${paths[@]}"; do
    if [[ "$path" == "$root" || "$path" == "$root/"* ]]; then
      allowed=true
      break
    fi
  done
  if [ "$allowed" = false ]; then
    echo "Unexpected staged path: $path" >&2
    exit 1
  fi
done < <(git diff --cached --name-only -z)

git add -- "${paths[@]}"
if git diff --cached --quiet; then
  echo 'Maintenance output already matches; no proposal branch created.'
  exit 0
fi

branch="maintenance/$1-${GITHUB_RUN_ID}"
git switch -c "$branch"
# Explicit command-scoped identity also overrides inherited author variables.
GIT_AUTHOR_NAME='github-actions[bot]' \
GIT_AUTHOR_EMAIL='41898282+github-actions[bot]@users.noreply.github.com' \
GIT_COMMITTER_NAME='github-actions[bot]' \
GIT_COMMITTER_EMAIL='41898282+github-actions[bot]@users.noreply.github.com' \
  git commit -m "$message"
git push origin "HEAD:refs/heads/$branch"
echo "Open a PR and require its checks before merge: https://github.com/${GITHUB_REPOSITORY}/compare/bootstrap/core-v1...${branch}?expand=1" >> "$GITHUB_STEP_SUMMARY"
