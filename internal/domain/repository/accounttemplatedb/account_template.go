package accounttemplatedb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type AccountTemplate struct {
	Id           int64           `json:"id" xorm:"pk autoincr 'id'"`
	ProviderCode string          `json:"providerCode" xorm:"'provider_code'"`
	Name         string          `json:"name" xorm:"'name'"`
	Type         string          `json:"type" xorm:"'type'"`
	Icon         string          `json:"icon" xorm:"'icon'"`
	Color        string          `json:"color" xorm:"'color'"`
	Sort         int             `json:"sort" xorm:"'sort'"`
	Enabled      bool            `json:"enabled" xorm:"'enabled'"`
	CreatedAt    model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt    model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

func (AccountTemplate) TableName() string {
	return "account_template"
}

type AccountTemplateQuery struct {
	Enabled *bool
}

func ListAccountTemplates(query AccountTemplateQuery) ([]AccountTemplate, error) {
	templates := make([]AccountTemplate, 0)
	session := infrastructure.Mysql.NewSession()
	defer session.Close()

	applyQuery(session, query)
	err := session.Asc("sort", "id").Find(&templates)
	return templates, err
}

func GetAccountTemplate(id int64) (*AccountTemplate, error) {
	template := &AccountTemplate{}
	has, err := infrastructure.Mysql.ID(id).Get(template)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("account template not found")
	}
	return template, nil
}

func GetAccountTemplateByCode(providerCode string) (*AccountTemplate, error) {
	template := &AccountTemplate{}
	has, err := infrastructure.Mysql.Where("provider_code = ?", providerCode).Get(template)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("account template not found")
	}
	return template, nil
}

func InsertAccountTemplate(template *AccountTemplate) error {
	_, err := infrastructure.Mysql.InsertOne(template)
	return err
}

func UpdateAccountTemplate(template *AccountTemplate) error {
	_, err := infrastructure.Mysql.ID(template.Id).Omit("CreatedAt").Update(template)
	return err
}

func applyQuery(session *xorm.Session, query AccountTemplateQuery) {
	if query.Enabled != nil {
		session.And("enabled = ?", *query.Enabled)
	}
}

// CheckCodeExists 检查模板代码是否已存在
func CheckCodeExists(providerCode string, excludeId int64) (bool, error) {
	session := infrastructure.Mysql.Where("provider_code = ?", providerCode)
	if excludeId > 0 {
		session.And("id != ?", excludeId)
	}
	return session.Exist(&AccountTemplate{})
}

// BuildTemplateMap 构建 code -> template 的映射
func BuildTemplateMap() (map[string]*AccountTemplate, error) {
	templates, err := ListAccountTemplates(AccountTemplateQuery{})
	if err != nil {
		return nil, err
	}
	result := make(map[string]*AccountTemplate)
	for i := range templates {
		result[strings.ToUpper(templates[i].ProviderCode)] = &templates[i]
	}
	return result, nil
}
