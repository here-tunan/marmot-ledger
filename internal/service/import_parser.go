package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"marmot-ledger/internal/domain/entity/importconfig"
	"strings"

	"github.com/xuri/excelize/v2"
)

// MaxImportPreviewRows 单次解析返回的最大数据行数；超出截断并在 warning 提示。
const MaxImportPreviewRows = 10000

// ParseImportTable 解析上传文件，按配置定位表头行，返回表头与数据行（按列序对齐）。
// truncated 为 true 表示文件行数超过 MaxImportPreviewRows，只返回前 N 行。
func ParseImportTable(config *importconfig.ImportConfig, data []byte) (headerRow []string, dataRows [][]string, truncated bool, err error) {
	fileType := strings.ToLower(strings.TrimSpace(config.FileType))
	switch fileType {
	case "csv":
		return parseCsv(config, data)
	case "xlsx":
		return parseXlsx(config, data)
	default:
		return nil, nil, false, fmt.Errorf("unsupported file type: %s", fileType)
	}
}

func parseCsv(config *importconfig.ImportConfig, data []byte) (headerRow []string, dataRows [][]string, truncated bool, err error) {
	// utf-8 BOM 容错
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	reader := csv.NewReader(bytes.NewReader(data))
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // 允许列数不齐，后续按表头长度对齐
	allRows, err := reader.ReadAll()
	if err != nil {
		return nil, nil, false, fmt.Errorf("read csv failed: %w", err)
	}
	return splitHeaderAndData(config, allRows)
}

func parseXlsx(config *importconfig.ImportConfig, data []byte) (headerRow []string, dataRows [][]string, truncated bool, err error) {
	file, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, nil, false, fmt.Errorf("open xlsx failed: %w", err)
	}
	defer file.Close()

	sheet := strings.TrimSpace(config.SheetName)
	if sheet == "" {
		sheets := file.GetSheetList()
		if len(sheets) == 0 {
			return nil, nil, false, fmt.Errorf("xlsx has no sheets")
		}
		sheet = sheets[0]
	}
	rows, err := file.GetRows(sheet)
	if err != nil {
		return nil, nil, false, fmt.Errorf("read xlsx sheet %q failed: %w", sheet, err)
	}
	return splitHeaderAndData(config, rows)
}

// splitHeaderAndData 按 config.HeaderRow（1-based）切分表头行与数据行，空行剔除，按 MaxImportPreviewRows 截断。
func splitHeaderAndData(config *importconfig.ImportConfig, rows [][]string) (headerRow []string, dataRows [][]string, truncated bool, err error) {
	headerRowNum := config.HeaderRow
	if headerRowNum < 1 {
		headerRowNum = 1
	}
	if len(rows) < headerRowNum {
		return nil, nil, false, fmt.Errorf("header row %d not found, file only has %d rows", headerRowNum, len(rows))
	}
	headerRow = rows[headerRowNum-1]
	for i := headerRowNum; i < len(rows); i++ {
		if isBlankRow(rows[i]) {
			continue
		}
		if len(dataRows) >= MaxImportPreviewRows {
			truncated = true
			break
		}
		dataRows = append(dataRows, rows[i])
	}
	return headerRow, dataRows, truncated, nil
}

func isBlankRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
