package record

import (
	"marmot-ledger/internal/domain/entity/financialevent"

	"github.com/shopspring/decimal"
)

type RecordRequest struct {
	Scenario                string          `json:"scenario"`
	BucketId                int64           `json:"bucketId"`
	FromBucketId            int64           `json:"fromBucketId"`
	ToBucketId              int64           `json:"toBucketId"`
	Amount                  decimal.Decimal `json:"amount"`
	Currency                string          `json:"currency"`
	Description             string          `json:"description"`
	EventTime               string          `json:"eventTime"`
	RelatedFinancialEventId int64           `json:"relatedFinancialEventId"`
	Remark                  string          `json:"remark"`
}

type RecordResponse struct {
	FinancialEvent financialevent.FinancialEvent `json:"financialEvent"`
}
