package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/category"
	"marmot-ledger/internal/domain/repository/categorydb"
	"marmot-ledger/internal/domain/repository/categorygroupdb"
	"marmot-ledger/internal/infrastructure"
	"strings"
)

func CreateCategory(userId int64, categoryInfo *category.Category) (*category.Category, error) {
	if err := validateCategory(userId, categoryInfo); err != nil {
		return nil, err
	}

	categoryDb := toCategoryDb(userId, categoryInfo)
	categoryDb.IsActive = true
	if err := categorydb.InsertCategory(categoryDb); err != nil {
		return nil, err
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
	return &entity, nil
}

func UpdateCategory(userId int64, id int64, categoryInfo *category.Category) (*category.Category, error) {
	categoryInfo.Id = id
	if err := validateCategory(userId, categoryInfo); err != nil {
		return nil, err
	}

	categoryDb := toCategoryDb(userId, categoryInfo)
	categoryDb.Id = id
	if err := categorydb.UpdateCategory(categoryDb); err != nil {
		return nil, err
	}
	return GetCategory(userId, id)
}

func DeleteCategory(userId int64, id int64) error {
	return categorydb.SoftDeleteCategory(id, userId)
}

func validateCategory(userId int64, categoryInfo *category.Category) error {
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
	if categoryInfo.CategoryGroupId == 0 {
		return errors.New("category group is required")
	}

	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	group, err := categorygroupdb.GetCategoryGroupById(session, categoryInfo.CategoryGroupId)
	if err != nil {
		return err
	}
	if !group.Enabled {
		return errors.New("category group is disabled")
	}
	if group.Type != categoryType {
		return errors.New("category group type does not match category type")
	}
	return nil
}

func toCategoryDb(userId int64, categoryInfo *category.Category) *categorydb.Category {
	return &categorydb.Category{
		Id:              categoryInfo.Id,
		UserId:          userId,
		Name:            strings.TrimSpace(categoryInfo.Name),
		Type:            strings.TrimSpace(categoryInfo.Type),
		CategoryGroupId: categoryInfo.CategoryGroupId,
		IsActive:        categoryInfo.IsActive,
	}
}

func toCategoryEntity(categoryDb *categorydb.CategoryView) category.Category {
	return category.Category{
		Id:                 categoryDb.Id,
		UserId:             categoryDb.UserId,
		Name:               categoryDb.Name,
		Type:               categoryDb.Type,
		CategoryGroupId:    categoryDb.CategoryGroupId,
		CategoryGroupCode:  categoryDb.CategoryGroupCode,
		CategoryGroupName:  categoryDb.CategoryGroupName,
		CategoryGroupColor: categoryDb.CategoryGroupColor,
		CategoryGroupIcon:  categoryDb.CategoryGroupIcon,
		IsActive:           categoryDb.IsActive,
	}
}
