package config

import (
	"fmt"
	"inttest-runtime/pkg/utils"
	"regexp"
	"strings"
	"time"

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

var strTimeIntervalRegex = regexp.MustCompile(`^\d+(s|ms|ns)$`)

func (i StrTimeInterval) IsValid() bool {
	return strTimeIntervalRegex.MatchString(string(i))
}

func (i StrTimeInterval) ToDuration() time.Duration {
	s := strings.ToLower(string(i))
	modulus, strMult := 0, ""
	fmt.Sscanf(s, "%d%s", &modulus, &strMult)
	multiplier := lo.Switch[string, time.Duration](strMult).
		Case("s", time.Second).
		Case("ms", time.Millisecond).
		Case("ns", time.Nanosecond).
		Default(time.Second)

	return multiplier * time.Duration(modulus)
}
