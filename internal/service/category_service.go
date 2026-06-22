package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/category"
	"marmot-ledger/internal/domain/repository/categorydb"
	"marmot-ledger/internal/domain/repository/familycategorygroupdb"
	"strings"
)

func CreateCategory(userId int64, categoryInfo *category.Category) (*category.Category, error) {
	if err := validateCategory(categoryInfo); err != nil {
		return nil, err
	}

	categoryDb := toCategoryDb(userId, categoryInfo)
	categoryDb.IsActive = true
	if err := categorydb.InsertCategory(categoryDb); err != nil {
		return nil, err
	}

	// 如果指定了家庭分组，添加成员关系
	if len(categoryInfo.GroupIds) > 0 {
		for _, groupId := range categoryInfo.GroupIds {
			member := &familycategorygroupdb.FamilyCategoryGroupMember{
				FamilyGroupId: groupId,
				CategoryId:    categoryDb.Id,
			}
			_ = familycategorygroupdb.InsertGroupMember(member)
		}
	}

	return GetCategory(userId, categoryDb.Id)
}

func ListCategories(userId int64, categoryType string, isActive *bool) ([]category.Category, error) {
	categories, err := categorydb.ListCategories(userId, categorydb.CategoryQuery{
		Type:     categoryType,
		IsActive: isActive,
	})
	if err != nil {
		return nil, err
	}

	result := make([]category.Category, 0, len(categories))
	for _, item := range categories {
		result = append(result, toCategoryEntity(&item))
	}
	return result, nil
}

func GetCategory(userId int64, id int64) (*category.Category, error) {
	categoryDb, err := categorydb.GetCategory(id, userId)
	if err != nil {
		return nil, err
	}
	entity := toCategoryEntity(categoryDb)

	// 填充所属家庭分组ID
	groupIds, err := familycategorygroupdb.ListGroupsForCategory(id)
	if err == nil {
		entity.GroupIds = groupIds
	}

	return &entity, nil
}

func UpdateCategory(userId int64, id int64, categoryInfo *category.Category) (*category.Category, error) {
	categoryInfo.Id = id
	if err := validateCategory(categoryInfo); err != nil {
		return nil, err
	}

	categoryDb := toCategoryDb(userId, categoryInfo)
	categoryDb.Id = id
	if err := categorydb.UpdateCategory(categoryDb); err != nil {
		return nil, err
	}

	// 更新家庭分组成员关系（先删后加）
	// TODO: 更高效的更新方式：对比差异只改变化的

	return GetCategory(userId, id)
}

// DeleteCategory 删除分类，返回被影响的账单数量
func DeleteCategory(userId int64, id int64) (int64, error) {
	// 先统计该分类下的账单数量
	count, err := categorydb.CountEventsByCategory(userId, id)
	if err != nil {
		return 0, err
	}

	// 执行软删除
	err = categorydb.SoftDeleteCategory(id, userId)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CheckCategoryUsage 检查分类的使用情况
func CheckCategoryUsage(userId int64, id int64) (int64, error) {
	return categorydb.CountEventsByCategory(userId, id)
}

func validateCategory(categoryInfo *category.Category) error {
	if categoryInfo == nil {
		return errors.New("category is required")
	}
	categoryType := strings.TrimSpace(categoryInfo.Type)
	if strings.TrimSpace(categoryInfo.Name) == "" {
		return errors.New("category name is required")
	}
	if categoryType != EventTypeIncome && categoryType != EventTypeExpense {
		return errors.New("category type is invalid")
	}
	return nil
}

func toCategoryDb(userId int64, categoryInfo *category.Category) *categorydb.Category {
	return &categorydb.Category{
		Id:       categoryInfo.Id,
		UserId:   userId,
		Name:     strings.TrimSpace(categoryInfo.Name),
		Type:     strings.TrimSpace(categoryInfo.Type),
		Icon:     categoryInfo.Icon,
		Color:    categoryInfo.Color,
		IsActive: categoryInfo.IsActive,
	}
}

func toCategoryEntity(categoryDb *categorydb.CategoryView) category.Category {
	return category.Category{
		Id:       categoryDb.Id,
		UserId:   categoryDb.UserId,
		Name:     categoryDb.Name,
		Type:     categoryDb.Type,
		Icon:     categoryDb.Icon,
		Color:    categoryDb.Color,
		IsActive: categoryDb.IsActive,
	}
}
