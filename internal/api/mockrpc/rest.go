package mockrpc

import (
	"errors"
	"fmt"
	domainTypes "inttest-runtime/internal/domain/types"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RestMockApi struct {
	httpRpcApi
	service IRestApiService
}

func NewRestMockApi(service IRestApiService) *RestMockApi {
	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		if !strings.Contains(ctx.Get(fiber.HeaderContentType), fiber.MIMEApplicationJSON) {
			return errors.New("REST API only accepts application/json")
		}
		return ctx.Next()
	})
	return &RestMockApi{
		httpRpcApi: httpRpcApi{
			app: app,
		},
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

func (api *RestMockApi) makeHttpMockHandlerImpl(urlPattern, method string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// we need to pass only parsed values into business logic
		var model any
		if err := ctx.BodyParser(&model); err != nil {
			return err
		}
		resp, err := api.service.HandleRestRequest(ctx.Context(), urlPattern, method, domainTypes.RestClientRequestParams{
			UrlParams:   ctx.AllParams(),
			QueryParams: ctx.Queries(),
			Headers:     getHeaderMap(ctx),
			Body:        model,
		})
		if err != nil {
			return err
		}
		for header, headerVal := range resp.Headers {
			ctx.Set(header, headerVal)
		}
		return ctx.Status(resp.StatusCode).JSON(resp.Body)
	}
}
