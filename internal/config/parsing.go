package config

import (
	"inttest-runtime/pkg/utils"
	"strings"

	"github.com/samber/lo"
)

func (r ParametrizedRestRoute) Params() []RestRouteParam {
	rS := string(r)
	colonPos := lo.FilterMap(utils.S2B(rS), func(c byte, i int) (int, bool) {
		return i, c == ':'
	})
	return lo.Map(strings.Split(rS, ":")[1:], func(pName string, i int) RestRouteParam {
		before, _, _ := strings.Cut(pName, "/")
		return RestRouteParam{
			Name: before,
			Pos:  colonPos[i],
		}
	})
}
