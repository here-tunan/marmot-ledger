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

func DeleteTransactionAccount(id int64) error {
	account := &moneydb.TransactionAccount{Id: id}
	return account.Delete()
}

func CheckAccountUsage(accountId int64) (int64, error) {
	return moneydb.CountTransactionsByAccount(accountId)
}
