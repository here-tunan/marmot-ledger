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
)

func CreateInvestmentSellRecord(userId int64, req *record.RecordRequest) (*record.RecordResponse, error) {
	if req == nil {
		return nil, errors.New("record request is required")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return nil, errors.New("currency is required")
	}
	if req.FromBucketId == 0 || req.ToBucketId == 0 {
		return nil, errors.New("invest bucket and cash bucket are required")
	}
	if req.FromBucketId == req.ToBucketId {
		return nil, errors.New("invest and cash buckets must differ")
	}
	if !req.Amount.IsPositive() {
		return nil, errors.New("received amount must be greater than 0")
	}
	if req.RemainingMarketValue.IsNegative() {
		return nil, errors.New("remaining market value must be >= 0")
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

	investBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.FromBucketId, userId)
	if err != nil {
		return nil, err
	}
	if err := validateBucketForRecord(investBucket, currency); err != nil {
		return nil, err
	}
	if err := validateInvestmentBucket(investBucket); err != nil {
		return nil, err
	}

	cashBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.ToBucketId, userId)
	if err != nil {
		return nil, err
	}
	if err := validateBucketForRecord(cashBucket, currency); err != nil {
		return nil, err
	}
	if err := validateInvestmentCashSideBucket(cashBucket); err != nil {
		return nil, err
	}

	received := req.Amount
	remaining := req.RemainingMarketValue
	targetMarketValue := received.Add(remaining)
	revalueDelta := targetMarketValue.Sub(investBucket.Balance)

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

	var revalueEvent *financialeventdb.FinancialEvent
	if !revalueDelta.IsZero() {
		revalueDesc := "卖出前市值修正"
		revalueEvent = &financialeventdb.FinancialEvent{
			UserId:              userId,
			EventType:           EventTypeInvestmentRevalue,
			Description:         revalueDesc,
			EventTime:           eventTime,
			Currency:            currency,
			Amount:              revalueDelta,
			IncludeInStatistics: false,
			Source:              EventSourceManual,
			Status:              EventStatusActive,
			Remark:              remarkPtr,
		}
		if err := financialeventdb.InsertFinancialEvent(session, revalueEvent); err != nil {
			return nil, err
		}

		revalueRole := EntryRoleRevaluationGain
		if revalueDelta.IsNegative() {
			revalueRole = EntryRoleRevaluationLoss
		}
		revalueEntry := &ledgerentrydb.LedgerEntry{
			FinancialEventId: revalueEvent.Id,
			UserId:           userId,
			BucketId:         investBucket.Id,
			Currency:         currency,
			Amount:           revalueDelta,
			BalanceAfter:     targetMarketValue,
			EntryRole:        revalueRole,
		}
		if err := ledgerentrydb.InsertLedgerEntry(session, revalueEntry); err != nil {
			return nil, err
		}
	}

	investAfter := targetMarketValue.Sub(received)
	cashAfter := cashBucket.Balance.Add(received)
	if err := validateBalanceAfter(investBucket, investAfter); err != nil {
		return nil, err
	}

	var relatedFinancialEventId *int64
	if req.RelatedFinancialEventId != 0 {
		relatedFinancialEventId = &req.RelatedFinancialEventId
	}
	sellEvent := &financialeventdb.FinancialEvent{
		UserId:                  userId,
		RelatedFinancialEventId: relatedFinancialEventId,
		EventType:               EventTypeInvestmentSell,
		Description:             strings.TrimSpace(req.Description),
		EventTime:               eventTime,
		Currency:                currency,
		Amount:                  received,
		IncludeInStatistics:     false,
		Source:                  EventSourceManual,
		Status:                  EventStatusActive,
		Remark:                  remarkPtr,
	}
	if err := financialeventdb.InsertFinancialEvent(session, sellEvent); err != nil {
		return nil, err
	}

	investEntry := &ledgerentrydb.LedgerEntry{
		FinancialEventId: sellEvent.Id,
		UserId:           userId,
		BucketId:         investBucket.Id,
		Currency:         currency,
		Amount:           received.Neg(),
		BalanceAfter:     investAfter,
		EntryRole:        EntryRoleInvestmentSell,
	}
	cashEntry := &ledgerentrydb.LedgerEntry{
		FinancialEventId: sellEvent.Id,
		UserId:           userId,
		BucketId:         cashBucket.Id,
		Currency:         currency,
		Amount:           received,
		BalanceAfter:     cashAfter,
		EntryRole:        EntryRoleCashLeg,
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, investEntry); err != nil {
		return nil, err
	}
	if err := ledgerentrydb.InsertLedgerEntry(session, cashEntry); err != nil {
		return nil, err
	}

	var groupId int64
	if revalueEvent != nil {
		groupId = revalueEvent.Id
	} else {
		groupId = sellEvent.Id
	}
	if revalueEvent != nil {
		if err := financialeventdb.UpdateEventGroupId(session, revalueEvent.Id, userId, groupId); err != nil {
			return nil, err
		}
		revalueEvent.EventGroupId = &groupId
	}
	if err := financialeventdb.UpdateEventGroupId(session, sellEvent.Id, userId, groupId); err != nil {
		return nil, err
	}
	sellEvent.EventGroupId = &groupId

	if err := bucketdb.UpdateBucketBalance(session, investBucket.Id, userId, investAfter); err != nil {
		return nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, cashBucket.Id, userId, cashAfter); err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true

	entryEntities := []ledgerentry.LedgerEntry{
		toLedgerEntryEntity(investEntry),
		toLedgerEntryEntity(cashEntry),
	}
	return &record.RecordResponse{FinancialEvent: *toFinancialEventEntity(sellEvent, entryEntities)}, nil
}
