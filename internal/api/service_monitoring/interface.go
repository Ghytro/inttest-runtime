package services

import (
	"context"
	"inttest-runtime/internal/domain"
	"inttest-runtime/internal/usecase/services"
)

type UseCase interface {
	GetStatus(ctx context.Context, serviceID domain.ServiceID) (services.ServiceStatusResp, error)
}
