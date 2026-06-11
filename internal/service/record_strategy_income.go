package service

import "marmot-ledger/internal/domain/repository/ledgerentrydb"

type IncomeRecordStrategy struct{}

func (IncomeRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	bucket := ctx.Buckets[ctx.Request.BucketId]
	delta := ctx.Amount
	balanceAfter := bucket.Balance.Add(delta)
	if err := validateBalanceAfter(bucket, balanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeIncome, ctx.Currency, true, ctx.Category)
	entry := &ledgerentrydb.LedgerEntry{UserId: ctx.UserId, BucketId: bucket.Id, Currency: ctx.Currency, Amount: delta, BalanceAfter: balanceAfter, EntryRole: EntryRoleIncome}
	return &RecordBuildResult{Event: event, Entries: []*ledgerentrydb.LedgerEntry{entry}, BucketBalances: singleBalance(bucket.Id, balanceAfter)}, nil
}
