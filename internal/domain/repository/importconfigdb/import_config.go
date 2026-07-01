package importconfigdb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"

	"xorm.io/xorm"
)

type ImportConfig struct {
	Id        int64           `json:"id" xorm:"pk autoincr 'id'"`
	UserId    int64           `json:"userId" xorm:"'user_id'"`
	Name      string          `json:"name" xorm:"'name'"`
	FileType  string          `json:"fileType" xorm:"'file_type'"`
	SheetName string          `json:"sheetName" xorm:"'sheet_name'"`
	HeaderRow int             `json:"headerRow" xorm:"'header_row'"`
	Mappings  string          `json:"mappings" xorm:"'mappings'"` // JSON 文本
	Filters   string          `json:"filters" xorm:"'filters'"`   // JSON 文本
	Icon      string          `json:"icon" xorm:"'icon'"`
	Sort      int             `json:"sort" xorm:"'sort'"`
	IsActive  bool            `json:"isActive" xorm:"'is_active'"`
	IsDeleted bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type ImportConfigQuery struct {
	IsActive *bool
}

func (ImportConfig) TableName() string {
	return "import_config"
}

func InsertImportConfig(config *ImportConfig) error {
	_, err := infrastructure.Mysql.InsertOne(config)
	return err
}

func ListImportConfigs(userId int64, query ImportConfigQuery) ([]ImportConfig, error) {
	configs := make([]ImportConfig, 0)
	session := infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0)
	applyQuery(session, query)
	err := session.Asc("sort", "id").Find(&configs)
	return configs, err
}

func GetImportConfig(id int64, userId int64) (*ImportConfig, error) {
	config := &ImportConfig{}
	has, err := infrastructure.Mysql.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).Get(config)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("import config not found")
	}
	return config, nil
}

func UpdateImportConfig(config *ImportConfig) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", config.Id, config.UserId, 0).
		Cols("name", "file_type", "sheet_name", "header_row", "mappings", "filters", "icon", "sort", "is_active").
		Update(config)
	return err
}

func SoftDeleteImportConfig(id int64, userId int64) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).
		Cols("is_deleted").
		Update(&ImportConfig{IsDeleted: true})
	return err
}

func applyQuery(session *xorm.Session, query ImportConfigQuery) {
	if query.IsActive != nil {
		session.And("is_active = ?", *query.IsActive)
	}
}
