package familycategorygroupdb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type FamilyCategoryGroup struct {
	Id              int64           `json:"id" xorm:"pk autoincr 'id'"`
	FamilyId        int64           `json:"familyId" xorm:"'family_id'"`
	Name            string          `json:"name" xorm:"'name'"`
	Type            string          `json:"type" xorm:"'type'"`
	Icon            string          `json:"icon" xorm:"'icon'"`
	Color           string          `json:"color" xorm:"'color'"`
	CreatedByUserId int64           `json:"createdByUserId" xorm:"'created_by_user_id'"`
	Sort            int             `json:"sort" xorm:"'sort'"`
	IsActive        bool            `json:"isActive" xorm:"'is_active'"`
	IsDeleted       bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt       model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt       model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

func (FamilyCategoryGroup) TableName() string {
	return "family_category_group"
}

type FamilyCategoryGroupQuery struct {
	FamilyId int64
	Type     string
	IsActive *bool
}

func ListFamilyCategoryGroups(query FamilyCategoryGroupQuery) ([]FamilyCategoryGroup, error) {
	groups := make([]FamilyCategoryGroup, 0)
	session := infrastructure.Mysql.NewSession()
	defer session.Close()

	session.And("family_id = ? AND is_deleted = 0", query.FamilyId)
	if strings.TrimSpace(query.Type) != "" {
		session.And("type = ?", strings.TrimSpace(query.Type))
	}
	if query.IsActive != nil {
		session.And("is_active = ?", *query.IsActive)
	}

	err := session.Asc("sort", "id").Find(&groups)
	return groups, err
}

func GetFamilyCategoryGroup(id int64, familyId int64) (*FamilyCategoryGroup, error) {
	group := &FamilyCategoryGroup{}
	has, err := infrastructure.Mysql.Where("id = ? AND family_id = ? AND is_deleted = 0", id, familyId).Get(group)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("family category group not found")
	}
	return group, nil
}

func GetFamilyCategoryGroupById(session *xorm.Session, id int64) (*FamilyCategoryGroup, error) {
	group := &FamilyCategoryGroup{}
	has, err := session.Where("id = ? AND is_deleted = 0", id).Get(group)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("family category group not found")
	}
	return group, nil
}

func InsertFamilyCategoryGroup(group *FamilyCategoryGroup) error {
	_, err := infrastructure.Mysql.InsertOne(group)
	return err
}

func InsertFamilyCategoryGroupTx(session *xorm.Session, group *FamilyCategoryGroup) error {
	_, err := session.InsertOne(group)
	return err
}

func UpdateFamilyCategoryGroup(group *FamilyCategoryGroup) error {
	_, err := infrastructure.Mysql.ID(group.Id).Omit("CreatedAt", "FamilyId", "CreatedByUserId").Update(group)
	return err
}

func SoftDeleteFamilyCategoryGroup(id int64, familyId int64) error {
	_, err := infrastructure.Mysql.Where("id = ? AND family_id = ?", id, familyId).Cols("is_deleted", "is_active").Update(&FamilyCategoryGroup{IsDeleted: true, IsActive: false})
	return err
}

// CheckNameExists 检查家庭内分组名称是否已存在
func CheckNameExists(familyId int64, name string, excludeId int64) (bool, error) {
	session := infrastructure.Mysql.Where("family_id = ? AND name = ? AND is_deleted = 0", familyId, name)
	if excludeId > 0 {
		session.And("id != ?", excludeId)
	}
	return session.Exist(&FamilyCategoryGroup{})
}
