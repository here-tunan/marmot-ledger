package service

import (
	"errors"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
	"strings"
)

type ExchangeRecordStrategy struct{}

func (ExchangeRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	fromBucket := ctx.Buckets[ctx.Request.FromBucketId]
	toBucket := ctx.Buckets[ctx.Request.ToBucketId]
	if fromBucket == nil || toBucket == nil {
		return nil, errors.New("exchange buckets not loaded")
	}
	if fromBucket.Id == toBucket.Id {
		return nil, errors.New("exchange buckets must be different")
	}

	fromAmount := ctx.Amount
	toAmount := ctx.Request.ToAmount
	toCurrency := strings.ToUpper(strings.TrimSpace(ctx.Request.ToCurrency))

	fromBalanceAfter := fromBucket.Balance.Sub(fromAmount)
	toBalanceAfter := toBucket.Balance.Add(toAmount)
	if err := validateBalanceAfter(fromBucket, fromBalanceAfter); err != nil {
		return nil, err
	}
	if err := validateBalanceAfter(toBucket, toBalanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeExchange, ctx.Currency, false, nil)

	fromEntry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     fromBucket.Id,
		Currency:     ctx.Currency,
		Amount:       fromAmount.Neg(),
		BalanceAfter: fromBalanceAfter,
		EntryRole:    EntryRoleExchangeOut,
	}
	toEntry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     toBucket.Id,
		Currency:     toCurrency,
		Amount:       toAmount,
		BalanceAfter: toBalanceAfter,
		EntryRole:    EntryRoleExchangeIn,
	}
	return &RecordBuildResult{
		Event:          event,
		Entries:        []*ledgerentrydb.LedgerEntry{fromEntry, toEntry},
		BucketBalances: mapBalance(fromBucket.Id, fromBalanceAfter, toBucket.Id, toBalanceAfter),
	}, nil
}
