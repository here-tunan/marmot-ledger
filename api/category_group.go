package api

import (
	"marmot-ledger/internal/domain/entity/categorygroup"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CategoryGroupMount() *fiber.App {
	app := fiber.New()

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]categorygroup.CategoryGroup]{}

		var enabled *bool
		if value := ctx.Query("enabled"); value != "" {
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
			}
			enabled = &parsed
		} else {
			parsed := true
			enabled = &parsed
		}

		groups, err := service.ListCategoryGroups(ctx.Query("type"), enabled)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(groups))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[categorygroup.CategoryGroup]{}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		group, err := service.GetCategoryGroup(id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*group))
	})

	return app
}
