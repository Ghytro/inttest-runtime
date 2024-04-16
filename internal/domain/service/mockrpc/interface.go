package rest

import (
	"context"
	"inttest-runtime/internal/config"
)

type Repository interface {
	GetHttpServiceBehaviorByUrlMethod(ctx context.Context, serviceType config.RpcServiceType, urlPattern, method string) ([]config.HttpHandlerBehaviorItem, error)
}
