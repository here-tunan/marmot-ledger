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
	ChannelId               *int64          `json:"channelId" xorm:"'channel_id'"`
	EventTime               model.LocalTime `json:"eventTime" xorm:"'event_time'"`
	Currency                string          `json:"currency" xorm:"'currency'"`
	Amount                  decimal.Decimal `json:"amount" xorm:"'amount'"`
	IncludeInStatistics     bool            `json:"includeInStatistics" xorm:"'include_in_statistics'"`
	Source                  string          `json:"source" xorm:"'source'"`
	Status                  string          `json:"status" xorm:"'status'"`
	Remark                  *string         `json:"remark" xorm:"'remark'"`
	IsDeleted               bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt               model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt               model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type FinancialEventQuery struct {
	EventType           string
	StartTime           string
	EndTime             string
	Currency            string
	CategoryId          int64
	ChannelId           int64
	BucketId            int64
	Keyword             string
	IncludeInStatistics *bool
	Page                int
	PageSize            int
}

type GroupPageItem struct {
	GroupKey int64           `xorm:"'group_key'"`
	MaxTime  model.LocalTime `xorm:"'max_time'"`
}

func (FinancialEvent) TableName() string {
	return "financial_event"
}

func InsertFinancialEvent(session *xorm.Session, event *FinancialEvent) error {
	_, err := session.InsertOne(event)
	return err
}

func ListFinancialEvents(userId int64, query FinancialEventQuery) ([]FinancialEvent, int64, error) {
	return ListFinancialEventsByUserIds([]int64{userId}, query)
}

func ListFinancialEventsByUserIds(userIds []int64, query FinancialEventQuery) ([]FinancialEvent, int64, error) {
	events := make([]FinancialEvent, 0)
	inClause, params := userIdsInClause(userIds)
	if len(userIds) == 0 {
		return events, 0, nil
	}

	page := query.Page
	pageSize := query.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	countEvents := make([]FinancialEvent, 0)
	countSession := infrastructure.Mysql.Where("user_id IN ("+inClause+") AND is_deleted = ?", append(params, 0)...)
	applyFinancialEventQuery(countSession, query, userIds)
	if err := countSession.Find(&countEvents); err != nil {
		return nil, 0, err
	}
	total := int64(len(countEvents))

	session := infrastructure.Mysql.Where("user_id IN ("+inClause+") AND is_deleted = ?", append(params, 0)...)
	applyFinancialEventQuery(session, query, userIds)
	err := session.Desc("event_time", "id").Limit(pageSize, (page-1)*pageSize).Find(&events)
	return events, total, err
}

func ListAllFinancialEventsByUserIds(userIds []int64, query FinancialEventQuery) ([]FinancialEvent, error) {
	events := make([]FinancialEvent, 0)
	inClause, params := userIdsInClause(userIds)
	if len(userIds) == 0 {
		return events, nil
	}
	session := infrastructure.Mysql.Where("user_id IN ("+inClause+") AND is_deleted = ?", append(params, 0)...)
	applyFinancialEventQuery(session, query, userIds)
	return events, session.Desc("event_time", "id").Find(&events)
}

func ListFinancialEventsByGroupKeys(userIds []int64, groupKeys []int64) ([]FinancialEvent, error) {
	events := make([]FinancialEvent, 0)
	if len(userIds) == 0 || len(groupKeys) == 0 {
		return events, nil
	}
	userInClause, userParams := userIdsInClause(userIds)
	groupPlaceholders := make([]string, 0, len(groupKeys))
	params := make([]any, 0, len(userParams)+len(groupKeys))
	params = append(params, userParams...)
	for _, key := range groupKeys {
		groupPlaceholders = append(groupPlaceholders, "?")
		params = append(params, key)
	}
	sql := "SELECT * FROM financial_event WHERE user_id IN (" + userInClause + ") AND is_deleted = 0 AND COALESCE(event_group_id, id) IN (" + strings.Join(groupPlaceholders, ",") + ") ORDER BY event_time ASC, id ASC"
	err := infrastructure.Mysql.SQL(sql, params...).Find(&events)
	return events, err
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

func GetFinancialEventForUpdate(session *xorm.Session, id int64, userId int64) (*FinancialEvent, error) {
	event := &FinancialEvent{}
	has, err := session.SQL("SELECT * FROM financial_event WHERE id = ? AND user_id = ? AND is_deleted = 0 FOR UPDATE", id, userId).Get(event)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("financial event not found")
	}
	return event, nil
}

func UpdateFinancialEvent(session *xorm.Session, event *FinancialEvent) error {
	_, err := session.
		Where("id = ? AND user_id = ? AND is_deleted = ?", event.Id, event.UserId, 0).
		Cols("event_type", "description", "category_id", "channel_id", "event_time", "currency", "amount", "include_in_statistics", "source", "status", "remark", "related_financial_event_id").
		Update(event)
	return err
}

func UpdateEventGroupId(session *xorm.Session, id int64, userId int64, groupId int64) error {
	_, err := session.
		Where("id = ? AND user_id = ?", id, userId).
		Cols("event_group_id").
		Update(&FinancialEvent{EventGroupId: &groupId})
	return err
}

// ListEventsByGroupId 列出同 group 内所有未删除事件，按时间升序。
func ListEventsByGroupId(userId int64, groupId int64) ([]FinancialEvent, error) {
	events := make([]FinancialEvent, 0)
	err := infrastructure.Mysql.
		Where("user_id = ? AND event_group_id = ? AND is_deleted = 0", userId, groupId).
		Asc("event_time", "id").
		Find(&events)
	return events, err
}

// ListEventsByBucketAndType 反查指定 bucket 上某类型的事件，用于 outstanding 计算。
// 通过 ledger_entry.bucket_id 反查 financial_event_id。
func ListEventsByBucketAndType(userId int64, bucketId int64, eventType string) ([]FinancialEvent, error) {
	events := make([]FinancialEvent, 0)
	err := infrastructure.Mysql.SQL(
		"SELECT DISTINCT f.* FROM financial_event f JOIN ledger_entry l ON l.financial_event_id = f.id WHERE f.user_id = ? AND f.event_type = ? AND f.is_deleted = 0 AND l.bucket_id = ? ORDER BY f.event_time DESC, f.id DESC",
		userId, eventType, bucketId,
	).Find(&events)
	return events, err
}

// ListEventsByTypeAcrossBuckets 跨桶查询某类型 create 事件（不限定 bucket）。
func ListEventsByTypeAcrossBuckets(userId int64, eventType string) ([]FinancialEvent, error) {
	events := make([]FinancialEvent, 0)
	err := infrastructure.Mysql.
		Where("user_id = ? AND event_type = ? AND is_deleted = 0", userId, eventType).
		Desc("event_time", "id").
		Find(&events)
	return events, err
}

// EventBucketInfo 记录某条 event 关联的 counterparty bucket_id 和 name（用于 outstanding 列表渲染）。
type EventBucketInfo struct {
	EventId    int64
	BucketId   int64
	BucketName string
}

// ListEventBucketInfo 一次性查多条 event 对应的 counterparty bucket（按 entry_role 反查）。
// counterEntryRoles 支持多个：receivable_create / deposit_create / loan_out / cash_leg 等。
func ListEventBucketInfo(userId int64, eventIds []int64, counterEntryRoles []string) ([]EventBucketInfo, error) {
	infos := make([]EventBucketInfo, 0)
	if len(eventIds) == 0 || len(counterEntryRoles) == 0 {
		return infos, nil
	}
	rows, err := infrastructure.Mysql.DB().Query(
		"SELECT l.financial_event_id, l.bucket_id, b.name FROM ledger_entry l JOIN bucket b ON b.id = l.bucket_id WHERE l.user_id = ? AND l.financial_event_id IN ("+placeholders(len(eventIds))+") AND l.entry_role IN ("+placeholders(len(counterEntryRoles))+")",
		buildArgs(userId, eventIds, counterEntryRoles)...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var info EventBucketInfo
		if err := rows.Scan(&info.EventId, &info.BucketId, &info.BucketName); err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = "?"
	}
	return strings.Join(parts, ",")
}

func buildArgs(userId int64, eventIds []int64, roles []string) []any {
	args := make([]any, 0, 1+len(eventIds)+len(roles))
	args = append(args, userId)
	for _, id := range eventIds {
		args = append(args, id)
	}
	for _, r := range roles {
		args = append(args, r)
	}
	return args
}

// ListEventsByRelatedIds 列出 related_financial_event_id 落在给定 ids 内的所有事件，用于 outstanding 累加。
func ListEventsByRelatedIds(userId int64, relatedIds []int64) ([]FinancialEvent, error) {
	events := make([]FinancialEvent, 0)
	if len(relatedIds) == 0 {
		return events, nil
	}
	err := infrastructure.Mysql.
		Where("user_id = ? AND is_deleted = 0", userId).
		In("related_financial_event_id", relatedIds).
		Find(&events)
	return events, err
}

func SoftDeleteFinancialEvent(session *xorm.Session, id int64, userId int64) error {
	_, err := session.
		Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).
		Cols("is_deleted").
		Update(&FinancialEvent{IsDeleted: true})
	return err
}

func applyFinancialEventQuery(session *xorm.Session, query FinancialEventQuery, userIds []int64) {
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
	if strings.TrimSpace(query.Currency) != "" {
		session.And("currency = ?", strings.ToUpper(strings.TrimSpace(query.Currency)))
	}
	if query.CategoryId != 0 {
		session.And("category_id = ?", query.CategoryId)
	}
	if query.ChannelId != 0 {
		session.And("channel_id = ?", query.ChannelId)
	}
	if query.BucketId != 0 {
		inClause, params := userIdsInClause(userIds)
		args := append([]any{query.BucketId}, params...)
		session.And("id IN (SELECT financial_event_id FROM ledger_entry WHERE bucket_id = ? AND user_id IN ("+inClause+"))", args...)
	}
	if strings.TrimSpace(query.Keyword) != "" {
		keyword := "%" + strings.TrimSpace(query.Keyword) + "%"
		session.And("(description LIKE ? OR remark LIKE ?)", keyword, keyword)
	}
	if query.IncludeInStatistics != nil {
		session.And("include_in_statistics = ?", *query.IncludeInStatistics)
	}
}

func userIdsInClause(userIds []int64) (string, []any) {
	if len(userIds) == 0 {
		return "0", []any{}
	}
	placeholders := make([]string, 0, len(userIds))
	params := make([]any, 0, len(userIds))
	for _, id := range userIds {
		placeholders = append(placeholders, "?")
		params = append(params, id)
	}
	return strings.Join(placeholders, ","), params
}

func mustLocation() *time.Location {
	location, err := time.LoadLocation(model.Loc)
	if err != nil {
		return time.Local
	}
	return location
}
