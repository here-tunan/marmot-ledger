package api

import (
	"marmot-ledger/internal/domain/entity/financialevent"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FinancialEventMount() *fiber.App {
	app := fiber.New()

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[service.PageResult[financialevent.FinancialEvent]]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		query, err := parseFinancialEventQuery(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		events, err := service.ListFinancialEvents(userId, query)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*events))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[financialevent.FinancialEvent]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		eventInfo, err := service.GetFinancialEvent(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*eventInfo))
	})

	return app
}

func parseFinancialEventQuery(ctx *fiber.Ctx) (financialevent.FinancialEventQuery, error) {
	query := financialevent.FinancialEventQuery{
		EventType: ctx.Query("eventType"),
		StartTime: ctx.Query("startTime"),
		EndTime:   ctx.Query("endTime"),
	}

	if page := ctx.Query("page"); page != "" {
		parsed, err := strconv.Atoi(page)
		if err != nil {
			return query, err
		}
		query.Page = parsed
	}
	if pageSize := ctx.Query("pageSize"); pageSize != "" {
		parsed, err := strconv.Atoi(pageSize)
		if err != nil {
			return query, err
		}
		query.PageSize = parsed
	}

	return query, nil
}
