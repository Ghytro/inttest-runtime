package services

import (
	"inttest-runtime/internal/api"
	"inttest-runtime/internal/domain"
	"inttest-runtime/internal/errors/internalErr"

	"github.com/gofiber/fiber/v2"
)

type API struct {
	service UseCase
}

func NewAPI(service UseCase) *API {
	return &API{
		service: service,
	}
}

func (a *API) MountRouters(app fiber.Router) {
	r := fiber.New()

	r.Get("/service/:serviceId", a.getServiceStatus)

	app.Mount("/monitoring", r)
}

func (a *API) getServiceStatus(ctx *fiber.Ctx) error {
	const location = "error while getting service status"

	serviceID, err := domain.ParseServiceID(ctx.Params("serviceId"))
	if err != nil {
		return api.NewErrResponse(
			err,
			api.ErrReasonBadRequest,
			location,
			api.ErrRespWithCode(internalErr.ErrCodeJSONParsing),
			api.ErrRespWithHttpCode(fiber.StatusBadRequest),
		)
	}

	resp, err := a.service.GetStatus(ctx.Context(), serviceID)
	if err != nil {
		return api.NewErrResponse(
			err,
			api.ErrReasonBadRequest,
			location,
			api.ErrRespWithCode(internalErr.ErrCodeServiceGet),
			api.ErrRespWithHttpCode(fiber.StatusBadRequest),
		)
	}

	return ctx.JSON(resp)
}
