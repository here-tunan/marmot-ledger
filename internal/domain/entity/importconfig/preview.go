package importconfig

// PreviewRow 一行数据按配置映射后的结果。Errors 为空且 Filtered=false 表示该行可入账。
type PreviewRow struct {
	RowIndex     int      `json:"rowIndex"`
	Scenario     string   `json:"scenario"`
	Amount       string   `json:"amount"` // 规范化字符串，便于可编辑展示
	Currency     string   `json:"currency"`
	EventTime    string   `json:"eventTime"` // 归一化为 "2006-01-02 15:04:05"
	CategoryId   int64    `json:"categoryId"`
	ChannelId    int64    `json:"channelId"`
	BucketId     int64    `json:"bucketId"`   // 单桶场景 = 资金桶；transfer = 转出桶
	ToBucketId   int64    `json:"toBucketId"` // 仅 transfer 使用
	Description  string   `json:"description"`
	Remark       string   `json:"remark"`
	Errors       []string `json:"errors"`
	Filtered     bool     `json:"filtered"` // 命中 drop 规则被过滤掉
	FilterReason string   `json:"filterReason"`
}

// ImportPreview 解析 + 规则映射后的整体预览结果。
type ImportPreview struct {
	Headers          []string     `json:"headers"`
	Rows             []PreviewRow `json:"rows"`
	Warnings         []string     `json:"warnings"`
	TotalRows        int          `json:"totalRows"`
	ErrorRowCount    int          `json:"errorRowCount"`
	FilteredRowCount int          `json:"filteredRowCount"`
}
