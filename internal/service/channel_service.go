package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/channel"
	"marmot-ledger/internal/domain/repository/channeldb"
	"marmot-ledger/internal/domain/repository/chantemplatedb"
	"marmot-ledger/internal/infrastructure"
	"strings"
)

func ListChannels(userId int64, query channel.ChannelQuery) ([]channel.Channel, error) {
	channels, err := channeldb.ListChannels(userId, channeldb.ChannelQuery{ChannelType: query.ChannelType, EventType: query.EventType, IsActive: query.IsActive})
	if err != nil {
		return nil, err
	}
	result := make([]channel.Channel, 0, len(channels))
	for _, item := range channels {
		result = append(result, toChannelEntity(&item))
	}
	return result, nil
}

func CreateChannel(userId int64, channelInfo *channel.Channel) (*channel.Channel, error) {
	if err := validateChannel(channelInfo); err != nil {
		return nil, err
	}
	channelDb := toChannelDb(userId, channelInfo)
	channelDb.IsActive = true
	if err := channeldb.InsertChannel(channelDb); err != nil {
		return nil, err
	}
	return GetChannel(userId, channelDb.Id)
}

func GetChannel(userId int64, id int64) (*channel.Channel, error) {
	channelDb, err := channeldb.GetChannel(id, userId)
	if err != nil {
		return nil, err
	}
	entity := toChannelEntity(channelDb)
	return &entity, nil
}

func UpdateChannel(userId int64, id int64, channelInfo *channel.Channel) (*channel.Channel, error) {
	if err := validateChannel(channelInfo); err != nil {
		return nil, err
	}
	if _, err := channeldb.GetChannel(id, userId); err != nil {
		return nil, err
	}
	channelDb := toChannelDb(userId, channelInfo)
	channelDb.Id = id
	if err := channeldb.UpdateChannel(channelDb); err != nil {
		return nil, err
	}
	return GetChannel(userId, id)
}

func DeleteChannel(userId int64, id int64) (int64, error) {
	count, err := channeldb.CountEventsByChannel(userId, id)
	if err != nil {
		return 0, err
	}
	if err := channeldb.SoftDeleteChannel(id, userId); err != nil {
		return 0, err
	}
	return count, nil
}

func CheckChannelUsage(userId int64, id int64) (int64, error) {
	return channeldb.CountEventsByChannel(userId, id)
}

func ImportChannelTemplates(userId int64, templateIds []int64) ([]channel.Channel, error) {
	if len(templateIds) == 0 {
		return nil, errors.New("no templates selected")
	}
	existing, err := channeldb.ExistingTemplateIds(userId, templateIds)
	if err != nil {
		return nil, err
	}
	result := make([]channel.Channel, 0, len(templateIds))
	for _, templateId := range templateIds {
		if existing[templateId] {
			continue
		}
		templateDb, err := chantemplatedb.GetChannelTemplate(templateId)
		if err != nil {
			return nil, err
		}
		created, err := CreateChannel(userId, &channel.Channel{
			ChannelTemplateId:   templateDb.Id,
			Name:                templateDb.Name,
			ChannelType:         templateDb.ChannelType,
			ProviderCode:        templateDb.ProviderCode,
			SupportedEventTypes: templateDb.SupportedEventTypes,
			Icon:                templateDb.Icon,
			Sort:                templateDb.Sort,
			IsActive:            true,
		})
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate") || strings.Contains(err.Error(), "duplicate") {
				continue
			}
			return nil, err
		}
		result = append(result, *created)
	}
	return result, nil
}

func EnsureDefaultChannels(userId int64) error {
	count, err := channeldb.CountChannels(userId)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	enabled := true
	templates, err := chantemplatedb.ListChannelTemplates(chantemplatedb.ChannelTemplateQuery{Enabled: &enabled})
	if err != nil {
		return err
	}
	items := make([]channeldb.Channel, 0, len(templates))
	for _, template := range templates {
		templateId := template.Id
		items = append(items, channeldb.Channel{
			UserId:              userId,
			ChannelTemplateId:   &templateId,
			Name:                template.Name,
			ChannelType:         strings.ToLower(strings.TrimSpace(template.ChannelType)),
			ProviderCode:        strings.TrimSpace(template.ProviderCode),
			SupportedEventTypes: normalizeEventTypes(template.SupportedEventTypes),
			Icon:                template.Icon,
			Sort:                template.Sort,
			IsActive:            true,
			IsDeleted:           false,
		})
	}
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	return channeldb.InsertChannels(session, items)
}

func validateChannel(channelInfo *channel.Channel) error {
	if channelInfo == nil {
		return errors.New("channel is required")
	}
	if strings.TrimSpace(channelInfo.Name) == "" {
		return errors.New("channel name is required")
	}
	if strings.TrimSpace(channelInfo.ChannelType) == "" {
		return errors.New("channel type is required")
	}
	return nil
}

func toChannelDb(userId int64, channelInfo *channel.Channel) *channeldb.Channel {
	var templateId *int64
	if channelInfo.ChannelTemplateId > 0 {
		id := channelInfo.ChannelTemplateId
		templateId = &id
	}
	return &channeldb.Channel{
		Id:                  channelInfo.Id,
		UserId:              userId,
		ChannelTemplateId:   templateId,
		Name:                strings.TrimSpace(channelInfo.Name),
		ChannelType:         strings.ToLower(strings.TrimSpace(channelInfo.ChannelType)),
		ProviderCode:        strings.TrimSpace(channelInfo.ProviderCode),
		SupportedEventTypes: normalizeEventTypes(channelInfo.SupportedEventTypes),
		Icon:                strings.TrimSpace(channelInfo.Icon),
		Sort:                channelInfo.Sort,
		IsActive:            channelInfo.IsActive,
	}
}

func toChannelEntity(channelDb *channeldb.Channel) channel.Channel {
	var templateId int64
	if channelDb.ChannelTemplateId != nil {
		templateId = *channelDb.ChannelTemplateId
	}
	return channel.Channel{
		Id:                  channelDb.Id,
		UserId:              channelDb.UserId,
		ChannelTemplateId:   templateId,
		Name:                channelDb.Name,
		ChannelType:         channelDb.ChannelType,
		ProviderCode:        channelDb.ProviderCode,
		SupportedEventTypes: channelDb.SupportedEventTypes,
		Icon:                channelDb.Icon,
		Sort:                channelDb.Sort,
		IsActive:            channelDb.IsActive,
	}
}
