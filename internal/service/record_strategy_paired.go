package service

import (
	"errors"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
)

type pairedDirection int

const (
	pairedCashToCounter pairedDirection = iota // create receivable / deposit / loan_out
	pairedCounterToCash                        // collect / refund
)

func isPairedScenario(scenario string) bool {
	switch scenario {
	case EventTypeReceivableCreate, EventTypeReceivableCollect,
		EventTypeDepositCreate, EventTypeDepositRefund,
		EventTypeLoanOut, EventTypeLoanCollect:
		return true
	}
	return false
}

func counterpartyBucketTypeFor(scenario string) string {
	switch scenario {
	case EventTypeReceivableCreate, EventTypeReceivableCollect:
		return "receivable"
	case EventTypeDepositCreate, EventTypeDepositRefund:
		return "deposit"
	case EventTypeLoanOut, EventTypeLoanCollect:
		return "loan_out"
	}
	return ""
}

func relatedCreateTypeFor(scenario string) string {
	switch scenario {
	case EventTypeReceivableCollect:
		return EventTypeReceivableCreate
	case EventTypeDepositRefund:
		return EventTypeDepositCreate
	case EventTypeLoanCollect:
		return EventTypeLoanOut
	}
	return ""
}

func directionFor(scenario string) pairedDirection {
	switch scenario {
	case EventTypeReceivableCreate, EventTypeDepositCreate, EventTypeLoanOut:
		return pairedCashToCounter
	}
	return pairedCounterToCash
}

func validateCashSideBucket(bucket *bucketdb.Bucket) error {
	if !cashSideBucketTypes[bucket.BucketType] {
		return errors.New("cash bucket type is invalid for this scenario")
	}
	if bucket.BucketNature != BucketNatureAsset {
		return errors.New("cash bucket must be an asset bucket")
	}
	return nil
}

func validateCounterpartyBucket(bucket *bucketdb.Bucket, expectedType string) error {
	if bucket.BucketType != expectedType {
		return errors.New("counterparty bucket type does not match scenario")
	}
	if bucket.BucketNature != BucketNatureAsset {
		return errors.New("counterparty bucket must be an asset bucket")
	}
	return nil
}

func buildPairedRecord(ctx *RecordContext, eventType string, counterEntryRole string) (*RecordBuildResult, error) {
	cashBucket := ctx.Buckets[ctx.Request.FromBucketId]
	counterBucket := ctx.Buckets[ctx.Request.ToBucketId]
	if cashBucket == nil || counterBucket == nil {
		return nil, errors.New("paired buckets not loaded")
	}
	if cashBucket.Id == counterBucket.Id {
		return nil, errors.New("cash bucket and counterparty bucket must differ")
	}
	if err := validateCashSideBucket(cashBucket); err != nil {
		return nil, err
	}
	if err := validateCounterpartyBucket(counterBucket, counterpartyBucketTypeFor(eventType)); err != nil {
		return nil, err
	}

	direction := directionFor(eventType)
	cashDelta := ctx.Amount.Neg()
	counterDelta := ctx.Amount
	if direction == pairedCounterToCash {
		cashDelta = ctx.Amount
		counterDelta = ctx.Amount.Neg()
	}

	cashBalanceAfter := cashBucket.Balance.Add(cashDelta)
	counterBalanceAfter := counterBucket.Balance.Add(counterDelta)
	if err := validateBalanceAfter(cashBucket, cashBalanceAfter); err != nil {
		return nil, err
	}
	if err := validateBalanceAfter(counterBucket, counterBalanceAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, eventType, ctx.Currency, false, nil)

	if direction == pairedCounterToCash && ctx.Request.RelatedFinancialEventId != 0 {
		related, err := financialeventdb.GetFinancialEvent(ctx.Request.RelatedFinancialEventId, ctx.UserId)
		if err != nil {
			return nil, err
		}
		expectedType := relatedCreateTypeFor(eventType)
		if related.EventType != expectedType {
			return nil, errors.New("related event type does not match scenario")
		}
		if related.Currency != ctx.Currency {
			return nil, errors.New("related event currency does not match")
		}
		// 校验 outstanding ≥ ctx.Amount
		outstandingItems, err := ListOutstandingForBucket(ctx.UserId, counterBucket.Id, expectedType)
		if err != nil {
			return nil, err
		}
		var matched bool
		for _, item := range outstandingItems {
			if item.Id == related.Id {
				matched = true
				if ctx.Amount.GreaterThan(item.OutstandingAmount) {
					return nil, errors.New("amount exceeds outstanding of related event")
				}
				break
			}
		}
		if !matched {
			return nil, errors.New("related event has no outstanding amount left")
		}
		// 继承 event_group_id
		if related.EventGroupId != nil {
			event.EventGroupId = related.EventGroupId
		} else {
			relatedId := related.Id
			event.EventGroupId = &relatedId
		}
	}

	cashEntry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     cashBucket.Id,
		Currency:     ctx.Currency,
		Amount:       cashDelta,
		BalanceAfter: cashBalanceAfter,
		EntryRole:    EntryRoleCashLeg,
	}
	counterEntry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     counterBucket.Id,
		Currency:     ctx.Currency,
		Amount:       counterDelta,
		BalanceAfter: counterBalanceAfter,
		EntryRole:    counterEntryRole,
	}
	return &RecordBuildResult{
		Event:          event,
		Entries:        []*ledgerentrydb.LedgerEntry{cashEntry, counterEntry},
		BucketBalances: mapBalance(cashBucket.Id, cashBalanceAfter, counterBucket.Id, counterBalanceAfter),
	}, nil
}

type ReceivableCreateRecordStrategy struct{}

func (ReceivableCreateRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	return buildPairedRecord(ctx, EventTypeReceivableCreate, EntryRoleReceivableCreate)
}

type ReceivableCollectRecordStrategy struct{}

func (ReceivableCollectRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	return buildPairedRecord(ctx, EventTypeReceivableCollect, EntryRoleReceivableCollect)
}

type DepositCreateRecordStrategy struct{}

func (DepositCreateRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	return buildPairedRecord(ctx, EventTypeDepositCreate, EntryRoleDepositCreate)
}

type DepositRefundRecordStrategy struct{}

func (DepositRefundRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	return buildPairedRecord(ctx, EventTypeDepositRefund, EntryRoleDepositRefund)
}

type LoanOutRecordStrategy struct{}

func (LoanOutRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	return buildPairedRecord(ctx, EventTypeLoanOut, EntryRoleLoanOut)
}

type LoanCollectRecordStrategy struct{}

func (LoanCollectRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	return buildPairedRecord(ctx, EventTypeLoanCollect, EntryRoleLoanCollect)
}
