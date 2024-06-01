package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"inttest-runtime/internal/config"
	"inttest-runtime/pkg/embedded"
	"inttest-runtime/pkg/utils"
	"inttest-runtime/pkg/xmltree"
	"log"
	"maps"
	"reflect"

	"github.com/samber/lo"
)

type MockLogicExecutor struct {
	pyExecutor *PyPrecompiledExecutor
}

func NewMockLogicExecutor(exec *PyPrecompiledExecutor) *MockLogicExecutor {
	return &MockLogicExecutor{
		pyExecutor: exec,
	}
}

type RestClientRequestParams struct {
	UrlParams   map[string]string
	QueryParams map[string]string
	Headers     map[string]string
	Body        any // parsed json either map[string]any or []any
}

type SoapClientRequestParams struct {
	UrlParams   map[string]string
	QueryParams map[string]string
	Headers     map[string]string
	Body        *xmltree.Node
}

type RestLogicResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       any
}

type SoapLogicResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       *xmltree.Node
}

func (exec *MockLogicExecutor) PerformRestLogic(
	clientReqParams RestClientRequestParams,
	prioritisedBehaviors []config.HttpHandlerBehaviorItem,
) (*RestLogicResponse, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть
	for _, b := range prioritisedBehaviors {
		var (
			resp      *RestLogicResponse
			performed bool
			err       error
		)
		switch b.Type {
		case config.RestHandlerBehaviorType_STUB:
			resp, performed, err = exec.restStubLogic(clientReqParams, b.HttpHandlerBehaviorUnion.HttpStubBehavior)

		case config.RestHandlerBehaviorType_MOCK:
			resp, performed, err = exec.restMockLogic(clientReqParams, b.HttpHandlerBehaviorUnion.HttpMockBehavior)

		default:
			return nil, fmt.Errorf("unknown behavior type: %s", b.Type)
		}

		if err != nil {
			return nil, err
		}
		if !performed {
			continue
		}
		resp.Body, err = toJsonSerializeable(resp.Body)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	return nil, errors.New("behavior was not set, check correctness of config")
}

func toJsonSerializeable(v any) (any, error) {
	switch t := v.(type) {
	case map[any]any:
		result := make(map[string]any, len(t))
		for mk, mv := range t {
			kStr, ok := mk.(string)
			if !ok {
				return nil, errors.New("all json keys must be strings")
			}
			resVal, err := toJsonSerializeable(mv)
			if err != nil {
				return nil, err
			}
			result[kStr] = resVal
		}
		return result, nil

	case []any:
		result := make([]any, 0, len(t))
		for _, sv := range t {
			resItem, err := toJsonSerializeable(sv)
			if err != nil {
				return nil, err
			}
			result = append(result, resItem)
		}
		return result, nil
	}

	return v, nil
}

func clientParamsFitBehavior(url, query, header map[string]string, behavParams config.HttpStubBehaviorParams) bool {
	mapIncludes := func(super, sub map[string]string) bool {
		for kSub, vSub := range sub {
			if vSuper, ok := super[kSub]; !ok || vSuper != vSub {
				return false
			}
		}
		return true
	}

	return mapIncludes(url, behavParams.Url) &&
		mapIncludes(query, behavParams.Query) &&
		mapIncludes(header, behavParams.Headers)
}

func (exec *MockLogicExecutor) restStubLogic(reqParams RestClientRequestParams, behavior config.HttpStubBehavior) (*RestLogicResponse, bool, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть

	// check if all the request parameters fit the behavior logic
	if !clientParamsFitBehavior(
		reqParams.UrlParams,
		reqParams.QueryParams,
		reqParams.Headers,
		behavior.Params,
	) {
		return nil, false, nil
	}

	// parse json body and compare to behavior
	// client body fits the behavior body if it has similar type (either list or object)
	// and equals according to the rules:
	// 1) list: (todo think about the order strictness rule) all the objects must be equal
	// 2) object: all the keys present in client's request should be present in behavior body
	const (
		jsonObjType_OBJECT = iota + 1
		jsonObjType_LIST
		jsonObjType_NULL
		jsonObjType_PRIMITIVE
	)
	getObjectType := func(obj any) int {
		if _, ok := obj.(map[string]any); ok {
			return jsonObjType_OBJECT
		}
		if _, ok := obj.([]any); ok {
			return jsonObjType_LIST
		}
		if obj == nil {
			return jsonObjType_NULL
		}
		return jsonObjType_PRIMITIVE
	}
	clientT := getObjectType(reqParams.Body)
	var behavObj any
	if err := json.Unmarshal(utils.S2B(behavior.Params.Body), &behavObj); err != nil {
		return nil, false, err
	}
	behavT := getObjectType(behavObj)
	if clientT != behavT {
		// type mismatch means objects are not equal
		return nil, false, nil
	}
	anyMapsEqual := func(obj1, obj2 any) bool {
		m1, ok := obj1.(map[string]any)
		if !ok {
			log.Println("runtime error: unable to cast to map in rest stub comparison")
			return false
		}
		m2, ok := obj2.(map[string]any)
		if !ok {
			log.Println("runtime error: unable to cast to map in rest stub comparison")
			return false
		}
		return reflect.DeepEqual(m1, m2)
	}
	var respBody any
	if err := json.Unmarshal(utils.S2B(behavior.Response.Body), &respBody); err != nil {
		return nil, false, err
	}
	domainResponse := &RestLogicResponse{
		StatusCode: int(behavior.Response.Status),
		Headers:    maps.Clone(behavior.Response.Headers),
		Body:       respBody,
	}
	switch clientT {
	case jsonObjType_OBJECT:
		if anyMapsEqual(reqParams.Body, behavObj) {
			return domainResponse, true, nil
		}
	case jsonObjType_LIST:
		// todo: include order strictness in comparison
		if reflect.DeepEqual(reqParams.Body, behavObj) {
			return domainResponse, true, nil
		}
	case jsonObjType_NULL:
		if reqParams.Body == nil && behavObj == nil {
			return domainResponse, true, nil
		}

	case jsonObjType_PRIMITIVE:
		if reqParams.Body == behavObj {
			return domainResponse, true, nil
		}
	}

	return nil, false, nil
}

func (exec *MockLogicExecutor) restMockLogic(reqParams RestClientRequestParams, behavior config.HttpMockBehavior) (*RestLogicResponse, bool, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть

	// construct py arguments before calling py func
	// немного не нравится доступ до интерпретатора
	args, err := NewRestMockFuncArgBuilder(exec.pyExecutor.pyCtx).
		SetHeaders(reqParams.Headers).
		SetQueryParams(reqParams.QueryParams).
		SetUrlParams(reqParams.UrlParams).
		SetBody(reqParams.Body).
		Build()

	if err != nil {
		return nil, false, err
	}

	result, err := exec.pyExecutor.ExecFunc(behavior.Impl, args...)
	if err != nil {
		return nil, false, err
	}
	if result.IsNone() {
		return nil, false, nil
	}

	var jsonRes any
	if mapRes, err := result.ToMap(); err == nil {
		var ok bool
		jsonRes, ok = goPyDict(mapRes).toJsonDict()
		if !ok {
			return nil, false, errors.New("python error: result dict is not json serializeable")
		}
	}
	if sliceRes, err := result.ToSlice(); err == nil {
		jsonResSlice := make([]any, 0, len(sliceRes))
		for _, item := range sliceRes {
			item := item
			if mapItem, ok := item.(map[any]any); ok {
				item, ok = goPyDict(mapItem).toJsonDict()
				if !ok {
					return nil, false, errors.New("python error: result dict is not json serializeable")
				}
			}
			jsonResSlice = append(jsonResSlice, item)
		}
		jsonRes = jsonResSlice
	}
	if primitiveRes, err := result.ToAny(); err == nil {
		jsonRes = primitiveRes
	}

	return &RestLogicResponse{
		StatusCode: 200, // todo
		Headers:    nil, // todo
		Body:       jsonRes,
	}, true, nil
}

type goPyDict map[any]any

func (d goPyDict) toJsonDict() (map[string]any, bool) {
	toStrMapIfTypeMatch := func(v any) (result any, ok bool) {
		if vMap, ok := v.(map[any]any); ok {
			result, ok = goPyDict(vMap).toJsonDict()
			if !ok {
				return nil, false
			}
			return result, true
		}
		return v, true
	}

	result := make(map[string]any, len(d))
	for k, v := range d {
		kStr, ok := k.(string)
		if !ok {
			return nil, false
		}
		var newVal any = v
		newVal, ok = toStrMapIfTypeMatch(v)
		if !ok {
			return nil, false
		}
		if vSlice, ok := v.([]any); ok {
			newSlice := make([]any, 0, len(vSlice))
			for _, item := range vSlice {
				newItem, ok := toStrMapIfTypeMatch(item)
				if !ok {
					return nil, false
				}
				newSlice = append(newSlice, newItem)
			}
			newVal = newSlice
		}

		result[kStr] = newVal
	}
	return result, true
}

func (exec *MockLogicExecutor) PerformSoapLogic(reqParams SoapClientRequestParams, behavior []config.HttpHandlerBehaviorItem) (*SoapLogicResponse, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть
	for _, b := range behavior {
		switch b.Type {
		case config.RestHandlerBehaviorType_STUB:
			resp, performed, err := exec.soapStubLogic(reqParams, b.HttpHandlerBehaviorUnion.HttpStubBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return resp, nil
		case config.RestHandlerBehaviorType_MOCK:
			resp, performed, err := exec.soapMockLogic(reqParams, b.HttpHandlerBehaviorUnion.HttpMockBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return resp, nil
		}
	}

	return nil, errors.New("no response")
}

func (exec *MockLogicExecutor) soapStubLogic(reqParams SoapClientRequestParams, behavior config.HttpStubBehavior) (*SoapLogicResponse, bool, error) {
	// check if all the request parameters fit the behavior logic
	if !clientParamsFitBehavior(
		reqParams.UrlParams,
		reqParams.QueryParams,
		reqParams.Headers,
		behavior.Params,
	) {
		return nil, false, nil
	}

	// check if body fits the behavior
	behavReqBody, err := xmltree.FromBytes(utils.S2B(behavior.Params.Body))
	if err != nil {
		return nil, false, err
	}
	if !reqParams.Body.Equal(behavReqBody) {
		return nil, false, nil
	}

	// if the body fits serialize response
	behavRespBody, err := xmltree.FromBytes(utils.S2B(behavior.Response.Body))
	if err != nil {
		return nil, false, err
	}
	return &SoapLogicResponse{
		StatusCode: int(behavior.Response.Status),
		Headers:    maps.Clone(behavior.Response.Headers),
		Body:       behavRespBody,
	}, true, nil
}

func (exec *MockLogicExecutor) soapMockLogic(reqParams SoapClientRequestParams, behavior config.HttpMockBehavior) (*SoapLogicResponse, bool, error) {
	return nil, false, errors.New("python code executor must be implemented")
}

type RestMockFuncArgBuilder struct {
	pyCtx       *embedded.PyRuntime
	urlParams   map[any]any
	headers     map[any]any
	queryParams map[any]any
	body        any
}

func NewRestMockFuncArgBuilder(interpreter *embedded.PyRuntime) *RestMockFuncArgBuilder {
	return &RestMockFuncArgBuilder{
		pyCtx: interpreter,
	}
}

func (b *RestMockFuncArgBuilder) SetUrlParams(params map[string]string) *RestMockFuncArgBuilder {
	b.urlParams = lo.MapEntries(params, func(k string, v string) (any, any) { return k, v })
	return b
}

func (b *RestMockFuncArgBuilder) SetHeaders(headers map[string]string) *RestMockFuncArgBuilder {
	b.headers = lo.MapEntries(headers, func(k string, v string) (any, any) { return k, v })
	return b
}

func (b *RestMockFuncArgBuilder) SetQueryParams(params map[string]string) *RestMockFuncArgBuilder {
	b.queryParams = lo.MapEntries(params, func(k string, v string) (any, any) { return k, v })
	return b
}

func (b *RestMockFuncArgBuilder) SetBody(body any) *RestMockFuncArgBuilder {
	b.body = body
	return b
}

func (b *RestMockFuncArgBuilder) Build() (result []embedded.PyValue, err error) {
	headers, err := b.pyCtx.NewDict(b.headers)
	if err != nil {
		return nil, err
	}
	result = append(result, headers)

	query, err := b.pyCtx.NewDict(b.queryParams)
	if err != nil {
		return nil, err
	}
	result = append(result, query)

	urlParams, err := b.pyCtx.NewDict(b.urlParams)
	if err != nil {
		return nil, err
	}
	result = append(result, urlParams)

	var body embedded.PyValue
	switch t := b.body.(type) {
	case []any:
		body, err = b.pyCtx.NewList(t)
	case map[string]any:
		body, err = b.pyCtx.NewDict(lo.MapKeys(t, func(_ any, k string) any { return k }))
	default:
		body, err = b.pyCtx.NewPrimitive(t)
	}
	if err != nil {
		return nil, err
	}

	return []embedded.PyValue{
		urlParams, query, headers, body,
	}, nil
}
