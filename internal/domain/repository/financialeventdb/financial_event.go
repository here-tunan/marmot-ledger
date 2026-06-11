package financialeventdb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"xorm.io/xorm"
)

type FinancialEvent struct {
	Id                      int64           `json:"id" xorm:"pk autoincr 'id'"`
	UserId                  int64           `json:"userId" xorm:"'user_id'"`
	EventGroupId            *int64          `json:"eventGroupId" xorm:"'event_group_id'"`
	RelatedFinancialEventId *int64          `json:"relatedFinancialEventId" xorm:"'related_financial_event_id'"`
	EventType               string          `json:"eventType" xorm:"'event_type'"`
	Description             string          `json:"description" xorm:"'description'"`
	CategoryId              *int64          `json:"categoryId" xorm:"'category_id'"`
	CategoryGroupId         *int64          `json:"categoryGroupId" xorm:"'category_group_id'"`
	ChannelType             *string         `json:"channelType" xorm:"'channel_type'"`
	ChannelAccountId        *int64          `json:"channelAccountId" xorm:"'channel_account_id'"`
	EventTime               model.LocalTime `json:"eventTime" xorm:"'event_time'"`
	Currency                string          `json:"currency" xorm:"'currency'"`
	Amount                  decimal.Decimal `json:"amount" xorm:"'amount'"`
	BaseCurrency            string          `json:"baseCurrency" xorm:"'base_currency'"`
	BaseAmount              decimal.Decimal `json:"baseAmount" xorm:"'base_amount'"`
	ExchangeRate            decimal.Decimal `json:"exchangeRate" xorm:"'exchange_rate'"`
	IncludeInStatistics     bool            `json:"includeInStatistics" xorm:"'include_in_statistics'"`
	Source                  string          `json:"source" xorm:"'source'"`
	Status                  string          `json:"status" xorm:"'status'"`
	Remark                  *string         `json:"remark" xorm:"'remark'"`
	IsDeleted               bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt               model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt               model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type FinancialEventQuery struct {
	EventType string
	StartTime string
	EndTime   string
	Page      int
	PageSize  int
}

func (FinancialEvent) TableName() string {
	return "financial_event"
}

func InsertFinancialEvent(session *xorm.Session, event *FinancialEvent) error {
	_, err := session.InsertOne(event)
	return err
}

func ListFinancialEvents(userId int64, query FinancialEventQuery) ([]FinancialEvent, int64, error) {
	events := make([]FinancialEvent, 0)

	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	countSession := infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0)
	applyFinancialEventQuery(countSession, query)
	total, err := countSession.Count(&FinancialEvent{})
	if err != nil {
		return nil, 0, err
	}

	session := infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0)
	applyFinancialEventQuery(session, query)
	err = session.Desc("event_time", "id").Limit(pageSize, (page-1)*pageSize).Find(&events)
	return events, total, err
}

func GetFinancialEvent(id int64, userId int64) (*FinancialEvent, error) {
	events := make([]FinancialEvent, 0, 1)
	err := infrastructure.Mysql.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).Limit(1).Find(&events)
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return nil, errors.New("financial event not found")
	}
	return &events[0], nil
}

func applyFinancialEventQuery(session *xorm.Session, query FinancialEventQuery) {
	if strings.TrimSpace(query.EventType) != "" {
		session.And("event_type = ?", strings.TrimSpace(query.EventType))
	}
	if strings.TrimSpace(query.StartTime) != "" {
		if startTime, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(query.StartTime), mustLocation()); err == nil {
			session.And("event_time >= ?", startTime)
		}
	}
	if strings.TrimSpace(query.EndTime) != "" {
		if endTime, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(query.EndTime), mustLocation()); err == nil {
			session.And("event_time <= ?", endTime)
		}
	}
}

func mustLocation() *time.Location {
	location, err := time.LoadLocation(model.Loc)
	if err != nil {
		return time.Local
	}
	return location
}
