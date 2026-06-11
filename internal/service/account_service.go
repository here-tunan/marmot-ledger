package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/account"
	"marmot-ledger/internal/domain/repository/accountdb"
	"strings"
)

func CreateAccount(userId int64, accountInfo *account.Account) (*account.Account, error) {
	if err := validateAccount(accountInfo); err != nil {
		return nil, err
	}

	accountDb := toAccountDb(userId, accountInfo)
	accountDb.IsActive = true
	if err := accountdb.InsertAccount(accountDb); err != nil {
		return nil, err
	}

	return toAccountEntity(accountDb), nil
}

func ListAccounts(userId int64, accountType string, isActive *bool) ([]account.Account, error) {
	accounts, err := accountdb.ListAccounts(userId, accountdb.AccountQuery{
		Type:     accountType,
		IsActive: isActive,
	})
	if err != nil {
		return nil, err
	}

	result := make([]account.Account, 0, len(accounts))
	for _, item := range accounts {
		result = append(result, *toAccountEntity(&item))
	}
	return result, nil
}

func GetAccount(userId int64, id int64) (*account.Account, error) {
	accountDb, err := accountdb.GetAccount(id, userId)
	if err != nil {
		return nil, err
	}
	return toAccountEntity(accountDb), nil
}

func UpdateAccount(userId int64, id int64, accountInfo *account.Account) (*account.Account, error) {
	if err := validateAccount(accountInfo); err != nil {
		return nil, err
	}

	accountDb := toAccountDb(userId, accountInfo)
	accountDb.Id = id
	if err := accountdb.UpdateAccount(accountDb); err != nil {
		return nil, err
	}

	return GetAccount(userId, id)
}

func DeleteAccount(userId int64, id int64) error {
	count, err := accountdb.CountActiveBuckets(id, userId)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("account has active buckets")
	}
	return accountdb.SoftDeleteAccount(id, userId)
}

func validateAccount(accountInfo *account.Account) error {
	if accountInfo == nil {
		return errors.New("account is required")
	}
	if strings.TrimSpace(accountInfo.Name) == "" {
		return errors.New("account name is required")
	}
	if strings.TrimSpace(accountInfo.Type) == "" {
		return errors.New("account type is required")
	}
	return nil
}

func toAccountDb(userId int64, accountInfo *account.Account) *accountdb.Account {
	var standardAccountId *int64
	if accountInfo.StandardAccountId != 0 {
		standardAccountId = &accountInfo.StandardAccountId
	}

	return &accountdb.Account{
		Id:                accountInfo.Id,
		UserId:            userId,
		Name:              strings.TrimSpace(accountInfo.Name),
		Type:              strings.TrimSpace(accountInfo.Type),
		ProviderCode:      strings.TrimSpace(accountInfo.ProviderCode),
		AccountGroupKey:   strings.TrimSpace(accountInfo.AccountGroupKey),
		StandardAccountId: standardAccountId,
		BankCode:          strings.TrimSpace(accountInfo.BankCode),
		DisplayName:       strings.TrimSpace(accountInfo.DisplayName),
		Icon:              strings.TrimSpace(accountInfo.Icon),
		Color:             strings.TrimSpace(accountInfo.Color),
		IsActive:          accountInfo.IsActive,
	}
}

func toAccountEntity(accountDb *accountdb.Account) *account.Account {
	var standardAccountId int64
	if accountDb.StandardAccountId != nil {
		standardAccountId = *accountDb.StandardAccountId
	}

	return &account.Account{
		Id:                accountDb.Id,
		UserId:            accountDb.UserId,
		Name:              accountDb.Name,
		Type:              accountDb.Type,
		ProviderCode:      accountDb.ProviderCode,
		AccountGroupKey:   accountDb.AccountGroupKey,
		StandardAccountId: standardAccountId,
		BankCode:          accountDb.BankCode,
		DisplayName:       accountDb.DisplayName,
		Icon:              accountDb.Icon,
		Color:             accountDb.Color,
		IsActive:          accountDb.IsActive,
	}
}
