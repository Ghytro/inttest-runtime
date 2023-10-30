package service_repo

import (
	"context"
	"inttest-runtime/internal/domain"
	"log/slog"
)

type ServiceRepository struct {
	log *slog.Logger
	reg *ServiceRegistry
}

func NewServiceRepository() *ServiceRepository {
	return &ServiceRepository{
		log: slog.With("location", "service_repo"),
		reg: NewServiceRegistry(),
	}
}

func (r *ServiceRepository) CreateRpcService(ctx context.Context, service *domain.RpcService) error {
	return r.reg.WithLock(func(r *ServiceRegistry) error {
		return r.Register(service)
	})
}

func (r *ServiceRepository) GetRpcService(ctx context.Context, id domain.ServiceID) (res *domain.RpcService, err error) {
	err = r.reg.WithLock(func(r *ServiceRegistry) error {
		res, err = r.Get(id)
		return err
	})
	if err != nil {
		return nil, err
	}
	return
}

type IServiceRepository interface {
	CreateRpcService(ctx context.Context, service *domain.RpcService) error
	GetRpcService(ctx context.Context, id domain.ServiceID) (*domain.RpcService, error)
}
