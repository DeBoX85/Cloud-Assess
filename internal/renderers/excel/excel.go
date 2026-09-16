package excel

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/DeBoX85/Cloud-Assess/internal/branding"
	"github.com/DeBoX85/Cloud-Assess/internal/renderers/tables"
	"github.com/DeBoX85/Cloud-Assess/internal/result"
	"github.com/xuri/excelize/v2"
)

const maxInventoryRows = 1_048_566

type styleCache struct {
	title  int
	header int
	blue   int
	white  int
}

// WriteFile renders a human-oriented XLSX report from the canonical assessment result.
func WriteFile(data *result.AssessmentResult, filename string, opts tables.Options) error {
	if data == nil {
		return fmt.Errorf("assessment result is nil")
	}
	if filename == "" {
		return fmt.Errorf("Excel output filename is empty")
	}
	if err := validateInventorySize(len(data.Resources)); err != nil {
		return err
	}

	projected, err := tables.Build(data, opts)
	if err != nil {
		return err
	}
	file := excelize.NewFile()
	styles, err := createStyles(file)
	if err != nil {
		_ = file.Close()
		return err
	}

	first := true
	for _, table := range projected {
		if !tables.ShouldRender(data, table) {
			continue
		}
		if first {
			if err := file.SetSheetName("Sheet1", table.SheetName); err != nil {
				_ = file.Close()
				return fmt.Errorf("rename first Excel sheet: %w", err)
			}
			first = false
		} else if _, err := file.NewSheet(table.SheetName); err != nil {
			_ = file.Close()
			return fmt.Errorf("create Excel sheet %q: %w", table.SheetName, err)
		}
		if err := writeSheet(file, table, styles); err != nil {
			_ = file.Close()
			return err
		}
	}

	if first {
		_ = file.Close()
		return fmt.Errorf("assessment produced no Excel sheets")
	}
	if err := file.SaveAs(filename); err != nil {
		_ = file.Close()
		return fmt.Errorf("save Excel report %q: %w", filename, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close Excel report %q: %w", filename, err)
	}
	if err := os.Chmod(filename, 0o600); err != nil {
		return fmt.Errorf("set Excel report permissions %q: %w", filename, err)
	}
	return nil
}

func validateInventorySize(count int) error {
	if count > maxInventoryRows {
		return fmt.Errorf("inventory contains %d resources, exceeding the XLSX limit of %d resources", count, maxInventoryRows)
	}
	return nil
}

func createStyles(file *excelize.File) (*styleCache, error) {
	title, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Size: 16}})
	if err != nil {
		return nil, fmt.Errorf("create Excel title style: %w", err)
	}
	header, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#CAEDFB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "top", WrapText: true},
	})
	if err != nil {
		return nil, fmt.Errorf("create Excel header style: %w", err)
	}
	blue, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#EAF6FB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Vertical: "top", WrapText: true},
	})
	if err != nil {
		return nil, fmt.Errorf("create Excel alternating-row style: %w", err)
	}
	white, err := file.NewStyle(&excelize.Style{Alignment: &excelize.Alignment{Vertical: "top", WrapText: true}})
	if err != nil {
		return nil, fmt.Errorf("create Excel row style: %w", err)
	}
	return &styleCache{title: title, header: header, blue: blue, white: white}, nil
}

func writeSheet(file *excelize.File, table tables.Table, styles *styleCache) error {
	if len(table.Rows) == 0 {
		return nil
	}
	writer, err := file.NewStreamWriter(table.SheetName)
	if err != nil {
		return fmt.Errorf("create stream writer for %q: %w", table.SheetName, err)
	}

	widths := computeWidths(table.Rows, 1000)
	for i, width := range widths {
		if width < 8 {
			width = 8
		}
		if err := writer.SetColWidth(i+1, i+1, float64(width)); err != nil {
			return fmt.Errorf("set column width for %q: %w", table.SheetName, err)
		}
	}

	brand := branding.Default()
	if err := writer.SetRow("A1", []interface{}{excelize.Cell{Value: brand.ReportTitle, StyleID: styles.title}}); err != nil {
		return fmt.Errorf("write report title for %q: %w", table.SheetName, err)
	}
	if err := writer.SetRow("A2", []interface{}{excelize.Cell{Value: table.SheetName}}); err != nil {
		return fmt.Errorf("write sheet title for %q: %w", table.SheetName, err)
	}

	headers := make([]interface{}, len(table.Rows[0]))
	for i, header := range table.Rows[0] {
		headers[i] = excelize.Cell{Value: header, StyleID: styles.header}
	}
	if err := writer.SetRow("A4", headers, excelize.RowOpts{StyleID: styles.header}); err != nil {
		return fmt.Errorf("write headers for %q: %w", table.SheetName, err)
	}

	hyperlinkColumn := hyperlinkColumn(table.Key)
	lastRow := 4
	for i, row := range table.Rows[1:] {
		lastRow = i + 5
		styleID := styles.white
		if lastRow%2 == 0 {
			styleID = styles.blue
		}
		cells := make([]interface{}, len(row))
		for column, value := range row {
			if hyperlinkColumn > 0 && column == hyperlinkColumn-1 && isHTTPURL(value) {
				cells[column] = excelize.Cell{Formula: hyperlinkFormula(value), StyleID: styleID}
			} else {
				cells[column] = excelize.Cell{Value: value, StyleID: styleID}
			}
		}
		if err := writer.SetRow("A"+strconv.Itoa(lastRow), cells, excelize.RowOpts{StyleID: styleID}); err != nil {
			return fmt.Errorf("write data row for %q: %w", table.SheetName, err)
		}
	}

	if lastCell, err := excelize.CoordinatesToCellName(len(table.Rows[0]), max(4, lastRow)); err == nil {
		if err := file.AutoFilter(table.SheetName, fmt.Sprintf("A4:%s", lastCell), nil); err != nil {
			return fmt.Errorf("set autofilter for %q: %w", table.SheetName, err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush Excel sheet %q: %w", table.SheetName, err)
	}
	return nil
}

func hyperlinkColumn(key string) int {
	switch key {
	case "recommendations":
		return 11
	case "impacted":
		return 18
	case "defenderRecommendations":
		return 11
	default:
		return 0
	}
}

func hyperlinkFormula(value string) string {
	escaped := strings.ReplaceAll(value, `"`, `""`)
	return `HYPERLINK("` + escaped + `","` + escaped + `")`
}

func isHTTPURL(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://")
}

func computeWidths(records [][]string, maxSampleRows int) []int {
	if len(records) == 0 {
		return nil
	}
	if maxSampleRows <= 0 {
		maxSampleRows = 1000
	}
	widths := make([]int, len(records[0]))
	limit := min(len(records), maxSampleRows)
	for rowIndex := 0; rowIndex < limit; rowIndex++ {
		for column, value := range records[rowIndex] {
			if column >= len(widths) {
				break
			}
			width := len(value) + 3
			if width > 120 {
				width = 120
			}
			if width > widths[column] {
				widths[column] = width
			}
		}
	}
	return widths
}
