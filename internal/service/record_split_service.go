package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/categorydb"
	"marmot-ledger/internal/domain/repository/currencydb"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const splitMaxShares = 100

func CreateSplitRecord(userId int64, req *record.RecordRequest) (*record.RecordResponse, error) {
	if err := validateSplitRequest(req); err != nil {
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

	cashBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.CashBucketId, userId)
	if err != nil {
		return nil, err
	}
	if err := validateBucketForRecord(cashBucket, currency); err != nil {
		return nil, err
	}
	if err := validateCashSideBucket(cashBucket); err != nil {
		return nil, err
	}

	receivableBucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, req.ReceivableBucketId, userId)
	if err != nil {
		return nil, err
	}
	if err := validateBucketForRecord(receivableBucket, currency); err != nil {
		return nil, err
	}
	if err := validateCounterpartyBucket(receivableBucket, "receivable"); err != nil {
		return nil, err
	}

	// 解析 categoryView (self share 用)
	var categoryView *categorydb.CategoryView
	if req.CategoryId != 0 {
		categoryView, err = resolveRecordCategory(session, userId, req.CategoryId, EventTypeExpense)
		if err != nil {
			return nil, err
		}
	}

	eventTime := model.LocalTime(time.Now())
	if strings.TrimSpace(req.EventTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(req.EventTime), mustServiceLocation()); err == nil {
			eventTime = model.LocalTime(parsed)
		}
	}

	cashRunning := cashBucket.Balance
	counterRunning := receivableBucket.Balance

	type splitOutput struct {
		event   *financialeventdb.FinancialEvent
		entries []*ledgerentrydb.LedgerEntry
	}
	outputs := make([]splitOutput, 0, len(req.Shares))

	for _, share := range req.Shares {
		if !share.Amount.IsPositive() {
			return nil, errors.New("split share amount must be greater than 0")
		}
		desc := strings.TrimSpace(share.Description)
		if desc == "" {
			desc = strings.TrimSpace(req.Description)
		}

		var remarkPtr *string
		if rk := strings.TrimSpace(req.Remark); rk != "" {
			remarkPtr = &rk
		}

		eventType := EventTypeReceivableCreate
		entryRoleCounter := EntryRoleReceivableCreate
		entryRoleCash := EntryRoleCashLeg
		includeStats := false
		var catId *int64

		if share.IsSelf {
			eventType = EventTypeExpense
			entryRoleCash = EntryRoleExpense
			includeStats = true
			if categoryView != nil {
				catId = &categoryView.Id
			}
		}

		cashRunning = cashRunning.Sub(share.Amount)
		if cashRunning.IsNegative() {
			return nil, errors.New("cash bucket balance cannot be negative")
		}
		var counterAfter decimal.Decimal
		if !share.IsSelf {
			counterRunning = counterRunning.Add(share.Amount)
			counterAfter = counterRunning
		}

		event := &financialeventdb.FinancialEvent{
			UserId:              userId,
			EventType:           eventType,
			Description:         desc,
			CategoryId:          catId,
			EventTime:           eventTime,
			Currency:            currency,
			Amount:              share.Amount,
			IncludeInStatistics: includeStats,
			Source:              EventSourceManual,
			Status:              EventStatusActive,
			Remark:              remarkPtr,
			IsDeleted:           false,
		}
		if err := financialeventdb.InsertFinancialEvent(session, event); err != nil {
			return nil, err
		}

		entries := make([]*ledgerentrydb.LedgerEntry, 0, 2)
		cashEntry := &ledgerentrydb.LedgerEntry{
			FinancialEventId: event.Id,
			UserId:           userId,
			BucketId:         cashBucket.Id,
			Currency:         currency,
			Amount:           share.Amount.Neg(),
			BalanceAfter:     cashRunning,
			EntryRole:        entryRoleCash,
		}
		if err := ledgerentrydb.InsertLedgerEntry(session, cashEntry); err != nil {
			return nil, err
		}
		entries = append(entries, cashEntry)
		if !share.IsSelf {
			counterEntry := &ledgerentrydb.LedgerEntry{
				FinancialEventId: event.Id,
				UserId:           userId,
				BucketId:         receivableBucket.Id,
				Currency:         currency,
				Amount:           share.Amount,
				BalanceAfter:     counterAfter,
				EntryRole:        entryRoleCounter,
			}
			if err := ledgerentrydb.InsertLedgerEntry(session, counterEntry); err != nil {
				return nil, err
			}
			entries = append(entries, counterEntry)
		}
		outputs = append(outputs, splitOutput{event: event, entries: entries})
	}

	// 用第一条 event.Id 作为 group id，回写每条
	if len(outputs) == 0 {
		return nil, errors.New("split must have at least one share")
	}
	groupId := outputs[0].event.Id
	for _, out := range outputs {
		if err := financialeventdb.UpdateEventGroupId(session, out.event.Id, userId, groupId); err != nil {
			return nil, err
		}
		out.event.EventGroupId = &groupId
	}

	if err := bucketdb.UpdateBucketBalance(session, cashBucket.Id, userId, cashRunning); err != nil {
		return nil, err
	}
	if err := bucketdb.UpdateBucketBalance(session, receivableBucket.Id, userId, counterRunning); err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true

	// 选返回的主事件：优先 self expense，否则第一条
	primary := outputs[0]
	for _, out := range outputs {
		if out.event.EventType == EventTypeExpense {
			primary = out
			break
		}
	}

	entryEntities := make([]ledgerentry.LedgerEntry, 0, len(primary.entries))
	for _, item := range primary.entries {
		entryEntities = append(entryEntities, toLedgerEntryEntity(item))
	}
	return &record.RecordResponse{FinancialEvent: *toFinancialEventEntity(primary.event, entryEntities)}, nil
}

func validateSplitRequest(req *record.RecordRequest) error {
	if req == nil {
		return errors.New("record request is required")
	}
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}
	if req.CashBucketId == 0 || req.ReceivableBucketId == 0 {
		return errors.New("cash bucket and receivable bucket are required")
	}
	if req.CashBucketId == req.ReceivableBucketId {
		return errors.New("cash and receivable buckets must differ")
	}
	if len(req.Shares) == 0 {
		return errors.New("split must have at least one share")
	}
	if len(req.Shares) > splitMaxShares {
		return errors.New("split shares exceed limit")
	}
	hasFriend := false
	total := decimal.Zero
	for _, share := range req.Shares {
		if !share.Amount.IsPositive() {
			return errors.New("split share amount must be greater than 0")
		}
		if !share.IsSelf {
			hasFriend = true
		}
		total = total.Add(share.Amount)
	}
	if !hasFriend {
		return errors.New("split needs at least one friend share; otherwise use plain expense")
	}
	if !req.Amount.IsPositive() {
		return errors.New("total amount must be greater than 0")
	}
	if !total.Equal(req.Amount) {
		return errors.New("split shares total must equal request amount")
	}
	return nil
}
