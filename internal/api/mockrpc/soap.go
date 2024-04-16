package mockrpc

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type SoapMockApi struct {
	addr    string
	app     *fiber.App
	service ISoapApiService
}

func NewSoapMockApi(addr string, service ISoapApiService) *SoapMockApi {
	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)
		return ctx.Next()
	})

	return &SoapMockApi{
		addr:    addr,
		app:     app,
		service: service,
	}
}

func (api SoapMockApi) Register(route, method string) error {
	var registrator func(route string, handlers ...fiber.Handler) fiber.Router
	switch method {
	case fiber.MethodGet:
		registrator = api.app.Get
	case fiber.MethodPost:
		registrator = api.app.Post
	default:
		return fmt.Errorf("unknown http method: %s", method)
	}

	registrator(route, api.makeSoapMockHandlerImpl(route, method))
	return nil
}

func (api SoapMockApi) makeSoapMockHandlerImpl(route, method string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		params := extractReqParams(ctx)
		resp, err := api.service.HandleSoapRequest(ctx.Context(), route, method, params)
		if err != nil {
			return err // todo: error handler
		}

		for header, headerVal := range resp.Headers {
			ctx.Set(header, headerVal)
		}
		return ctx.Status(int(resp.Status)).SendString(resp.Body)
	}
}
