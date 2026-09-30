//go:build windows

package reportfile

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

// Copy the existing effective DACL before writing report bytes. Protect it from
// broader parent inheritance on replacement. New reports inherit the directory.
func preserveReplacementPrivacy(destination, staged string) error {
	if _, err := os.Lstat(destination); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	descriptor, err := windows.GetNamedSecurityInfo(destination, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return fmt.Errorf("read existing report DACL: %w", err)
	}
	dacl, _, err := descriptor.DACL()
	if err != nil {
		return fmt.Errorf("read existing report DACL: %w", err)
	}
	if dacl == nil {
		return fmt.Errorf("existing report has no DACL; refusing replacement")
	}
	err = windows.SetNamedSecurityInfo(staged, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, dacl, nil)
	runtime.KeepAlive(descriptor)
	if err != nil {
		return fmt.Errorf("preserve existing report DACL: %w", err)
	}
	return nil
}
