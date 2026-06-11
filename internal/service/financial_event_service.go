package service

import (
	"marmot-ledger/internal/domain/entity/financialevent"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/domain/repository/ledgerentrydb"
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
		CategoryId:          query.CategoryId,
		CategoryGroupId:     query.CategoryGroupId,
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
	var categoryGroupId int64
	if eventDb.CategoryGroupId != nil {
		categoryGroupId = *eventDb.CategoryGroupId
	}
	var channelType string
	if eventDb.ChannelType != nil {
		channelType = *eventDb.ChannelType
	}
	var channelAccountId int64
	if eventDb.ChannelAccountId != nil {
		channelAccountId = *eventDb.ChannelAccountId
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
		CategoryGroupId:         categoryGroupId,
		ChannelType:             channelType,
		ChannelAccountId:        channelAccountId,
		EventTime:               eventDb.EventTime.String(),
		Currency:                eventDb.Currency,
		Amount:                  eventDb.Amount,
		BaseCurrency:            eventDb.BaseCurrency,
		BaseAmount:              eventDb.BaseAmount,
		ExchangeRate:            eventDb.ExchangeRate,
		IncludeInStatistics:     eventDb.IncludeInStatistics,
		Source:                  eventDb.Source,
		Status:                  eventDb.Status,
		Remark:                  remark,
		LedgerEntries:           entries,
	}
}
