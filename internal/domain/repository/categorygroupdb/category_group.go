package categorygroupdb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type CategoryGroup struct {
	Id        int64           `json:"id" xorm:"pk autoincr 'id'"`
	GroupCode string          `json:"groupCode" xorm:"'group_code'"`
	Name      string          `json:"name" xorm:"'name'"`
	Type      string          `json:"type" xorm:"'type'"`
	Icon      string          `json:"icon" xorm:"'icon'"`
	Color     string          `json:"color" xorm:"'color'"`
	Sort      int             `json:"sort" xorm:"'sort'"`
	Enabled   bool            `json:"enabled" xorm:"'enabled'"`
	CreatedAt model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type CategoryGroupQuery struct {
	Type    string
	Enabled *bool
}

func (CategoryGroup) TableName() string {
	return "category_group"
}

func ListCategoryGroups(query CategoryGroupQuery) ([]CategoryGroup, error) {
	groups := make([]CategoryGroup, 0)
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	applyCategoryGroupQuery(session, query)
	err := session.Asc("sort", "id").Find(&groups)
	return groups, err
}

func GetCategoryGroupById(session *xorm.Session, id int64) (*CategoryGroup, error) {
	group := &CategoryGroup{}
	has, err := session.Where("id = ?", id).Get(group)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("category group not found")
	}
	return group, nil
}

func GetCategoryGroup(id int64) (*CategoryGroup, error) {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	return GetCategoryGroupById(session, id)
}

func applyCategoryGroupQuery(session *xorm.Session, query CategoryGroupQuery) {
	if strings.TrimSpace(query.Type) != "" {
		session.And("type = ?", strings.TrimSpace(query.Type))
	}
	if query.Enabled != nil {
		session.And("enabled = ?", *query.Enabled)
	}
}
