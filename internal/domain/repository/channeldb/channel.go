package channeldb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type Channel struct {
	Id                  int64           `json:"id" xorm:"pk autoincr 'id'"`
	UserId              int64           `json:"userId" xorm:"'user_id'"`
	ChannelTemplateId   *int64          `json:"channelTemplateId" xorm:"'channel_template_id'"`
	Name                string          `json:"name" xorm:"'name'"`
	ChannelType         string          `json:"channelType" xorm:"'channel_type'"`
	ProviderCode        string          `json:"providerCode" xorm:"'provider_code'"`
	SupportedEventTypes string          `json:"supportedEventTypes" xorm:"'supported_event_types'"`
	Icon                string          `json:"icon" xorm:"'icon'"`
	Sort                int             `json:"sort" xorm:"'sort'"`
	IsActive            bool            `json:"isActive" xorm:"'is_active'"`
	IsDeleted           bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt           model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt           model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type ChannelQuery struct {
	ChannelType string
	EventType   string
	IsActive    *bool
}

func (Channel) TableName() string {
	return "personal_channel"
}

func InsertChannel(channel *Channel) error {
	_, err := infrastructure.Mysql.InsertOne(channel)
	return err
}

func InsertChannels(session *xorm.Session, channels []Channel) error {
	if len(channels) == 0 {
		return nil
	}
	_, err := session.Insert(&channels)
	return err
}

func ListChannels(userId int64, query ChannelQuery) ([]Channel, error) {
	channels := make([]Channel, 0)
	session := infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0)
	applyQuery(session, query)
	err := session.Asc("sort", "id").Find(&channels)
	return channels, err
}

func CountChannels(userId int64) (int64, error) {
	return infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0).Count(&Channel{})
}

func GetChannel(id int64, userId int64) (*Channel, error) {
	channel := &Channel{}
	has, err := infrastructure.Mysql.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).Get(channel)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("channel not found")
	}
	return channel, nil
}

func UpdateChannel(channel *Channel) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", channel.Id, channel.UserId, 0).
		Cols("name", "channel_type", "provider_code", "supported_event_types", "icon", "sort", "is_active").
		Update(channel)
	return err
}

func SoftDeleteChannel(id int64, userId int64) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).
		Cols("is_deleted").
		Update(&Channel{IsDeleted: true})
	return err
}

func CountEventsByChannel(userId int64, channelId int64) (int64, error) {
	type row struct {
		Count int64 `xorm:"'count'"`
	}
	result := &row{}
	_, err := infrastructure.Mysql.SQL("SELECT COUNT(*) AS count FROM financial_event WHERE user_id = ? AND channel_id = ? AND is_deleted = 0", userId, channelId).Get(result)
	return result.Count, err
}

func ExistingTemplateIds(userId int64, templateIds []int64) (map[int64]bool, error) {
	result := map[int64]bool{}
	if len(templateIds) == 0 {
		return result, nil
	}
	channels := make([]Channel, 0)
	if err := infrastructure.Mysql.Where("user_id = ? AND is_deleted = 0", userId).In("channel_template_id", templateIds).Find(&channels); err != nil {
		return nil, err
	}
	for _, item := range channels {
		if item.ChannelTemplateId != nil {
			result[*item.ChannelTemplateId] = true
		}
	}
	return result, nil
}

func applyQuery(session *xorm.Session, query ChannelQuery) {
	if strings.TrimSpace(query.ChannelType) != "" {
		session.And("channel_type = ?", strings.ToLower(strings.TrimSpace(query.ChannelType)))
	}
	if strings.TrimSpace(query.EventType) != "" {
		eventType := strings.ToLower(strings.TrimSpace(query.EventType))
		session.And("FIND_IN_SET(?, supported_event_types)", eventType)
	}
	if query.IsActive != nil {
		session.And("is_active = ?", *query.IsActive)
	}
}
