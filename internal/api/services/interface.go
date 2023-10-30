package services

import (
	"context"
	"inttest-runtime/internal/domain"
	"inttest-runtime/internal/useCase/services"
)

type UseCase interface {
	GetStatus(ctx context.Context, serviceID domain.ServiceID) (services.ServiceStatusResp, error)
}
