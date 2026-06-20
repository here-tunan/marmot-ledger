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

	"xorm.io/xorm"
)

func CreateFamilyTransferRecord(userId int64, req *record.RecordRequest) (*record.RecordResponse, error) {
	if req == nil {
		return nil, errors.New("record request is required")
	}
	if req.FamilyId == 0 {
		return nil, errors.New("family id is required")
	}
	if req.FromBucketId == 0 || req.ToBucketId == 0 {
		return nil, errors.New("from bucket and to bucket are required")
	}
	if req.FromBucketId == req.ToBucketId {
		return nil, errors.New("from and to buckets must differ")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return nil, errors.New("currency is required")
	}
	if !req.Amount.IsPositive() {
		return nil, errors.New("amount must be greater than 0")
	}

	activeUserIds, err := ActiveFamilyUserIds(userId, req.FamilyId)
	if err != nil {
		return nil, err
	}
	activeSet := make(map[int64]bool, len(activeUserIds))
	for _, id := range activeUserIds {
		activeSet[id] = true
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

	fromBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.FromBucketId, userId)
	if err != nil {
		return nil, err
	}
	if err := validateBucketForRecord(fromBucket, currency); err != nil {
		return nil, err
	}

	toBucket, err := getBucketForUpdateById(session, req.ToBucketId)
	if err != nil {
		return nil, err
	}
	if !activeSet[toBucket.UserId] {
		return nil, errors.New("target bucket owner is not an active family member")
	}
	if err := validateBucketForRecord(toBucket, currency); err != nil {
		return nil, err
	}

	fromAfter := fromBucket.Balance.Sub(req.Amount)
	toAfter := toBucket.Balance.Add(req.Amount)
	if err := validateBalanceAfter(fromBucket, fromAfter); err != nil {
		return nil, err
	}
	if err := validateBalanceAfter(toBucket, toAfter); err != nil {
		return nil, err
	}

	eventTime := model.LocalTime(time.Now())
	if strings.TrimSpace(req.EventTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(req.EventTime), mustServiceLocation()); err == nil {
			eventTime = model.LocalTime(parsed)
		}
	}
	var remarkPtr *string
	if rk := strings.TrimSpace(req.Remark); rk != "" {
		remarkPtr = &rk
	}

	event := &financialeventdb.FinancialEvent{
		UserId: userId, EventType: EventTypeFamilyTransfer, Description: strings.TrimSpace(req.Description),
		EventTime: eventTime, Currency: currency, Amount: req.Amount, IncludeInStatistics: false,
		Source: EventSourceManual, Status: EventStatusActive, Remark: remarkPtr,
	}
	if err := financialeventdb.InsertFinancialEvent(session, event); err != nil {
		return nil, err
	}

	outEntry := &ledgerentrydb.LedgerEntry{FinancialEventId: event.Id, UserId: userId, BucketId: fromBucket.Id, Currency: currency, Amount: req.Amount.Neg(), BalanceAfter: fromAfter, EntryRole: EntryRoleFamilyTransferOut}
	inEntry := &ledgerentrydb.LedgerEntry{FinancialEventId: event.Id, UserId: toBucket.UserId, BucketId: toBucket.Id, Currency: currency, Amount: req.Amount, BalanceAfter: toAfter, EntryRole: EntryRoleFamilyTransferIn}
	if err := ledgerentrydb.InsertLedgerEntry(session, outEntry); err != nil {
		return nil, err
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, inEntry); err != nil {
		return nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, fromBucket.Id, userId, fromAfter); err != nil {
		return nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, toBucket.Id, toBucket.UserId, toAfter); err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true
	entries := []ledgerentry.LedgerEntry{toLedgerEntryEntity(outEntry), toLedgerEntryEntity(inEntry)}
	return &record.RecordResponse{FinancialEvent: *toFinancialEventEntity(event, entries)}, nil
}

func getBucketForUpdateById(session *xorm.Session, id int64) (*bucketdb.Bucket, error) {
	bucket := &bucketdb.Bucket{}
	has, err := session.SQL("SELECT * FROM bucket WHERE id = ? AND is_deleted = 0 FOR UPDATE", id).Get(bucket)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("bucket not found")
	}
	return bucket, nil
}
