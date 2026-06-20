package api

import (
	"marmot-ledger/api/middleware"
	"marmot-ledger/internal/domain/entity/category"
	"marmot-ledger/internal/domain/entity/categorytemplate"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CategoryTemplateMount() *fiber.App {
	app := fiber.New()

	// 用户接口 - 获取启用的模板列表
	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]categorytemplate.CategoryTemplate]{}
		typeFilter := ctx.Query("type", "")
		templates, err := service.ListCategoryTemplatesForUser(typeFilter)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(templates))
	})

	// 用户批量导入模板分类
	app.Post("/import", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]category.Category]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		req := &categorytemplate.ImportRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		if len(req.TemplateIds) == 0 {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "请选择要导入的分类模板"))
		}
		imported, err := service.ImportTemplatesForUser(userId, req.TemplateIds)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(imported))
	})

	// 管理员接口
	adminGroup := app.Group("", middleware.AdminRequired)

	// 管理员获取完整列表
	adminGroup.Get("/admin", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]categorytemplate.CategoryTemplate]{}
		typeFilter := ctx.Query("type", "")
		templates, err := service.ListCategoryTemplatesForAdmin(typeFilter)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(templates))
	})

	// 管理员获取单个模板
	adminGroup.Get("/admin/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[categorytemplate.CategoryTemplate]{}
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的ID"))
		}
		template, err := service.GetCategoryTemplate(id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.NotFound), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	// 管理员创建模板
	adminGroup.Post("/admin", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[categorytemplate.CategoryTemplate]{}
		req := &categorytemplate.CreateTemplateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		template, err := service.CreateCategoryTemplate(req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	// 管理员更新模板
	adminGroup.Put("/admin/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[categorytemplate.CategoryTemplate]{}
		id, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的ID"))
		}
		req := &categorytemplate.UpdateTemplateRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		template, err := service.UpdateCategoryTemplate(id, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*template))
	})

	return app
}
