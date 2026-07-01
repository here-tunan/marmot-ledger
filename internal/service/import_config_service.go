package service

import (
	"encoding/json"
	"errors"
	"marmot-ledger/internal/domain/entity/importconfig"
	"marmot-ledger/internal/domain/repository/importconfigdb"
	"strings"
)

var importConfigTargetFields = map[string]bool{
	"amount":      true,
	"description": true,
	"remark":      true,
	"eventTime":   true,
	"currency":    true,
	"scenario":    true,
	"category":    true,
	"channel":     true,
	"bucket":      true,
}

var importConfigOperators = map[string]bool{
	"contains":       true,
	"equals":         true,
	"notContains":    true,
	"notEquals":      true,
	"containsAny":    true,
	"notContainsAny": true,
	"equalsAny":      true,
	"notEqualsAny":   true,
}

var importFilterActions = map[string]bool{
	"drop": true,
	"keep": true,
}

func ListImportConfigs(userId int64, query importconfig.ImportConfigQuery) ([]importconfig.ImportConfig, error) {
	configs, err := importconfigdb.ListImportConfigs(userId, importconfigdb.ImportConfigQuery{IsActive: query.IsActive})
	if err != nil {
		return nil, err
	}
	result := make([]importconfig.ImportConfig, 0, len(configs))
	for _, item := range configs {
		entity, err := toImportConfigEntity(&item)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
	}
	return result, nil
}

func CreateImportConfig(userId int64, info *importconfig.ImportConfig) (*importconfig.ImportConfig, error) {
	if err := validateImportConfig(info); err != nil {
		return nil, err
	}
	configDb, err := toImportConfigDb(userId, info)
	if err != nil {
		return nil, err
	}
	configDb.IsActive = true
	if err := importconfigdb.InsertImportConfig(configDb); err != nil {
		return nil, err
	}
	return GetImportConfig(userId, configDb.Id)
}

func GetImportConfig(userId int64, id int64) (*importconfig.ImportConfig, error) {
	configDb, err := importconfigdb.GetImportConfig(id, userId)
	if err != nil {
		return nil, err
	}
	entity, err := toImportConfigEntity(configDb)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func UpdateImportConfig(userId int64, id int64, info *importconfig.ImportConfig) (*importconfig.ImportConfig, error) {
	if err := validateImportConfig(info); err != nil {
		return nil, err
	}
	if _, err := importconfigdb.GetImportConfig(id, userId); err != nil {
		return nil, err
	}
	configDb, err := toImportConfigDb(userId, info)
	if err != nil {
		return nil, err
	}
	configDb.Id = id
	if err := importconfigdb.UpdateImportConfig(configDb); err != nil {
		return nil, err
	}
	return GetImportConfig(userId, id)
}

func DeleteImportConfig(userId int64, id int64) error {
	if _, err := importconfigdb.GetImportConfig(id, userId); err != nil {
		return err
	}
	return importconfigdb.SoftDeleteImportConfig(id, userId)
}

func validateImportConfig(info *importconfig.ImportConfig) error {
	if info == nil {
		return errors.New("import config is required")
	}
	if strings.TrimSpace(info.Name) == "" {
		return errors.New("import config name is required")
	}
	fileType := strings.ToLower(strings.TrimSpace(info.FileType))
	if fileType != "xlsx" && fileType != "csv" {
		return errors.New("file type must be xlsx or csv")
	}
	for _, mapping := range info.Mappings {
		if !importConfigTargetFields[mapping.TargetField] {
			return errors.New("invalid target field: " + mapping.TargetField)
		}
		for _, rule := range mapping.Rules {
			if !importConfigOperators[rule.Operator] {
				return errors.New("invalid rule operator: " + rule.Operator)
			}
		}
	}
	for _, filter := range info.Filters {
		if !importConfigOperators[filter.Operator] {
			return errors.New("invalid filter operator: " + filter.Operator)
		}
		if !importFilterActions[filter.Action] {
			return errors.New("invalid filter action: " + filter.Action)
		}
	}
	return nil
}

func toImportConfigDb(userId int64, info *importconfig.ImportConfig) (*importconfigdb.ImportConfig, error) {
	mappings := info.Mappings
	if mappings == nil {
		mappings = []importconfig.FieldMapping{}
	}
	mappingsJson, err := json.Marshal(mappings)
	if err != nil {
		return nil, err
	}
	filters := info.Filters
	if filters == nil {
		filters = []importconfig.FilterRule{}
	}
	filtersJson, err := json.Marshal(filters)
	if err != nil {
		return nil, err
	}
	headerRow := info.HeaderRow
	if headerRow <= 0 {
		headerRow = 1
	}
	return &importconfigdb.ImportConfig{
		Id:        info.Id,
		UserId:    userId,
		Name:      strings.TrimSpace(info.Name),
		FileType:  strings.ToLower(strings.TrimSpace(info.FileType)),
		SheetName: strings.TrimSpace(info.SheetName),
		HeaderRow: headerRow,
		Mappings:  string(mappingsJson),
		Filters:   string(filtersJson),
		Icon:      strings.TrimSpace(info.Icon),
		Sort:      info.Sort,
		IsActive:  info.IsActive,
	}, nil
}

func toImportConfigEntity(configDb *importconfigdb.ImportConfig) (importconfig.ImportConfig, error) {
	mappings := make([]importconfig.FieldMapping, 0)
	if strings.TrimSpace(configDb.Mappings) != "" {
		if err := json.Unmarshal([]byte(configDb.Mappings), &mappings); err != nil {
			return importconfig.ImportConfig{}, err
		}
	}
	filters := make([]importconfig.FilterRule, 0)
	if strings.TrimSpace(configDb.Filters) != "" {
		if err := json.Unmarshal([]byte(configDb.Filters), &filters); err != nil {
			return importconfig.ImportConfig{}, err
		}
	}
	return importconfig.ImportConfig{
		Id:        configDb.Id,
		UserId:    configDb.UserId,
		Name:      configDb.Name,
		FileType:  configDb.FileType,
		SheetName: configDb.SheetName,
		HeaderRow: configDb.HeaderRow,
		Icon:      configDb.Icon,
		Sort:      configDb.Sort,
		IsActive:  configDb.IsActive,
		Mappings:  mappings,
		Filters:   filters,
	}, nil
}
