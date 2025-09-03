package service

import (
	"go-my-life/internal/domain/repository/moneydb"
)

func PutTransactionAccount(account *moneydb.TransactionAccount) error {
	var err error
	if account.Id == 0 {
		err = account.Insert()
	} else {
		err = account.Update()
	}
	return err
}

func AllAccounts() ([]moneydb.TransactionAccount, error) {
	return moneydb.AllAccounts()
}

func AllAccountsByUser(userId int64) ([]moneydb.TransactionAccount, error) {
	return moneydb.AllAccountsByUser(userId)
}

func DeleteTransactionAccount(id int64) error {
	account := &moneydb.TransactionAccount{Id: id}
	return account.Delete()
}

func CheckAccountUsage(accountId int64) (int64, error) {
	return moneydb.CountTransactionsByAccount(accountId)
}

// FindOrCreateAccount 查找或创建用户账户
func FindOrCreateAccount(userId int64, accountName string) (int64, error) {
	// 先查找是否已存在
	accounts, err := moneydb.AllAccountsByUser(userId)
	if err != nil {
		return 0, err
	}

	for _, account := range accounts {
		if account.Name == accountName {
			return account.Id, nil
		}
	}

	// 不存在则创建
	newAccount := &moneydb.TransactionAccount{
		UserId: userId,
		Name:   accountName,
		Desc:   "导入时自动创建",
	}

	err = newAccount.Insert()
	if err != nil {
		return 0, err
	}

	return newAccount.Id, nil
}
