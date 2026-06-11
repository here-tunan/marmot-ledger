package service

import "marmot-ledger/internal/domain/repository/ledgerentrydb"

type ExpenseRecordStrategy struct{}

func (ExpenseRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	bucket := ctx.Buckets[ctx.Request.BucketId]
	delta := signedDeltaForSingleBucket(bucket, ctx.Amount, EventTypeExpense)
	balanceAfter := bucket.Balance.Add(delta)
	if err := validateBalanceAfter(bucket, balanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeExpense, ctx.Currency, true, ctx.Category)
	entry := &ledgerentrydb.LedgerEntry{UserId: ctx.UserId, BucketId: bucket.Id, Currency: ctx.Currency, Amount: delta, BalanceAfter: balanceAfter, EntryRole: EntryRoleExpense}
	return &RecordBuildResult{Event: event, Entries: []*ledgerentrydb.LedgerEntry{entry}, BucketBalances: singleBalance(bucket.Id, balanceAfter)}, nil
}
