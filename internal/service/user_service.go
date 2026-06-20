package service

import (
	"errors"
	"marmot-ledger/internal/domain/entity/user"
	"marmot-ledger/internal/domain/repository/userdb"
)

func RegisterUser(account, password, name string) (*user.User, error) {
	// 检查是否已存在
	existsUser := &userdb.User{Account: account}
	err := existsUser.GetUser()
	if err == nil {
		return nil, errors.New("该账号已被注册")
	}

	displayName := name
	if displayName == "" {
		displayName = account
	}

	newUser := &user.User{
		Account:  account,
		Password: password,
		Name:     displayName,
	}

	err = PutUser(newUser)
	if err != nil {
		return nil, errors.New("注册失败: " + err.Error())
	}

	return newUser, nil
}

func GetUser(user *userdb.User) (*userdb.User, error) {
	err := user.GetUser()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func PutUser(user *user.User) error {
	userDb := &userdb.User{
		Id:       user.Id,
		Account:  user.Account,
		Password: user.Password,
		Name:     user.Name,
		Desc:     user.Desc,
		Avatar:   user.Avatar,
		Extra:    user.Extra,
	}

	err := userDb.PutUser()
	user.Id = userDb.Id
	if err != nil {
		return err
	}
	return nil
}
