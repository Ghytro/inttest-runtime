package configparser

import (
	"errors"
	"inttest-runtime/internal/config"
	"inttest-runtime/pkg/utils"
	"strings"

	"github.com/samber/lo"
)

type ParamName string

func (n ParamName) String() string {
	return string(n)
}

type ParsedRoute struct {
	PathParams  []ParamName
	QueryParams map[ParamName]string
	Headers     map[ParamName]string
}

func ParseRESTRoute(route config.RestRoute) (res ParsedRoute, err error) {
	res.PathParams, err = findRestUrlPathParams(route.Route)
	if err != nil {
		return res, err
	}
	res.QueryParams, err = findRestPathQueryParams(route.Route)
	if err != nil {
		return res, err
	}
	res.Headers, err = findRestHeaderParams(route.Headers)
	if err != nil {
		return res, err
	}
	return
}

func findRestUrlPathParams(url string) ([]ParamName, error) {
	type intPair = utils.Pair[int, int]

	urlPath := strings.Split(url, "?")[0]
	var paramBraces []intPair
	i := 0
	for i < len(urlPath) {
		if url[i] == '{' {
			j := i + 1
			for j < len(url) && url[j] != '}' {
				j++
			}
			if j == len(url) {
				return nil, errors.New("incorrect parameter definition")
			}
			paramBraces = append(paramBraces, intPair{i, j})
			i = j + 1
			continue
		}
		i++
	}
	return lo.Map(paramBraces, func(p intPair, _ int) ParamName {
		return ParamName(urlPath[p.First+1 : p.Second])
	}), nil
}

func findRestPathQueryParams(url string) (result map[ParamName]string, err error) {
	queryParams := strings.Split(url, "?")[1]
	for _, queryAssignment := range strings.Split(queryParams, "&") {
		eqIdx := strings.Index(queryAssignment, "=")
		if v := queryAssignment[eqIdx+1:]; isBraced(v) {
			result[ParamName(unbraced(v))] = queryAssignment[:eqIdx]
		}
	}
	return result, nil
}

func findRestHeaderParams(headers config.HttpHeaders) (result map[ParamName]string, err error) {
	return lo.MapEntries(
		lo.PickBy(headers, func(_, value string) bool {
			return isBraced(value)
		}),
		func(k, v string) (ParamName, string) {
			return ParamName(unbraced(v)), k
		},
	), nil
}

// todo: json body param parser

func isBraced(s string) bool {
	return s[0] == '{' && s[len(s)-1] == '}'
}

func unbraced(s string) string {
	return s[1 : len(s)-1]
}
