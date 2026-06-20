package familycategorygroupdb

import (
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"

	"xorm.io/xorm"
)

type FamilyCategoryGroupMember struct {
	Id            int64           `json:"id" xorm:"pk autoincr 'id'"`
	FamilyGroupId int64           `json:"familyGroupId" xorm:"'family_group_id'"`
	CategoryId    int64           `json:"categoryId" xorm:"'category_id'"`
	AddedByUserId int64           `json:"addedByUserId" xorm:"'added_by_user_id'"`
	AddedAt       model.LocalTime `json:"addedAt" xorm:"created 'added_at'"`
}

func (FamilyCategoryGroupMember) TableName() string {
	return "family_category_group_member"
}

// InsertGroupMember 添加成员到分组
func InsertGroupMember(member *FamilyCategoryGroupMember) error {
	_, err := infrastructure.Mysql.InsertOne(member)
	return err
}

func InsertGroupMemberTx(session *xorm.Session, member *FamilyCategoryGroupMember) error {
	_, err := session.InsertOne(member)
	return err
}

// DeleteGroupMember 从分组移除成员
func DeleteGroupMember(familyGroupId int64, categoryId int64) error {
	_, err := infrastructure.Mysql.Where("family_group_id = ? AND category_id = ?", familyGroupId, categoryId).Delete(&FamilyCategoryGroupMember{})
	return err
}

func DeleteGroupMemberTx(session *xorm.Session, familyGroupId int64, categoryId int64) error {
	_, err := session.Where("family_group_id = ? AND category_id = ?", familyGroupId, categoryId).Delete(&FamilyCategoryGroupMember{})
	return err
}

// DeleteAllMembersByGroupId 移除分组下所有成员
func DeleteAllMembersByGroupId(familyGroupId int64) error {
	_, err := infrastructure.Mysql.Where("family_group_id = ?", familyGroupId).Delete(&FamilyCategoryGroupMember{})
	return err
}

// ListGroupMemberIds 获取分组下所有分类ID
func ListGroupMemberIds(familyGroupId int64) ([]int64, error) {
	members := make([]FamilyCategoryGroupMember, 0)
	err := infrastructure.Mysql.Where("family_group_id = ?", familyGroupId).Find(&members)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.CategoryId)
	}
	return ids, nil
}

// ListGroupsForCategory 获取分类所属的所有家庭分组ID
func ListGroupsForCategory(categoryId int64) ([]int64, error) {
	members := make([]FamilyCategoryGroupMember, 0)
	err := infrastructure.Mysql.Where("category_id = ?", categoryId).Find(&members)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(members))
	for _, m := range members {
		ids = append(ids, m.FamilyGroupId)
	}
	return ids, nil
}

// CheckMemberExists 检查分类是否已在分组中
func CheckMemberExists(familyGroupId int64, categoryId int64) (bool, error) {
	return infrastructure.Mysql.Where("family_group_id = ? AND category_id = ?", familyGroupId, categoryId).Exist(&FamilyCategoryGroupMember{})
}

// GetAllGroupMembers 获取所有分组成员关系（用于统计聚合）
func GetAllGroupMembers() ([]FamilyCategoryGroupMember, error) {
	members := make([]FamilyCategoryGroupMember, 0)
	err := infrastructure.Mysql.Find(&members)
	return members, err
}
