package main

import (
	"fmt"
	"log"
	_ "marmot-ledger/env"
	"marmot-ledger/internal/domain/entity/account"
	"marmot-ledger/internal/domain/entity/bucket"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/internal/service"
	"time"

	"github.com/shopspring/decimal"
)

const userID = 1

func main() {
	defer infrastructure.Mysql.Close()

	ccb := createAccount("建设银行房贷", "bank", "CreditCard", "#1d4ed8")
	huabei := createAccount("蚂蚁花呗", "liability", "Warning", "#06b6d4")
	cmbCredit := findAccountByName("招行信用卡")

	mortgage := createLiabilityBucket(ccb.Id, "建设银行房贷", "liability", "CNY", "720000")
	huabeiBucket := createLiabilityBucket(huabei.Id, "花呗账单", "liability", "CNY", "8000")
	creditBucket := findBucketByName("招行信用卡")
	if creditBucket == nil {
		creditBucket = createLiabilityBucket(cmbCredit.Id, "招行信用卡", "credit", "CNY", "3000")
	}

	start := time.Now().AddDate(-2, 0, 0)
	count := 0

	// 房贷：每月还本金约 2400-3100，余额逐月下降。
	for i := 0; i < 24; i++ {
		date := start.AddDate(0, i, 12)
		principal := 2400 + (i%5)*170
		createAdjustment(mortgage.Id, mortgage.Currency, decimal.NewFromInt(int64(-principal)), date, fmt.Sprintf("房贷本金归还 %d", principal))
		count++
	}

	// 花呗：每月消费后增加，下月还款后下降；制造短期负债波动。
	for i := 0; i < 24; i++ {
		date := start.AddDate(0, i, 8)
		spend := 1200 + (i%6)*180
		createAdjustment(huabeiBucket.Id, huabeiBucket.Currency, decimal.NewFromInt(int64(spend)), date, fmt.Sprintf("花呗当月账单 %d", spend))
		count++
		payDate := start.AddDate(0, i, 20)
		pay := spend - 150 + (i%3)*50
		if pay < 0 {
			pay = spend
		}
		createAdjustment(huabeiBucket.Id, huabeiBucket.Currency, decimal.NewFromInt(int64(-pay)), payDate, fmt.Sprintf("花呗还款 %d", pay))
		count++
	}

	// 信用卡：每月账单 + 还款，制造信用卡负债波动。
	for i := 0; i < 24; i++ {
		date := start.AddDate(0, i, 15)
		bill := 1800 + (i%7)*260
		createAdjustment(creditBucket.Id, creditBucket.Currency, decimal.NewFromInt(int64(bill)), date, fmt.Sprintf("信用卡账单 %d", bill))
		count++
		payDate := start.AddDate(0, i, 27)
		pay := bill - 200 + (i%4)*80
		createAdjustment(creditBucket.Id, creditBucket.Currency, decimal.NewFromInt(int64(-pay)), payDate, fmt.Sprintf("信用卡还款 %d", pay))
		count++
	}

	log.Printf("Created %d liability adjustment records", count)
	showBalances()
}

func createAccount(name, accType, icon, color string) *account.Account {
	if existing := findAccountByName(name); existing != nil {
		return existing
	}
	created, err := service.CreateAccount(userID, &account.Account{Name: name, Type: accType, Icon: icon, Color: color, IsActive: true})
	must(err, "create account "+name)
	return created
}

func createLiabilityBucket(accountID int64, name, bucketType, currency, initial string) *bucket.Bucket {
	if existing := findBucketByName(name); existing != nil {
		return existing
	}
	amount, _ := decimal.NewFromString(initial)
	created, err := service.CreateBucket(userID, &bucket.Bucket{AccountId: accountID, Name: name, Currency: currency, InitialBalance: amount, BucketType: bucketType, BucketNature: "liability"})
	must(err, "create bucket "+name)
	return created
}

func createAdjustment(bucketID int64, currency string, delta decimal.Decimal, eventTime time.Time, desc string) {
	_, err := service.CreateRecord(userID, &record.RecordRequest{
		Scenario:    "balance_adjustment",
		BucketId:    bucketID,
		Amount:      delta,
		Currency:    currency,
		Description: desc,
		EventTime:   eventTime.Format("2006-01-02 15:04:05"),
		Remark:      "负债样例数据",
	})
	must(err, desc)
}

func findAccountByName(name string) *account.Account {
	row := &account.Account{}
	has, err := infrastructure.Mysql.SQL("SELECT id, user_id, name, type, icon, color, is_active FROM account WHERE user_id = ? AND name = ? AND is_deleted = 0 LIMIT 1", userID, name).Get(row)
	must(err, "find account "+name)
	if !has {
		return nil
	}
	return row
}

func findBucketByName(name string) *bucket.Bucket {
	row := &bucket.Bucket{}
	has, err := infrastructure.Mysql.SQL("SELECT id, user_id, account_id, name, currency, balance, initial_balance, bucket_type, bucket_nature, is_active FROM bucket WHERE user_id = ? AND name = ? AND is_deleted = 0 LIMIT 1", userID, name).Get(row)
	must(err, "find bucket "+name)
	if !has {
		return nil
	}
	return row
}

func showBalances() {
	rows, err := infrastructure.Mysql.QueryString("SELECT name, currency, balance FROM bucket WHERE user_id = 1 AND bucket_nature = 'liability' AND is_deleted = 0 ORDER BY name")
	must(err, "show balances")
	for _, r := range rows {
		log.Printf("%s %s %s", r["name"], r["currency"], r["balance"])
	}
}

func must(err error, ctx string) {
	if err != nil {
		log.Fatalf("FAIL %s: %v", ctx, err)
	}
}
