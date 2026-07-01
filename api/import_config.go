package api

import (
	"fmt"
	"io"
	"marmot-ledger/internal/domain/entity/importconfig"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func ImportConfigMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[importconfig.ImportConfig]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		info := &importconfig.ImportConfig{}
		if err := ctx.BodyParser(info); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		created, err := service.CreateImportConfig(userId, info)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]importconfig.ImportConfig]{}
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
		configs, err := service.ListImportConfigs(userId, importconfig.ImportConfigQuery{IsActive: isActive})
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(configs))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[importconfig.ImportConfig]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		info, err := service.GetImportConfig(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*info))
	})

	app.Put("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[importconfig.ImportConfig]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		info := &importconfig.ImportConfig{}
		if err := ctx.BodyParser(info); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		updated, err := service.UpdateImportConfig(userId, id, info)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*updated))
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
		if err := service.DeleteImportConfig(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(map[string]interface{}{"success": true}))
	})

	app.Post("/:id/preview", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[importconfig.ImportPreview]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		configInfo, err := service.GetImportConfig(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		fileHeader, err := ctx.FormFile("file")
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "file is required"))
		}
		if fileHeader.Size > 10*1024*1024 {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "file is too large (max 10MB)"))
		}
		src, err := fileHeader.Open()
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		defer src.Close()
		data, err := io.ReadAll(src)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		var defaultBucketId int64
		if value := ctx.FormValue("defaultBucketId"); value != "" {
			if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
				defaultBucketId = parsed
			}
		}
		headerRow, dataRows, truncated, err := service.ParseImportTable(configInfo, data)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		preview, err := service.BuildPreviewRows(configInfo, headerRow, dataRows, defaultBucketId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		if truncated {
			preview.Warnings = append(preview.Warnings, fmt.Sprintf("文件超过 %d 行，已截断仅展示前 %d 行；如需导入全部行，请拆分文件。", service.MaxImportPreviewRows, service.MaxImportPreviewRows))
		}
		return ctx.JSON(result.OK(*preview))
	})

	app.Post("/:id/commit", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[service.ImportCommitResult]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		// 校验配置归属，避免用别人配置的 id 提交（虽然 rows 里没有 configId 依赖，但保持权限一致）
		if _, err := service.GetImportConfig(userId, id); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		req := &struct {
			Rows []service.ImportCommitRow `json:"rows"`
		}{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		commitResult, err := service.CommitImportRecords(userId, req.Rows)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*commitResult))
	})

	return app
}
