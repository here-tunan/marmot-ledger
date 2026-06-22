package bucket

import "github.com/shopspring/decimal"

type Bucket struct {
	Id                      int64           `json:"id"`
	UserId                  int64           `json:"userId"`
	AccountId               int64           `json:"accountId"`
	Name                    string          `json:"name"`
	Currency                string          `json:"currency"`
	Balance                 decimal.Decimal `json:"balance"`
	InitialBalance          decimal.Decimal `json:"initialBalance"`
	BucketType              string          `json:"bucketType"`
	BucketNature            string          `json:"bucketNature"`
	IsActive                bool            `json:"isActive"`
	InitialFinancialEventId int64           `json:"initialFinancialEventId,omitempty"`
	InitialLedgerEntryId    int64           `json:"initialLedgerEntryId,omitempty"`
}

type BucketQuery struct {
	AccountId    int64
	Currency     string
	BucketType   string
	BucketNature string
	IsActive     *bool
}
