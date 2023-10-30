package api

import (
	"errors"
	"inttest-runtime/internal/errors/internalErr"

	"github.com/gofiber/fiber/v2"
)

func ErrHandler(ctx *fiber.Ctx, err error) error {
	var t *errResponse
	if errors.As(err, &t) {
		httpCode := t.httpCode
		if httpCode == 0 {
			httpCode = fiber.StatusInternalServerError
		}
		ctx.Response().Header.Set(_headerContentType, _contentTypeApplicationJson)
		return ctx.Status(httpCode).JSON(t)
	}
	return fiber.DefaultErrorHandler(ctx, err)
}

func ResponseContentTypeMiddleware(ctx *fiber.Ctx) error {
	ctx.Set(_headerContentType, _contentTypeApplicationJsonUtf8)
	return ctx.Next()
}

func RequestContentTypeMiddleware(ctx *fiber.Ctx) error {
	contentTypeHeader := ctx.GetReqHeaders()[_headerContentType]
	if len(contentTypeHeader) == 0 || (len(contentTypeHeader) != 0 && contentTypeHeader[0] != _contentTypeApplicationJson) {
		return NewErrResponse(
			errors.New("data needs to be json encoded"),
			ErrReasonBadRequest,
			"json parsing",
			ErrRespWithHttpCode(fiber.StatusBadRequest),
			ErrRespWithCode(internalErr.ErrCodeJSONParsing),
		)
	}
	return ctx.Next()
}
