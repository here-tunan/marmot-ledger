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

func GetStatisticsSummaries(userId int64, query statisticsdb.StatisticsQuery) ([]statistics.Summary, error) {
	return statisticsdb.GetSummariesByUserIds([]int64{userId}, query)
}

func GetStatisticsCategoryGroups(userId int64, query statisticsdb.StatisticsQuery) ([]statistics.CategoryGroupStats, error) {
	return statisticsdb.GetCategoryGroupStatsListByUserIds([]int64{userId}, query)
}

func GetStatisticsTrend(userId int64, query statisticsdb.StatisticsQuery, granularity string) ([]statistics.TrendPoint, error) {
	return statisticsdb.GetTrendByUserIds([]int64{userId}, query, granularity)
}

func GetInvestmentSummaries(userId int64, query statisticsdb.StatisticsQuery) ([]statistics.InvestmentSummary, error) {
	return statisticsdb.GetInvestmentSummary(userId, query)
}

func GetNetWorthTrend(userId int64, query statisticsdb.StatisticsQuery, granularity string) ([]statistics.NetWorthTrendPoint, error) {
	return statisticsdb.GetNetWorthTrendByUserIds([]int64{userId}, query, granularity)
}
