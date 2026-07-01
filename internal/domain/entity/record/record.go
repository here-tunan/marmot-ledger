package record

import (
	"marmot-ledger/internal/domain/entity/financialevent"

	"github.com/shopspring/decimal"
)

type SplitShare struct {
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	IsSelf      bool            `json:"isSelf"`
}

type RecordRequest struct {
	Scenario                string          `json:"scenario"`
	BucketId                int64           `json:"bucketId"`
	FromBucketId            int64           `json:"fromBucketId"`
	ToBucketId              int64           `json:"toBucketId"`
	CashBucketId            int64           `json:"cashBucketId"`
	ReceivableBucketId      int64           `json:"receivableBucketId"`
	FamilyId                int64           `json:"familyId"`
	CategoryId              int64           `json:"categoryId"`
	Amount                  decimal.Decimal `json:"amount"`
	Currency                string          `json:"currency"`
	ToAmount                decimal.Decimal `json:"toAmount"`
	ToCurrency              string          `json:"toCurrency"`
	RemainingMarketValue    decimal.Decimal `json:"remainingMarketValue"`
	Shares                  []SplitShare    `json:"shares"`
	Description             string          `json:"description"`
	EventTime               string          `json:"eventTime"`
	RelatedFinancialEventId int64           `json:"relatedFinancialEventId"`
	ChannelId               int64           `json:"channelId"`
	Remark                  string          `json:"remark"`
	Source                  string          `json:"source"` // 空则回退 manual；"import" 表示批量导入
}

type RecordResponse struct {
	FinancialEvent financialevent.FinancialEvent `json:"financialEvent"`
}
