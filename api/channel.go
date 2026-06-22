package api

import (
	"marmot-ledger/internal/domain/entity/channel"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ChannelMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[channel.Channel]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		channelInfo := &channel.Channel{}
		if err := ctx.BodyParser(channelInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		created, err := service.CreateChannel(userId, channelInfo)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]channel.Channel]{}
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
		channels, err := service.ListChannels(userId, channel.ChannelQuery{ChannelType: ctx.Query("channelType"), EventType: ctx.Query("eventType"), IsActive: isActive})
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(channels))
	})

	app.Post("/import", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]channel.Channel]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		req := &channel.ImportRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		items, err := service.ImportChannelTemplates(userId, req.TemplateIds)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(items))
	})

	app.Post("/initialize-defaults", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[map[string]interface{}]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		if err := service.EnsureDefaultChannels(userId); err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(map[string]interface{}{"success": true}))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[channel.Channel]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		channelInfo, err := service.GetChannel(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*channelInfo))
	})

	app.Put("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[channel.Channel]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		channelInfo := &channel.Channel{}
		if err := ctx.BodyParser(channelInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		updated, err := service.UpdateChannel(userId, id, channelInfo)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*updated))
	})

	app.Get("/:id/usage", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[map[string]interface{}]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		count, err := service.CheckChannelUsage(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(map[string]interface{}{"eventCount": count}))
	})

	app.Delete("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[map[string]interface{}]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		count, err := service.DeleteChannel(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(map[string]interface{}{"success": true, "affectedCount": count}))
	})

	return app
}
