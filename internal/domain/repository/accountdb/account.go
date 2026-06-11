package accountdb

import (
	"errors"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"

	"xorm.io/xorm"
)

type Account struct {
	Id                int64           `json:"id" xorm:"pk autoincr 'id'"`
	UserId            int64           `json:"userId" xorm:"'user_id'"`
	Name              string          `json:"name" xorm:"'name'"`
	Type              string          `json:"type" xorm:"'type'"`
	ProviderCode      string          `json:"providerCode" xorm:"'provider_code'"`
	AccountGroupKey   string          `json:"accountGroupKey" xorm:"'account_group_key'"`
	StandardAccountId *int64          `json:"standardAccountId" xorm:"'standard_account_id'"`
	BankCode          string          `json:"bankCode" xorm:"'bank_code'"`
	DisplayName       string          `json:"displayName" xorm:"'display_name'"`
	Icon              string          `json:"icon" xorm:"'icon'"`
	Color             string          `json:"color" xorm:"'color'"`
	IsActive          bool            `json:"isActive" xorm:"'is_active'"`
	IsDeleted         bool            `json:"isDeleted" xorm:"'is_deleted'"`
	CreatedAt         model.LocalTime `json:"createdAt" xorm:"created 'created_at'"`
	UpdatedAt         model.LocalTime `json:"updatedAt" xorm:"updated 'updated_at'"`
}

type AccountQuery struct {
	Type     string
	IsActive *bool
}

func (Account) TableName() string {
	return "account"
}

func InsertAccount(account *Account) error {
	_, err := infrastructure.Mysql.InsertOne(account)
	return err
}

func ListAccounts(userId int64, query AccountQuery) ([]Account, error) {
	accounts := make([]Account, 0)
	session := infrastructure.Mysql.Where("user_id = ? AND is_deleted = ?", userId, 0)

	if strings.TrimSpace(query.Type) != "" {
		session.And("type = ?", strings.TrimSpace(query.Type))
	}
	if query.IsActive != nil {
		session.And("is_active = ?", *query.IsActive)
	}

	err := session.Asc("id").Find(&accounts)
	return accounts, err
}

func GetAccountByIdForUser(session *xorm.Session, id int64, userId int64) (*Account, error) {
	account := &Account{}
	query := session.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0)
	has, err := query.Get(account)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("account not found")
	}
	return account, nil
}

func GetAccount(id int64, userId int64) (*Account, error) {
	session := infrastructure.Mysql.NewSession()
	defer session.Close()
	return GetAccountByIdForUser(session, id, userId)
}

func UpdateAccount(account *Account) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", account.Id, account.UserId, 0).
		Cols("name", "type", "provider_code", "account_group_key", "standard_account_id", "bank_code", "display_name", "icon", "color", "is_active").
		Update(account)
	return err
}

func CountActiveBuckets(accountId int64, userId int64) (int64, error) {
	return infrastructure.Mysql.
		Table("bucket").
		Where("account_id = ? AND user_id = ? AND is_deleted = ?", accountId, userId, 0).
		Count()
}

func SoftDeleteAccount(id int64, userId int64) error {
	_, err := infrastructure.Mysql.
		Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, 0).
		Cols("is_deleted").
		Update(&Account{IsDeleted: true})
	return err
}
