package main

import (
	"fmt"
	"log"
	_ "marmot-ledger/env"
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/domain/repository/financialeventdb"
	"marmot-ledger/internal/infrastructure"
	"marmot-ledger/internal/service"
	"time"

	"github.com/shopspring/decimal"
)

const userID = 1

func main() {
	defer infrastructure.Mysql.Close()

	cashID := bucketID("招行储蓄卡")
	recvID := bucketID("应收款")
	foodCatID := categoryID("餐厅吃饭")
	investID := bucketID("盈透 QQQ 持仓")
	investCashID := bucketID("盈透现金账户")

	now := time.Now()
	stamp := now.Format("150405")

	// 1) split demo: 4 人晚餐
	splitReq := &record.RecordRequest{
		Scenario:           "split",
		CashBucketId:       cashID,
		ReceivableBucketId: recvID,
		Currency:           "CNY",
		Amount:             decimal.NewFromInt(480),
		CategoryId:         foodCatID,
		Description:        "Group Demo · 4 人晚餐 " + stamp,
		EventTime:          now.Add(-2 * time.Hour).Format("2006-01-02 15:04:05"),
		Shares: []record.SplitShare{
			{IsSelf: true, Amount: decimal.NewFromInt(120), Description: "我的晚餐份额"},
			{IsSelf: false, Amount: decimal.NewFromInt(120), Description: "小张晚餐分摊 " + stamp},
			{IsSelf: false, Amount: decimal.NewFromInt(120), Description: "李四晚餐分摊 " + stamp},
			{IsSelf: false, Amount: decimal.NewFromInt(120), Description: "小王晚餐分摊 " + stamp},
		},
	}
	_, err := service.CreateRecord(userID, splitReq)
	must(err, "create split demo")
	log.Println("created split group demo")

	// 2) collect one receivable, related to 小张's receivable_create
	zhangEventID := eventIDByDescription("小张晚餐分摊 " + stamp)
	collectReq := &record.RecordRequest{
		Scenario:                "receivable_collect",
		FromBucketId:            cashID,
		ToBucketId:              recvID,
		Currency:                "CNY",
		Amount:                  decimal.NewFromInt(120),
		RelatedFinancialEventId: zhangEventID,
		Description:             "小张晚餐还款 " + stamp,
		EventTime:               now.Add(-1 * time.Hour).Format("2006-01-02 15:04:05"),
	}
	_, err = service.CreateRecord(userID, collectReq)
	must(err, "create collect demo")
	log.Printf("created related collect demo, related=%d", zhangEventID)

	// 3) investment sell demo: make automatic revalue + sell group
	current := bucketBalance(investID)
	if current.LessThan(decimal.NewFromInt(300)) {
		buyReq := &record.RecordRequest{
			Scenario:     "investment_buy",
			FromBucketId: investCashID,
			ToBucketId:   investID,
			Currency:     "USD",
			Amount:       decimal.NewFromInt(500),
			Description:  "Group Demo · QQQ 买入补仓 " + stamp,
			EventTime:    now.Add(-90 * time.Minute).Format("2006-01-02 15:04:05"),
		}
		_, err = service.CreateRecord(userID, buyReq)
		must(err, "create investment buy before sell")
		current = bucketBalance(investID)
	}

	received := decimal.NewFromInt(260)
	remaining := current.Add(decimal.NewFromInt(80)).Sub(received) // 自动产生 +80 revalue，再卖出 260
	if remaining.IsNegative() {
		remaining = decimal.Zero
	}
	sellReq := &record.RecordRequest{
		Scenario:             "investment_sell",
		FromBucketId:         investID,
		ToBucketId:           investCashID,
		Currency:             "USD",
		Amount:               received,
		RemainingMarketValue: remaining,
		Description:          "Group Demo · QQQ 部分卖出 " + stamp,
		EventTime:            now.Add(-30 * time.Minute).Format("2006-01-02 15:04:05"),
	}
	_, err = service.CreateRecord(userID, sellReq)
	must(err, "create investment sell demo")
	log.Printf("created investment sell group demo, received=%s remaining=%s", received, remaining)

	fmt.Println("Done. Open Records and search keyword: Group Demo or 展示 latest records.")
}

func bucketID(name string) int64 {
	var row struct{ Id int64 }
	has, err := infrastructure.Mysql.SQL("SELECT id FROM bucket WHERE user_id=? AND name=? AND is_deleted=0 LIMIT 1", userID, name).Get(&row)
	must(err, "find bucket "+name)
	if !has {
		log.Fatalf("bucket not found: %s", name)
	}
	return row.Id
}

func bucketBalance(id int64) decimal.Decimal {
	var row struct{ Balance decimal.Decimal }
	has, err := infrastructure.Mysql.SQL("SELECT balance FROM bucket WHERE id=? LIMIT 1", id).Get(&row)
	must(err, "find bucket balance")
	if !has {
		return decimal.Zero
	}
	return row.Balance
}

func categoryID(name string) int64 {
	var row struct{ Id int64 }
	has, err := infrastructure.Mysql.SQL("SELECT id FROM category WHERE user_id=? AND name=? AND is_active=1 LIMIT 1", userID, name).Get(&row)
	must(err, "find category "+name)
	if !has {
		log.Fatalf("category not found: %s", name)
	}
	return row.Id
}

func eventIDByDescription(desc string) int64 {
	var row financialeventdb.FinancialEvent
	has, err := infrastructure.Mysql.SQL("SELECT * FROM financial_event WHERE user_id=? AND description=? AND event_type='receivable_create' AND is_deleted=0 ORDER BY id DESC LIMIT 1", userID, desc).Get(&row)
	must(err, "find event "+desc)
	if !has {
		log.Fatalf("event not found: %s", desc)
	}
	return row.Id
}

func must(err error, ctx string) {
	if err != nil {
		log.Fatalf("%s: %v", ctx, err)
	}
}
