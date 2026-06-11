package categorydb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type Category struct {
	Id              int64           `json:"id" xorm:"pk autoincr 'id'"`
	UserId          int64           `json:"userId" xorm:"'user_id'"`
	Name            string          `json:"name" xorm:"'name'"`
	Type            string          `json:"type" xorm:"'type'"`
	CategoryGroupId int64           `json:"categoryGroupId" xorm:"'category_group_id'"`
	IsActive        bool            `json:"isActive" xorm:"'is_active'"`
	IsDeleted       bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt       model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt       model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type CategoryView struct {
	Category           `xorm:"extends"`
	CategoryGroupCode  string `json:"categoryGroupCode" xorm:"'category_group_code'"`
	CategoryGroupName  string `json:"categoryGroupName" xorm:"'category_group_name'"`
	CategoryGroupColor string `json:"categoryGroupColor" xorm:"'category_group_color'"`
	CategoryGroupIcon  string `json:"categoryGroupIcon" xorm:"'category_group_icon'"`
}

type CategoryQuery struct {
	Type     string
	IsActive *bool
}

func (Category) TableName() string {
	return "category"
}

func InsertCategory(category *Category) error {
	_, err := infrastructure.Mysql.InsertOne(category)
	return err
}

func ListCategories(userId int64, query CategoryQuery) ([]CategoryView, error) {
	categories := make([]CategoryView, 0)
	session := infrastructure.Mysql.Table("category").Alias("c").
		Select("c.*, cg.group_code AS category_group_code, cg.name AS category_group_name, cg.color AS category_group_color, cg.icon AS category_group_icon").
		Join("LEFT", "category_group cg", "c.category_group_id = cg.id").
		Where("c.user_id = ? AND c.is_deleted = ?", userId, 0)

	if strings.TrimSpace(query.Type) != "" {
		session.And("c.type = ?", strings.TrimSpace(query.Type))
	}
	if query.IsActive != nil {
		session.And("c.is_active = ?", *query.IsActive)
	}

	err := session.Asc("c.type", "cg.sort", "c.id").Find(&categories)
	return categories, err
}

func GetCategoryByIdForUser(session *xorm.Session, id int64, userId int64) (*CategoryView, error) {
	categories := make([]CategoryView, 0, 1)
	err := session.Table("category").Alias("c").
		Select("c.*, cg.group_code AS category_group_code, cg.name AS category_group_name, cg.color AS category_group_color, cg.icon AS category_group_icon").
		Join("LEFT", "category_group cg", "c.category_group_id = cg.id").
		Where("c.id = ? AND c.user_id = ? AND c.is_deleted = ?", id, userId, 0).
		Limit(1).
		Find(&categories)
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, errors.New("category not found")
	}
	return &categories[0], nil
}

func GetCategory(id int64, userId int64) (*CategoryView, error) {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	return GetCategoryByIdForUser(session, id, userId)
}

func UpdateCategory(category *Category) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", category.Id, category.UserId, 0).
		Cols("name", "type", "category_group_id", "is_active").
		Update(category)
	return err
}

func SoftDeleteCategory(id int64, userId int64) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).
		Cols("is_deleted").
		Update(&Category{IsDeleted: true})
	return err
}
