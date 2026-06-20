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

	app.Get("/summaries", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]statistics.Summary]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		summaries, err := service.GetStatisticsSummaries(userId, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(summaries))
	})

	app.Get("/category-groups", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]statistics.CategoryGroupStats]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		stats, err := service.GetStatisticsCategoryGroups(userId, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(stats))
	})

	app.Get("/trend", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]statistics.TrendPoint]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		granularity := ctx.Query("granularity")
		if granularity != "week" {
			granularity = "month"
		}
		points, err := service.GetStatisticsTrend(userId, parseStatisticsQuery(ctx), granularity)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(points))
	})

	app.Get("/investment", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]statistics.InvestmentSummary]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		items, err := service.GetInvestmentSummaries(userId, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(items))
	})

	app.Get("/net-worth-trend", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]statistics.NetWorthTrendPoint]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		granularity := ctx.Query("granularity")
		if granularity != "week" {
			granularity = "month"
		}
		points, err := service.GetNetWorthTrend(userId, parseStatisticsQuery(ctx), granularity)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(points))
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
