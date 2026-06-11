package api

import (
	"marmot-ledger/internal/domain/entity/category"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CategoryMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[category.Category]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		categoryInfo := &category.Category{}
		if err := ctx.BodyParser(categoryInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		created, err := service.CreateCategory(userId, categoryInfo)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]category.Category]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		var isActive *bool
		if value := ctx.Query("isActive"); value != "" {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
			}
			isActive = &parsed
		}

		categories, err := service.ListCategories(userId, ctx.Query("type"), isActive)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(categories))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[category.Category]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		categoryInfo, err := service.GetCategory(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*categoryInfo))
	})

	app.Put("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[category.Category]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		categoryInfo := &category.Category{}
		if err := ctx.BodyParser(categoryInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		updated, err := service.UpdateCategory(userId, id, categoryInfo)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*updated))
	})

	app.Delete("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bool]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		if err := service.DeleteCategory(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	return app
}
