package mockrpc

import (
	"fmt"
	domainTypes "inttest-runtime/internal/domain/types"
	"inttest-runtime/pkg/xmltree"

	"github.com/gofiber/fiber/v2"
)

type SoapMockApi struct {
	httpRpcApi
	service ISoapApiService
}

func NewSoapMockApi(service ISoapApiService) *SoapMockApi {
	app := fiber.New()
	app.Use(func(ctx *fiber.Ctx) error {
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationXML)
		return ctx.Next()
	})

	return &SoapMockApi{
		httpRpcApi: httpRpcApi{
			app: app,
		},
		service: service,
	}
}

func (api *SoapMockApi) Register(route, method string) error {
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

func (api *SoapMockApi) makeSoapMockHandlerImpl(route, method string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		xmlBytes := ctx.Body()
		model, err := xmltree.FromBytes(xmlBytes)
		if err != nil {
			return err
		}
		resp, err := api.service.HandleSoapRequest(ctx.Context(), route, method, domainTypes.SoapClientRequestParams{
			UrlParams:   ctx.AllParams(),
			QueryParams: ctx.Queries(),
			Headers:     getHeaderMap(ctx),
			Body:        model,
		})
		if err != nil {
			return err
		}
		resultPayload, err := resp.Body.Marshal()
		if err != nil {
			return err
		}
		for header, headerVal := range resp.Headers {
			ctx.Set(header, headerVal)
		}
		return ctx.Status(resp.StatusCode).Send(resultPayload)
	}
}
