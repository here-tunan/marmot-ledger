package chantemplatedb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type ChannelTemplate struct {
	Id                  int64           `json:"id" xorm:"pk autoincr 'id'"`
	ChannelCode         string          `json:"channelCode" xorm:"'channel_code'"`
	Name                string          `json:"name" xorm:"'name'"`
	ChannelType         string          `json:"channelType" xorm:"'channel_type'"`
	ProviderCode        string          `json:"providerCode" xorm:"'provider_code'"`
	SupportedEventTypes string          `json:"supportedEventTypes" xorm:"'supported_event_types'"`
	Icon                string          `json:"icon" xorm:"'icon'"`
	Sort                int             `json:"sort" xorm:"'sort'"`
	Enabled             bool            `json:"enabled" xorm:"'enabled'"`
	CreatedAt           model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt           model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

func (ChannelTemplate) TableName() string {
	return "channel_template"
}

type ChannelTemplateQuery struct {
	Enabled *bool
}

func ListChannelTemplates(query ChannelTemplateQuery) ([]ChannelTemplate, error) {
	templates := make([]ChannelTemplate, 0)
	session := infrastructure.Mysql.NewSession()
	defer session.Close()

	applyQuery(session, query)
	err := session.Asc("sort", "id").Find(&templates)
	return templates, err
}

func GetChannelTemplate(id int64) (*ChannelTemplate, error) {
	template := &ChannelTemplate{}
	has, err := infrastructure.Mysql.ID(id).Get(template)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("channel template not found")
	}
	return template, nil
}

func GetChannelTemplateByCode(code string) (*ChannelTemplate, error) {
	template := &ChannelTemplate{}
	has, err := infrastructure.Mysql.Where("channel_code = ?", code).Get(template)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("channel template not found")
	}
	return template, nil
}

func InsertChannelTemplate(template *ChannelTemplate) error {
	_, err := infrastructure.Mysql.InsertOne(template)
	return err
}

func UpdateChannelTemplate(template *ChannelTemplate) error {
	_, err := infrastructure.Mysql.ID(template.Id).Omit("CreatedAt").Update(template)
	return err
}

func applyQuery(session *xorm.Session, query ChannelTemplateQuery) {
	if query.Enabled != nil {
		session.And("enabled = ?", *query.Enabled)
	}
}

// CheckCodeExists 检查模板代码是否已存在
func CheckCodeExists(channelCode string, excludeId int64) (bool, error) {
	session := infrastructure.Mysql.Where("channel_code = ?", channelCode)
	if excludeId > 0 {
		session.And("id != ?", excludeId)
	}
	return session.Exist(&ChannelTemplate{})
}

// BuildTemplateMap 构建 code -> template 的映射
func BuildTemplateMap() (map[string]*ChannelTemplate, error) {
	templates, err := ListChannelTemplates(ChannelTemplateQuery{})
	if err != nil {
		return nil, err
	}
	result := make(map[string]*ChannelTemplate)
	for i := range templates {
		result[strings.ToUpper(templates[i].ChannelCode)] = &templates[i]
	}
	return result, nil
}
