package service

import (
	"errors"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
)

type BalanceAdjustmentRecordStrategy struct{}

func (BalanceAdjustmentRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	if ctx.Amount.IsZero() {
		return nil, errors.New("adjustment amount must not be 0")
	}

	bucket := ctx.Buckets[ctx.Request.BucketId]
	balanceAfter := bucket.Balance.Add(ctx.Amount)
	if err := validateBalanceAfter(bucket, balanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeBalanceAdjustment, ctx.Currency, false, nil)
	entry := &ledgerentrydb.LedgerEntry{UserId: ctx.UserId, BucketId: bucket.Id, Currency: ctx.Currency, Amount: ctx.Amount, BalanceAfter: balanceAfter, EntryRole: EntryRoleAdjustment}
	return &RecordBuildResult{Event: event, Entries: []*ledgerentrydb.LedgerEntry{entry}, BucketBalances: singleBalance(bucket.Id, balanceAfter)}, nil
}
