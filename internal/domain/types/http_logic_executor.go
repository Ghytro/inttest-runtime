package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"inttest-runtime/internal/config"
	"inttest-runtime/pkg/utils"
	"inttest-runtime/pkg/xmltree"
	"log"
	"maps"
	"reflect"
)

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

func PerformRestLogic(
	clientReqParams RestClientRequestParams,
	prioritisedBehaviors []config.HttpHandlerBehaviorItem,
) (*RestLogicResponse, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть
	for _, b := range prioritisedBehaviors {
		switch b.Type {
		case config.RestHandlerBehaviorType_STUB:
			resp, performed, err := restStubLogic(clientReqParams, b.HttpHandlerBehaviorUnion.HttpStubBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return resp, nil

		case config.RestHandlerBehaviorType_MOCK:
			resp, performed, err := restMockLogic(clientReqParams, b.HttpHandlerBehaviorUnion.HttpMockBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return resp, nil

		default:
			return nil, fmt.Errorf("unknown behavior type: %s", b.Type)
		}
	}

	return nil, errors.New("behavior was not set, check correctness of config")
}

func restStubLogic(reqParams RestClientRequestParams, behavior config.HttpStubBehavior) (*RestLogicResponse, bool, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть

	// check if all the request parameters fit the behavior logic
	if !maps.Equal(reqParams.Headers, behavior.Params.Headers) ||
		!maps.Equal(reqParams.QueryParams, behavior.Params.Query) ||
		!maps.Equal(reqParams.UrlParams, behavior.Params.Url) {
		// that just means that the behavior is wrong
		// and we need to try another one
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

func restMockLogic(reqParams RestClientRequestParams, behavior config.HttpMockBehavior) (*RestLogicResponse, bool, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть
	return nil, false, errors.New("python code executor not implemented")
}

func PerformSoapLogic(reqParams SoapClientRequestParams, behavior []config.HttpHandlerBehaviorItem) (*SoapLogicResponse, error) {
	// todo: зависимость доменной логики на пакет с конфигом выглядит плохо, переделать если время есть
	for _, b := range behavior {
		switch b.Type {
		case config.RestHandlerBehaviorType_STUB:
			resp, performed, err := soapStubLogic(reqParams, b.HttpHandlerBehaviorUnion.HttpStubBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return resp, nil
		case config.RestHandlerBehaviorType_MOCK:
			resp, performed, err := soapMockLogic(reqParams, b.HttpHandlerBehaviorUnion.HttpMockBehavior)
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

func soapStubLogic(reqParams SoapClientRequestParams, behavior config.HttpStubBehavior) (*SoapLogicResponse, bool, error) {
	// check if all the request parameters fit the behavior logic
	if !maps.Equal(reqParams.Headers, behavior.Params.Headers) ||
		!maps.Equal(reqParams.QueryParams, behavior.Params.Query) ||
		!maps.Equal(reqParams.UrlParams, behavior.Params.Url) {
		// that just means that the behavior is wrong
		// and we need to try another one
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

func soapMockLogic(reqParams SoapClientRequestParams, behavior config.HttpMockBehavior) (*SoapLogicResponse, bool, error) {
	return nil, false, errors.New("python code executor must be implemented")
}
