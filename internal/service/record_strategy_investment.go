package service

import (
	"errors"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
)

func validateInvestmentBucket(bucket *bucketdb.Bucket) error {
	if !investmentBucketTypes[bucket.BucketType] {
		return errors.New("bucket must be an investment bucket")
	}
	if bucket.BucketNature != BucketNatureAsset {
		return errors.New("investment bucket must be an asset bucket")
	}
	return nil
}

func isInvestmentScenario(scenario string) bool {
	switch scenario {
	case EventTypeInvestmentBuy, EventTypeInvestmentSell, EventTypeInvestmentIncome, EventTypeInvestmentRevalue:
		return true
	}
	return false
}

func validateInvestmentCashSideBucket(bucket *bucketdb.Bucket) error {
	if !investmentCashSideBucketTypes[bucket.BucketType] {
		return errors.New("cash bucket type is invalid for this scenario")
	}
	if bucket.BucketNature != BucketNatureAsset {
		return errors.New("cash bucket must be an asset bucket")
	}
	return nil
}

type InvestmentBuyRecordStrategy struct{}

func (InvestmentBuyRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	cashBucket := ctx.Buckets[ctx.Request.FromBucketId]
	investBucket := ctx.Buckets[ctx.Request.ToBucketId]
	if cashBucket == nil || investBucket == nil {
		return nil, errors.New("investment buckets not loaded")
	}
	if cashBucket.Id == investBucket.Id {
		return nil, errors.New("cash bucket and investment bucket must differ")
	}
	if err := validateInvestmentCashSideBucket(cashBucket); err != nil {
		return nil, err
	}
	if err := validateInvestmentBucket(investBucket); err != nil {
		return nil, err
	}

	cashAfter := cashBucket.Balance.Sub(ctx.Amount)
	investAfter := investBucket.Balance.Add(ctx.Amount)
	if err := validateBalanceAfter(cashBucket, cashAfter); err != nil {
		return nil, err
	}
	if err := validateBalanceAfter(investBucket, investAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeInvestmentBuy, ctx.Currency, false, nil)
	cashEntry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     cashBucket.Id,
		Currency:     ctx.Currency,
		Amount:       ctx.Amount.Neg(),
		BalanceAfter: cashAfter,
		EntryRole:    EntryRoleCashLeg,
	}
	investEntry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     investBucket.Id,
		Currency:     ctx.Currency,
		Amount:       ctx.Amount,
		BalanceAfter: investAfter,
		EntryRole:    EntryRoleInvestmentBuy,
	}
	return &RecordBuildResult{
		Event:          event,
		Entries:        []*ledgerentrydb.LedgerEntry{cashEntry, investEntry},
		BucketBalances: mapBalance(cashBucket.Id, cashAfter, investBucket.Id, investAfter),
	}, nil
}

type InvestmentIncomeRecordStrategy struct{}

func (InvestmentIncomeRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	cashBucket := ctx.Buckets[ctx.Request.BucketId]
	if cashBucket == nil {
		return nil, errors.New("cash bucket not loaded")
	}
	if err := validateInvestmentCashSideBucket(cashBucket); err != nil {
		return nil, err
	}

	cashAfter := cashBucket.Balance.Add(ctx.Amount)
	if err := validateBalanceAfter(cashBucket, cashAfter); err != nil {
		return nil, err
	}

	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeInvestmentIncome, ctx.Currency, false, nil)
	entry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     cashBucket.Id,
		Currency:     ctx.Currency,
		Amount:       ctx.Amount,
		BalanceAfter: cashAfter,
		EntryRole:    EntryRoleInvestmentIncome,
	}
	return &RecordBuildResult{
		Event:          event,
		Entries:        []*ledgerentrydb.LedgerEntry{entry},
		BucketBalances: singleBalance(cashBucket.Id, cashAfter),
	}, nil
}

type InvestmentRevalueRecordStrategy struct{}

func (InvestmentRevalueRecordStrategy) Build(ctx *RecordContext) (*RecordBuildResult, error) {
	if ctx.Amount.IsZero() {
		return nil, errors.New("revalue amount must not be 0")
	}
	bucket := ctx.Buckets[ctx.Request.BucketId]
	if bucket == nil {
		return nil, errors.New("investment bucket not loaded")
	}
	if err := validateInvestmentBucket(bucket); err != nil {
		return nil, err
	}

	balanceAfter := bucket.Balance.Add(ctx.Amount)
	if err := validateBalanceAfter(bucket, balanceAfter); err != nil {
		return nil, err
	}

	role := EntryRoleRevaluationGain
	if ctx.Amount.IsNegative() {
		role = EntryRoleRevaluationLoss
	}
	event := buildRecordEvent(ctx.UserId, ctx.Request, EventTypeInvestmentRevalue, ctx.Currency, false, nil)
	entry := &ledgerentrydb.LedgerEntry{
		UserId:       ctx.UserId,
		BucketId:     bucket.Id,
		Currency:     ctx.Currency,
		Amount:       ctx.Amount,
		BalanceAfter: balanceAfter,
		EntryRole:    role,
	}
	return &RecordBuildResult{
		Event:          event,
		Entries:        []*ledgerentrydb.LedgerEntry{entry},
		BucketBalances: singleBalance(bucket.Id, balanceAfter),
	}, nil
}
