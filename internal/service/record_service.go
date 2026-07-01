package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/domain/repository/bucketdb"
	"marmot-ledger/internal/domain/repository/categorydb"
	"marmot-ledger/internal/domain/repository/channeldb"
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
	if strings.TrimSpace(req.Scenario) == "split" {
		return CreateSplitRecord(userId, req)
	}
	if strings.TrimSpace(req.Scenario) == EventTypeInvestmentSell {
		return CreateInvestmentSellRecord(userId, req)
	}
	if strings.TrimSpace(req.Scenario) == EventTypeFamilyTransfer {
		return CreateFamilyTransferRecord(userId, req)
	}
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

	ctx, err := buildRecordContext(session, userId, req, currency)
	if err != nil {
		return nil, err
	}

	strategy, ok := recordStrategies[strings.TrimSpace(req.Scenario)]
	if !ok {
		return nil, errors.New("record scenario is unsupported")
	}

	buildResult, err := strategy.Build(ctx)
	if err != nil {
		return nil, err
	}
	if err := persistRecordBuildResult(session, userId, buildResult); err != nil {
		return nil, err
	}

	eventDb := buildResult.Event
	entries := buildResult.Entries

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

func buildRecordContext(session *xorm.Session, userId int64, req *record.RecordRequest, currency string) (*RecordContext, error) {
	buckets := make(map[int64]*bucketdb.Bucket)
	loadBucket := func(id int64) error {
		if id == 0 {
			return nil
		}
		if _, ok := buckets[id]; ok {
			return nil
		}
		bucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, id, userId)
		if err != nil {
			return err
		}
		if err := validateBucketForRecord(bucket, currency); err != nil {
			return err
		}
		buckets[id] = bucket
		return nil
	}

	if req.Scenario == EventTypeTransfer {
		if err := loadBucket(req.FromBucketId); err != nil {
			return nil, err
		}
		if err := loadBucket(req.ToBucketId); err != nil {
			return nil, err
		}
	} else if req.Scenario == EventTypeExchange {
		if err := loadBucketForCurrency(session, userId, buckets, req.FromBucketId, currency); err != nil {
			return nil, err
		}
		toCurrency := strings.ToUpper(strings.TrimSpace(req.ToCurrency))
		if _, err := currencydb.GetEnabledCurrency(session, toCurrency); err != nil {
			return nil, err
		}
		if err := loadBucketForCurrency(session, userId, buckets, req.ToBucketId, toCurrency); err != nil {
			return nil, err
		}
	} else if isPairedScenario(req.Scenario) {
		if err := loadBucket(req.FromBucketId); err != nil {
			return nil, err
		}
		if err := loadBucket(req.ToBucketId); err != nil {
			return nil, err
		}
	} else if req.Scenario == EventTypeInvestmentBuy {
		if err := loadBucket(req.FromBucketId); err != nil {
			return nil, err
		}
		if err := loadBucket(req.ToBucketId); err != nil {
			return nil, err
		}
	} else {
		if err := loadBucket(req.BucketId); err != nil {
			return nil, err
		}
	}

	categoryView, err := resolveRecordCategory(session, userId, req.CategoryId, req.Scenario)
	if err != nil {
		return nil, err
	}
	channelId, err := resolveRecordChannel(userId, req.ChannelId, req.Scenario)
	if err != nil {
		return nil, err
	}

	eventTime := model.LocalTime(time.Now())
	if strings.TrimSpace(req.EventTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(req.EventTime), mustServiceLocation()); err == nil {
			eventTime = model.LocalTime(parsed)
		}
	}

	return &RecordContext{UserId: userId, Request: req, Session: session, Amount: req.Amount, Currency: currency, EventTime: eventTime, Buckets: buckets, Category: categoryView, ChannelId: channelId}, nil
}

func persistRecordBuildResult(session *xorm.Session, userId int64, buildResult *RecordBuildResult) error {
	if err := financialeventdb.InsertFinancialEvent(session, buildResult.Event); err != nil {
		return err
	}
	return persistLedgerEntriesAndBalances(session, userId, buildResult)
}

func persistLedgerEntriesAndBalances(session *xorm.Session, userId int64, buildResult *RecordBuildResult) error {
	for _, entry := range buildResult.Entries {
		entry.FinancialEventId = buildResult.Event.Id
		if err := ledgerentrydb.InsertLedgerEntry(session, entry); err != nil {
			return err
		}
	}
	for bucketId, balance := range buildResult.BucketBalances {
		if err := bucketdb.UpdateBucketBalance(session, bucketId, userId, balance); err != nil {
			return err
		}
	}
	return nil
}

func UpdateRecord(userId int64, id int64, req *record.RecordRequest) (*record.RecordResponse, error) {
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

	oldEvent, err := financialeventdb.GetFinancialEventForUpdate(session, id, userId)
	if err != nil {
		return nil, err
	}
	if _, ok := recordStrategies[oldEvent.EventType]; !ok {
		return nil, errors.New("financial event type does not support edit")
	}

	oldEntries, err := ledgerentrydb.ListLedgerEntriesByEventInSession(session, id, userId)
	if err != nil {
		return nil, err
	}
	if err := rollbackLedgerEntries(session, userId, oldEntries); err != nil {
		return nil, err
	}
	if err := ledgerentrydb.DeleteLedgerEntriesByEvent(session, id, userId); err != nil {
		return nil, err
	}

	currency := strings.ToUpper(strings.TrimSpace(req.Currency))
	if _, err := currencydb.GetEnabledCurrency(session, currency); err != nil {
		return nil, err
	}
	ctx, err := buildRecordContext(session, userId, req, currency)
	if err != nil {
		return nil, err
	}
	strategy, ok := recordStrategies[strings.TrimSpace(req.Scenario)]
	if !ok {
		return nil, errors.New("record scenario is unsupported")
	}
	buildResult, err := strategy.Build(ctx)
	if err != nil {
		return nil, err
	}
	buildResult.Event.Id = id
	if err := financialeventdb.UpdateFinancialEvent(session, buildResult.Event); err != nil {
		return nil, err
	}
	if err := persistLedgerEntriesAndBalances(session, userId, buildResult); err != nil {
		return nil, err
	}

	if err := session.Commit(); err != nil {
		return nil, err
	}
	committed = true

	entryEntities := make([]ledgerentry.LedgerEntry, 0, len(buildResult.Entries))
	for _, item := range buildResult.Entries {
		entryEntities = append(entryEntities, toLedgerEntryEntity(item))
	}
	return &record.RecordResponse{FinancialEvent: *toFinancialEventEntity(buildResult.Event, entryEntities)}, nil
}

func DeleteRecord(userId int64, id int64) error {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	committed := false
	defer func() {
		if !committed {
			_ = session.Rollback()
		}
	}()

	event, err := financialeventdb.GetFinancialEventForUpdate(session, id, userId)
	if err != nil {
		return err
	}
	if _, ok := recordStrategies[event.EventType]; !ok {
		return errors.New("financial event type does not support delete")
	}
	entries, err := ledgerentrydb.ListLedgerEntriesByEventInSession(session, id, userId)
	if err != nil {
		return err
	}
	if err := rollbackLedgerEntries(session, userId, entries); err != nil {
		return err
	}
	if err := ledgerentrydb.DeleteLedgerEntriesByEvent(session, id, userId); err != nil {
		return err
	}
	if err := financialeventdb.SoftDeleteFinancialEvent(session, id, userId); err != nil {
		return err
	}
	if err := session.Commit(); err != nil {
		return err
	}
	committed = true
	return nil
}

func rollbackLedgerEntries(session *xorm.Session, userId int64, entries []ledgerentrydb.LedgerEntry) error {
	for _, entry := range entries {
		bucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, entry.BucketId, userId)
		if err != nil {
			return err
		}
		newBalance := bucket.Balance.Sub(entry.Amount)
		if err := validateBalanceAfter(bucket, newBalance); err != nil {
			return err
		}
		if err := bucketdb.UpdateBucketBalance(session, bucket.Id, userId, newBalance); err != nil {
			return err
		}
	}
	return nil
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

	categoryView, err := resolveRecordCategory(session, userId, req.CategoryId, eventType)
	if err != nil {
		return nil, nil, err
	}

	eventDb := buildRecordEvent(userId, req, eventType, currency, eventIncludeInStatistics(eventType), categoryView)
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

	eventDb := buildRecordEvent(userId, req, EventTypeTransfer, currency, false, nil)
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
	if strings.TrimSpace(req.Currency) == "" {
		return errors.New("currency is required")
	}

	scenario := strings.TrimSpace(req.Scenario)
	if scenario == EventTypeBalanceAdjustment || scenario == EventTypeInvestmentRevalue {
		if req.Amount.IsZero() {
			return errors.New("amount must not be 0")
		}
	} else if !req.Amount.IsPositive() {
		return errors.New("amount must be greater than 0")
	}

	switch scenario {
	case EventTypeIncome, EventTypeExpense, EventTypeRefund, EventTypeBalanceAdjustment:
		if req.BucketId == 0 {
			return errors.New("bucket id is required")
		}
	case EventTypeTransfer:
		if req.FromBucketId == 0 || req.ToBucketId == 0 {
			return errors.New("from bucket and to bucket are required")
		}
	case EventTypeExchange:
		if req.FromBucketId == 0 || req.ToBucketId == 0 {
			return errors.New("from bucket and to bucket are required")
		}
		if req.FromBucketId == req.ToBucketId {
			return errors.New("exchange buckets must be different")
		}
		if !req.ToAmount.IsPositive() {
			return errors.New("to amount must be greater than 0")
		}
		if strings.TrimSpace(req.ToCurrency) == "" {
			return errors.New("to currency is required")
		}
		if strings.EqualFold(strings.TrimSpace(req.Currency), strings.TrimSpace(req.ToCurrency)) {
			return errors.New("exchange currencies must differ; use transfer when both sides share the currency")
		}
	case EventTypeReceivableCreate, EventTypeReceivableCollect,
		EventTypeDepositCreate, EventTypeDepositRefund,
		EventTypeLoanOut, EventTypeLoanCollect:
		if req.FromBucketId == 0 || req.ToBucketId == 0 {
			return errors.New("cash bucket and counterparty bucket are required")
		}
		if req.FromBucketId == req.ToBucketId {
			return errors.New("cash bucket and counterparty bucket must differ")
		}
	case EventTypeInvestmentBuy:
		if req.FromBucketId == 0 || req.ToBucketId == 0 {
			return errors.New("cash bucket and investment bucket are required")
		}
		if req.FromBucketId == req.ToBucketId {
			return errors.New("cash bucket and investment bucket must differ")
		}
	case EventTypeInvestmentIncome:
		if req.BucketId == 0 {
			return errors.New("cash bucket is required")
		}
	case EventTypeInvestmentRevalue:
		if req.BucketId == 0 {
			return errors.New("investment bucket is required")
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

func loadBucketForCurrency(session *xorm.Session, userId int64, buckets map[int64]*bucketdb.Bucket, bucketId int64, currency string) error {
	if bucketId == 0 {
		return errors.New("bucket id is required")
	}
	if _, ok := buckets[bucketId]; ok {
		return nil
	}
	bucket, err := bucketdb.GetBucketByIdForUserForUpdate(session, bucketId, userId)
	if err != nil {
		return err
	}
	if err := validateBucketForRecord(bucket, currency); err != nil {
		return err
	}
	buckets[bucketId] = bucket
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

func resolveRecordCategory(session *xorm.Session, userId int64, categoryId int64, eventType string) (*categorydb.CategoryView, error) {
	if categoryId == 0 || eventType == EventTypeTransfer || eventType == EventTypeExchange || eventType == EventTypeBalanceAdjustment || isPairedScenario(eventType) || isInvestmentScenario(eventType) {
		return nil, nil
	}

	categoryView, err := categorydb.GetCategoryByIdForUser(session, categoryId, userId)
	if err != nil {
		return nil, err
	}
	if !categoryView.IsActive {
		return nil, errors.New("category is inactive")
	}
	if eventType == EventTypeIncome && categoryView.Type != EventTypeIncome {
		return nil, errors.New("category type does not match income")
	}
	if (eventType == EventTypeExpense || eventType == EventTypeRefund) && categoryView.Type != EventTypeExpense {
		return nil, errors.New("category type does not match expense")
	}
	return categoryView, nil
}

func resolveRecordChannel(userId int64, channelId int64, scenario string) (*int64, error) {
	if channelId == 0 {
		return nil, nil
	}
	channelDb, err := channeldb.GetChannel(channelId, userId)
	if err != nil {
		return nil, err
	}
	if !channelDb.IsActive {
		return nil, errors.New("channel is inactive")
	}
	eventType := strings.ToLower(strings.TrimSpace(scenario))
	if eventType != "" && channelDb.SupportedEventTypes != "" && !supportsEventType(channelDb.SupportedEventTypes, eventType) {
		return nil, errors.New("channel does not support this record type")
	}
	return &channelDb.Id, nil
}

func supportsEventType(supported string, eventType string) bool {
	for _, item := range strings.Split(supported, ",") {
		if strings.ToLower(strings.TrimSpace(item)) == eventType {
			return true
		}
	}
	return false
}

func buildRecordEvent(userId int64, req *record.RecordRequest, eventType string, currency string, includeInStatistics bool, categoryView *categorydb.CategoryView) *financialeventdb.FinancialEvent {
	var relatedFinancialEventId *int64
	if req.RelatedFinancialEventId != 0 {
		relatedFinancialEventId = &req.RelatedFinancialEventId
	}
	var remark *string
	if strings.TrimSpace(req.Remark) != "" {
		trimmed := strings.TrimSpace(req.Remark)
		remark = &trimmed
	}

	var channelId *int64
	if req.ChannelId > 0 {
		channelId = &req.ChannelId
	}

	eventTime := model.LocalTime(time.Now())
	if strings.TrimSpace(req.EventTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(req.EventTime), mustServiceLocation()); err == nil {
			eventTime = model.LocalTime(parsed)
		}
	}

	var categoryId *int64
	if categoryView != nil {
		categoryId = &categoryView.Id
	}

	source := strings.TrimSpace(req.Source)
	if source == "" {
		source = EventSourceManual
	}

	return &financialeventdb.FinancialEvent{
		UserId:                  userId,
		RelatedFinancialEventId: relatedFinancialEventId,
		EventType:               eventType,
		Description:             strings.TrimSpace(req.Description),
		CategoryId:              categoryId,
		ChannelId:               channelId,
		EventTime:               eventTime,
		Currency:                currency,
		Amount:                  req.Amount,
		IncludeInStatistics:     includeInStatistics,
		Source:                  source,
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
