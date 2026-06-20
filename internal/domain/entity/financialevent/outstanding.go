package financialevent

import "github.com/shopspring/decimal"

type OutstandingItem struct {
	Id                int64           `json:"id"`
	EventTime         string          `json:"eventTime"`
	EventType         string          `json:"eventType"`
	Description       string          `json:"description"`
	Currency          string          `json:"currency"`
	Amount            decimal.Decimal `json:"amount"`
	OutstandingAmount decimal.Decimal `json:"outstandingAmount"`
	EventGroupId      int64           `json:"eventGroupId"`
	BucketId          int64           `json:"bucketId"`
	BucketName        string          `json:"bucketName"`
}

type OutstandingSummary struct {
	Receivables []OutstandingItem `json:"receivables"`
	Deposits    []OutstandingItem `json:"deposits"`
	LoansOut    []OutstandingItem `json:"loansOut"`
}
