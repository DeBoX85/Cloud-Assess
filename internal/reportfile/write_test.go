package reportfile

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFailedRenderPreservesExistingReportAndRemovesTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")
	if err := os.WriteFile(path, []byte("previous complete report"), 0o600); err != nil {
		t.Fatal(err)
	}
	failure := errors.New("simulated output failure")
	err := Write(path, func(file *os.File) error {
		if _, err := file.WriteString("partial report"); err != nil {
			return err
		}
		return failure
	})
	if !errors.Is(err, failure) {
		t.Fatalf("error=%v, want render failure", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "previous complete report" {
		t.Fatalf("existing report changed: %q, %v", data, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary files remain: %v, %v", entries, err)
	}
}

func TestSuccessfulReplacementAndInvalidDestination(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.json")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteBytes(path, []byte("new")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "new" {
		t.Fatalf("replacement=%q error=%v", data, err)
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Fatalf("replacement must use private Unix permissions: %v %v", info, err)
		}
	}
	if err := WriteBytes(filepath.Join(path, "child"), []byte("data")); err == nil {
		t.Fatal("invalid destination reported success")
	}
}

func TestSymbolicLinkDestinationIsRejected(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows symlink creation requires environment-specific privileges")
	}
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	link := filepath.Join(dir, "report")
	if err := os.WriteFile(target, []byte("untouched"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := WriteBytes(link, []byte("new")); err == nil {
		t.Fatal("symlink destination accepted")
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "untouched" {
		t.Fatalf("symlink target modified: %q %v", data, err)
	}
}
