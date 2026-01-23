package moneydb

import (
	"fmt"
	"go-my-life/internal/infrastructure"
	"time"
)

type TransactionCategory struct {
	// id
	Id int64 `json:"id"`
	// 用户ID
	UserId int64 `json:"userId"`
	// 名称
	Name string `json:"name"`
	// 收入/支出
	Type int `json:"type"`
	// 描述
	Desc string `json:"desc"`
	// 是否删除
	IsDeleted bool `json:"isDeleted"`
	// 系统创建时间
	GmtCreate time.Time `json:"gmtCreate" xorm:"updated"`
	// 系统更新时间
	GmtModified time.Time `json:"gmtModified" xorm:"updated"`
}

func AllCategory() ([]TransactionCategory, error) {
	session := infrastructure.Mysql.Where("is_deleted = 0")
	var categories []TransactionCategory
	err := session.Find(&categories)
	if err != nil {
		return nil, err
	}
	return categories, err
}

func AllCategoriesByUser(userId int64) ([]TransactionCategory, error) {
	session := infrastructure.Mysql.Where("is_deleted = 0 AND user_id = ?", userId)
	var categories []TransactionCategory
	err := session.Find(&categories)
	if err != nil {
		return nil, err
	}
	return categories, err
}

func AllCategoriesByUserIds(userIds []int64) ([]TransactionCategory, error) {
	if len(userIds) == 0 {
		return []TransactionCategory{}, nil
	}

	session := infrastructure.Mysql.In("user_id", userIds).And("is_deleted = 0")
	var categories []TransactionCategory
	err := session.Find(&categories)
	if err != nil {
		return nil, err
	}
	return categories, err
}

func (category *TransactionCategory) Insert() error {
	_, err := infrastructure.Mysql.InsertOne(category)
	return err
}

func (category *TransactionCategory) Update() error {
	affected, err := infrastructure.Mysql.ID(category.Id).Omit("GmtCreate").Update(category)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("更新失败，未找到ID为%d的记录", category.Id)
	}
	return nil
}

func (category *TransactionCategory) Delete() error {
	// 软删除：设置 is_deleted = true
	category.IsDeleted = true
	_, err := infrastructure.Mysql.ID(category.Id).Cols("is_deleted").Update(category)
	return err
}
