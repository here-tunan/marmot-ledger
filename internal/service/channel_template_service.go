package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/chantemplate"
	"marmot-ledger/internal/domain/repository/chantemplatedb"
	"strings"
)

// ListChannelTemplatesForUser 获取用户可见的渠道模板列表
func ListChannelTemplatesForUser() ([]chantemplate.ChannelTemplate, error) {
	enabled := true
	templates, err := chantemplatedb.ListChannelTemplates(chantemplatedb.ChannelTemplateQuery{
		Enabled: &enabled,
	})
	if err != nil {
		return nil, err
	}
	result := make([]chantemplate.ChannelTemplate, 0, len(templates))
	for _, t := range templates {
		result = append(result, toChannelTemplateEntity(&t))
	}
	return result, nil
}

// ListChannelTemplatesForAdmin 获取管理员可见的完整模板列表
func ListChannelTemplatesForAdmin() ([]chantemplate.ChannelTemplate, error) {
	templates, err := chantemplatedb.ListChannelTemplates(chantemplatedb.ChannelTemplateQuery{})
	if err != nil {
		return nil, err
	}
	result := make([]chantemplate.ChannelTemplate, 0, len(templates))
	for _, t := range templates {
		result = append(result, toChannelTemplateEntity(&t))
	}
	return result, nil
}

// CreateChannelTemplate 创建渠道模板（管理员）
func CreateChannelTemplate(req *chantemplate.CreateTemplateRequest) (*chantemplate.ChannelTemplate, error) {
	if strings.TrimSpace(req.ChannelCode) == "" {
		return nil, errors.New("渠道代码不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("渠道名称不能为空")
	}

	// 检查代码是否已存在
	exists, err := chantemplatedb.CheckCodeExists(strings.ToUpper(req.ChannelCode), 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("渠道代码已存在")
	}

	templateDb := &chantemplatedb.ChannelTemplate{
		ChannelCode:         strings.ToUpper(strings.TrimSpace(req.ChannelCode)),
		Name:                strings.TrimSpace(req.Name),
		ChannelType:         strings.ToLower(strings.TrimSpace(req.ChannelType)),
		ProviderCode:        strings.TrimSpace(req.ProviderCode),
		SupportedEventTypes: normalizeEventTypes(req.SupportedEventTypes),
		Icon:                strings.TrimSpace(req.Icon),
		Sort:                req.Sort,
		Enabled:             true,
		Remark:              strings.TrimSpace(req.Remark),
	}

	err = chantemplatedb.InsertChannelTemplate(templateDb)
	if err != nil {
		return nil, err
	}

	entity := toChannelTemplateEntity(templateDb)
	return &entity, nil
}

// UpdateChannelTemplate 更新渠道模板（管理员）
func UpdateChannelTemplate(id int64, req *chantemplate.UpdateTemplateRequest) (*chantemplate.ChannelTemplate, error) {
	templateDb, err := chantemplatedb.GetChannelTemplate(id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) != "" {
		templateDb.Name = strings.TrimSpace(req.Name)
	}
	if strings.TrimSpace(req.ChannelType) != "" {
		templateDb.ChannelType = strings.ToLower(strings.TrimSpace(req.ChannelType))
	}
	templateDb.ProviderCode = strings.TrimSpace(req.ProviderCode)
	templateDb.SupportedEventTypes = normalizeEventTypes(req.SupportedEventTypes)
	templateDb.Icon = strings.TrimSpace(req.Icon)
	templateDb.Sort = req.Sort
	templateDb.Remark = strings.TrimSpace(req.Remark)
	if req.Enabled != nil {
		templateDb.Enabled = *req.Enabled
	}

	err = chantemplatedb.UpdateChannelTemplate(templateDb)
	if err != nil {
		return nil, err
	}

	entity := toChannelTemplateEntity(templateDb)
	return &entity, nil
}

// GetChannelTemplate 获取单个模板
func GetChannelTemplate(id int64) (*chantemplate.ChannelTemplate, error) {
	templateDb, err := chantemplatedb.GetChannelTemplate(id)
	if err != nil {
		return nil, err
	}
	entity := toChannelTemplateEntity(templateDb)
	return &entity, nil
}

func normalizeEventTypes(value string) string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(part))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, ",")
}

func toChannelTemplateEntity(db *chantemplatedb.ChannelTemplate) chantemplate.ChannelTemplate {
	return chantemplate.ChannelTemplate{
		Id:                  db.Id,
		ChannelCode:         db.ChannelCode,
		Name:                db.Name,
		ChannelType:         db.ChannelType,
		ProviderCode:        db.ProviderCode,
		SupportedEventTypes: db.SupportedEventTypes,
		Icon:                db.Icon,
		Sort:                db.Sort,
		Enabled:             db.Enabled,
		Remark:              db.Remark,
	}
}
