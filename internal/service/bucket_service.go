package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/bucket"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/domain/repository/accountdb"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/currencydb"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"
	"time"
)

func CreateBucket(userId int64, bucketInfo *bucket.Bucket) (*bucket.Bucket, error) {
	if err := validateBucket(bucketInfo); err != nil {
		return nil, err
	}

	session := infrastructure.Mysql.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = session.Rollback()
		}
	}()

	if _, err := accountdb.GetAccountByIdForUser(session, bucketInfo.AccountId, userId); err != nil {
		return nil, err
	}

	currency := strings.ToUpper(strings.TrimSpace(bucketInfo.Currency))
	if _, err := currencydb.GetEnabledCurrency(session, currency); err != nil {
		return nil, err
	}

	initialBalance := bucketInfo.InitialBalance
	bucketDb := &bucketdb.Bucket{
		UserId:         userId,
		AccountId:      bucketInfo.AccountId,
		Name:           strings.TrimSpace(bucketInfo.Name),
		Currency:       currency,
		Balance:        initialBalance,
		InitialBalance: initialBalance,
		BucketType:     strings.TrimSpace(bucketInfo.BucketType),
		BucketNature:   strings.TrimSpace(bucketInfo.BucketNature),
		BucketGroupKey: strings.TrimSpace(bucketInfo.BucketGroupKey),
		IsActive:       true,
		IsDeleted:      false,
	}
	if err := bucketdb.InsertBucket(session, bucketDb); err != nil {
		return nil, err
	}

	now := model.LocalTime(time.Now())
	remark := "初始化余额"
	eventDb := &financialeventdb.FinancialEvent{
		UserId:              userId,
		EventType:           EventTypeBalanceAdjustment,
		Description:         "初始化余额",
		EventTime:           now,
		Currency:            currency,
		Amount:              initialBalance,
		IncludeInStatistics: false,
		Source:              EventSourceManual,
		Status:              EventStatusActive,
		Remark:              &remark,
		IsDeleted:           false,
	}
	if err := financialeventdb.InsertFinancialEvent(session, eventDb); err != nil {
		return nil, err
	}

	entryDb := &ledgerentrydb.LedgerEntry{
		FinancialEventId: eventDb.Id,
		UserId:           userId,
		BucketId:         bucketDb.Id,
		Currency:         currency,
		Amount:           initialBalance,
		BalanceAfter:     initialBalance,
		EntryRole:        EntryRoleAdjustment,
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, entryDb); err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true

	created := toBucketEntity(bucketDb)
	created.InitialFinancialEventId = eventDb.Id
	created.InitialLedgerEntryId = entryDb.Id
	return created, nil
}

func ListBuckets(userId int64, query bucket.BucketQuery) ([]bucket.Bucket, error) {
	buckets, err := bucketdb.ListBuckets(userId, bucketdb.BucketQuery{
		AccountId:    query.AccountId,
		Currency:     query.Currency,
		BucketType:   query.BucketType,
		BucketNature: query.BucketNature,
		IsActive:     query.IsActive,
	})
	if err != nil {
		return nil, err
	}

	result := make([]bucket.Bucket, 0, len(buckets))
	for _, item := range buckets {
		result = append(result, *toBucketEntity(&item))
	}
	return result, nil
}

func GetBucket(userId int64, id int64) (*bucket.Bucket, error) {
	bucketDb, err := bucketdb.GetBucket(id, userId)
	if err != nil {
		return nil, err
	}
	return toBucketEntity(bucketDb), nil
}

func UpdateBucket(userId int64, id int64, bucketInfo *bucket.Bucket) (*bucket.Bucket, error) {
	if err := validateBucketUpdate(bucketInfo); err != nil {
		return nil, err
	}

	if _, err := bucketdb.GetBucket(id, userId); err != nil {
		return nil, err
	}

	bucketDb := &bucketdb.Bucket{
		Id:             id,
		UserId:         userId,
		Name:           strings.TrimSpace(bucketInfo.Name),
		BucketType:     strings.TrimSpace(bucketInfo.BucketType),
		BucketNature:   strings.TrimSpace(bucketInfo.BucketNature),
		BucketGroupKey: strings.TrimSpace(bucketInfo.BucketGroupKey),
		IsActive:       bucketInfo.IsActive,
	}
	if err := bucketdb.UpdateBucket(bucketDb); err != nil {
		return nil, err
	}

	return GetBucket(userId, id)
}

func ListBucketLedgerEntries(userId int64, bucketId int64) ([]ledgerentry.LedgerEntry, error) {
	if _, err := bucketdb.GetBucket(bucketId, userId); err != nil {
		return nil, err
	}

	entries, err := ledgerentrydb.ListLedgerEntriesByBucket(bucketId, userId)
	if err != nil {
		return nil, err
	}

	result := make([]ledgerentry.LedgerEntry, 0, len(entries))
	for _, item := range entries {
		result = append(result, toLedgerEntryEntity(&item))
	}
	return result, nil
}

func validateBucket(bucketInfo *bucket.Bucket) error {
	if err := validateBucketUpdate(bucketInfo); err != nil {
		return err
	}
	if bucketInfo.AccountId == 0 {
		return errors.New("account id is required")
	}
	if strings.TrimSpace(bucketInfo.Currency) == "" {
		return errors.New("currency is required")
	}
	if bucketInfo.InitialBalance.IsNegative() {
		return errors.New("initial balance must be greater than or equal to 0")
	}
	return nil
}

func validateBucketUpdate(bucketInfo *bucket.Bucket) error {
	if bucketInfo == nil {
		return errors.New("bucket is required")
	}
	if strings.TrimSpace(bucketInfo.Name) == "" {
		return errors.New("bucket name is required")
	}
	bucketType := strings.TrimSpace(bucketInfo.BucketType)
	if !validBucketTypes[bucketType] {
		return errors.New("bucket type is invalid")
	}
	bucketNature := strings.TrimSpace(bucketInfo.BucketNature)
	if bucketNature != BucketNatureAsset && bucketNature != BucketNatureLiability {
		return errors.New("bucket nature is invalid")
	}
	return nil
}

func toBucketEntity(bucketDb *bucketdb.Bucket) *bucket.Bucket {
	return &bucket.Bucket{
		Id:             bucketDb.Id,
		UserId:         bucketDb.UserId,
		AccountId:      bucketDb.AccountId,
		Name:           bucketDb.Name,
		Currency:       bucketDb.Currency,
		Balance:        bucketDb.Balance,
		InitialBalance: bucketDb.InitialBalance,
		BucketType:     bucketDb.BucketType,
		BucketNature:   bucketDb.BucketNature,
		BucketGroupKey: bucketDb.BucketGroupKey,
		IsActive:       bucketDb.IsActive,
	}
}

func toLedgerEntryEntity(entryDb *ledgerentrydb.LedgerEntry) ledgerentry.LedgerEntry {
	return ledgerentry.LedgerEntry{
		Id:               entryDb.Id,
		FinancialEventId: entryDb.FinancialEventId,
		UserId:           entryDb.UserId,
		BucketId:         entryDb.BucketId,
		Currency:         entryDb.Currency,
		Amount:           entryDb.Amount,
		BalanceAfter:     entryDb.BalanceAfter,
		EntryRole:        entryDb.EntryRole,
		CreatedAt:        entryDb.CreatedAt.String(),
	}
}
