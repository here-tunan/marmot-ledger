package api

import (
	"context"
	"go-my-life/internal/domain/entity/money"
	"go-my-life/internal/domain/repository/moneydb"
	"go-my-life/internal/infrastructure"
	service "go-my-life/internal/service/money"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// MoneyMount 收入支出记账的web层
func MoneyMount() *fiber.App {
	app := fiber.New()

	app.Post("/transaction/list", func(c *fiber.Ctx) error {
		param := &moneydb.TransactionQueryParam{}
		// 解析JSON RequestBody
		err := c.BodyParser(param)
		if err != nil {
			return err
		}
		param.UserIds = append(param.UserIds, c.Locals("userId").(int64))
		// 获取结果
		res, total, err := service.QueryTransactions(*param)
		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"code":    "500",
				"error":   err.Error(),
			})
		}
		return c.JSON(&fiber.Map{
			"success": true,
			"data":    res,
			"total":   total,
		})
	})

	app.Post("/transaction/analysis", func(c *fiber.Ctx) error {
		param := &service.TransactionsAnalysisParam{}

		// 解析JSON RequestBody
		err := c.BodyParser(param)
		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}
		param.UserIds = append(param.UserIds, c.Locals("userId").(int64))
		results, err := service.TransactionsAnalysis(param)
		if err != nil {
			return c.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}

		return c.JSON(&fiber.Map{
			"success": true,
			"data":    results,
		})

	})

	app.Delete("/transaction", func(ctx *fiber.Ctx) error {
		id := ctx.QueryInt("id")
		if id == 0 {
			return ctx.JSON(&fiber.Map{
				"success": false,
			})
		}
		err := service.DeleteTransaction(int64(id))
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"code":    "500",
				"error":   err,
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
		})
	})

	// transaction put
	app.Put("/transaction", func(ctx *fiber.Ctx) error {
		transaction := &money.Transaction{}
		err := ctx.BodyParser(transaction)
		if err != nil {
			log.Println(err)
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "param parse failed!",
			})
		}
		err = service.PutTransaction(ctx.Locals("userId").(int64), transaction)
		success := true
		if err != nil {
			success = false
		}
		return ctx.JSON(&fiber.Map{
			"success": success,
			"error":   err,
		})
	})

	// batch put
	app.Put("/transaction/batch", func(ctx *fiber.Ctx) error {
		var transactions []*moneydb.Transaction
		err := ctx.BodyParser(&transactions)
		if err != nil {
			log.Println(err)
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "param parse failed!",
			})
		}

		// 获取用户id
		userId := ctx.Locals("userId").(int64)

		num, err := service.BatchPutTransaction(userId, transactions)
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err,
				"data":    num,
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
			"data":    num,
		})
	})

	app.Get("/transactionCategory", func(ctx *fiber.Ctx) error {
		userId := ctx.Locals("userId").(int64)
		res, err := service.AllCategoriesByUser(userId)
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
			"data":    res,
		})
	})

	app.Put("/transactionCategory", func(ctx *fiber.Ctx) error {
		param := &moneydb.TransactionCategory{}
		err := ctx.BodyParser(param)
		if err != nil {
			log.Printf("解析请求参数失败: %v", err)
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "请求参数格式错误: " + err.Error(),
				"code":    "400",
			})
		}

		// 设置用户ID
		userId := ctx.Locals("userId").(int64)
		param.UserId = userId

		// 参数验证
		if strings.TrimSpace(param.Name) == "" {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "类型名称不能为空",
				"code":    "400",
			})
		}

		if param.Type != 1 && param.Type != 2 {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "收支类型必须是1（收入）或2（支出）",
				"code":    "400",
			})
		}

		err = service.PutTransactionCategory(param)
		if err != nil {
			log.Printf("保存交易类型失败: %v", err)
			errorMsg := "保存失败"
			if err.Error() != "" {
				errorMsg = err.Error()
			}
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   errorMsg,
				"code":    "500",
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Delete("/transactionCategory", func(ctx *fiber.Ctx) error {
		id := ctx.QueryInt("id")
		if id == 0 {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "ID parameter is required",
				"code":    "400",
			})
		}
		err := service.DeleteTransactionCategory(int64(id))
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Get("/transactionAccount", func(ctx *fiber.Ctx) error {
		userId := ctx.Locals("userId").(int64)
		res, err := service.AllAccountsByUser(userId)
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
			"data":    res,
		})
	})

	app.Put("/transactionAccount", func(ctx *fiber.Ctx) error {
		param := &moneydb.TransactionAccount{}
		err := ctx.BodyParser(param)
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "400",
			})
		}

		// 设置用户ID
		userId := ctx.Locals("userId").(int64)
		param.UserId = userId

		err = service.PutTransactionAccount(param)
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
		})
	})

	app.Delete("/transactionAccount", func(ctx *fiber.Ctx) error {
		id := ctx.QueryInt("id")
		if id == 0 {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "ID parameter is required",
				"code":    "400",
			})
		}
		err := service.DeleteTransactionAccount(int64(id))
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
		})
	})

	// 检查账户是否被交易记录使用
	app.Get("/transactionAccount/usage-check", func(ctx *fiber.Ctx) error {
		accountId := ctx.QueryInt("accountId")
		if accountId == 0 {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "accountId parameter is required",
				"code":    "400",
			})
		}

		count, err := service.CheckAccountUsage(int64(accountId))
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}

		return ctx.JSON(&fiber.Map{
			"success": true,
			"data": &fiber.Map{
				"isUsed": count > 0,
				"count":  count,
			},
		})
	})

	// 检查分类是否被交易记录使用
	app.Get("/transactionCategory/usage-check", func(ctx *fiber.Ctx) error {
		categoryId := ctx.QueryInt("categoryId")
		if categoryId == 0 {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   "categoryId parameter is required",
				"code":    "400",
			})
		}

		count, err := service.CheckCategoryUsage(int64(categoryId))
		if err != nil {
			return ctx.JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
				"code":    "500",
			})
		}

		return ctx.JSON(&fiber.Map{
			"success": true,
			"data": &fiber.Map{
				"isUsed": count > 0,
				"count":  count,
			},
		})
	})

	// 微信XLSX账单导入接口
	app.Post("/transaction/import/wechat", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "没有找到文件",
			})
		}

		// 验证文件扩展名
		if !strings.HasSuffix(strings.ToLower(file.Filename), ".xlsx") {
			return ctx.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "微信账单仅支持XLSX格式文件",
			})
		}

		xlsxFile, err := file.Open()
		if err != nil {
			return ctx.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "无法打开文件",
			})
		}
		defer func(xlsxFile multipart.File) {
			err := xlsxFile.Close()
			if err != nil {
				log.Print(err)
			}
		}(xlsxFile)

		result, err := service.ProcessWeChatXLSX(ctx.Locals("userId").(int64), xlsxFile)
		if err != nil {
			return ctx.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return ctx.JSON(&fiber.Map{
			"success":     true,
			"data":        result.Transactions,
			"warnings":    result.Warnings,
			"encoding":    result.Encoding,
			"total_count": len(result.Transactions),
			"platform":    "wechat",
		})
	})

	// 支付宝CSV账单导入接口
	app.Post("/transaction/import/alipay", func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "没有找到文件",
			})
		}

		// 验证文件扩展名
		if !strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
			return ctx.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   "支付宝账单仅支持CSV格式文件",
			})
		}

		csvFile, err := file.Open()
		if err != nil {
			return ctx.Status(500).JSON(&fiber.Map{
				"success": false,
				"error":   "无法打开文件",
			})
		}
		defer func(csvFile multipart.File) {
			err := csvFile.Close()
			if err != nil {
				log.Print(err)
			}
		}(csvFile)

		result, err := service.ProcessAlipayCSV(ctx.Locals("userId").(int64), csvFile)
		if err != nil {
			return ctx.Status(400).JSON(&fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		return ctx.JSON(&fiber.Map{
			"success":     true,
			"data":        result.Transactions,
			"warnings":    result.Warnings,
			"encoding":    result.Encoding,
			"total_count": len(result.Transactions),
			"platform":    "alipay",
		})
	})

	app.Get("/transaction/category_analysis", func(ctx *fiber.Ctx) error {
		desc := ctx.Query("desc")
		transactionType := ctx.QueryInt("type", 1) // 默认为收入类型

		res := service.AnalysisCategory(desc, transactionType)
		return ctx.JSON(&fiber.Map{
			"success": true,
			"data":    res,
		})
	})

	app.Get("/transaction/category_immigrate", func(ctx *fiber.Ctx) error {
		param := moneydb.TransactionQueryParam{
			PageSize:  100,
			PageIndex: 1,
			StartTime: "2020-01-01",
			EndTime:   "2025-01-01",
		}

		transactions, _, _ := service.QueryTransactions(param)
		for transactions != nil && len(transactions) > 0 {
			for _, transaction := range transactions {
				doc := moneydb.TransactionIndex{
					Id:          transaction.Id,
					Description: transaction.Description,
					Type:        transaction.Type,
					Category:    transaction.Category,
				}
				infrastructure.EsClient.Index(moneydb.EsIndex).Id(strconv.FormatInt(transaction.Id, 10)).Request(doc).Do(context.Background())
			}
			param.PageIndex = param.PageIndex + 1
			transactions, _, _ = service.QueryTransactions(param)
		}
		return ctx.JSON(&fiber.Map{
			"success": true,
		})
	})
	return app
}
