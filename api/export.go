package api

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"marmot-ledger/internal/domain/entity/bucket"
	"marmot-ledger/internal/domain/entity/financialevent"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ExportMount() *fiber.App {
	app := fiber.New()

	app.Get("/records.csv", func(ctx *fiber.Ctx) error {
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.Status(401).SendString(myerror.Unauthorized.String())
		}
		query, err := parseFinancialEventQuery(ctx)
		if err != nil {
			return ctx.Status(400).SendString(err.Error())
		}
		events, err := service.ExportFinancialEvents(userId, query)
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		rows := make([][]string, 0, len(events))
		for _, ev := range events {
			rows = append(rows, []string{
				strconv.FormatInt(ev.Id, 10), ev.EventTime, ev.EventType, ev.Description, ev.Currency, ev.Amount.String(), boolText(ev.IncludeInStatistics),
				strconv.FormatInt(ev.CategoryId, 10), strconv.FormatInt(ev.ChannelId, 10), strconv.FormatInt(ev.RelatedFinancialEventId, 10), strconv.FormatInt(ev.EventGroupId, 10), ev.Remark,
			})
		}
		return writeCSV(ctx, fileName("marmot-records"), []string{"ID", "时间", "类型", "描述", "币种", "金额", "计入统计", "分类ID", "渠道ID", "关联事件ID", "事件组ID", "备注"}, rows)
	})

	app.Get("/buckets.csv", func(ctx *fiber.Ctx) error {
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.Status(401).SendString(myerror.Unauthorized.String())
		}
		items, err := service.ListBuckets(userId, bucket.BucketQuery{})
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		rows := make([][]string, 0, len(items))
		for _, b := range items {
			rows = append(rows, []string{
				strconv.FormatInt(b.Id, 10), b.Name, b.Currency, b.Balance.String(), b.InitialBalance.String(), b.BucketType, b.BucketNature, boolText(b.IsActive), strconv.FormatInt(b.AccountId, 10),
			})
		}
		return writeCSV(ctx, fileName("marmot-buckets"), []string{"ID", "名称", "币种", "当前余额", "初始余额", "类型", "性质", "是否启用", "账户ID"}, rows)
	})

	app.Get("/outstanding.csv", func(ctx *fiber.Ctx) error {
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.Status(401).SendString(myerror.Unauthorized.String())
		}
		summary, err := service.GetOutstandingSummary(userId)
		if err != nil {
			return ctx.Status(500).SendString(err.Error())
		}
		rows := make([][]string, 0)
		appendRows := func(kind string, items []financialevent.OutstandingItem) {
			for _, item := range items {
				rows = append(rows, []string{kind, strconv.FormatInt(item.Id, 10), item.EventTime, item.Description, item.Currency, item.Amount.String(), item.OutstandingAmount.String(), strconv.FormatInt(item.BucketId, 10), item.BucketName, strconv.FormatInt(item.EventGroupId, 10)})
			}
		}
		appendRows("receivable", summary.Receivables)
		appendRows("deposit", summary.Deposits)
		appendRows("loan_out", summary.LoansOut)
		return writeCSV(ctx, fileName("marmot-outstanding"), []string{"类型", "原事件ID", "时间", "描述", "币种", "原始金额", "未结清金额", "资金桶ID", "资金桶名称", "事件组ID"}, rows)
	})

	return app
}

func boolText(v bool) string {
	if v {
		return "是"
	}
	return "否"
}
func fileName(prefix string) string {
	return fmt.Sprintf("%s-%s.csv", prefix, time.Now().Format("20060102-150405"))
}

func writeCSV(ctx *fiber.Ctx, filename string, headers []string, rows [][]string) error {
	buf := &bytes.Buffer{}
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(buf)
	if err := writer.Write(headers); err != nil {
		return err
	}
	if err := writer.WriteAll(rows); err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	ctx.Set("Content-Type", "text/csv; charset=utf-8")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	return ctx.Send(buf.Bytes())
}
