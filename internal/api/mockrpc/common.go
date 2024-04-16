package mockrpc

import (
	domainTypes "inttest-runtime/internal/domain/types"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func extractReqParams(ctx *fiber.Ctx) domainTypes.HttpClientRequestParams {
	return domainTypes.HttpClientRequestParams{
		Headers: lo.MapValues(ctx.GetReqHeaders(), func(v []string, _ string) string {
			if len(v) == 0 {
				log.Println("empty slice in http headers")
				return ""
			}
			return v[0]
		}),
		Url:   ctx.AllParams(),
		Query: ctx.Queries(),
		Body:  string(ctx.Body()),
	}
}
