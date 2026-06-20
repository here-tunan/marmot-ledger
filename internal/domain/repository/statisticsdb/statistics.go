package statisticsdb

import (
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/pkg/model"
	"strings"
	"time"

	"github.com/shopspring/decimal"
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
	timeWhere, timeParams := buildTimeWhere(query, "")
	params = append(params, timeParams...)

	sql := `
SELECT
  ? AS currency,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN amount ELSE 0 END), 0) AS income,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount ELSE 0 END), 0) AS gross_expense,
  COALESCE(SUM(CASE WHEN event_type = 'refund' THEN amount ELSE 0 END), 0) AS refund,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount WHEN event_type = 'refund' THEN -amount ELSE 0 END), 0) AS net_expense,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount WHEN event_type = 'refund' THEN -amount ELSE 0 END), 0) AS expense,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN amount WHEN event_type = 'expense' THEN -amount WHEN event_type = 'refund' THEN amount ELSE 0 END), 0) AS net,
  COUNT(*) AS event_count
FROM financial_event
WHERE user_id IN (` + inClause + `)
  AND is_deleted = 0
  AND status = 'active'
  AND include_in_statistics = 1
  AND event_type IN ('income', 'expense', 'refund')
  AND currency = ?` + timeWhere

	result := &statistics.Summary{}
	_, err := infrastructure.Mysql.SQL(sql, append([]any{currency}, params...)...).Get(result)
	return result, err
}

func GetSummariesByUserIds(userIds []int64, query StatisticsQuery) ([]statistics.Summary, error) {
	inClause, userParams := userIdsInClause(userIds)
	params := append([]any{}, userParams...)
	currencyWhere := ""
	if strings.TrimSpace(query.Currency) != "" {
		currencyWhere = " AND currency = ?"
		params = append(params, strings.ToUpper(strings.TrimSpace(query.Currency)))
	}
	timeWhere, timeParams := buildTimeWhere(query, "")
	params = append(params, timeParams...)

	sql := `
SELECT
  currency,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN amount ELSE 0 END), 0) AS income,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount ELSE 0 END), 0) AS gross_expense,
  COALESCE(SUM(CASE WHEN event_type = 'refund' THEN amount ELSE 0 END), 0) AS refund,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount WHEN event_type = 'refund' THEN -amount ELSE 0 END), 0) AS net_expense,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount WHEN event_type = 'refund' THEN -amount ELSE 0 END), 0) AS expense,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN amount WHEN event_type = 'expense' THEN -amount WHEN event_type = 'refund' THEN amount ELSE 0 END), 0) AS net,
  COUNT(*) AS event_count
FROM financial_event
WHERE user_id IN (` + inClause + `)
  AND is_deleted = 0
  AND status = 'active'
  AND include_in_statistics = 1
  AND event_type IN ('income', 'expense', 'refund')` + currencyWhere + timeWhere + `
GROUP BY currency
ORDER BY currency ASC`

	items := make([]statistics.Summary, 0)
	err := infrastructure.Mysql.SQL(sql, params...).Find(&items)
	return items, err
}

func GetCategoryGroupStats(userId int64, query StatisticsQuery) (*statistics.CategoryGroupStats, error) {
	return GetCategoryGroupStatsByUserIds([]int64{userId}, query)
}

func GetCategoryGroupStatsByUserIds(userIds []int64, query StatisticsQuery) (*statistics.CategoryGroupStats, error) {
	currency := normalizeCurrency(query.Currency)
	statsType := normalizeStatsType(query.Type)

	inClause, userParams := userIdsInClause(userIds)
	params := append([]any{statsType}, userParams...)
	params = append(params, currency)
	timeWhere, timeParams := buildTimeWhere(query, "fe")
	params = append(params, timeParams...)
	eventWhere, amountExpr := categoryEventWhereAndAmount(statsType)

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
  AND fe.currency = ?` + eventWhere + timeWhere + `
GROUP BY cg.id, cg.group_code, cg.name, cg.color, cg.icon
ORDER BY amount DESC, count DESC`

	items := make([]statistics.CategoryGroupItem, 0)
	err := infrastructure.Mysql.SQL(sql, params...).Find(&items)
	return &statistics.CategoryGroupStats{Currency: currency, Type: statsType, Items: items}, err
}

func GetCategoryGroupStatsListByUserIds(userIds []int64, query StatisticsQuery) ([]statistics.CategoryGroupStats, error) {
	statsType := normalizeStatsType(query.Type)
	inClause, userParams := userIdsInClause(userIds)
	params := append([]any{statsType}, userParams...)
	currencyWhere := ""
	if strings.TrimSpace(query.Currency) != "" {
		currencyWhere = " AND fe.currency = ?"
		params = append(params, strings.ToUpper(strings.TrimSpace(query.Currency)))
	}
	timeWhere, timeParams := buildTimeWhere(query, "fe")
	params = append(params, timeParams...)
	eventWhere, amountExpr := categoryEventWhereAndAmount(statsType)

	type row struct {
		Currency                     string `xorm:"'currency'"`
		statistics.CategoryGroupItem `xorm:"extends"`
	}
	rows := make([]row, 0)
	sql := `
SELECT
  fe.currency AS currency,
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
  AND fe.include_in_statistics = 1` + currencyWhere + eventWhere + timeWhere + `
GROUP BY fe.currency, cg.id, cg.group_code, cg.name, cg.color, cg.icon
ORDER BY fe.currency ASC, amount DESC, count DESC`
	if err := infrastructure.Mysql.SQL(sql, params...).Find(&rows); err != nil {
		return nil, err
	}

	result := make([]statistics.CategoryGroupStats, 0)
	index := make(map[string]int)
	for _, item := range rows {
		pos, ok := index[item.Currency]
		if !ok {
			index[item.Currency] = len(result)
			result = append(result, statistics.CategoryGroupStats{Currency: item.Currency, Type: statsType, Items: make([]statistics.CategoryGroupItem, 0)})
			pos = len(result) - 1
		}
		result[pos].Items = append(result[pos].Items, item.CategoryGroupItem)
	}
	return result, nil
}

type trendPeriod struct {
	Label string
	End   time.Time
}

func buildTrendPeriods(query StatisticsQuery, granularity string) []trendPeriod {
	location, err := time.LoadLocation(model.Loc)
	if err != nil {
		location = time.Local
	}
	now := time.Now().In(location)
	start := time.Date(now.Year()-1, now.Month(), now.Day(), 0, 0, 0, 0, location)
	end := now
	if strings.TrimSpace(query.StartTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(query.StartTime), location); err == nil {
			start = parsed
		}
	}
	if strings.TrimSpace(query.EndTime) != "" {
		if parsed, err := time.ParseInLocation(model.TimeFormat, strings.TrimSpace(query.EndTime), location); err == nil {
			end = parsed
		}
	}

	periods := make([]trendPeriod, 0)
	if granularity == "week" {
		cursor := start
		for !cursor.After(end) {
			periodEnd := cursor.AddDate(0, 0, 6)
			if periodEnd.After(end) {
				periodEnd = end
			}
			year, week := periodEnd.ISOWeek()
			periods = append(periods, trendPeriod{Label: decimal.NewFromInt(int64(year)).String() + "-W" + pad2(week), End: periodEnd})
			cursor = cursor.AddDate(0, 0, 7)
		}
		return periods
	}

	cursor := time.Date(start.Year(), start.Month(), 1, 23, 59, 59, 0, location)
	for !cursor.After(end) {
		periodEnd := time.Date(cursor.Year(), cursor.Month()+1, 0, 23, 59, 59, 0, location)
		if periodEnd.After(end) {
			periodEnd = end
		}
		periods = append(periods, trendPeriod{Label: cursor.Format("2006-01"), End: periodEnd})
		cursor = cursor.AddDate(0, 1, 0)
	}
	return periods
}

func pad2(v int) string {
	if v < 10 {
		return "0" + decimal.NewFromInt(int64(v)).String()
	}
	return decimal.NewFromInt(int64(v)).String()
}

func GetNetWorthTrendByUserIds(userIds []int64, query StatisticsQuery, granularity string) ([]statistics.NetWorthTrendPoint, error) {
	periods := buildTrendPeriods(query, granularity)
	if len(periods) == 0 {
		return []statistics.NetWorthTrendPoint{}, nil
	}
	inClause, userParams := userIdsInClause(userIds)
	currency := strings.ToUpper(strings.TrimSpace(query.Currency))
	result := make([]statistics.NetWorthTrendPoint, 0, len(periods))

	for _, p := range periods {
		params := append([]any{}, userParams...)
		params = append(params, p.End.Format(model.TimeFormat))
		params = append(params, userParams...)
		currencyWhere := ""
		if currency != "" {
			currencyWhere = " AND b.currency = ?"
			params = append(params, currency)
		}

		type row struct {
			Currency     string `xorm:"'currency'"`
			BucketNature string `xorm:"'bucket_nature'"`
			BalanceAfter string `xorm:"'balance_after'"`
		}
		rows := make([]row, 0)
		sql := `
SELECT
  b.currency AS currency,
  b.bucket_nature AS bucket_nature,
  le.balance_after AS balance_after
FROM bucket b
JOIN ledger_entry le ON le.bucket_id = b.id
JOIN financial_event fe ON fe.id = le.financial_event_id
JOIN (
  SELECT le2.bucket_id, MAX(le2.id) AS max_entry_id
  FROM ledger_entry le2
  JOIN financial_event fe2 ON fe2.id = le2.financial_event_id
  WHERE le2.user_id IN (` + inClause + `)
    AND fe2.event_time <= ?
    AND fe2.is_deleted = 0
    AND fe2.status = 'active'
  GROUP BY le2.bucket_id
) latest ON latest.max_entry_id = le.id
WHERE b.user_id IN (` + inClause + `)
  AND b.is_deleted = 0
  AND b.is_active = 1` + currencyWhere
		if err := infrastructure.Mysql.SQL(sql, params...).Find(&rows); err != nil {
			return nil, err
		}
		byCurrency := map[string]*statistics.NetWorthTrendPoint{}
		for _, r := range rows {
			balance, _ := decimal.NewFromString(r.BalanceAfter)
			point := byCurrency[r.Currency]
			if point == nil {
				point = &statistics.NetWorthTrendPoint{Period: p.Label, Currency: r.Currency}
				byCurrency[r.Currency] = point
			}
			if r.BucketNature == "liability" {
				point.Liability = point.Liability.Add(balance)
			} else {
				point.Asset = point.Asset.Add(balance)
			}
		}
		for _, point := range byCurrency {
			point.NetWorth = point.Asset.Sub(point.Liability)
			result = append(result, *point)
		}
	}
	return result, nil
}

func GetInvestmentSummary(userId int64, query StatisticsQuery) ([]statistics.InvestmentSummary, error) {
	type bucketRow struct {
		BucketId    int64  `xorm:"'bucket_id'"`
		BucketName  string `xorm:"'bucket_name'"`
		BucketType  string `xorm:"'bucket_type'"`
		Currency    string `xorm:"'currency'"`
		MarketValue string `xorm:"'market_value'"`
	}
	buckets := make([]bucketRow, 0)
	if err := infrastructure.Mysql.SQL(`
SELECT id AS bucket_id, name AS bucket_name, bucket_type, currency, balance AS market_value
FROM bucket
WHERE user_id = ?
  AND is_deleted = 0
  AND is_active = 1
  AND bucket_type IN ('investment_cash', 'investment_asset')
ORDER BY currency ASC, name ASC`, userId).Find(&buckets); err != nil {
		return nil, err
	}

	type flowRow struct {
		BucketId      int64  `xorm:"'bucket_id'"`
		Currency      string `xorm:"'currency'"`
		BuyAmount     string `xorm:"'buy_amount'"`
		RevalueAmount string `xorm:"'revalue_amount'"`
	}
	timeWhere, timeParams := buildTimeWhere(query, "fe")
	flowRows := make([]flowRow, 0)
	flowParams := append([]any{userId}, timeParams...)
	if err := infrastructure.Mysql.SQL(`
SELECT
  le.bucket_id AS bucket_id,
  le.currency AS currency,
  COALESCE(SUM(CASE WHEN le.entry_role = 'investment_buy' THEN le.amount ELSE 0 END), 0) AS buy_amount,
  COALESCE(SUM(CASE WHEN le.entry_role IN ('revaluation_gain', 'revaluation_loss') THEN le.amount ELSE 0 END), 0) AS revalue_amount
FROM ledger_entry le
JOIN financial_event fe ON fe.id = le.financial_event_id
JOIN bucket b ON b.id = le.bucket_id
WHERE le.user_id = ?
  AND fe.is_deleted = 0
  AND fe.status = 'active'
  AND b.bucket_type IN ('investment_cash', 'investment_asset')
  AND le.entry_role IN ('investment_buy', 'revaluation_gain', 'revaluation_loss')`+timeWhere+`
GROUP BY le.bucket_id, le.currency`, flowParams...).Find(&flowRows); err != nil {
		return nil, err
	}

	type incomeRow struct {
		Currency string `xorm:"'currency'"`
		Income   string `xorm:"'income_amount'"`
	}
	incomeRows := make([]incomeRow, 0)
	incomeParams := append([]any{userId}, timeParams...)
	if err := infrastructure.Mysql.SQL(`
SELECT currency, COALESCE(SUM(amount), 0) AS income_amount
FROM financial_event fe
WHERE user_id = ?
  AND is_deleted = 0
  AND status = 'active'
  AND event_type = 'investment_income'`+timeWhere+`
GROUP BY currency`, incomeParams...).Find(&incomeRows); err != nil {
		return nil, err
	}

	summaryMap := map[string]*statistics.InvestmentSummary{}
	bucketStats := map[int64]*statistics.InvestmentBucketStats{}

	ensureSummary := func(currency string) *statistics.InvestmentSummary {
		item, ok := summaryMap[currency]
		if !ok {
			item = &statistics.InvestmentSummary{Currency: currency, Buckets: make([]statistics.InvestmentBucketStats, 0)}
			summaryMap[currency] = item
		}
		return item
	}

	for _, row := range buckets {
		market, _ := decimal.NewFromString(row.MarketValue)
		bucket := &statistics.InvestmentBucketStats{
			BucketId:    row.BucketId,
			BucketName:  row.BucketName,
			BucketType:  row.BucketType,
			Currency:    row.Currency,
			MarketValue: market,
		}
		bucketStats[row.BucketId] = bucket
		summary := ensureSummary(row.Currency)
		summary.MarketValue = summary.MarketValue.Add(market)
	}

	for _, row := range flowRows {
		bucket := bucketStats[row.BucketId]
		if bucket == nil {
			continue
		}
		buy, _ := decimal.NewFromString(row.BuyAmount)
		revalue, _ := decimal.NewFromString(row.RevalueAmount)
		bucket.BuyAmount = buy
		bucket.RevalueAmount = revalue
		bucket.ProfitLoss = revalue
		if !buy.IsZero() {
			bucket.ReturnRate = bucket.ProfitLoss.Div(buy)
		}
		summary := ensureSummary(row.Currency)
		summary.BuyAmount = summary.BuyAmount.Add(buy)
		summary.RevalueAmount = summary.RevalueAmount.Add(revalue)
	}

	for _, row := range incomeRows {
		income, _ := decimal.NewFromString(row.Income)
		summary := ensureSummary(row.Currency)
		summary.IncomeAmount = summary.IncomeAmount.Add(income)
	}

	result := make([]statistics.InvestmentSummary, 0, len(summaryMap))
	for currency, summary := range summaryMap {
		for _, bucket := range bucketStats {
			if bucket.Currency == currency {
				summary.Buckets = append(summary.Buckets, *bucket)
			}
		}
		summary.ProfitLoss = summary.RevalueAmount.Add(summary.IncomeAmount)
		if !summary.BuyAmount.IsZero() {
			summary.ReturnRate = summary.ProfitLoss.Div(summary.BuyAmount)
		}
		result = append(result, *summary)
	}
	return result, nil
}

func GetTrendByUserIds(userIds []int64, query StatisticsQuery, granularity string) ([]statistics.TrendPoint, error) {
	currency := normalizeCurrency(query.Currency)
	periodExpr := "DATE_FORMAT(event_time, '%Y-%m')"
	if granularity == "week" {
		periodExpr = "DATE_FORMAT(event_time, '%x-W%v')"
	}
	inClause, userParams := userIdsInClause(userIds)
	params := append([]any{}, userParams...)
	params = append(params, currency)
	timeWhere, timeParams := buildTimeWhere(query, "")
	params = append(params, timeParams...)

	sql := `
SELECT
  ` + periodExpr + ` AS period,
  COALESCE(SUM(CASE WHEN event_type = 'income' THEN amount ELSE 0 END), 0) AS income,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount ELSE 0 END), 0) AS gross_expense,
  COALESCE(SUM(CASE WHEN event_type = 'refund' THEN amount ELSE 0 END), 0) AS refund,
  COALESCE(SUM(CASE WHEN event_type = 'expense' THEN amount WHEN event_type = 'refund' THEN -amount ELSE 0 END), 0) AS net_expense
FROM financial_event
WHERE user_id IN (` + inClause + `)
  AND is_deleted = 0
  AND status = 'active'
  AND include_in_statistics = 1
  AND event_type IN ('income', 'expense', 'refund')
  AND currency = ?` + timeWhere + `
GROUP BY ` + periodExpr + `
ORDER BY period ASC`

	type row struct {
		Period       string `xorm:"'period'"`
		Income       string `xorm:"'income'"`
		GrossExpense string `xorm:"'gross_expense'"`
		Refund       string `xorm:"'refund'"`
		NetExpense   string `xorm:"'net_expense'"`
	}
	rows := make([]row, 0)
	if err := infrastructure.Mysql.SQL(sql, params...).Find(&rows); err != nil {
		return nil, err
	}

	result := make([]statistics.TrendPoint, 0, len(rows))
	for _, r := range rows {
		income, _ := decimal.NewFromString(r.Income)
		grossExpense, _ := decimal.NewFromString(r.GrossExpense)
		refund, _ := decimal.NewFromString(r.Refund)
		netExpense, _ := decimal.NewFromString(r.NetExpense)
		result = append(result, statistics.TrendPoint{
			Period:       r.Period,
			Income:       income,
			GrossExpense: grossExpense,
			Refund:       refund,
			NetExpense:   netExpense,
		})
	}
	return result, nil
}

func categoryEventWhereAndAmount(statsType string) (string, string) {
	if statsType == "expense" {
		return " AND fe.event_type IN ('expense', 'refund')", "COALESCE(SUM(CASE WHEN fe.event_type = 'expense' THEN fe.amount WHEN fe.event_type = 'refund' THEN -fe.amount ELSE 0 END), 0)"
	}
	return " AND fe.event_type = 'income'", "COALESCE(SUM(fe.amount), 0)"
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

func normalizeStatsType(statsType string) string {
	if strings.TrimSpace(statsType) == "" {
		return "expense"
	}
	return strings.TrimSpace(statsType)
}

func buildTimeWhere(query StatisticsQuery, alias string) (string, []any) {
	prefix := ""
	if strings.TrimSpace(alias) != "" {
		prefix = strings.TrimSpace(alias) + "."
	}
	where := ""
	params := make([]any, 0, 2)
	if strings.TrimSpace(query.StartTime) != "" {
		where += " AND " + prefix + "event_time >= ?"
		params = append(params, strings.TrimSpace(query.StartTime))
	}
	if strings.TrimSpace(query.EndTime) != "" {
		where += " AND " + prefix + "event_time <= ?"
		params = append(params, strings.TrimSpace(query.EndTime))
	}
	return where, params
}
