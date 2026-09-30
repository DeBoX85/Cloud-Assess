// Package reportfile stages each report before replacing its destination.
package reportfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// Write preserves an existing report if rendering or closing the staged file fails.
// Replacement uses os.Rename; atomic replacement is not promised on every platform.
// Unix files are private; Windows access follows the destination directory's ACL.
// A multi-file report is not a transaction. Callers must retain successful paths.
func Write(filename string, render func(*os.File) error) error {
	if info, err := os.Lstat(filename); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("report destination is a symbolic link: %s", filename)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(filename), ".cloud-assess-report-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if err := preserveReplacementPrivacy(filename, file.Name()); err != nil {
		return err
	}
	if err := render(file); err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), filename)
}

func WriteBytes(filename string, data []byte) error {
	return Write(filename, func(file *os.File) error {
		_, err := file.Write(data)
		return err
	})
}
