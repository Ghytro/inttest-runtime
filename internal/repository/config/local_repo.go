package config

import (
	"context"
	"errors"
	"inttest-runtime/internal/config"

	domainTypes "inttest-runtime/internal/domain/types"

	"github.com/samber/lo"
)

type LocalRepository struct {
	config domainTypes.Config
}

func NewLocalRepository(conf config.Config) *LocalRepository {
	return &LocalRepository{
		config: domainTypes.Config(conf),
	}
}

func (r LocalRepository) Get(ctx context.Context) (*domainTypes.Config, error) {
	return lo.ToPtr(r.config), nil
}

func (r LocalRepository) GetHttpServiceBehaviorByUrlMethod(ctx context.Context, serviceType config.RpcServiceType, urlPattern, method string) ([]config.HttpHandlerBehaviorItem, error) {
	result, ok := r.config.RpcServiceBehavsByUrlMethod(urlPattern, method, serviceType)
	if !ok {
		return nil, errors.New("behavior undefined")
	}
	return result, nil
}
