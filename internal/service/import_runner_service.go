package service

import (
	"fmt"
	"marmot-ledger/internal/domain/entity/importconfig"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

var importPreviewScenarios = map[string]bool{
	"income": true, "expense": true, "refund": true, "transfer": true,
}

// BuildPreviewRows 把解析出的原始行按配置映射成预览行，逐行收集错误但不中断。
func BuildPreviewRows(config *importconfig.ImportConfig, headerRow []string, dataRows [][]string, defaultBucketId int64) (*importconfig.ImportPreview, error) {
	headerIndex := buildHeaderIndex(headerRow)
	mappingByField := map[string]importconfig.FieldMapping{}
	for _, m := range config.Mappings {
		mappingByField[m.TargetField] = m
	}

	preview := &importconfig.ImportPreview{
		Headers: cloneStrings(headerRow),
		Rows:    make([]importconfig.PreviewRow, 0, len(dataRows)),
	}
	preview.Warnings = collectWarnings(config, headerIndex)

	for i, row := range dataRows {
		cellByColumn := buildCellLookup(headerRow, row)
		previewRow := buildPreviewRow(mappingByField, cellByColumn, i, defaultBucketId)
		applyFilters(config.Filters, cellByColumn, &previewRow)
		preview.Rows = append(preview.Rows, previewRow)
	}
	preview.TotalRows = len(preview.Rows)
	for _, r := range preview.Rows {
		if r.Filtered {
			preview.FilteredRowCount++
			continue
		}
		if len(r.Errors) > 0 {
			preview.ErrorRowCount++
		}
	}
	return preview, nil
}

// applyFilters 按顺序首条命中决定丢/留；都不命中默认保留。
func applyFilters(filters []importconfig.FilterRule, cellByColumn map[string]string, row *importconfig.PreviewRow) {
	for _, f := range filters {
		cell := strings.TrimSpace(cellByColumn[strings.TrimSpace(f.MatchColumn)])
		if cell == "" {
			continue
		}
		if !matchRule(f.Operator, cell, f.MatchValue) {
			continue
		}
		if f.Action == "drop" {
			row.Filtered = true
			row.FilterReason = fmt.Sprintf("列「%s」%s「%s」", f.MatchColumn, operatorLabel(f.Operator), f.MatchValue)
		}
		return // keep 或 drop 都终止评估
	}
}

func operatorLabel(op string) string {
	switch op {
	case "contains":
		return "包含"
	case "equals":
		return "等于"
	case "notContains":
		return "不包含"
	case "notEquals":
		return "不等于"
	case "containsAny":
		return "包含任意"
	case "notContainsAny":
		return "都不包含"
	case "equalsAny":
		return "等于任意"
	case "notEqualsAny":
		return "都不等于"
	default:
		return op
	}
}

func buildHeaderIndex(headerRow []string) map[string]int {
	index := map[string]int{}
	for i, cell := range headerRow {
		name := strings.TrimSpace(cell)
		if name == "" {
			continue
		}
		if _, exists := index[name]; !exists {
			index[name] = i
		}
	}
	return index
}

// buildCellLookup 把一行单元格按列名索引（取该列名首次出现的列号）。
func buildCellLookup(headerRow []string, row []string) map[string]string {
	lookup := map[string]string{}
	for i, name := range headerRow {
		key := strings.TrimSpace(name)
		if key == "" {
			continue
		}
		if _, exists := lookup[key]; exists {
			continue
		}
		if i < len(row) {
			lookup[key] = row[i]
		} else {
			lookup[key] = ""
		}
	}
	return lookup
}

func collectWarnings(config *importconfig.ImportConfig, headerIndex map[string]int) []string {
	seen := map[string]bool{}
	var warnings []string
	add := func(column, field string) {
		column = strings.TrimSpace(column)
		if column == "" {
			return
		}
		if _, ok := headerIndex[column]; ok {
			return
		}
		key := field + "|" + column
		if seen[key] {
			return
		}
		seen[key] = true
		warnings = append(warnings, fmt.Sprintf("配置字段 %s 引用了文件中不存在的列「%s」", field, column))
	}
	for _, m := range config.Mappings {
		add(m.SourceColumn, m.TargetField)
		for _, r := range m.Rules {
			add(r.MatchColumn, m.TargetField+"(rule)")
		}
	}
	for _, f := range config.Filters {
		add(f.MatchColumn, "filter")
	}
	return warnings
}

func buildPreviewRow(mappingByField map[string]importconfig.FieldMapping, cellByColumn map[string]string, rowIndex int, defaultBucketId int64) importconfig.PreviewRow {
	row := importconfig.PreviewRow{RowIndex: rowIndex + 1, Errors: []string{}}

	// 直接映射字段
	row.Amount = strings.TrimSpace(cellByColumn[strings.TrimSpace(mappingByField["amount"].SourceColumn)])
	row.Description = strings.TrimSpace(cellByColumn[strings.TrimSpace(mappingByField["description"].SourceColumn)])
	row.Remark = strings.TrimSpace(cellByColumn[strings.TrimSpace(mappingByField["remark"].SourceColumn)])
	row.EventTime = strings.TrimSpace(cellByColumn[strings.TrimSpace(mappingByField["eventTime"].SourceColumn)])

	// 规则字段
	row.Currency = resolveFieldValue(mappingByField["currency"], cellByColumn)
	row.Scenario = resolveFieldValue(mappingByField["scenario"], cellByColumn)
	if row.Scenario == "" {
		row.Scenario = "expense"
	}
	row.CategoryId = parseInt64Value(mappingByField["category"], cellByColumn, "category", &row)
	row.ChannelId = parseInt64Value(mappingByField["channel"], cellByColumn, "channel", &row)

	bucketRaw := resolveFieldValue(mappingByField["bucket"], cellByColumn)
	bucketId, bucketErr := strconv.ParseInt(strings.TrimSpace(bucketRaw), 10, 64)
	if bucketErr != nil || bucketId == 0 {
		bucketId = defaultBucketId
	}
	row.BucketId = bucketId
	if row.Scenario == "transfer" {
		if bucketId == 0 {
			row.Errors = append(row.Errors, "未指定转出资金桶")
		}
		// 转入桶不参与规则匹配，由用户在预览手动选
		if row.ToBucketId == 0 {
			row.Errors = append(row.Errors, "未指定转入资金桶")
		}
	} else if bucketId == 0 {
		row.Errors = append(row.Errors, "未指定资金桶")
	}

	// amount 校验 + 规范化
	if row.Amount == "" {
		row.Errors = append(row.Errors, "金额为空")
	} else {
		normalized, ok := parseAmount(row.Amount)
		if !ok {
			row.Errors = append(row.Errors, fmt.Sprintf("金额无法解析「%s」", row.Amount))
		} else {
			row.Amount = normalized
		}
	}

	// eventTime 校验 + 归一化
	if row.EventTime == "" {
		row.Errors = append(row.Errors, "交易时间为空")
	} else {
		normalized, ok := parseEventTime(row.EventTime)
		if !ok {
			row.Errors = append(row.Errors, fmt.Sprintf("交易时间无法解析「%s」", row.EventTime))
		} else {
			row.EventTime = normalized
		}
	}

	// scenario 合法性
	if !importPreviewScenarios[row.Scenario] {
		row.Errors = append(row.Errors, fmt.Sprintf("收支类型不合法「%s」", row.Scenario))
	}

	// currency 兜底
	if row.Currency == "" {
		row.Currency = "CNY"
	}

	return row
}

// resolveFieldValue 按规则顺序取第一个命中的 resultValue，都不命中用 defaultValue。
func resolveFieldValue(mapping importconfig.FieldMapping, cellByColumn map[string]string) string {
	for _, rule := range mapping.Rules {
		cell := strings.TrimSpace(cellByColumn[strings.TrimSpace(rule.MatchColumn)])
		if cell == "" {
			continue // 空单元格不参与规则匹配
		}
		if matchRule(rule.Operator, cell, rule.MatchValue) {
			return rule.ResultValue
		}
	}
	return mapping.DefaultValue
}

func matchRule(operator, cell, matchValue string) bool {
	cell = strings.TrimSpace(cell)
	values := strings.Split(matchValue, ",")
	for i := range values {
		values[i] = strings.TrimSpace(values[i])
	}
	switch operator {
	case "contains":
		return strings.Contains(cell, matchValue)
	case "equals":
		return cell == matchValue
	case "notContains":
		return !strings.Contains(cell, matchValue)
	case "notEquals":
		return cell != matchValue
	case "containsAny":
		for _, v := range values {
			if v == "" {
				continue
			}
			if strings.Contains(cell, v) {
				return true
			}
		}
		return false
	case "notContainsAny":
		for _, v := range values {
			if v == "" {
				continue
			}
			if strings.Contains(cell, v) {
				return false
			}
		}
		return true
	case "equalsAny":
		for _, v := range values {
			if v == "" {
				continue
			}
			if cell == v {
				return true
			}
		}
		return false
	case "notEqualsAny":
		for _, v := range values {
			if v == "" {
				continue
			}
			if cell == v {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func parseInt64Value(mapping importconfig.FieldMapping, cellByColumn map[string]string, label string, row *importconfig.PreviewRow) int64 {
	raw := resolveFieldValue(mapping, cellByColumn)
	id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
	if err != nil {
		if strings.TrimSpace(raw) == "" {
			return 0 // 空结果视为未指定，不算错误
		}
		row.Errors = append(row.Errors, fmt.Sprintf("%s id 无法解析「%s」", label, raw))
		return 0
	}
	return id
}

// parseAmount 清洗金额文本（去货币符号/千分位逗号/空格/全角），再解析为 decimal，返回规范字符串。
func parseAmount(raw string) (string, bool) {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, "，", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "¥", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, "€", "")
	s = strings.ReplaceAll(s, " ", "")
	d, err := decimal.NewFromString(s)
	if err != nil {
		return "", false
	}
	return d.StringFixed(2), true
}

// parseEventTime 尝试多种常见日期布局 + Excel 序列号，成功归一化为 "2006-01-02 15:04:05"。
func parseEventTime(raw string) (string, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", false
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"2006年01月02日 15:04:05",
		"2006年01月02日 15:04",
		"2006年01月02日",
		"2006.01.02 15:04:05",
		"2006.01.02 15:04",
		"2006.01.02",
		"20060102",
		"2006/1/2 15:04",
		"2006/1/2",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return parsed.Format("2006-01-02 15:04:05"), true
		}
	}
	// Excel 序列号（纯整数，自 1899-12-30 起的天数）
	if n, err := strconv.Atoi(s); err == nil && n > 0 && n < 2000000 {
		base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)
		return base.AddDate(0, 0, n).Format("2006-01-02 15:04:05"), true
	}
	return "", false
}

func cloneStrings(src []string) []string {
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
