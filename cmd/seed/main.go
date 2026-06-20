package main

import (
	"fmt"
	"log"
	_ "marmot-ledger/env"
	"marmot-ledger/internal/domain/entity/account"
	"marmot-ledger/internal/domain/entity/bucket"
	"marmot-ledger/internal/domain/entity/category"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/internal/service"
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

const (
	userID  = 1
	rngSeed = 42
)

func main() {
	defer infrastructure.Mysql.Close()
	rand.Seed(rngSeed)

	log.Println("Seeding accounts...")
	accounts := seedAccounts()
	for k, v := range accounts {
		log.Printf("  %s -> id=%d", k, v.Id)
	}

	log.Println("Seeding buckets...")
	buckets := seedBuckets(accounts)
	for k, v := range buckets {
		log.Printf("  %s -> id=%d (initial=%s)", k, v.Id, v.InitialBalance)
	}

	log.Println("Seeding categories...")
	categories := seedCategories()
	for k, v := range categories {
		log.Printf("  %s -> id=%d", k, v.Id)
	}

	log.Println("Seeding records (24 months)...")
	count := seedRecords(buckets, categories)
	log.Printf("Created %d records", count)
	log.Println("Done!")
}

func seedAccounts() map[string]*account.Account {
	type accSpec struct {
		key, name, accType, icon, color string
	}
	specs := []accSpec{
		{"cash", "现金钱包", "cash", "Money", "#f59e0b"},
		{"wechat", "微信钱包", "wallet", "Wallet", "#22c55e"},
		{"alipay", "支付宝", "wallet", "Wallet", "#1677ff"},
		{"cmb", "招商银行储蓄卡", "bank", "CreditCard", "#ef4444"},
		{"icbc", "工行储蓄卡", "bank", "CreditCard", "#f97316"},
		{"cmb_credit", "招行信用卡", "credit", "CreditCard", "#ec4899"},
		{"ibkr", "盈透证券", "investment", "TrendCharts", "#1f2933"},
		{"futu", "富途证券", "investment", "TrendCharts", "#06b6d4"},
	}
	out := make(map[string]*account.Account, len(specs))
	for _, s := range specs {
		acc := &account.Account{
			Name:     s.name,
			Type:     s.accType,
			Icon:     s.icon,
			Color:    s.color,
			IsActive: true,
		}
		created, err := service.CreateAccount(userID, acc)
		must(err, "account "+s.name)
		out[s.key] = created
	}
	return out
}

func seedBuckets(accounts map[string]*account.Account) map[string]*bucket.Bucket {
	type bSpec struct {
		key, name, currency, bucketType, nature string
		accountKey                              string
		initial                                 string
	}
	specs := []bSpec{
		{"cash_cny", "人民币现金", "CNY", "cash", "asset", "cash", "3000"},
		{"wechat_cny", "微信零钱", "CNY", "wallet", "asset", "wechat", "5000"},
		{"alipay_cny", "支付宝余额", "CNY", "wallet", "asset", "alipay", "8000"},
		{"cmb_cny", "招行储蓄卡", "CNY", "bank", "asset", "cmb", "80000"},
		{"cmb_usd", "招行美元账户", "USD", "bank", "asset", "cmb", "12000"},
		{"icbc_cny", "工行储蓄卡", "CNY", "bank", "asset", "icbc", "45000"},
		{"cmb_credit", "招行信用卡", "CNY", "credit", "liability", "cmb_credit", "0"},
		{"ibkr_cash", "盈透现金账户", "USD", "investment_cash", "asset", "ibkr", "8000"},
		{"ibkr_qqq", "盈透 QQQ 持仓", "USD", "investment_asset", "asset", "ibkr", "0"},
		{"ibkr_voo", "盈透 VOO 持仓", "USD", "investment_asset", "asset", "ibkr", "0"},
		{"futu_hkd", "富途现金账户", "HKD", "investment_cash", "asset", "futu", "60000"},
		{"futu_hk_etf", "富途盈富基金 2800", "HKD", "investment_asset", "asset", "futu", "0"},
		{"recv", "应收款", "CNY", "receivable", "asset", "alipay", "0"},
		{"deposit", "押金", "CNY", "deposit", "asset", "cmb", "0"},
	}
	out := make(map[string]*bucket.Bucket, len(specs))
	for _, s := range specs {
		acc := accounts[s.accountKey]
		if acc == nil {
			log.Fatalf("account key not found: %s", s.accountKey)
		}
		initial, _ := decimal.NewFromString(s.initial)
		bk := &bucket.Bucket{
			AccountId:      acc.Id,
			Name:           s.name,
			Currency:       s.currency,
			InitialBalance: initial,
			BucketType:     s.bucketType,
			BucketNature:   s.nature,
		}
		created, err := service.CreateBucket(userID, bk)
		must(err, "bucket "+s.name)
		out[s.key] = created
	}
	return out
}

func seedCategories() map[string]*category.Category {
	type catSpec struct {
		key, name, catType, groupCode string
	}
	specs := []catSpec{
		// expense
		{"food_meal", "外卖", "expense", "FOOD"},
		{"food_dining", "餐厅吃饭", "expense", "FOOD"},
		{"food_coffee", "咖啡奶茶", "expense", "FOOD"},
		{"grocery", "买菜日用", "expense", "GROCERY"},
		{"transport_subway", "地铁公交", "expense", "TRANSPORT"},
		{"transport_taxi", "打车", "expense", "TRANSPORT"},
		{"housing_rent", "房租", "expense", "HOUSING"},
		{"utilities", "水电燃气", "expense", "UTILITIES"},
		{"shopping_clothes", "服装鞋包", "expense", "SHOPPING"},
		{"shopping_electronics", "数码电子", "expense", "SHOPPING"},
		{"entertainment_movie", "电影演出", "expense", "ENTERTAINMENT"},
		{"entertainment_subscription", "订阅服务", "expense", "ENTERTAINMENT"},
		{"healthcare", "医疗健康", "expense", "HEALTHCARE"},
		{"education", "教育学习", "expense", "EDUCATION"},
		{"travel", "旅行", "expense", "TRAVEL"},
		// income
		{"salary", "工资", "income", "SALARY"},
		{"bonus", "奖金", "income", "BONUS"},
		{"reimbursement", "公司报销", "income", "REIMBURSEMENT"},
	}
	groups := loadCategoryGroups()
	out := make(map[string]*category.Category, len(specs))
	for _, s := range specs {
		groupID, ok := groups[s.groupCode]
		if !ok {
			log.Fatalf("category group not found: %s", s.groupCode)
		}
		_ = groupID
		c := &category.Category{
			Name:     s.name,
			Type:     s.catType,
			IsActive: true,
		}
		created, err := service.CreateCategory(userID, c)
		must(err, "category "+s.name)
		out[s.key] = created
	}
	return out
}

func loadCategoryGroups() map[string]int64 {
	groups, err := infrastructure.Mysql.QueryString("SELECT id, group_code FROM category_group")
	must(err, "load category groups")
	out := make(map[string]int64, len(groups))
	for _, row := range groups {
		var id int64
		fmt.Sscan(row["id"], &id)
		out[row["group_code"]] = id
	}
	return out
}

func seedRecords(buckets map[string]*bucket.Bucket, cats map[string]*category.Category) int {
	now := time.Now()
	startDate := now.AddDate(-2, 0, 0)
	count := 0

	cnyCash := []string{"alipay_cny", "wechat_cny", "cash_cny", "cmb_cny", "cmb_cny", "cmb_cny"}
	_ = cnyCash

	// helpers
	rec := func(scenario string, eventTime time.Time, amount string, currency string, fromBucket, toBucket, plainBucket, catKey, desc string) {
		req := &record.RecordRequest{
			Scenario:    scenario,
			Currency:    currency,
			Description: desc,
			EventTime:   eventTime.Format("2006-01-02 15:04:05"),
		}
		amt, _ := decimal.NewFromString(amount)
		req.Amount = amt
		if plainBucket != "" {
			req.BucketId = buckets[plainBucket].Id
		}
		if fromBucket != "" {
			req.FromBucketId = buckets[fromBucket].Id
		}
		if toBucket != "" {
			req.ToBucketId = buckets[toBucket].Id
		}
		if catKey != "" {
			req.CategoryId = cats[catKey].Id
		}
		_, err := service.CreateRecord(userID, req)
		if err != nil {
			log.Printf("WARN record skipped (%s %s %s): %v", scenario, amount, desc, err)
			return
		}
		count++
	}

	recExchange := func(eventTime time.Time, fromBucket, toBucket, fromAmount, toAmount, fromCurr, toCurr, desc string) {
		req := &record.RecordRequest{
			Scenario:    "exchange",
			Currency:    fromCurr,
			ToCurrency:  toCurr,
			Description: desc,
			EventTime:   eventTime.Format("2006-01-02 15:04:05"),
		}
		fa, _ := decimal.NewFromString(fromAmount)
		ta, _ := decimal.NewFromString(toAmount)
		req.Amount = fa
		req.ToAmount = ta
		req.FromBucketId = buckets[fromBucket].Id
		req.ToBucketId = buckets[toBucket].Id
		_, err := service.CreateRecord(userID, req)
		if err != nil {
			log.Printf("WARN exchange skipped: %v", err)
			return
		}
		count++
	}

	recRevalue := func(eventTime time.Time, bucketKey, delta, currency, desc string) {
		req := &record.RecordRequest{
			Scenario:    "investment_revalue",
			Currency:    currency,
			Description: desc,
			EventTime:   eventTime.Format("2006-01-02 15:04:05"),
		}
		amt, _ := decimal.NewFromString(delta)
		req.Amount = amt
		req.BucketId = buckets[bucketKey].Id
		_, err := service.CreateRecord(userID, req)
		if err != nil {
			log.Printf("WARN revalue skipped: %v", err)
			return
		}
		count++
	}

	// 24 months loop
	for monthOffset := 0; monthOffset < 24; monthOffset++ {
		baseDate := startDate.AddDate(0, monthOffset, 0)
		year, month, _ := baseDate.Date()
		loc := baseDate.Location()

		// salary: 5th
		salary := 12000 + rand.Intn(2000)
		rec("income", time.Date(year, month, 5, 9, 30, 0, 0, loc),
			fmt.Sprintf("%d", salary), "CNY",
			"", "", "cmb_cny", "salary",
			fmt.Sprintf("%d月工资", month))

		// occasional bonus
		if monthOffset%3 == 0 {
			bonus := 3000 + rand.Intn(5000)
			rec("income", time.Date(year, month, 15, 14, 0, 0, 0, loc),
				fmt.Sprintf("%d", bonus), "CNY",
				"", "", "cmb_cny", "bonus",
				"季度奖金")
		}

		// rent: 1st
		rent := 5500
		rec("expense", time.Date(year, month, 1, 10, 0, 0, 0, loc),
			fmt.Sprintf("%d", rent), "CNY",
			"", "", "cmb_cny", "housing_rent",
			fmt.Sprintf("%d月房租", month))

		// utilities: ~8th
		util := 200 + rand.Intn(150)
		rec("expense", time.Date(year, month, 8, 19, 0, 0, 0, loc),
			fmt.Sprintf("%d", util), "CNY",
			"", "", "cmb_cny", "utilities",
			"水电燃气")

		// food expenses (8-12 per month)
		mealCount := 8 + rand.Intn(5)
		for i := 0; i < mealCount; i++ {
			day := 1 + rand.Intn(28)
			isMeal := rand.Intn(3)
			amt := 0
			catKey := ""
			desc := ""
			switch isMeal {
			case 0:
				amt = 25 + rand.Intn(50)
				catKey = "food_meal"
				desc = "美团外卖"
			case 1:
				amt = 80 + rand.Intn(200)
				catKey = "food_dining"
				desc = "餐厅"
			default:
				amt = 18 + rand.Intn(25)
				catKey = "food_coffee"
				desc = "咖啡"
			}
			payBucket := pick(cnyCash, rand.Intn(len(cnyCash)))
			rec("expense", time.Date(year, month, day, 12+rand.Intn(8), rand.Intn(60), 0, 0, loc),
				fmt.Sprintf("%d.%d", amt, rand.Intn(100)), "CNY",
				"", "", payBucket, catKey, desc)
		}

		// grocery 3-4 times
		for i := 0; i < 3+rand.Intn(2); i++ {
			day := 2 + rand.Intn(27)
			amt := 80 + rand.Intn(180)
			rec("expense", time.Date(year, month, day, 19, 0, 0, 0, loc),
				fmt.Sprintf("%d.%d", amt, rand.Intn(100)), "CNY",
				"", "", pick(cnyCash, rand.Intn(len(cnyCash))), "grocery",
				"超市买菜")
		}

		// transport 5-8 times
		for i := 0; i < 5+rand.Intn(4); i++ {
			day := 1 + rand.Intn(28)
			isTaxi := rand.Intn(3) == 0
			catKey := "transport_subway"
			amt := 4 + rand.Intn(8)
			desc := "地铁"
			if isTaxi {
				catKey = "transport_taxi"
				amt = 25 + rand.Intn(40)
				desc = "滴滴打车"
			}
			rec("expense", time.Date(year, month, day, 8, 30, 0, 0, loc),
				fmt.Sprintf("%d", amt), "CNY",
				"", "", "wechat_cny", catKey, desc)
		}

		// shopping 1-3 times
		for i := 0; i < 1+rand.Intn(3); i++ {
			day := 5 + rand.Intn(20)
			isElec := rand.Intn(4) == 0
			amt := 80 + rand.Intn(400)
			catKey := "shopping_clothes"
			desc := "买衣服"
			if isElec {
				amt = 300 + rand.Intn(2000)
				catKey = "shopping_electronics"
				desc = "数码"
			}
			rec("expense", time.Date(year, month, day, 14, 0, 0, 0, loc),
				fmt.Sprintf("%d", amt), "CNY",
				"", "", "alipay_cny", catKey, desc)
		}

		// entertainment / subscription
		rec("expense", time.Date(year, month, 10, 20, 0, 0, 0, loc),
			fmt.Sprintf("%d", 30+rand.Intn(80)), "CNY",
			"", "", "alipay_cny", "entertainment_subscription",
			"流媒体订阅")
		if monthOffset%2 == 0 {
			rec("expense", time.Date(year, month, 18, 19, 0, 0, 0, loc),
				fmt.Sprintf("%d", 60+rand.Intn(140)), "CNY",
				"", "", pick(cnyCash, rand.Intn(len(cnyCash))), "entertainment_movie",
				"看电影")
		}

		// healthcare quarterly
		if monthOffset%4 == 1 {
			rec("expense", time.Date(year, month, 22, 10, 0, 0, 0, loc),
				fmt.Sprintf("%d", 200+rand.Intn(800)), "CNY",
				"", "", "cmb_cny", "healthcare", "医院")
		}

		// education monthly
		if monthOffset%3 == 1 {
			rec("expense", time.Date(year, month, 12, 21, 0, 0, 0, loc),
				fmt.Sprintf("%d", 99+rand.Intn(400)), "CNY",
				"", "", "alipay_cny", "education", "课程")
		}

		// travel half yearly
		if monthOffset%6 == 4 {
			rec("expense", time.Date(year, month, 20, 9, 0, 0, 0, loc),
				fmt.Sprintf("%d", 1500+rand.Intn(3500)), "CNY",
				"", "", "cmb_cny", "travel", "周边旅行")
		}

		// transfer between cash buckets monthly - replenish wechat/alipay
		rec("transfer", time.Date(year, month, 6, 11, 0, 0, 0, loc),
			"3000", "CNY",
			"cmb_cny", "wechat_cny", "", "",
			"补微信零花")
		rec("transfer", time.Date(year, month, 6, 11, 5, 0, 0, loc),
			"3000", "CNY",
			"cmb_cny", "alipay_cny", "", "",
			"补支付宝零花")

		// reimbursement: every 4 months
		if monthOffset%4 == 2 {
			rec("income", time.Date(year, month, 28, 14, 0, 0, 0, loc),
				fmt.Sprintf("%d", 500+rand.Intn(1500)), "CNY",
				"", "", "cmb_cny", "reimbursement",
				"公司报销")
		}

		// monthly investment buy - QQQ + VOO (USD), 2800 (HKD)
		// month 0..23: dollar cost average
		dcaMonth := 4
		if monthOffset%dcaMonth == 0 {
			// US: 200 USD into QQQ
			rec("investment_buy", time.Date(year, month, 7, 10, 0, 0, 0, loc),
				"200", "USD",
				"ibkr_cash", "ibkr_qqq", "", "",
				"QQQ 定投")
			// US: 200 USD into VOO
			rec("investment_buy", time.Date(year, month, 7, 10, 5, 0, 0, loc),
				"200", "USD",
				"ibkr_cash", "ibkr_voo", "", "",
				"VOO 定投")
			// HK: 2000 HKD into 2800
			rec("investment_buy", time.Date(year, month, 9, 10, 0, 0, 0, loc),
				"2000", "HKD",
				"futu_hkd", "futu_hk_etf", "", "",
				"2800 港股定投")
		}

		// quarterly market revalue
		if monthOffset%3 == 2 {
			delta := -50 + rand.Intn(220) // -50 to +169
			recRevalue(time.Date(year, month, 28, 16, 0, 0, 0, loc),
				"ibkr_qqq", fmt.Sprintf("%d", delta), "USD",
				"季度市值刷新")
			delta2 := -40 + rand.Intn(180)
			recRevalue(time.Date(year, month, 28, 16, 5, 0, 0, loc),
				"ibkr_voo", fmt.Sprintf("%d", delta2), "USD",
				"季度市值刷新")
			delta3 := -200 + rand.Intn(900)
			recRevalue(time.Date(year, month, 28, 16, 10, 0, 0, loc),
				"futu_hk_etf", fmt.Sprintf("%d", delta3), "HKD",
				"季度市值刷新")
		}

		// once a year: dividend
		if monthOffset%12 == 6 {
			rec("investment_income", time.Date(year, month, 15, 10, 0, 0, 0, loc),
				"45", "USD",
				"", "", "ibkr_cash", "",
				"QQQ + VOO 分红")
			rec("investment_income", time.Date(year, month, 15, 10, 5, 0, 0, loc),
				"180", "HKD",
				"", "", "futu_hkd", "",
				"2800 派息")
		}

		// occasional exchange CNY → USD (cmb_cny → cmb_usd)
		if monthOffset%5 == 3 {
			recExchange(time.Date(year, month, 22, 11, 0, 0, 0, loc),
				"cmb_cny", "cmb_usd", "2200", "300", "CNY", "USD",
				"换汇 CNY→USD")
			// 然后把 USD 转入 IBKR
			rec("transfer", time.Date(year, month, 22, 11, 30, 0, 0, loc),
				"300", "USD",
				"cmb_usd", "ibkr_cash", "", "",
				"美元入金 IBKR")
		}
	}

	return count
}

func pick(arr []string, i int) string {
	return arr[i%len(arr)]
}

func must(err error, ctx string) {
	if err != nil {
		log.Fatalf("FAIL %s: %v", ctx, err)
	}
}
