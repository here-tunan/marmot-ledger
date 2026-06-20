package categorytemplatedb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type CategoryTemplate struct {
	Id           int64           `json:"id" xorm:"pk autoincr 'id'"`
	TemplateCode string          `json:"templateCode" xorm:"'template_code'"`
	Name         string          `json:"name" xorm:"'name'"`
	Type         string          `json:"type" xorm:"'type'"`
	Icon         string          `json:"icon" xorm:"'icon'"`
	Color        string          `json:"color" xorm:"'color'"`
	Sort         int             `json:"sort" xorm:"'sort'"`
	Enabled      bool            `json:"enabled" xorm:"'enabled'"`
	CreatedAt    model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt    model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

func (CategoryTemplate) TableName() string {
	return "category_template"
}

type CategoryTemplateQuery struct {
	Enabled *bool
	Type    string
}

func ListCategoryTemplates(query CategoryTemplateQuery) ([]CategoryTemplate, error) {
	templates := make([]CategoryTemplate, 0)
	session := infrastructure.Mysql.NewSession()
	defer session.Close()

	applyQuery(session, query)
	err := session.Asc("type", "sort", "id").Find(&templates)
	return templates, err
}

func GetCategoryTemplate(id int64) (*CategoryTemplate, error) {
	template := &CategoryTemplate{}
	has, err := infrastructure.Mysql.ID(id).Get(template)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("category template not found")
	}
	return template, nil
}

func GetCategoryTemplateByCode(templateCode string) (*CategoryTemplate, error) {
	template := &CategoryTemplate{}
	has, err := infrastructure.Mysql.Where("template_code = ?", templateCode).Get(template)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("category template not found")
	}
	return template, nil
}

func InsertCategoryTemplate(template *CategoryTemplate) error {
	_, err := infrastructure.Mysql.InsertOne(template)
	return err
}

func UpdateCategoryTemplate(template *CategoryTemplate) error {
	_, err := infrastructure.Mysql.ID(template.Id).Omit("CreatedAt").Update(template)
	return err
}

func applyQuery(session *xorm.Session, query CategoryTemplateQuery) {
	if query.Enabled != nil {
		session.And("enabled = ?", *query.Enabled)
	}
	if strings.TrimSpace(query.Type) != "" {
		session.And("type = ?", strings.TrimSpace(query.Type))
	}
}

// CheckCodeExists 检查模板代码是否已存在
func CheckCodeExists(templateCode string, excludeId int64) (bool, error) {
	session := infrastructure.Mysql.Where("template_code = ?", templateCode)
	if excludeId > 0 {
		session.And("id != ?", excludeId)
	}
	return session.Exist(&CategoryTemplate{})
}

// BatchGetByIds 批量获取模板
func BatchGetByIds(ids []int64) ([]CategoryTemplate, error) {
	if len(ids) == 0 {
		return []CategoryTemplate{}, nil
	}
	templates := make([]CategoryTemplate, 0)
	err := infrastructure.Mysql.In("id", ids).Find(&templates)
	return templates, err
}
