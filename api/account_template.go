package api

import (
	"marmot-ledger/api/middleware"
	"marmot-ledger/internal/domain/entity/account"
	"marmot-ledger/internal/domain/entity/accounttemplate"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AccountTemplateMount() *fiber.App {
	app := fiber.New()

	// 用户接口 - 获取启用的模板列表
	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]accounttemplate.AccountTemplate]{}
		templates, err := service.ListAccountTemplatesForUser()
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(templates))
	})

	// 用户基于模板创建账户
	app.Post("/instantiate", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[account.Account]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		req := &accounttemplate.InstantiateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		if req.TemplateId == 0 {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "模板ID不能为空"))
		}
		newAccount, err := service.InstantiateAccountFromTemplate(userId, req.TemplateId, req.CustomName)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*newAccount))
	})

	// 管理员接口
	adminGroup := app.Group("", middleware.AdminRequired)

	// 管理员获取完整列表
	adminGroup.Get("/admin", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]accounttemplate.AccountTemplate]{}
		templates, err := service.ListAccountTemplatesForAdmin()
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(templates))
	})

	// 管理员获取单个模板
	adminGroup.Get("/admin/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[accounttemplate.AccountTemplate]{}
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的ID"))
		}
		template, err := service.GetAccountTemplate(id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.NotFound), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	// 管理员创建模板
	adminGroup.Post("/admin", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[accounttemplate.AccountTemplate]{}
		req := &accounttemplate.CreateTemplateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		template, err := service.CreateAccountTemplate(req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	// 管理员更新模板
	adminGroup.Put("/admin/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[accounttemplate.AccountTemplate]{}
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的ID"))
		}
		req := &accounttemplate.UpdateTemplateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		template, err := service.UpdateAccountTemplate(id, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	return app
}
