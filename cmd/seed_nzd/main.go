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
	"math/rand"
	"time"

	"github.com/shopspring/decimal"
)

const userID = 1

func main() {
	defer infrastructure.Mysql.Close()
	rand.Seed(20260615)

	wise := accountByName("Wise NZD Account")
	if wise == nil {
		created, err := service.CreateAccount(userID, &account.Account{Name: "Wise NZD Account", Type: "wallet", Icon: "Wallet", Color: "#14b8a6", IsActive: true})
		must(err, "create wise")
		wise = created
	}
	asb := accountByName("ASB Bank")
	if asb == nil {
		created, err := service.CreateAccount(userID, &account.Account{Name: "ASB Savings", Type: "bank", Icon: "CreditCard", Color: "#fbbf24", IsActive: true})
		must(err, "create asb")
		asb = created
	}

	wiseBucket := bucketByName("Wise NZD 余额")
	if wiseBucket == nil {
		wiseBucket = createBucket(wise.Id, "Wise NZD 余额", "wallet", "asset", "NZD", "4500")
	}
	asbBucket := bucketByName("ASB Savings")
	if asbBucket == nil {
		asbBucket = createBucket(asb.Id, "ASB Savings", "bank", "asset", "NZD", "8000")
	}
	recvBucket := bucketByName("NZD 应收款")
	if recvBucket == nil {
		recvBucket = createBucket(wise.Id, "NZD 应收款", "receivable", "asset", "NZD", "0")
	}

	cats := loadCategories()
	start := time.Now().AddDate(-1, 0, 0)
	count := 0
	rec := func(scenario string, when time.Time, amount string, bucketID, fromID, toID int64, catKey, desc string) {
		amt, _ := decimal.NewFromString(amount)
		req := &record.RecordRequest{Scenario: scenario, Currency: "NZD", Amount: amt, BucketId: bucketID, FromBucketId: fromID, ToBucketId: toID, Description: desc, EventTime: when.Format("2006-01-02 15:04:05")}
		if catKey != "" {
			req.CategoryId = cats[catKey]
		}
		_, err := service.CreateRecord(userID, req)
		if err != nil {
			log.Printf("WARN skip %s %s %s: %v", scenario, amount, desc, err)
			return
		}
		count++
	}

	for i := 0; i < 12; i++ {
		base := start.AddDate(0, i, 0)
		y, m, _ := base.Date()
		loc := base.Location()
		income := 1800 + rand.Intn(1400)
		rec("income", time.Date(y, m, 3, 10, 0, 0, 0, loc), fmt.Sprintf("%d", income), wiseBucket.Id, 0, 0, "salary", "海外项目收入")
		transfer := 500 + rand.Intn(900)
		rec("transfer", time.Date(y, m, 4, 11, 0, 0, 0, loc), fmt.Sprintf("%d", transfer), 0, wiseBucket.Id, asbBucket.Id, "", "转入 ASB 储蓄")
		for j := 0; j < 3+rand.Intn(3); j++ {
			day := 5 + rand.Intn(20)
			kind := rand.Intn(5)
			cat, desc, amount := "food_coffee", "Flat white", 6+rand.Intn(8)
			switch kind {
			case 1:
				cat, desc, amount = "grocery", "Countdown 超市", 40+rand.Intn(90)
			case 2:
				cat, desc, amount = "transport_subway", "AT HOP 交通", 8+rand.Intn(20)
			case 3:
				cat, desc, amount = "travel", "周末短途", 120+rand.Intn(400)
			case 4:
				cat, desc, amount = "entertainment_subscription", "Spotify NZ", 16
			}
			rec("expense", time.Date(y, m, day, 12+rand.Intn(8), 0, 0, 0, loc), fmt.Sprintf("%d", amount), wiseBucket.Id, 0, 0, cat, desc)
		}
		if i%4 == 1 {
			// NZD split-like receivable create/collect sample
			rec("receivable_create", time.Date(y, m, 18, 19, 0, 0, 0, loc), "80", 0, wiseBucket.Id, recvBucket.Id, "", "朋友 BBQ 分摊")
			rec("receivable_collect", time.Date(y, m, 25, 20, 0, 0, 0, loc), "80", 0, wiseBucket.Id, recvBucket.Id, "", "朋友 BBQ 还款")
		}
	}
	log.Printf("Created %d NZD records", count)
}

func createBucket(accountID int64, name, typ, nature, currency, initial string) *bucket.Bucket {
	amt, _ := decimal.NewFromString(initial)
	created, err := service.CreateBucket(userID, &bucket.Bucket{AccountId: accountID, Name: name, BucketType: typ, BucketNature: nature, Currency: currency, InitialBalance: amt})
	must(err, "create bucket "+name)
	return created
}
func accountByName(name string) *account.Account {
	row := &account.Account{}
	has, err := infrastructure.Mysql.SQL("SELECT * FROM account WHERE user_id=? AND name=? AND is_deleted=0 LIMIT 1", userID, name).Get(row)
	must(err, "find account")
	if !has {
		return nil
	}
	return row
}
func bucketByName(name string) *bucket.Bucket {
	row := &bucket.Bucket{}
	has, err := infrastructure.Mysql.SQL("SELECT * FROM bucket WHERE user_id=? AND name=? AND is_deleted=0 LIMIT 1", userID, name).Get(row)
	must(err, "find bucket")
	if !has {
		return nil
	}
	return row
}
func loadCategories() map[string]int64 {
	rows, err := infrastructure.Mysql.QueryString("SELECT id, name FROM category WHERE user_id=1 AND is_active=1")
	must(err, "categories")
	out := map[string]int64{}
	for _, r := range rows {
		var id int64
		fmt.Sscan(r["id"], &id)
		switch r["name"] {
		case "工资":
			out["salary"] = id
		case "咖啡奶茶":
			out["food_coffee"] = id
		case "买菜日用":
			out["grocery"] = id
		case "地铁公交":
			out["transport_subway"] = id
		case "旅行":
			out["travel"] = id
		case "订阅服务":
			out["entertainment_subscription"] = id
		}
	}
	return out
}
func must(err error, ctx string) {
	if err != nil {
		log.Fatalf("%s: %v", ctx, err)
	}
}
