package config

import (
	"errors"
	"inttest-runtime/pkg/utils"

	"github.com/samber/lo"
)

func mapPortGetter[T interface{ GetPort() int }](item T, _ int) int {
	return item.GetPort()
}

func (c Config) Validate() error {
	ports := lo.Interleave(
		lo.Map(c.RestServices, mapPortGetter),
		lo.Map(c.GrpcServices, mapPortGetter),
		lo.Map(c.Brokers, mapPortGetter),
	)
	if !utils.IsUniq(ports) {
		return errors.New("services' ports must be unique")
	}
	return nil
}
