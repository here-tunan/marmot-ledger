package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/category"
	"marmot-ledger/internal/domain/entity/categorytemplate"
	"marmot-ledger/internal/domain/repository/categorytemplatedb"
	"strings"
)

// ListCategoryTemplatesForUser 获取用户可见的分类模板列表
func ListCategoryTemplatesForUser(typeFilter string) ([]categorytemplate.CategoryTemplate, error) {
	enabled := true
	templates, err := categorytemplatedb.ListCategoryTemplates(categorytemplatedb.CategoryTemplateQuery{
		Enabled: &enabled,
		Type:    typeFilter,
	})
	if err != nil {
		return nil, err
	}
	result := make([]categorytemplate.CategoryTemplate, 0, len(templates))
	for _, t := range templates {
		result = append(result, toCategoryTemplateEntity(&t))
	}
	return result, nil
}

// ListCategoryTemplatesForAdmin 获取管理员可见的完整模板列表
func ListCategoryTemplatesForAdmin(typeFilter string) ([]categorytemplate.CategoryTemplate, error) {
	templates, err := categorytemplatedb.ListCategoryTemplates(categorytemplatedb.CategoryTemplateQuery{
		Type: typeFilter,
	})
	if err != nil {
		return nil, err
	}
	result := make([]categorytemplate.CategoryTemplate, 0, len(templates))
	for _, t := range templates {
		result = append(result, toCategoryTemplateEntity(&t))
	}
	return result, nil
}

// CreateCategoryTemplate 创建分类模板（管理员）
func CreateCategoryTemplate(req *categorytemplate.CreateTemplateRequest) (*categorytemplate.CategoryTemplate, error) {
	if strings.TrimSpace(req.TemplateCode) == "" {
		return nil, errors.New("模板代码不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("模板名称不能为空")
	}
	if strings.TrimSpace(req.Type) == "" {
		return nil, errors.New("分类类型不能为空")
	}

	// 检查代码是否已存在
	exists, err := categorytemplatedb.CheckCodeExists(strings.ToUpper(req.TemplateCode), 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("模板代码已存在")
	}

	templateDb := &categorytemplatedb.CategoryTemplate{
		TemplateCode: strings.ToUpper(req.TemplateCode),
		Name:         req.Name,
		Type:         req.Type,
		Icon:         req.Icon,
		Color:        req.Color,
		Sort:         req.Sort,
		Enabled:      true,
	}

	err = categorytemplatedb.InsertCategoryTemplate(templateDb)
	if err != nil {
		return nil, err
	}

	entity := toCategoryTemplateEntity(templateDb)
	return &entity, nil
}

// UpdateCategoryTemplate 更新分类模板（管理员）
func UpdateCategoryTemplate(id int64, req *categorytemplate.UpdateTemplateRequest) (*categorytemplate.CategoryTemplate, error) {
	templateDb, err := categorytemplatedb.GetCategoryTemplate(id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) != "" {
		templateDb.Name = req.Name
	}
	if strings.TrimSpace(req.Type) != "" {
		templateDb.Type = req.Type
	}
	if req.Icon != "" {
		templateDb.Icon = req.Icon
	}
	if req.Color != "" {
		templateDb.Color = req.Color
	}
	if req.Sort != 0 {
		templateDb.Sort = req.Sort
	}
	if req.Enabled != nil {
		templateDb.Enabled = *req.Enabled
	}

	err = categorytemplatedb.UpdateCategoryTemplate(templateDb)
	if err != nil {
		return nil, err
	}

	entity := toCategoryTemplateEntity(templateDb)
	return &entity, nil
}

// GetCategoryTemplate 获取单个模板
func GetCategoryTemplate(id int64) (*categorytemplate.CategoryTemplate, error) {
	templateDb, err := categorytemplatedb.GetCategoryTemplate(id)
	if err != nil {
		return nil, err
	}
	entity := toCategoryTemplateEntity(templateDb)
	return &entity, nil
}

// ImportTemplatesForUser 用户批量导入模板分类
func ImportTemplatesForUser(userId int64, templateIds []int64) ([]category.Category, error) {
	if len(templateIds) == 0 {
		return nil, errors.New("未选择任何模板")
	}

	templates, err := categorytemplatedb.BatchGetByIds(templateIds)
	if err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return nil, errors.New("未找到有效的模板")
	}

	result := make([]category.Category, 0, len(templates))
	for _, t := range templates {
		if !t.Enabled {
			continue
		}

		categoryInfo := &category.Category{
			Name:       t.Name,
			Type:       t.Type,
			TemplateId: t.Id,
			Icon:       t.Icon,
			Color:      t.Color,
			IsActive:   true,
		}

		created, err := CreateCategory(userId, categoryInfo)
		if err != nil {
			// 忽略重复名称错误，继续导入其他
			if strings.Contains(err.Error(), "already exists") {
				continue
			}
			return nil, err
		}
		result = append(result, *created)
	}

	return result, nil
}

func toCategoryTemplateEntity(db *categorytemplatedb.CategoryTemplate) categorytemplate.CategoryTemplate {
	return categorytemplate.CategoryTemplate{
		Id:           db.Id,
		TemplateCode: db.TemplateCode,
		Name:         db.Name,
		Type:         db.Type,
		Icon:         db.Icon,
		Color:        db.Color,
		Sort:         db.Sort,
		Enabled:      db.Enabled,
	}
}
