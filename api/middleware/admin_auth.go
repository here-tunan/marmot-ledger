package middleware

import (
	"marmot-ledger/internal/domain/repository/userdb"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"

	"github.com/gofiber/fiber/v2"
)

// AdminRequired 管理员权限中间件
func AdminRequired(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	if userId == nil {
		result := &myresult.MyResult[any]{}
		return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
	}

	role, err := userdb.GetUserRoleById(userId.(int64))
	if err != nil {
		result := &myresult.MyResult[any]{}
		return ctx.JSON(result.Err(int(myerror.Unauthorized), "获取用户角色失败"))
	}

	if role != "admin" {
		result := &myresult.MyResult[any]{}
		return ctx.JSON(result.Err(int(myerror.Forbidden), "需要管理员权限"))
	}

	return ctx.Next()
}
