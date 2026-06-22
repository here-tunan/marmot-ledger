package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/financialevent"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"

	"github.com/shopspring/decimal"
)

type PageResult[T any] struct {
	List     []T   `json:"list"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

func ListFinancialEvents(userId int64, query financialevent.FinancialEventQuery) (*PageResult[financialevent.FinancialEvent], error) {
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	query.Page = page
	query.PageSize = pageSize

	events, total, err := financialeventdb.ListFinancialEvents(userId, financialeventdb.FinancialEventQuery{
		EventType:           query.EventType,
		StartTime:           query.StartTime,
		EndTime:             query.EndTime,
		Currency:            query.Currency,
		CategoryId:          query.CategoryId,
		ChannelId:           query.ChannelId,
		BucketId:            query.BucketId,
		Keyword:             query.Keyword,
		IncludeInStatistics: query.IncludeInStatistics,
		Page:                query.Page,
		PageSize:            query.PageSize,
	})
	if err != nil {
		return nil, err
	}

	result := make([]financialevent.FinancialEvent, 0, len(events))
	for _, item := range events {
		result = append(result, *toFinancialEventEntity(&item, nil))
	}

	return &PageResult[financialevent.FinancialEvent]{
		List:     result,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func ListGroupedFinancialEvents(userId int64, query financialevent.FinancialEventQuery) (*PageResult[financialevent.FinancialEvent], error) {
	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	query.Page = 0
	query.PageSize = 0

	matched, err := financialeventdb.ListAllFinancialEventsByUserIds([]int64{userId}, financialeventdb.FinancialEventQuery{
		EventType:           query.EventType,
		StartTime:           query.StartTime,
		EndTime:             query.EndTime,
		Currency:            query.Currency,
		CategoryId:          query.CategoryId,
		ChannelId:           query.ChannelId,
		BucketId:            query.BucketId,
		Keyword:             query.Keyword,
		IncludeInStatistics: query.IncludeInStatistics,
	})
	if err != nil {
		return nil, err
	}

	seen := map[int64]bool{}
	groupKeys := make([]int64, 0)
	for _, ev := range matched {
		key := ev.Id
		if ev.EventGroupId != nil {
			key = *ev.EventGroupId
		}
		if !seen[key] {
			seen[key] = true
			groupKeys = append(groupKeys, key)
		}
	}
	total := int64(len(groupKeys))
	start := (page - 1) * pageSize
	if start >= len(groupKeys) {
		return &PageResult[financialevent.FinancialEvent]{List: []financialevent.FinancialEvent{}, Page: page, PageSize: pageSize, Total: total}, nil
	}
	end := start + pageSize
	if end > len(groupKeys) {
		end = len(groupKeys)
	}
	pageKeys := groupKeys[start:end]

	childrenDb, err := financialeventdb.ListFinancialEventsByGroupKeys([]int64{userId}, pageKeys)
	if err != nil {
		return nil, err
	}
	byGroup := map[int64][]financialeventdb.FinancialEvent{}
	for _, ev := range childrenDb {
		key := ev.Id
		if ev.EventGroupId != nil {
			key = *ev.EventGroupId
		}
		byGroup[key] = append(byGroup[key], ev)
	}

	result := make([]financialevent.FinancialEvent, 0, len(pageKeys))
	for _, key := range pageKeys {
		items := byGroup[key]
		if len(items) == 0 {
			continue
		}
		rep := chooseGroupRepresentative(items)
		children := make([]financialevent.FinancialEvent, 0, len(items))
		for _, item := range items {
			child := *toFinancialEventEntity(&item, nil)
			child.GroupKey = key
			child.GroupSize = len(items)
			child.DisplayAmount = child.Amount
			children = append(children, child)
		}
		entity := *toFinancialEventEntity(&rep, nil)
		entity.GroupKey = key
		entity.GroupSize = len(items)
		entity.DisplayAmount = groupDisplayAmount(items, rep.Amount)
		entity.Children = children
		result = append(result, entity)
	}
	return &PageResult[financialevent.FinancialEvent]{List: result, Page: page, PageSize: pageSize, Total: total}, nil
}

func chooseGroupRepresentative(items []financialeventdb.FinancialEvent) financialeventdb.FinancialEvent {
	for _, item := range items {
		if item.EventType == EventTypeExpense {
			return item
		}
	}
	for _, item := range items {
		if item.EventType == EventTypeInvestmentSell {
			return item
		}
	}
	return items[len(items)-1]
}

func groupDisplayAmount(items []financialeventdb.FinancialEvent, fallback decimal.Decimal) decimal.Decimal {
	hasExpense := false
	hasReceivable := false
	currency := ""
	total := decimal.Zero
	for _, item := range items {
		if currency == "" {
			currency = item.Currency
		}
		if item.Currency != currency {
			return fallback
		}
		if item.EventType == EventTypeExpense {
			hasExpense = true
			total = total.Add(item.Amount)
		}
		if item.EventType == EventTypeReceivableCreate {
			hasReceivable = true
			total = total.Add(item.Amount)
		}
	}
	if hasExpense && hasReceivable {
		return total
	}
	return fallback
}

func ExportFinancialEvents(userId int64, query financialevent.FinancialEventQuery) ([]financialevent.FinancialEvent, error) {
	events, err := financialeventdb.ListAllFinancialEventsByUserIds([]int64{userId}, financialeventdb.FinancialEventQuery{
		EventType:           query.EventType,
		StartTime:           query.StartTime,
		EndTime:             query.EndTime,
		Currency:            query.Currency,
		CategoryId:          query.CategoryId,
		ChannelId:           query.ChannelId,
		BucketId:            query.BucketId,
		Keyword:             query.Keyword,
		IncludeInStatistics: query.IncludeInStatistics,
	})
	if err != nil {
		return nil, err
	}
	result := make([]financialevent.FinancialEvent, 0, len(events))
	for _, item := range events {
		result = append(result, *toFinancialEventEntity(&item, nil))
	}
	return result, nil
}

func GetFinancialEvent(userId int64, id int64) (*financialevent.FinancialEvent, error) {
	eventDb, err := financialeventdb.GetFinancialEvent(id, userId)
	if err != nil {
		return nil, err
	}

	entryDbs, err := ledgerentrydb.ListLedgerEntriesByEvent(id, userId)
	if err != nil {
		return nil, err
	}

	entries := make([]ledgerentry.LedgerEntry, 0, len(entryDbs))
	for _, item := range entryDbs {
		entries = append(entries, toLedgerEntryEntity(&item))
	}

	return toFinancialEventEntity(eventDb, entries), nil
}

func toFinancialEventEntity(eventDb *financialeventdb.FinancialEvent, entries []ledgerentry.LedgerEntry) *financialevent.FinancialEvent {
	var eventGroupId int64
	if eventDb.EventGroupId != nil {
		eventGroupId = *eventDb.EventGroupId
	}
	var relatedFinancialEventId int64
	if eventDb.RelatedFinancialEventId != nil {
		relatedFinancialEventId = *eventDb.RelatedFinancialEventId
	}
	var categoryId int64
	if eventDb.CategoryId != nil {
		categoryId = *eventDb.CategoryId
	}
	var channelId int64
	if eventDb.ChannelId != nil {
		channelId = *eventDb.ChannelId
	}
	var remark string
	if eventDb.Remark != nil {
		remark = *eventDb.Remark
	}

	return &financialevent.FinancialEvent{
		Id:                      eventDb.Id,
		UserId:                  eventDb.UserId,
		EventGroupId:            eventGroupId,
		RelatedFinancialEventId: relatedFinancialEventId,
		EventType:               eventDb.EventType,
		Description:             eventDb.Description,
		CategoryId:              categoryId,
		ChannelId:               channelId,
		EventTime:               eventDb.EventTime.String(),
		Currency:                eventDb.Currency,
		Amount:                  eventDb.Amount,
		IncludeInStatistics:     eventDb.IncludeInStatistics,
		Source:                  eventDb.Source,
		Status:                  eventDb.Status,
		Remark:                  remark,
		LedgerEntries:           entries,
	}
}

// ListOutstandingForBucket 返回该 bucket 上 createEventType 类型尚未结清的事件。
// bucketId == 0 时跨所有桶。
func ListOutstandingForBucket(userId int64, bucketId int64, createEventType string) ([]financialevent.OutstandingItem, error) {
	if createEventType == "" {
		return nil, errors.New("event type is required")
	}

	var creates []financialeventdb.FinancialEvent
	var err error
	if bucketId == 0 {
		creates, err = financialeventdb.ListEventsByTypeAcrossBuckets(userId, createEventType)
	} else {
		creates, err = financialeventdb.ListEventsByBucketAndType(userId, bucketId, createEventType)
	}
	if err != nil {
		return nil, err
	}
	if len(creates) == 0 {
		return []financialevent.OutstandingItem{}, nil
	}

	createIds := make([]int64, 0, len(creates))
	for _, ev := range creates {
		createIds = append(createIds, ev.Id)
	}
	collects, err := financialeventdb.ListEventsByRelatedIds(userId, createIds)
	if err != nil {
		return nil, err
	}
	collectedSum := make(map[int64]decimal.Decimal, len(creates))
	for _, ev := range collects {
		if ev.RelatedFinancialEventId == nil {
			continue
		}
		key := *ev.RelatedFinancialEventId
		collectedSum[key] = collectedSum[key].Add(ev.Amount)
	}

	// 一次性查所有 create 事件的 counterparty bucket
	counterRole := counterEntryRoleForCreate(createEventType)
	bucketByEvent := map[int64]financialeventdb.EventBucketInfo{}
	if counterRole != "" {
		infos, err := financialeventdb.ListEventBucketInfo(userId, createIds, []string{counterRole})
		if err != nil {
			return nil, err
		}
		for _, info := range infos {
			bucketByEvent[info.EventId] = info
		}
	}

	result := make([]financialevent.OutstandingItem, 0, len(creates))
	for _, ev := range creates {
		outstanding := ev.Amount.Sub(collectedSum[ev.Id])
		if !outstanding.IsPositive() {
			continue
		}
		groupId := int64(0)
		if ev.EventGroupId != nil {
			groupId = *ev.EventGroupId
		}
		info := bucketByEvent[ev.Id]
		result = append(result, financialevent.OutstandingItem{
			Id:                ev.Id,
			EventTime:         ev.EventTime.String(),
			EventType:         ev.EventType,
			Description:       ev.Description,
			Currency:          ev.Currency,
			Amount:            ev.Amount,
			OutstandingAmount: outstanding,
			EventGroupId:      groupId,
			BucketId:          info.BucketId,
			BucketName:        info.BucketName,
		})
	}
	return result, nil
}

// counterEntryRoleForCreate 返回 create 类事件中"对手桶"那条 ledger_entry 的 entry_role。
func counterEntryRoleForCreate(createEventType string) string {
	switch createEventType {
	case "receivable_create":
		return "receivable_create"
	case "deposit_create":
		return "deposit_create"
	case "loan_out":
		return "loan_out"
	}
	return ""
}

// GetOutstandingSummary 一次性返回三类未结清事件。
func GetOutstandingSummary(userId int64) (*financialevent.OutstandingSummary, error) {
	receivables, err := ListOutstandingForBucket(userId, 0, "receivable_create")
	if err != nil {
		return nil, err
	}
	deposits, err := ListOutstandingForBucket(userId, 0, "deposit_create")
	if err != nil {
		return nil, err
	}
	loansOut, err := ListOutstandingForBucket(userId, 0, "loan_out")
	if err != nil {
		return nil, err
	}
	return &financialevent.OutstandingSummary{
		Receivables: receivables,
		Deposits:    deposits,
		LoansOut:    loansOut,
	}, nil
}

// GroupedEvent 同 group 视图中的一条事件，附带 outstanding 状态。
type GroupedEvent struct {
	Id                  int64           `json:"id"`
	EventType           string          `json:"eventType"`
	Description         string          `json:"description"`
	Currency            string          `json:"currency"`
	Amount              decimal.Decimal `json:"amount"`
	EventTime           string          `json:"eventTime"`
	IncludeInStatistics bool            `json:"includeInStatistics"`
	OutstandingAmount   decimal.Decimal `json:"outstandingAmount"`
	SettledAmount       decimal.Decimal `json:"settledAmount"`
	IsCurrent           bool            `json:"isCurrent"`
}

func isCreateLikeEvent(t string) bool {
	return t == "receivable_create" || t == "deposit_create" || t == "loan_out"
}
func isCollectLikeEvent(t string) bool {
	return t == "receivable_collect" || t == "deposit_refund" || t == "loan_collect"
}

// GetEventGroup 列出同 group 所有 event，create 类附带未结清状态。
func GetEventGroup(userId int64, groupId int64, currentEventId int64) ([]GroupedEvent, error) {
	if groupId == 0 {
		return []GroupedEvent{}, nil
	}
	events, err := financialeventdb.ListEventsByGroupId(userId, groupId)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return []GroupedEvent{}, nil
	}

	settledByCreate := make(map[int64]decimal.Decimal)
	for _, ev := range events {
		if !isCollectLikeEvent(ev.EventType) {
			continue
		}
		if ev.RelatedFinancialEventId == nil {
			continue
		}
		key := *ev.RelatedFinancialEventId
		settledByCreate[key] = settledByCreate[key].Add(ev.Amount)
	}

	result := make([]GroupedEvent, 0, len(events))
	for _, ev := range events {
		outstanding := decimal.Zero
		settled := decimal.Zero
		if isCreateLikeEvent(ev.EventType) {
			settled = settledByCreate[ev.Id]
			outstanding = ev.Amount.Sub(settled)
		}
		result = append(result, GroupedEvent{
			Id:                  ev.Id,
			EventType:           ev.EventType,
			Description:         ev.Description,
			Currency:            ev.Currency,
			Amount:              ev.Amount,
			EventTime:           ev.EventTime.String(),
			IncludeInStatistics: ev.IncludeInStatistics,
			OutstandingAmount:   outstanding,
			SettledAmount:       settled,
			IsCurrent:           ev.Id == currentEventId,
		})
	}
	return result, nil
}
