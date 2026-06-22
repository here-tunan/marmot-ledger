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

		var events *service.PageResult[financialevent.FinancialEvent]
		if ctx.QueryBool("groupMode", false) {
			events, err = service.ListGroupedFinancialEvents(userId, query)
		} else {
			events, err = service.ListFinancialEvents(userId, query)
		}
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*events))
	})

	app.Get("/outstanding", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]financialevent.OutstandingItem]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		bucketIdStr := ctx.Query("bucketId")
		var bucketId int64
		if bucketIdStr != "" {
			parsed, parseErr := strconv.ParseInt(bucketIdStr, 10, 64)
			if parseErr != nil {
				return ctx.JSON(result.Err(int(myerror.WrongParam), "bucketId is invalid"))
			}
			bucketId = parsed
		}
		eventType := ctx.Query("eventType")
		if eventType == "" {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "eventType is required"))
		}

		items, err := service.ListOutstandingForBucket(userId, bucketId, eventType)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(items))
	})

	app.Get("/group/:groupId", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]service.GroupedEvent]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		groupId, err := strconv.ParseInt(ctx.Params("groupId"), 10, 64)
		if err != nil || groupId == 0 {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "groupId is required"))
		}
		var currentId int64
		if currentIdStr := ctx.Query("currentId"); currentIdStr != "" {
			currentId, _ = strconv.ParseInt(currentIdStr, 10, 64)
		}

		events, err := service.GetEventGroup(userId, groupId, currentId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(events))
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
		Currency:  ctx.Query("currency"),
	}

	if categoryId := ctx.Query("categoryId"); categoryId != "" {
		parsed, err := strconv.ParseInt(categoryId, 10, 64)
		if err != nil {
			return query, err
		}
		query.CategoryId = parsed
	}
	if channelId := ctx.Query("channelId"); channelId != "" {
		parsed, err := strconv.ParseInt(channelId, 10, 64)
		if err != nil {
			return query, err
		}
		query.ChannelId = parsed
	}
	if bucketId := ctx.Query("bucketId"); bucketId != "" {
		parsed, err := strconv.ParseInt(bucketId, 10, 64)
		if err != nil {
			return query, err
		}
		query.BucketId = parsed
	}
	query.Keyword = ctx.Query("keyword")
	if includeInStatistics := ctx.Query("includeInStatistics"); includeInStatistics != "" {
		parsed, err := strconv.ParseBool(includeInStatistics)
		if err != nil {
			return query, err
		}
		query.IncludeInStatistics = &parsed
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
