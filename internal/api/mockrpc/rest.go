package mockrpc

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RestMockApi struct {
	addr    string
	app     *fiber.App
	service IRestApiService
}

func NewRestMockApi(addr string, service IRestApiService) *RestMockApi {
	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return ctx.Next()
	})
	return &RestMockApi{
		addr:    addr,
		app:     app,
		service: service,
	}
}

func (api *RestMockApi) Register(route, method string) error {
	var registrator func(route string, handlers ...fiber.Handler) fiber.Router
	switch method {
	case fiber.MethodGet:
		registrator = api.app.Get
	case fiber.MethodPost:
		registrator = api.app.Post
	case fiber.MethodPut:
		registrator = api.app.Put
	case fiber.MethodPatch:
		registrator = api.app.Patch
	case fiber.MethodDelete:
		registrator = api.app.Delete
	default:
		return fmt.Errorf("unknown http method: %s", method)
	}

	registrator(route, api.makeHttpMockHandlerImpl(route, method))

	return nil
}

func (api RestMockApi) makeHttpMockHandlerImpl(urlPattern, method string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		params := extractReqParams(ctx)
		resp, err := api.service.HandleRestRequest(ctx.Context(), urlPattern, method, params)
		if err != nil {
			return err // todo: error handler
		}

		for header, headerVal := range resp.Headers {
			ctx.Set(header, headerVal)
		}
		return ctx.Status(int(resp.Status)).SendString(resp.Body)
	}
}
