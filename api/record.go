package api

import (
	"marmot-ledger/internal/domain/entity/record"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"

	"github.com/gofiber/fiber/v2"
)

func RecordMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[record.RecordResponse]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		req := &record.RecordRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		created, err := service.CreateRecord(userId, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Put("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[record.RecordResponse]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		req := &record.RecordRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		updated, err := service.UpdateRecord(userId, id, req)
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
		if err := service.DeleteRecord(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	return app
}
