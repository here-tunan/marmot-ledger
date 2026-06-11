package service

import "marmot-ledger/internal/domain/repository/ledgerentrydb"

type RefundRecordStrategy struct{}

func (RefundRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	bucket := ctx.Buckets[ctx.Request.BucketId]
	delta := signedDeltaForSingleBucket(bucket, ctx.Amount, EventTypeRefund)
	balanceAfter := bucket.Balance.Add(delta)
	if err := validateBalanceAfter(bucket, balanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeRefund, ctx.Currency, true, ctx.Category)
	entry := &ledgerentrydb.LedgerEntry{UserId: ctx.UserId, BucketId: bucket.Id, Currency: ctx.Currency, Amount: delta, BalanceAfter: balanceAfter, EntryRole: EntryRoleRefund}
	return &RecordBuildResult{Event: event, Entries: []*ledgerentrydb.LedgerEntry{entry}, BucketBalances: singleBalance(bucket.Id, balanceAfter)}, nil
}
