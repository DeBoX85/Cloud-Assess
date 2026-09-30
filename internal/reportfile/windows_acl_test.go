package reportfile

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWindowsReportInheritsControlledDirectoryACL(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows ACL assertion")
	}
	dir := t.TempDir()
	run := func(script string) {
		t.Helper()
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "$ErrorActionPreference='Stop'; "+script)
		// Windows PowerShell must build its own module path, not inherit pwsh modules.
		for _, entry := range os.Environ() {
			if !strings.HasPrefix(strings.ToUpper(entry), "PSMODULEPATH=") {
				cmd.Env = append(cmd.Env, entry)
			}
		}
		cmd.Env = append(cmd.Env, "CLOUD_ASSESS_TEST_ACL_DIR="+dir)
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("ACL fixture failed: %v %s", err, output)
		}
	}
	run(`$path=$env:CLOUD_ASSESS_TEST_ACL_DIR; $sid=[System.Security.Principal.WindowsIdentity]::GetCurrent().User; $acl=Get-Acl -LiteralPath $path; $acl.SetAccessRuleProtection($true,$false); $rule=[System.Security.AccessControl.FileSystemAccessRule]::new($sid,'FullControl','ContainerInherit,ObjectInherit','None','Allow'); $acl.AddAccessRule($rule); Set-Acl -LiteralPath $path -AclObject $acl`)
	if err := WriteBytes(filepath.Join(dir, "report.json"), []byte("synthetic")); err != nil {
		t.Fatal(err)
	}
	run(`$sid=[System.Security.Principal.WindowsIdentity]::GetCurrent().User.Value; $acl=Get-Acl -LiteralPath (Join-Path $env:CLOUD_ASSESS_TEST_ACL_DIR 'report.json'); $allowed=@($acl.Access | Where-Object { $_.AccessControlType -eq 'Allow' }); if ($allowed.Count -eq 0) { throw 'No allow entries' }; foreach ($entry in $allowed) { if ($entry.IdentityReference.Translate([System.Security.Principal.SecurityIdentifier]).Value -ne $sid) { throw 'Unexpected principal can access report' } }`)
}
