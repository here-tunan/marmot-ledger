package api

import (
	"marmot-ledger/internal/domain/entity/financialevent"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"

	"github.com/gofiber/fiber/v2"
)

func OutstandingMount() *fiber.App {
	app := fiber.New()

	app.Get("/summary", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[financialevent.OutstandingSummary]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		summary, err := service.GetOutstandingSummary(userId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*summary))
	})

	return app
}
