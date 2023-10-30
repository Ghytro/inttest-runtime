package services

import (
	"context"
	"fmt"
	"inttest-runtime/internal/config"
	"inttest-runtime/internal/domain"
	"inttest-runtime/internal/errors/internalErr"
	"inttest-runtime/internal/repository"
	"inttest-runtime/internal/useCase/services/configparser"
	"inttest-runtime/pkg/utils"
	"inttest-runtime/pkg/worker"

	"github.com/pkg/errors"
)

type Service struct {
	repo repository.Transactioner
}

func NewServiceManager(repo repository.Transactioner) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetStatus(ctx context.Context, serviceID domain.ServiceID) (res ServiceStatusResp, err error) {
	var sEnt *domain.RpcService
	err = s.repo.WithLock(func(r repository.IRepository) error {
		sEnt, err = r.GetRpcService(ctx, serviceID)
		return err
	})
	if err != nil {
		return res, internalErr.WrapWithCode(err, internalErr.ErrCodeServiceGet)
	}

	res.ID = sEnt.ID
	res.LaunchedAt = sEnt.CreatedAt
	res.Type = sEnt.Type.String()
	res.CurrentStatus = ServiceStatus{
		Status:    sEnt.CurrentStatus.String(),
		UpdatedAt: sEnt.UpdatedAt,
	}
	return
}

type ServiceInitConfig struct {
	// MaxCompilers how many services may be compiling in parallel
	MaxCompilers int
	MaxCodegens  int
}

func (c ServiceInitConfig) copyApplyDefaults() (res ServiceInitConfig) {
	res.MaxCompilers = c.MaxCompilers
	if c.MaxCompilers == 0 {
		res.MaxCompilers = defaultMaxCompilers
	}

	res.MaxCodegens = c.MaxCodegens
	if c.MaxCodegens == 0 {
		res.MaxCodegens = defaultMaxCodegens
	}
	return res
}

const defaultMaxCompilers = 1
const defaultMaxCodegens = 1

func (s *Service) InitEnv(ctx context.Context, config *config.Config, initConf ServiceInitConfig) error {
	initConf = initConf.copyApplyDefaults()
	worker := worker.NewBoundedWorker(initConf.MaxCompilers)

	// create repository records about the services
	err := s.repo.WithLock(func(r repository.IRepository) error {
		if err := createServicesInRepo(ctx, r, config.RestServices); err != nil {
			return internalErr.WrapWithCode(err, internalErr.ErrCodeRestServiceCreate)
		}
		if err := createServicesInRepo(ctx, r, config.GrpcServices); err != nil {
			return internalErr.WrapWithCode(err, internalErr.ErrCodeGrpcServiceCreate)
		}
		return nil
	})
	if err != nil {
		return internalErr.WrapWithCode(err, internalErr.ErrCodeServiceCreate)
	}

	// codegen rest services
	for _, restService := range config.RestServices {
		for _, route := range restService.Handlers.Routes {
			params, err := configparser.ParseRESTRoute(route)
			if err != nil {
				err = errors.Wrap(err, "error while parsing REST parameters")
				return err
			}

		}
	}
}

type rpcServiceConstraint interface {
	config.RestService | config.GrpcService
}

func createServicesInRepo[T rpcServiceConstraint](
	ctx context.Context,
	repo repository.IRepository,
	services []T,
) error {
	// fixme: проект только начался а уже говнокодим
	switch utils.TypeInstAny[T]().(type) {
	case config.RestService:
		services := utils.SliceTypeAssert[config.RestService](services)
		for _, s := range services {
			if err := repo.CreateRpcService(ctx, &domain.RpcService{
				ID:            domain.ServiceID(s.ID),
				Type:          domain.RpcTypeRestApi,
				RemoteAddr:    fmt.Sprintf(":%d", s.Port),
				CurrentStatus: domain.ServiceUptimeStatus_STARTING,
			}); err != nil {
				return err
			}
		}

	case config.GrpcService:
		services := utils.SliceTypeAssert[config.GrpcService](services)
		for _, s := range services {
			if err := repo.CreateRpcService(ctx, &domain.RpcService{
				ID:            domain.ServiceID(s.ID),
				Type:          domain.RpcTypeGrpc,
				RemoteAddr:    fmt.Sprintf(":%d", s.Port),
				CurrentStatus: domain.ServiceUptimeStatus_STARTING,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func restCodegenJob()
