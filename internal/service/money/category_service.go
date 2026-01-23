package service

import (
	"context"
	"encoding/json"
	"go-my-life/internal/domain/entity/family"
	"go-my-life/internal/domain/repository/moneydb"
	"go-my-life/internal/infrastructure"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func PutTransactionCategory(category *moneydb.TransactionCategory) error {
	var err error
	if category.Id == 0 {
		err = category.Insert()
	} else {
		err = category.Update()
	}
	return err
}

func AllCategory() ([]moneydb.TransactionCategory, error) {
	return moneydb.AllCategory()
}

func AllCategoriesByUser(userId int64) ([]moneydb.TransactionCategory, error) {
	return moneydb.AllCategoriesByUser(userId)
}

// AllCategoriesByFamily 获取家庭所有成员的分类
func AllCategoriesByFamily(familyId int64) ([]moneydb.TransactionCategory, error) {
	// 获取家庭信息
	familyEntity := &family.Family{Id: familyId}
	err := familyEntity.GetFamily()
	if err != nil {
		return nil, err
	}

	// 收集家庭成员的用户ID
	var userIds []int64
	for _, member := range familyEntity.Members {
		userIds = append(userIds, member.UserId)
	}

	// 查询所有家庭成员的分类
	return moneydb.AllCategoriesByUserIds(userIds)
}

func DeleteTransactionCategory(id int64) error {
	category := &moneydb.TransactionCategory{Id: id}
	return category.Delete()
}

func CheckCategoryUsage(categoryId int64) (int64, error) {
	return moneydb.CountTransactionsByCategory(categoryId)
}

// AnalysisCategory 根据描述和类型进行分析
func AnalysisCategory(desc string, transactionType int, userId int64) int {
	size := 1
	res, err := infrastructure.EsClient.Search().
		Index(moneydb.EsIndex).
		Request(&search.Request{
			Query: &types.Query{
				Bool: &types.BoolQuery{
					Must: []types.Query{
						{
							Match: map[string]types.MatchQuery{
								"description": {Query: desc},
							},
						},
						{
							Term: map[string]types.TermQuery{
								"type": {Value: transactionType},
							},
						},
						{
							Term: map[string]types.TermQuery{
								"userId": {Value: userId},
							},
						},
					},
				},
			},
			Size: &size,
		}).Do(context.Background())
	if err != nil {
		return 0
	}

	if res.Hits.Total.Value > 0 {
		// 取第一个匹配的结果
		jsonRaw := res.Hits.Hits[0].Source_
		doc := &moneydb.TransactionIndex{}
		err := json.Unmarshal(jsonRaw, doc)
		if err != nil {
			return 1
		}
		return doc.Category
	}
	return 1
}
