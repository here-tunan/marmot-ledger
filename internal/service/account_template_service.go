package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/account"
	"marmot-ledger/internal/domain/entity/accounttemplate"
	"marmot-ledger/internal/domain/repository/accounttemplatedb"
	"strings"
)

// ListAccountTemplatesForUser 获取用户可见的账户模板列表
func ListAccountTemplatesForUser() ([]accounttemplate.AccountTemplate, error) {
	enabled := true
	templates, err := accounttemplatedb.ListAccountTemplates(accounttemplatedb.AccountTemplateQuery{
		Enabled: &enabled,
	})
	if err != nil {
		return nil, err
	}
	result := make([]accounttemplate.AccountTemplate, 0, len(templates))
	for _, t := range templates {
		result = append(result, toAccountTemplateEntity(&t))
	}
	return result, nil
}

// ListAccountTemplatesForAdmin 获取管理员可见的完整模板列表
func ListAccountTemplatesForAdmin() ([]accounttemplate.AccountTemplate, error) {
	templates, err := accounttemplatedb.ListAccountTemplates(accounttemplatedb.AccountTemplateQuery{})
	if err != nil {
		return nil, err
	}
	result := make([]accounttemplate.AccountTemplate, 0, len(templates))
	for _, t := range templates {
		result = append(result, toAccountTemplateEntity(&t))
	}
	return result, nil
}

// CreateAccountTemplate 创建账户模板（管理员）
func CreateAccountTemplate(req *accounttemplate.CreateTemplateRequest) (*accounttemplate.AccountTemplate, error) {
	if strings.TrimSpace(req.ProviderCode) == "" {
		return nil, errors.New("模板代码不能为空")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("模板名称不能为空")
	}
	if strings.TrimSpace(req.Type) == "" {
		return nil, errors.New("账户类型不能为空")
	}

	code := strings.ToUpper(strings.TrimSpace(req.ProviderCode))
	exists, err := accounttemplatedb.CheckCodeExists(code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("模板代码已存在")
	}

	accountType := strings.ToLower(strings.TrimSpace(req.Type))
	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		icon = defaultAccountIcon(accountType)
	}
	color := strings.TrimSpace(req.Color)
	if color == "" {
		color = defaultAccountColor(accountType)
	}

	templateDb := &accounttemplatedb.AccountTemplate{
		ProviderCode: code,
		Name:         strings.TrimSpace(req.Name),
		Type:         accountType,
		Icon:         icon,
		Color:        color,
		Sort:         req.Sort,
		Enabled:      true,
	}

	err = accounttemplatedb.InsertAccountTemplate(templateDb)
	if err != nil {
		return nil, err
	}

	entity := toAccountTemplateEntity(templateDb)
	return &entity, nil
}

// UpdateAccountTemplate 更新账户模板（管理员）
func UpdateAccountTemplate(id int64, req *accounttemplate.UpdateTemplateRequest) (*accounttemplate.AccountTemplate, error) {
	templateDb, err := accounttemplatedb.GetAccountTemplate(id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(req.Name) != "" {
		templateDb.Name = strings.TrimSpace(req.Name)
	}
	if strings.TrimSpace(req.Type) != "" {
		templateDb.Type = strings.ToLower(strings.TrimSpace(req.Type))
	}
	if req.Icon != "" {
		templateDb.Icon = strings.TrimSpace(req.Icon)
	}
	if req.Color != "" {
		templateDb.Color = strings.TrimSpace(req.Color)
	}
	if req.Sort != 0 {
		templateDb.Sort = req.Sort
	}
	if req.Enabled != nil {
		templateDb.Enabled = *req.Enabled
	}
	if strings.TrimSpace(templateDb.Icon) == "" {
		templateDb.Icon = defaultAccountIcon(templateDb.Type)
	}
	if strings.TrimSpace(templateDb.Color) == "" {
		templateDb.Color = defaultAccountColor(templateDb.Type)
	}

	err = accounttemplatedb.UpdateAccountTemplate(templateDb)
	if err != nil {
		return nil, err
	}

	entity := toAccountTemplateEntity(templateDb)
	return &entity, nil
}

// GetAccountTemplate 获取单个模板
func GetAccountTemplate(id int64) (*accounttemplate.AccountTemplate, error) {
	templateDb, err := accounttemplatedb.GetAccountTemplate(id)
	if err != nil {
		return nil, err
	}
	entity := toAccountTemplateEntity(templateDb)
	return &entity, nil
}

// InstantiateAccountFromTemplate 用户基于模板创建账户
func InstantiateAccountFromTemplate(userId int64, templateId int64, customName string) (*account.Account, error) {
	templateDb, err := accounttemplatedb.GetAccountTemplate(templateId)
	if err != nil {
		return nil, err
	}
	if !templateDb.Enabled {
		return nil, errors.New("该模板已禁用")
	}

	accountName := templateDb.Name
	if strings.TrimSpace(customName) != "" {
		accountName = strings.TrimSpace(customName)
	}

	accountInfo := &account.Account{
		Name:     accountName,
		Type:     templateDb.Type,
		Icon:     templateDb.Icon,
		Color:    templateDb.Color,
		IsActive: true,
	}

	return CreateAccount(userId, accountInfo)
}

func toAccountTemplateEntity(db *accounttemplatedb.AccountTemplate) accounttemplate.AccountTemplate {
	return accounttemplate.AccountTemplate{
		Id:           db.Id,
		ProviderCode: db.ProviderCode,
		Name:         db.Name,
		Type:         db.Type,
		Icon:         db.Icon,
		Color:        db.Color,
		Sort:         db.Sort,
		Enabled:      db.Enabled,
	}
}
