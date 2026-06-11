package financialevent

import (
	"marmot-ledger/internal/domain/entity/ledgerentry"

	"github.com/shopspring/decimal"
)

type FinancialEvent struct {
	Id                      int64                     `json:"id"`
	UserId                  int64                     `json:"userId"`
	EventGroupId            int64                     `json:"eventGroupId"`
	RelatedFinancialEventId int64                     `json:"relatedFinancialEventId"`
	EventType               string                    `json:"eventType"`
	Description             string                    `json:"description"`
	CategoryId              int64                     `json:"categoryId"`
	CategoryGroupId         int64                     `json:"categoryGroupId"`
	ChannelType             string                    `json:"channelType"`
	ChannelAccountId        int64                     `json:"channelAccountId"`
	EventTime               string                    `json:"eventTime"`
	Currency                string                    `json:"currency"`
	Amount                  decimal.Decimal           `json:"amount"`
	BaseCurrency            string                    `json:"baseCurrency"`
	BaseAmount              decimal.Decimal           `json:"baseAmount"`
	ExchangeRate            decimal.Decimal           `json:"exchangeRate"`
	IncludeInStatistics     bool                      `json:"includeInStatistics"`
	Source                  string                    `json:"source"`
	Status                  string                    `json:"status"`
	Remark                  string                    `json:"remark"`
	LedgerEntries           []ledgerentry.LedgerEntry `json:"ledgerEntries,omitempty"`
}

type FinancialEventQuery struct {
	EventType           string `json:"eventType"`
	StartTime           string `json:"startTime"`
	EndTime             string `json:"endTime"`
	CategoryId          int64  `json:"categoryId"`
	CategoryGroupId     int64  `json:"categoryGroupId"`
	BucketId            int64  `json:"bucketId"`
	Keyword             string `json:"keyword"`
	IncludeInStatistics *bool  `json:"includeInStatistics"`
	Page                int    `json:"page"`
	PageSize            int    `json:"pageSize"`
}
