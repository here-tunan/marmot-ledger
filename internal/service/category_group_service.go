package service

import (
	"marmot-ledger/internal/domain/entity/categorygroup"
	"marmot-ledger/internal/domain/repository/categorygroupdb"
)

func ListCategoryGroups(groupType string, enabled *bool) ([]categorygroup.CategoryGroup, error) {
	groups, err := categorygroupdb.ListCategoryGroups(categorygroupdb.CategoryGroupQuery{
		Type:    groupType,
		Enabled: enabled,
	})
	if err != nil {
		return nil, err
	}

	result := make([]categorygroup.CategoryGroup, 0, len(groups))
	for _, item := range groups {
		result = append(result, toCategoryGroupEntity(&item))
	}
	return result, nil
}

func GetCategoryGroup(id int64) (*categorygroup.CategoryGroup, error) {
	group, err := categorygroupdb.GetCategoryGroup(id)
	if err != nil {
		return nil, err
	}
	entity := toCategoryGroupEntity(group)
	return &entity, nil
}

func toCategoryGroupEntity(group *categorygroupdb.CategoryGroup) categorygroup.CategoryGroup {
	return categorygroup.CategoryGroup{
		Id:        group.Id,
		GroupCode: group.GroupCode,
		Name:      group.Name,
		Type:      group.Type,
		Icon:      group.Icon,
		Color:     group.Color,
		Sort:      group.Sort,
		Enabled:   group.Enabled,
	}
}
