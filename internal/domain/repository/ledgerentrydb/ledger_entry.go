package ledgerentrydb

import (
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

type LedgerEntry struct {
	Id               int64           `json:"id" xorm:"pk autoincr 'id'"`
	FinancialEventId int64           `json:"financialEventId" xorm:"'financial_event_id'"`
	UserId           int64           `json:"userId" xorm:"'user_id'"`
	BucketId         int64           `json:"bucketId" xorm:"'bucket_id'"`
	Currency         string          `json:"currency" xorm:"'currency'"`
	Amount           decimal.Decimal `json:"amount" xorm:"'amount'"`
	BalanceAfter     decimal.Decimal `json:"balanceAfter" xorm:"'balance_after'"`
	EntryRole        string          `json:"entryRole" xorm:"'entry_role'"`
	CreatedAt        model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
}

func (LedgerEntry) TableName() string {
	return "ledger_entry"
}

func InsertLedgerEntry(session *xorm.Session, entry *LedgerEntry) error {
	_, err := session.InsertOne(entry)
	return err
}

func ListLedgerEntriesByBucket(bucketId int64, userId int64) ([]LedgerEntry, error) {
	entries := make([]LedgerEntry, 0)
	err := infrastructure.Mysql.
		Where("bucket_id = ? AND user_id = ?", bucketId, userId).
		Desc("created_at", "id").
		Find(&entries)
	return entries, err
}

func ListLedgerEntriesByEvent(financialEventId int64, userId int64) ([]LedgerEntry, error) {
	entries := make([]LedgerEntry, 0)
	err := infrastructure.Mysql.
		Where("financial_event_id = ? AND user_id = ?", financialEventId, userId).
		Asc("id").
		Find(&entries)
	return entries, err
}
