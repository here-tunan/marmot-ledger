package api

import (
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/domain/repository/statisticsdb"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"

	"github.com/gofiber/fiber/v2"
)

func StatisticsMount() *fiber.App {
	app := fiber.New()

	app.Get("/summary", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[statistics.Summary]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		summary, err := service.GetStatisticsSummary(userId, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*summary))
	})

	app.Get("/category-group", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[statistics.CategoryGroupStats]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		stats, err := service.GetStatisticsCategoryGroup(userId, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*stats))
	})

	return app
}

func parseStatisticsQuery(ctx *fiber.Ctx) statisticsdb.StatisticsQuery {
	return statisticsdb.StatisticsQuery{
		StartTime: ctx.Query("startTime"),
		EndTime:   ctx.Query("endTime"),
		Currency:  ctx.Query("currency"),
		Type:      ctx.Query("type"),
	}
}
