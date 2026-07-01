package importconfig

// ImportConfig 表格导入配置：描述某个平台的表格如何映射到财务记录字段。
type ImportConfig struct {
	Id        int64          `json:"id"`
	UserId    int64          `json:"userId"`
	Name      string         `json:"name"`
	FileType  string         `json:"fileType"`  // xlsx | csv
	SheetName string         `json:"sheetName"` // xlsx 工作表名，空=第一个
	HeaderRow int            `json:"headerRow"` // 表头所在行，1-based
	Icon      string         `json:"icon"`
	Sort      int            `json:"sort"`
	IsActive  bool           `json:"isActive"`
	Mappings  []FieldMapping `json:"mappings"`
	Filters   []FilterRule   `json:"filters"` // 行过滤规则，按顺序首条命中决定保留/丢弃
}

// FilterRule 行过滤规则：某列内容满足 operator/matchValue 时执行 action。
type FilterRule struct {
	MatchColumn string `json:"matchColumn"`
	Operator    string `json:"operator"` // contains|equals|notContains|notEquals|containsAny|notContainsAny
	MatchValue  string `json:"matchValue"`
	Action      string `json:"action"` // drop | keep
}

// FieldMapping 单个目标字段的映射配置。
// 取值优先级：rules 按顺序取第一个命中的 resultValue，都不命中则用 defaultValue；
// 直接映射类字段（amount/description/remark/eventTime/currency）通常只设 sourceColumn。
type FieldMapping struct {
	TargetField  string      `json:"targetField"`  // amount|description|remark|eventTime|currency|scenario|category|channel|bucket
	SourceColumn string      `json:"sourceColumn"` // 表头列名
	DefaultValue string      `json:"defaultValue"`
	Rules        []FieldRule `json:"rules"`
}

// FieldRule 条件规则：当 matchColumn 列内容满足 operator/matchValue 时，命中 resultValue。
type FieldRule struct {
	MatchColumn string `json:"matchColumn"`
	Operator    string `json:"operator"` // contains|equals|notContains|notEquals
	MatchValue  string `json:"matchValue"`
	ResultValue string `json:"resultValue"`
}

type ImportConfigQuery struct {
	IsActive *bool
}
