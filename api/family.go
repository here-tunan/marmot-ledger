package api

import (
	"marmot-ledger/internal/domain/entity/family"
	"marmot-ledger/internal/domain/entity/statistics"
	"marmot-ledger/internal/domain/repository/statisticsdb"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FamilyMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[family.Family]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		req := &family.CreateFamilyRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		created, err := service.CreateFamily(userId, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]family.Family]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		items, err := service.ListFamilies(userId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(items))
	})

	app.Get("/invitations", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]family.Member]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		items, err := service.ListFamilyInvitations(userId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(items))
	})

	app.Post("/invitations/:id/accept", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bool]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		if err := service.AcceptFamilyInvitation(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	app.Post("/invitations/:id/reject", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bool]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		if err := service.RejectFamilyInvitation(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[family.Family]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		item, err := service.GetFamily(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*item))
	})

	app.Get("/:id/members", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]family.Member]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		includeInvited, _ := strconv.ParseBool(ctx.Query("includeInvited", "false"))
		items, err := service.ListFamilyMembers(userId, id, includeInvited)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(items))
	})

	app.Post("/:id/invitations", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[family.Member]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		req := &family.InviteRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		member, err := service.InviteFamilyMember(userId, id, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*member))
	})

	app.Get("/:id/statistics/summary", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[statistics.Summary]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		data, err := service.GetFamilyStatisticsSummary(userId, id, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*data))
	})

	app.Get("/:id/statistics/category-group", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[statistics.CategoryGroupStats]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		data, err := service.GetFamilyStatisticsCategoryGroup(userId, id, parseStatisticsQuery(ctx))
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*data))
	})

	_ = statisticsdb.StatisticsQuery{}
	return app
}
