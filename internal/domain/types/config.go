package types

import (
	"inttest-runtime/internal/config"

	"github.com/samber/lo"
)

type Config config.Config

func (c Config) RpcServiceBehavsByUrlMethod(urlPattern, method string, serviceTypes ...config.RpcServiceType) ([]config.HttpHandlerBehaviorItem, bool) {
	for _, service := range c.RpcServices {
		if len(serviceTypes) != 0 && !lo.Contains(serviceTypes, service.Type) {
			continue
		}

		route, ok := lo.Find(service.Routes, func(item config.HttpRoute) bool {
			return item.Route.String() == urlPattern && string(item.Method) == method
		})
		if ok {
			return route.Behavior, true
		}
	}
	return nil, false
}
