// brand-build creates a development executable with immutable presentation metadata.
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/DeBoX85/Cloud-Assess/internal/branding"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("brand-build", flag.ContinueOnError)
	profilePath := flags.String("profile", "", "Versioned JSON profile; omitted uses built-in defaults")
	output := flags.String("output", ".build/branding", "New executable's parent directory")
	version := flags.String("version", "dev", "Development candidate version")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected builder arguments")
	}
	if !regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,63}$`).MatchString(*version) {
		return fmt.Errorf("invalid version label")
	}
	if runtime.GOARCH != "amd64" || (runtime.GOOS != "linux" && runtime.GOOS != "windows") {
		return fmt.Errorf("builder supports native Linux/Windows amd64 only")
	}
	data := []byte(`{"schemaVersion":1}`)
	if *profilePath != "" {
		f, err := os.Open(*profilePath)
		if err != nil {
			return fmt.Errorf("open profile: %w", err)
		}
		data, err = io.ReadAll(io.LimitReader(f, branding.MaxProfileBytes+1))
		closeErr := f.Close()
		if err != nil {
			return fmt.Errorf("read profile: %w", err)
		}
		if closeErr != nil {
			return closeErr
		}
	}
	profile, err := branding.ParseProfile(data)
	if err != nil {
		return err
	}
	canonical, err := profile.Canonical()
	if err != nil {
		return err
	}
	// Explicit argv and base64url avoid shell/linker quoting of display text.
	encoded := base64.RawURLEncoding.EncodeToString(canonical)
	name := profile.CLIName
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	destination := filepath.Join(*output, name)
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return fmt.Errorf("refusing to replace existing executable")
	}
	if err := os.MkdirAll(*output, 0700); err != nil {
		return err
	}
	staging, err := os.MkdirTemp(*output, ".brand-build-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(staging)
	binary := filepath.Join(staging, name)
	cmd := exec.Command("go", "build", "-trimpath", "-mod=readonly", "-buildvcs=true", "-ldflags", "-X main.version="+*version+" -X github.com/DeBoX85/Cloud-Assess/internal/branding.embeddedProfile="+encoded, "-o", binary, "./cmd/cloud-assess")
	// Build for the native target regardless of inherited cross-compilation settings.
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS="+runtime.GOOS, "GOARCH="+runtime.GOARCH)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build branded executable: %w", err)
	}
	// Same-filesystem exclusive publication never overwrites another builder's output.
	if err := os.Link(binary, destination); err != nil {
		return fmt.Errorf("publish executable: %w", err)
	}
	fmt.Println(destination)
	return nil
}
