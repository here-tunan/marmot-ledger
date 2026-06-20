package api

import (
	"marmot-ledger/internal/domain/entity/familycategorygroup"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FamilyCategoryGroupMount() *fiber.App {
	app := fiber.New()

	// 获取家庭分类组列表
	app.Get("/family/:familyId/category-group", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]familycategorygroup.FamilyCategoryGroup]{}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		typeFilter := ctx.Query("type", "")
		groups, err := service.ListFamilyCategoryGroups(familyId, typeFilter)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(groups))
	})

	// 获取单个家庭分类组详情
	app.Get("/family/:familyId/category-group/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[familycategorygroup.FamilyCategoryGroup]{}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		groupId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分组ID"))
		}
		group, err := service.GetFamilyCategoryGroup(familyId, groupId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.NotFound), err.Error()))
		}
		return ctx.JSON(result.OK(*group))
	})

	// 创建家庭分类组
	app.Post("/family/:familyId/category-group", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[familycategorygroup.FamilyCategoryGroup]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		req := &familycategorygroup.CreateGroupRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		group, err := service.CreateFamilyCategoryGroup(familyId, userId, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*group))
	})

	// 更新家庭分类组
	app.Put("/family/:familyId/category-group/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[familycategorygroup.FamilyCategoryGroup]{}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		groupId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分组ID"))
		}
		req := &familycategorygroup.UpdateGroupRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		group, err := service.UpdateFamilyCategoryGroup(familyId, groupId, req)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(*group))
	})

	// 删除家庭分类组
	app.Delete("/family/:familyId/category-group/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bool]{}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		groupId, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分组ID"))
		}
		err = service.DeleteFamilyCategoryGroup(familyId, groupId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	// 批量添加分类到分组
	app.Post("/family/:familyId/category-group/:groupId/members", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bool]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		groupId, err := strconv.ParseInt(ctx.Params("groupId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分组ID"))
		}
		req := &familycategorygroup.AddMembersRequest{}
		if err := ctx.BodyParser(req); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		err = service.AddCategoriesToGroup(familyId, groupId, req.CategoryIds, userId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	// 从分组移除分类
	app.Delete("/family/:familyId/category-group/:groupId/members/:categoryId", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bool]{}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		groupId, err := strconv.ParseInt(ctx.Params("groupId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分组ID"))
		}
		categoryId, err := strconv.ParseInt(ctx.Params("categoryId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分类ID"))
		}
		err = service.RemoveCategoryFromGroup(familyId, groupId, categoryId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(true))
	})

	// 获取分组下所有分类ID
	app.Get("/family/:familyId/category-group/:groupId/members", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]int64]{}
		familyId, err := strconv.ParseInt(ctx.Params("familyId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的家庭ID"))
		}
		groupId, err := strconv.ParseInt(ctx.Params("groupId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分组ID"))
		}
		memberIds, err := service.GetGroupMemberIds(familyId, groupId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(memberIds))
	})

	// 获取分类所属的所有家庭分组
	app.Get("/category/:categoryId/family-groups", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]int64]{}
		categoryId, err := strconv.ParseInt(ctx.Params("categoryId"), 10, 64)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), "无效的分类ID"))
		}
		groupIds, err := service.GetGroupsForCategory(categoryId)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(groupIds))
	})

	return app
}
