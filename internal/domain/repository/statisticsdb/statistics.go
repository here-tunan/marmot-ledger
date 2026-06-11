package statisticsdb

import (
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/infrastructure"
	"strings"
)

type StatisticsQuery struct {
	StartTime string
	EndTime   string
	Currency  string
	Type      string
}

func GetSummary(userId int64, query StatisticsQuery) (*statistics.Summary, error) {
	return GetSummaryByUserIds([]int64{userId}, query)
}

func GetSummaryByUserIds(userIds []int64, query StatisticsQuery) (*statistics.Summary, error) {
	currency := normalizeCurrency(query.Currency)
	inClause, userParams := userIdsInClause(userIds)
	params := append(userParams, currency)
	timeWhere, timeParams := buildTimeWhere(query)
	params = append(params, timeParams...)

	sql := `
SELECT
  ? AS currency,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN base_amount ELSE 0 END), 0) AS income,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN base_amount ELSE 0 END), 0) AS gross_expense,
  COALESCE(SUM(CASE WHEN event_type = 'refund' THEN base_amount ELSE 0 END), 0) AS refund,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN base_amount WHEN event_type = 'refund' THEN -base_amount ELSE 0 END), 0) AS net_expense,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN base_amount WHEN event_type = 'refund' THEN -base_amount ELSE 0 END), 0) AS expense,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN base_amount WHEN event_type = 'expense' THEN -base_amount WHEN event_type = 'refund' THEN base_amount ELSE 0 END), 0) AS net,
  COUNT(*) AS event_count
FROM financial_event
WHERE user_id IN (` + inClause + `)
  AND is_deleted = 0
  AND status = 'active'
  AND include_in_statistics = 1
  AND event_type IN ('income', 'expense', 'refund')
  AND base_currency = ?` + timeWhere

	result := &statistics.Summary{}
	_, err := infrastructure.Mysql.SQL(sql, append([]any{currency}, params...)...).Get(result)
	return result, err
}

func GetCategoryGroupStats(userId int64, query StatisticsQuery) (*statistics.CategoryGroupStats, error) {
	return GetCategoryGroupStatsByUserIds([]int64{userId}, query)
}

func GetCategoryGroupStatsByUserIds(userIds []int64, query StatisticsQuery) (*statistics.CategoryGroupStats, error) {
	currency := normalizeCurrency(query.Currency)
	statsType := strings.TrimSpace(query.Type)
	if statsType == "" {
		statsType = "expense"
	}

	inClause, userParams := userIdsInClause(userIds)
	params := append([]any{statsType}, userParams...)
	params = append(params, currency)
	timeWhere, timeParams := buildTimeWhere(query)
	params = append(params, timeParams...)
	eventWhere := " AND fe.event_type = 'income'"
	amountExpr := "COALESCE(SUM(fe.base_amount), 0)"
	if statsType == "expense" {
		eventWhere = " AND fe.event_type IN ('expense', 'refund')"
		amountExpr = "COALESCE(SUM(CASE WHEN fe.event_type = 'expense' THEN fe.base_amount WHEN fe.event_type = 'refund' THEN -fe.base_amount ELSE 0 END), 0)"
	}

	sql := `
SELECT
  COALESCE(cg.id, 0) AS category_group_id,
  COALESCE(cg.group_code, 'UNCATEGORIZED') AS category_group_code,
  COALESCE(cg.name, '未分类') AS category_group_name,
  ? AS type,
  ` + amountExpr + ` AS amount,
  COUNT(*) AS count,
  COALESCE(cg.color, '') AS color,
  COALESCE(cg.icon, '') AS icon
FROM financial_event fe
LEFT JOIN category_group cg ON fe.category_group_id = cg.id
WHERE fe.user_id IN (` + inClause + `)
  AND fe.is_deleted = 0
  AND fe.status = 'active'
  AND fe.include_in_statistics = 1
  AND fe.base_currency = ?` + eventWhere + timeWhere + `
GROUP BY cg.id, cg.group_code, cg.name, cg.color, cg.icon
ORDER BY amount DESC, count DESC`

	items := make([]statistics.CategoryGroupItem, 0)
	err := infrastructure.Mysql.SQL(sql, params...).Find(&items)
	return &statistics.CategoryGroupStats{
		Currency: currency,
		Type:     statsType,
		Items:    items,
	}, err
}

func userIdsInClause(userIds []int64) (string, []any) {
	if len(userIds) == 0 {
		return "0", []any{}
	}
	placeholders := make([]string, 0, len(userIds))
	params := make([]any, 0, len(userIds))
	for _, id := range userIds {
		placeholders = append(placeholders, "?")
		params = append(params, id)
	}
	return strings.Join(placeholders, ","), params
}

func normalizeCurrency(currency string) string {
	if strings.TrimSpace(currency) == "" {
		return "CNY"
	}
	return strings.ToUpper(strings.TrimSpace(currency))
}

func buildTimeWhere(query StatisticsQuery) (string, []any) {
	where := ""
	params := make([]any, 0, 2)
	if strings.TrimSpace(query.StartTime) != "" {
		where += " AND event_time >= ?"
		params = append(params, strings.TrimSpace(query.StartTime))
	}
	if strings.TrimSpace(query.EndTime) != "" {
		where += " AND event_time <= ?"
		params = append(params, strings.TrimSpace(query.EndTime))
	}
	return where, params
}
