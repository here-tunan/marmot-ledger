package api

import (
	"marmot-ledger/api/middleware"
	"marmot-ledger/internal/domain/entity/chantemplate"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ChannelTemplateMount() *fiber.App {
	app := fiber.New()

	// 用户接口 - 获取启用的模板列表
	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]chantemplate.ChannelTemplate]{}
		templates, err := service.ListChannelTemplatesForUser()
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(templates))
	})

	// 管理员接口
	adminGroup := app.Group("", middleware.AdminRequired)

	// 管理员获取完整列表
	adminGroup.Get("/admin", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]chantemplate.ChannelTemplate]{}
		templates, err := service.ListChannelTemplatesForAdmin()
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(templates))
	})

	// 管理员获取单个模板
	adminGroup.Get("/admin/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[chantemplate.ChannelTemplate]{}
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的ID"))
		}
		template, err := service.GetChannelTemplate(id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.NotFound), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	// 管理员创建模板
	adminGroup.Post("/admin", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[chantemplate.ChannelTemplate]{}
		req := &chantemplate.CreateTemplateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		template, err := service.CreateChannelTemplate(req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	// 管理员更新模板
	adminGroup.Put("/admin/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[chantemplate.ChannelTemplate]{}
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的ID"))
		}
		req := &chantemplate.UpdateTemplateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		template, err := service.UpdateChannelTemplate(id, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	return app
}
