package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/familycategorygroup"
	"marmot-ledger/internal/domain/repository/familycategorygroupdb"
	"marmot-ledger/internal/infrastructure"
	"strings"
)

// ListFamilyCategoryGroups 获取家庭分类组列表
func ListFamilyCategoryGroups(familyId int64, typeFilter string) ([]familycategorygroup.FamilyCategoryGroup, error) {
	groups, err := familycategorygroupdb.ListFamilyCategoryGroups(familycategorygroupdb.FamilyCategoryGroupQuery{
		FamilyId: familyId,
		Type:     typeFilter,
	})
	if err != nil {
		return nil, err
	}
	result := make([]familycategorygroup.FamilyCategoryGroup, 0, len(groups))
	for _, g := range groups {
		result = append(result, toFamilyCategoryGroupEntity(&g))
	}
	return result, nil
}

// GetFamilyCategoryGroup 获取单个家庭分类组
func GetFamilyCategoryGroup(familyId int64, groupId int64) (*familycategorygroup.FamilyCategoryGroup, error) {
	groupDb, err := familycategorygroupdb.GetFamilyCategoryGroup(groupId, familyId)
	if err != nil {
		return nil, err
	}
	entity := toFamilyCategoryGroupEntity(groupDb)
	return &entity, nil
}

// CreateFamilyCategoryGroup 创建家庭分类组
func CreateFamilyCategoryGroup(familyId int64, userId int64, req *familycategorygroup.CreateGroupRequest) (*familycategorygroup.FamilyCategoryGroup, error) {
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("分组名称不能为空")
	}
	if strings.TrimSpace(req.Type) == "" {
		return nil, errors.New("分组类型不能为空")
	}

	// 检查家庭内名称是否重复
	exists, err := familycategorygroupdb.CheckNameExists(familyId, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("分组名称已存在")
	}

	groupDb := &familycategorygroupdb.FamilyCategoryGroup{
		FamilyId:        familyId,
		Name:            req.Name,
		Type:            req.Type,
		Icon:            req.Icon,
		Color:           req.Color,
		CreatedByUserId: userId,
		Sort:            req.Sort,
		IsActive:        true,
	}

	err = familycategorygroupdb.InsertFamilyCategoryGroup(groupDb)
	if err != nil {
		return nil, err
	}

	entity := toFamilyCategoryGroupEntity(groupDb)
	return &entity, nil
}

// UpdateFamilyCategoryGroup 更新家庭分类组
func UpdateFamilyCategoryGroup(familyId int64, groupId int64, req *familycategorygroup.UpdateGroupRequest) (*familycategorygroup.FamilyCategoryGroup, error) {
	groupDb, err := familycategorygroupdb.GetFamilyCategoryGroup(groupId, familyId)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) != "" {
		// 检查新名称是否重复
		if req.Name != groupDb.Name {
			exists, err := familycategorygroupdb.CheckNameExists(familyId, req.Name, groupId)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("分组名称已存在")
			}
		}
		groupDb.Name = req.Name
	}
	if strings.TrimSpace(req.Type) != "" {
		groupDb.Type = req.Type
	}
	if req.Icon != "" {
		groupDb.Icon = req.Icon
	}
	if req.Color != "" {
		groupDb.Color = req.Color
	}
	if req.Sort != 0 {
		groupDb.Sort = req.Sort
	}

	err = familycategorygroupdb.UpdateFamilyCategoryGroup(groupDb)
	if err != nil {
		return nil, err
	}

	entity := toFamilyCategoryGroupEntity(groupDb)
	return &entity, nil
}

// DeleteFamilyCategoryGroup 删除家庭分类组
func DeleteFamilyCategoryGroup(familyId int64, groupId int64) error {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	// 先删除所有成员关系
	err := familycategorygroupdb.DeleteAllMembersByGroupId(groupId)
	if err != nil {
		session.Rollback()
		return err
	}

	// 软删除分组
	err = familycategorygroupdb.SoftDeleteFamilyCategoryGroup(groupId, familyId)
	if err != nil {
		session.Rollback()
		return err
	}

	return session.Commit()
}

// AddCategoriesToGroup 批量添加分类到分组
func AddCategoriesToGroup(familyId int64, groupId int64, categoryIds []int64, addedByUserId int64) error {
	// 验证分组属于该家庭
	_, err := familycategorygroupdb.GetFamilyCategoryGroup(groupId, familyId)
	if err != nil {
		return err
	}

	if len(categoryIds) == 0 {
		return errors.New("未选择任何分类")
	}

	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	for _, categoryId := range categoryIds {
		// 检查是否已添加
		exists, err := familycategorygroupdb.CheckMemberExists(groupId, categoryId)
		if err != nil {
			session.Rollback()
			return err
		}
		if exists {
			continue
		}

		member := &familycategorygroupdb.FamilyCategoryGroupMember{
			FamilyGroupId: groupId,
			CategoryId:    categoryId,
			AddedByUserId: addedByUserId,
		}
		err = familycategorygroupdb.InsertGroupMemberTx(session, member)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	return session.Commit()
}

// RemoveCategoryFromGroup 从分组移除分类
func RemoveCategoryFromGroup(familyId int64, groupId int64, categoryId int64) error {
	// 验证分组属于该家庭
	_, err := familycategorygroupdb.GetFamilyCategoryGroup(groupId, familyId)
	if err != nil {
		return err
	}

	return familycategorygroupdb.DeleteGroupMember(groupId, categoryId)
}

// GetGroupMemberIds 获取分组下所有分类ID
func GetGroupMemberIds(familyId int64, groupId int64) ([]int64, error) {
	// 验证分组属于该家庭
	_, err := familycategorygroupdb.GetFamilyCategoryGroup(groupId, familyId)
	if err != nil {
		return nil, err
	}

	return familycategorygroupdb.ListGroupMemberIds(groupId)
}

// GetGroupsForCategory 获取分类所属的所有家庭分组
func GetGroupsForCategory(categoryId int64) ([]int64, error) {
	return familycategorygroupdb.ListGroupsForCategory(categoryId)
}

func toFamilyCategoryGroupEntity(db *familycategorygroupdb.FamilyCategoryGroup) familycategorygroup.FamilyCategoryGroup {
	return familycategorygroup.FamilyCategoryGroup{
		Id:        db.Id,
		FamilyId:  db.FamilyId,
		Name:      db.Name,
		Type:      db.Type,
		Icon:      db.Icon,
		Color:     db.Color,
		CreatedBy: db.CreatedByUserId,
		Sort:      db.Sort,
		IsActive:  db.IsActive,
	}
}
