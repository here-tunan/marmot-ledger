package api

import (
	"marmot-ledger/internal/domain/entity/account"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AccountMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[account.Account]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		accountInfo := &account.Account{}
		if err := ctx.BodyParser(accountInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		created, err := service.CreateAccount(userId, accountInfo)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]account.Account]{}
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

		accounts, err := service.ListAccounts(userId, ctx.Query("type"), isActive)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(accounts))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[account.Account]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		accountInfo, err := service.GetAccount(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*accountInfo))
	})

	app.Put("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[account.Account]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		accountInfo := &account.Account{}
		if err := ctx.BodyParser(accountInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		updated, err := service.UpdateAccount(userId, id, accountInfo)
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

		if err := service.DeleteAccount(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	return app
}

func getLoginUserId(ctx *fiber.Ctx) (int64, bool) {
	userId, ok := ctx.Locals("userId").(int64)
	return userId, ok
}

func parseIdParam(ctx *fiber.Ctx) (int64, error) {
	return strconv.ParseInt(ctx.Params("id"), 10, 64)
}
