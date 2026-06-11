package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/currencydb"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

func CreateRecord(userId int64, req *record.RecordRequest) (*record.RecordResponse, error) {
	if err := validateRecordRequest(req); err != nil {
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

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if _, err := currencydb.GetEnabledCurrency(session, currency); err != nil {
		return nil, err
	}

	eventDb, entries, err := createRecordInSession(session, userId, req, currency)
	if err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true

	entryEntities := make([]ledgerentry.LedgerEntry, 0, len(entries))
	for _, item := range entries {
		entryEntities = append(entryEntities, toLedgerEntryEntity(item))
	}

	return &record.RecordResponse{
		FinancialEvent: *toFinancialEventEntity(eventDb, entryEntities),
	}, nil
}

func createRecordInSession(session *xorm.Session, userId int64, req *record.RecordRequest, currency string) (*financialeventdb.FinancialEvent, []*ledgerentrydb.LedgerEntry, error) {
	scenario := strings.TrimSpace(req.Scenario)
	switch scenario {
	case EventTypeIncome:
		return createSingleBucketRecord(session, userId, req, currency, EventTypeIncome)
	case EventTypeExpense:
		return createSingleBucketRecord(session, userId, req, currency, EventTypeExpense)
	case EventTypeRefund:
		return createSingleBucketRecord(session, userId, req, currency, EventTypeRefund)
	case EventTypeTransfer:
		return createTransferRecord(session, userId, req, currency)
	default:
		return nil, nil, errors.New("record scenario is unsupported")
	}
}

func createSingleBucketRecord(session *xorm.Session, userId int64, req *record.RecordRequest, currency string, eventType string) (*financialeventdb.FinancialEvent, []*ledgerentrydb.LedgerEntry, error) {
	bucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.BucketId, userId)
	if err != nil {
		return nil, nil, err
	}
	if err := validateBucketForRecord(bucket, currency); err != nil {
		return nil, nil, err
	}

	delta := signedDeltaForSingleBucket(bucket, req.Amount, eventType)
	balanceAfter := bucket.Balance.Add(delta)
	if err := validateBalanceAfter(bucket, balanceAfter); err != nil {
		return nil, nil, err
	}

	eventDb := buildRecordEvent(userId, req, eventType, currency, eventIncludeInStatistics(eventType))
	if err := financialeventdb.InsertFinancialEvent(session, eventDb); err != nil {
		return nil, nil, err
	}

	entryDb := &ledgerentrydb.LedgerEntry{
		FinancialEventId: eventDb.Id,
		UserId:           userId,
		BucketId:         bucket.Id,
		Currency:         currency,
		Amount:           delta,
		BalanceAfter:     balanceAfter,
		EntryRole:        entryRoleForSingleBucket(eventType),
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, entryDb); err != nil {
		return nil, nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, bucket.Id, userId, balanceAfter); err != nil {
		return nil, nil, err
	}

	return eventDb, []*ledgerentrydb.LedgerEntry{entryDb}, nil
}

func createTransferRecord(session *xorm.Session, userId int64, req *record.RecordRequest, currency string) (*financialeventdb.FinancialEvent, []*ledgerentrydb.LedgerEntry, error) {
	fromBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.FromBucketId, userId)
	if err != nil {
		return nil, nil, err
	}
	toBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.ToBucketId, userId)
	if err != nil {
		return nil, nil, err
	}
	if fromBucket.Id == toBucket.Id {
		return nil, nil, errors.New("transfer buckets must be different")
	}
	if err := validateBucketForRecord(fromBucket, currency); err != nil {
		return nil, nil, err
	}
	if err := validateBucketForRecord(toBucket, currency); err != nil {
		return nil, nil, err
	}

	fromDelta := req.Amount.Neg()
	toDelta := req.Amount
	fromBalanceAfter := fromBucket.Balance.Add(fromDelta)
	toBalanceAfter := toBucket.Balance.Add(toDelta)
	if err := validateBalanceAfter(fromBucket, fromBalanceAfter); err != nil {
		return nil, nil, err
	}
	if err := validateBalanceAfter(toBucket, toBalanceAfter); err != nil {
		return nil, nil, err
	}

	eventDb := buildRecordEvent(userId, req, EventTypeTransfer, currency, false)
	if err := financialeventdb.InsertFinancialEvent(session, eventDb); err != nil {
		return nil, nil, err
	}

	fromEntry := &ledgerentrydb.LedgerEntry{
		FinancialEventId: eventDb.Id,
		UserId:           userId,
		BucketId:         fromBucket.Id,
		Currency:         currency,
		Amount:           fromDelta,
		BalanceAfter:     fromBalanceAfter,
		EntryRole:        EntryRoleTransferOut,
	}
	toEntry := &ledgerentrydb.LedgerEntry{
		FinancialEventId: eventDb.Id,
		UserId:           userId,
		BucketId:         toBucket.Id,
		Currency:         currency,
		Amount:           toDelta,
		BalanceAfter:     toBalanceAfter,
		EntryRole:        EntryRoleTransferIn,
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, fromEntry); err != nil {
		return nil, nil, err
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, toEntry); err != nil {
		return nil, nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, fromBucket.Id, userId, fromBalanceAfter); err != nil {
		return nil, nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, toBucket.Id, userId, toBalanceAfter); err != nil {
		return nil, nil, err
	}

	return eventDb, []*ledgerentrydb.LedgerEntry{fromEntry, toEntry}, nil
}

func validateRecordRequest(req *record.RecordRequest) error {
	if req == nil {
		return errors.New("record request is required")
	}
	if strings.TrimSpace(req.Scenario) == "" {
		return errors.New("record scenario is required")
	}
	if !req.Amount.IsPositive() {
		return errors.New("amount must be greater than 0")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}

	scenario := strings.TrimSpace(req.Scenario)
	switch scenario {
	case EventTypeIncome, EventTypeExpense, EventTypeRefund:
		if req.BucketId == 0 {
			return errors.New("bucket id is required")
		}
	case EventTypeTransfer:
		if req.FromBucketId == 0 || req.ToBucketId == 0 {
			return errors.New("from bucket and to bucket are required")
		}
	default:
		return errors.New("record scenario is unsupported")
	}
	return nil
}

func validateBucketForRecord(bucket *bucketdb.Bucket, currency string) error {
	if !bucket.IsActive {
		return errors.New("bucket is inactive")
	}
	if !strings.EqualFold(bucket.Currency, currency) {
		return errors.New("bucket currency does not match record currency")
	}
	return nil
}

func signedDeltaForSingleBucket(bucket *bucketdb.Bucket, amount decimal.Decimal, eventType string) decimal.Decimal {
	isLiability := bucket.BucketNature == BucketNatureLiability || bucket.BucketType == "credit" || bucket.BucketType == "liability"
	switch eventType {
	case EventTypeIncome:
		return amount
	case EventTypeExpense:
		if isLiability {
			return amount
		}
		return amount.Neg()
	case EventTypeRefund:
		if isLiability {
			return amount.Neg()
		}
		return amount
	default:
		return amount
	}
}

func entryRoleForSingleBucket(eventType string) string {
	switch eventType {
	case EventTypeIncome:
		return EntryRoleIncome
	case EventTypeExpense:
		return EntryRoleExpense
	case EventTypeRefund:
		return EntryRoleRefund
	default:
		return EntryRoleAdjustment
	}
}

func eventIncludeInStatistics(eventType string) bool {
	return eventType == EventTypeIncome || eventType == EventTypeExpense
}

func validateBalanceAfter(bucket *bucketdb.Bucket, balanceAfter decimal.Decimal) error {
	if balanceAfter.IsNegative() {
		return errors.New("bucket balance cannot be negative")
	}
	return nil
}

func buildRecordEvent(userId int64, req *record.RecordRequest, eventType string, currency string, includeInStatistics bool) *financialeventdb.FinancialEvent {
	var relatedFinancialEventId *int64
	if req.RelatedFinancialEventId != 0 {
		relatedFinancialEventId = &req.RelatedFinancialEventId
	}
	var remark *string
	if strings.TrimSpace(req.Remark) != "" {
		trimmed := strings.TrimSpace(req.Remark)
		remark = &trimmed
	}

	eventTime := model.LocalTime(time.Now())
	if strings.TrimSpace(req.EventTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(req.EventTime), mustServiceLocation()); err == nil {
			eventTime = model.LocalTime(parsed)
		}
	}

	return &financialeventdb.FinancialEvent{
		UserId:                  userId,
		RelatedFinancialEventId: relatedFinancialEventId,
		EventType:               eventType,
		Description:             strings.TrimSpace(req.Description),
		EventTime:               eventTime,
		Currency:                currency,
		Amount:                  req.Amount,
		BaseCurrency:            currency,
		BaseAmount:              req.Amount,
		ExchangeRate:            decimal.NewFromInt(1),
		IncludeInStatistics:     includeInStatistics,
		Source:                  EventSourceManual,
		Status:                  EventStatusActive,
		Remark:                  remark,
		IsDeleted:               false,
	}
}

func mustServiceLocation() *time.Location {
	location, err := time.LoadLocation(model.Loc)
	if err != nil {
		return time.Local
	}
	return location
}
