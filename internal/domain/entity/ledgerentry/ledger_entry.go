package ledgerentry

import "github.com/shopspring/decimal"

type LedgerEntry struct {
	Id               int64           `json:"id"`
	FinancialEventId int64           `json:"financialEventId"`
	UserId           int64           `json:"userId"`
	BucketId         int64           `json:"bucketId"`
	Currency         string          `json:"currency"`
	Amount           decimal.Decimal `json:"amount"`
	BalanceAfter     decimal.Decimal `json:"balanceAfter"`
	EntryRole        string          `json:"entryRole"`
	CreatedAt        string          `json:"createdAt"`
}
