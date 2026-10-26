package excel

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

// Child isolation bounds the OLD dependency's fatal stack overflow. It does not
// add workbook-input support to the product or change runtime stack settings.
func securityDependencyChild(t *testing.T, name string) bool {
	t.Helper()
	if os.Getenv("CLOUD_ASSESS_EXCELIZE_SECURITY_CHILD") == name {
		debug.SetMaxStack(1 << 20)
		return true
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^"+name+"$", "-test.count=1")
	cmd.Env = append(os.Environ(), "CLOUD_ASSESS_EXCELIZE_SECURITY_CHILD="+name)
	output, e := cmd.CombinedOutput()
	if e != nil || ctx.Err() != nil || !strings.Contains(string(output), "EXCELIZE_SECURITY_COMPLETED "+name) {
		if len(output) > 8192 {
			output = output[:8192]
		}
		t.Fatalf("dependency safety regression: %s: %v / %v\n%s", name, e, ctx.Err(), output)
	}
	return false
}
func TestExcelizeCircularArraySafety(t *testing.T) {
	const name = "TestExcelizeCircularArraySafety"
	if !securityDependencyChild(t, name) {
		return
	}
	f := excelize.NewFile()
	defer f.Close()
	typ, first, second := excelize.STCellFormulaTypeArray, "A1:A1", "B1:B1"
	if e := f.SetCellFormula("Sheet1", "A1", "_xlfn.ANCHORARRAY($B$1)", excelize.FormulaOpts{Ref: &first, Type: &typ}); e != nil {
		t.Fatal(e)
	}
	if e := f.SetCellFormula("Sheet1", "B1", "_xlfn.ANCHORARRAY($A$1)", excelize.FormulaOpts{Ref: &second, Type: &typ}); e != nil {
		t.Fatal(e)
	}
	got, e := f.CalcCellValue("Sheet1", "A1")
	if e != nil || got != "" {
		t.Fatal("circular array must terminate safely", got, e)
	}
	os.Stdout.WriteString("EXCELIZE_SECURITY_COMPLETED " + name + "\n")
}
func TestExcelizeMalformedEncryptionSafety(t *testing.T) {
	const name = "TestExcelizeMalformedEncryptionSafety"
	if !securityDependencyChild(t, name) {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("malformed encryption header caused panic: %v", p)
		}
	}()
	plain := bytes.Repeat([]byte("0123456789abcdef"), 512)
	opts := &excelize.Options{Password: "synthetic-fixture-password"}
	encrypted, e := excelize.Encrypt(plain, opts)
	if e != nil {
		t.Fatal(e)
	}
	original, e := excelize.Decrypt(encrypted, opts)
	if e != nil || !bytes.Equal(original, plain) {
		t.Fatalf("valid synthetic encryption roundtrip: err=%v got=%x want=%x", e, original, plain)
	}
	// ECMA standard version3.2, flags0x24, header length0xa4. Exact bytes
	// originate in the independently reviewed upstream standard writer.
	prefix := []byte{3, 0, 2, 0, 0x24, 0, 0, 0, 0xa4, 0, 0, 0}
	if bytes.Count(encrypted, prefix) != 1 {
		t.Fatal("synthetic standard header not uniquely identified")
	}
	offset := bytes.Index(encrypted, prefix)
	binary.LittleEndian.PutUint32(encrypted[offset+8:offset+12], 0xffffffff)
	got, e := excelize.Decrypt(encrypted, opts)
	if got != nil || !errors.Is(e, excelize.ErrWorkbookFileFormat) {
		t.Fatal("malformed header must fail closed", len(got), e)
	}
	os.Stdout.WriteString("EXCELIZE_SECURITY_COMPLETED " + name + "\n")
}

// Numeric-looking report values must remain exact strings despite upstream
// numeric precision changes; Unicode and formula-looking input stay data.
func TestExcelizeSecurityUpdateReportCompatibility(t *testing.T) {
	for _, value := range []string{"000012345678901234567890.1234500", "Norway æøå 日本語 😀", "=1+1", "+1+1", "@SUM(1,2)"} {
		t.Run(value, func(t *testing.T) {
			data := excelFixture()
			data.Resources[0].Name = value
			filename := filepath.Join(t.TempDir(), "security-compat.xlsx")
			if e := WriteFile(data, filename, tables.Options{}); e != nil {
				t.Fatal(e)
			}
			f, e := excelize.OpenFile(filename)
			if e != nil {
				t.Fatal(e)
			}
			defer f.Close()
			got, e := f.GetCellValue("Inventory", "E5", excelize.Options{RawCellValue: true})
			if e != nil || got != value {
				t.Fatal("report text/precision changed", got, value, e)
			}
			formula, e := f.GetCellFormula("Inventory", "E5")
			if e != nil || formula != "" {
				t.Fatal("report value became formula", formula, e)
			}
			typ, e := f.GetCellType("Inventory", "E5")
			if e != nil || typ != excelize.CellTypeInlineString {
				t.Fatal("report value no longer stored as text", typ, e)
			}
			style, e := f.GetCellStyle("Inventory", "E5")
			if e != nil || style == 0 {
				t.Fatal("report style lost", style, e)
			}
			header, e := f.GetCellValue("Inventory", "E4")
			if e != nil || header != "Resource Name" {
				t.Fatal("report header changed", header, e)
			}
		})
	}
}
