package api

import (
	"marmot-ledger/internal/domain/entity/bucket"
	"marmot-ledger/internal/domain/entity/ledgerentry"
	"marmot-ledger/internal/service"
	"marmot-ledger/pkg/myerror"
	"marmot-ledger/pkg/myresult"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func BucketMount() *fiber.App {
	app := fiber.New()

	app.Post("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bucket.Bucket]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		bucketInfo := &bucket.Bucket{}
		if err := ctx.BodyParser(bucketInfo); err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		created, err := service.CreateBucket(userId, bucketInfo)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*created))
	})

	app.Get("", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]bucket.Bucket]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		query, err := parseBucketQuery(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		buckets, err := service.ListBuckets(userId, query)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.InternalError), err.Error()))
		}
		return ctx.JSON(result.OK(buckets))
	})

	app.Get("/:id", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[bucket.Bucket]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		bucketInfo, err := service.GetBucket(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(*bucketInfo))
	})

	app.Get("/:id/ledger-entry", func(ctx *fiber.Ctx) error {
		result := &myresult.MyResult[[]ledgerentry.LedgerEntry]{}
		userId, ok := getLoginUserId(ctx)
		if !ok {
			return ctx.JSON(result.Err(int(myerror.Unauthorized), myerror.Unauthorized.String()))
		}

		id, err := parseIdParam(ctx)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}

		entries, err := service.ListBucketLedgerEntries(userId, id)
		if err != nil {
			return ctx.JSON(result.Err(int(myerror.WrongParam), err.Error()))
		}
		return ctx.JSON(result.OK(entries))
	})

	return app
}

func parseBucketQuery(ctx *fiber.Ctx) (bucket.BucketQuery, error) {
	query := bucket.BucketQuery{
		Currency:     ctx.Query("currency"),
		BucketType:   ctx.Query("bucketType"),
		BucketNature: ctx.Query("bucketNature"),
	}

	if accountId := ctx.Query("accountId"); accountId != "" {
		parsed, err := strconv.ParseInt(accountId, 10, 64)
		if err != nil {
			return query, err
		}
		query.AccountId = parsed
	}
	if isActive := ctx.Query("isActive"); isActive != "" {
		parsed, err := strconv.ParseBool(isActive)
		if err != nil {
			return query, err
		}
		query.IsActive = &parsed
	}

	return query, nil
}
