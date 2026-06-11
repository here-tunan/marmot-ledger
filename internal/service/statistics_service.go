package service

import (
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/domain/repository/statisticsdb"
)

func GetStatisticsSummary(userId int64, query statisticsdb.StatisticsQuery) (*statistics.Summary, error) {
	return statisticsdb.GetSummary(userId, query)
}

func GetStatisticsCategoryGroup(userId int64, query statisticsdb.StatisticsQuery) (*statistics.CategoryGroupStats, error) {
	return statisticsdb.GetCategoryGroupStats(userId, query)
}
