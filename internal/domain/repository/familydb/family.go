package familydb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"

	"xorm.io/xorm"
)

type Family struct {
	Id          int64           `json:"id" xorm:"pk autoincr 'id'"`
	Name        string          `json:"name" xorm:"'name'"`
	OwnerUserId int64           `json:"ownerUserId" xorm:"'owner_user_id'"`
	IsDeleted   bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt   model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt   model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type FamilyView struct {
	Family `xorm:"extends"`
	Role   string `json:"role" xorm:"'role'"`
}

func (Family) TableName() string { return "family" }

func InsertFamily(session *xorm.Session, family *Family) error {
	_, err := session.InsertOne(family)
	return err
}

func ListFamilies(userId int64) ([]FamilyView, error) {
	items := make([]FamilyView, 0)
	err := infrastructure.Mysql.Table("family").Alias("f").
		Select("f.*, fm.role AS role").
		Join("INNER", "family_member fm", "f.id = fm.family_id").
		Where("fm.user_id = ? AND fm.status = 'active' AND f.is_deleted = 0", userId).
		Asc("f.id").Find(&items)
	return items, err
}

func GetFamily(id int64) (*Family, error) {
	family := &Family{}
	has, err := infrastructure.Mysql.Where("id = ? AND is_deleted = 0", id).Get(family)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("family not found")
	}
	return family, nil
}

func IsActiveMember(familyId int64, userId int64) (bool, error) {
	return infrastructure.Mysql.Table("family_member").Where("family_id = ? AND user_id = ? AND status = 'active'", familyId, userId).Exist()
}
