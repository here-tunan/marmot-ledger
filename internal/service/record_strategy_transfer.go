package service

import (
	"errors"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
)

type TransferRecordStrategy struct{}

func (TransferRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	fromBucket := ctx.Buckets[ctx.Request.FromBucketId]
	toBucket := ctx.Buckets[ctx.Request.ToBucketId]
	if fromBucket.Id == toBucket.Id {
		return nil, errors.New("transfer buckets must be different")
	}
	fromDelta := ctx.Amount.Neg()
	toDelta := ctx.Amount
	fromBalanceAfter := fromBucket.Balance.Add(fromDelta)
	toBalanceAfter := toBucket.Balance.Add(toDelta)
	if err := validateBalanceAfter(fromBucket, fromBalanceAfter); err != nil {
		return nil, err
	}
	if err := validateBalanceAfter(toBucket, toBalanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeTransfer, ctx.Currency, false, nil)
	fromEntry := &ledgerentrydb.LedgerEntry{UserId: ctx.UserId, BucketId: fromBucket.Id, Currency: ctx.Currency, Amount: fromDelta, BalanceAfter: fromBalanceAfter, EntryRole: EntryRoleTransferOut}
	toEntry := &ledgerentrydb.LedgerEntry{UserId: ctx.UserId, BucketId: toBucket.Id, Currency: ctx.Currency, Amount: toDelta, BalanceAfter: toBalanceAfter, EntryRole: EntryRoleTransferIn}
	return &RecordBuildResult{Event: event, Entries: []*ledgerentrydb.LedgerEntry{fromEntry, toEntry}, BucketBalances: mapBalance(fromBucket.Id, fromBalanceAfter, toBucket.Id, toBalanceAfter)}, nil
}
