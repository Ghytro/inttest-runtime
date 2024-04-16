package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"inttest-runtime/internal/config"
	"inttest-runtime/pkg/utils"
	"log"
	"maps"
	"reflect"

	"github.com/samber/lo"
)

type HttpClientRequestParams config.HttpStubBehaviorParams

type HttpLogicResponse struct {
	Stub *HttpResp
	Mock *HttpResp
}

type HttpResp config.HttpStubBehaviorResponse

func PerformHttpLogicItem(
	clientReqParams HttpClientRequestParams,
	protocol config.RpcServiceType,
	prioritisedBehaviors []config.HttpHandlerBehaviorItem,
) (*HttpLogicResponse, error) {
	executor, ok := makeLogicExecutor(protocol, clientReqParams)
	if !ok {
		return nil, fmt.Errorf("unknown rpc service type: %s", protocol)
	}
	for _, b := range prioritisedBehaviors {
		switch b.Type {
		case config.RestHandlerBehaviorType_STUB:
			resp, performed, err := executor.performStubLogic(b.HttpHandlerBehaviorUnion.HttpStubBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return &HttpLogicResponse{
				Stub: resp,
			}, nil

		case config.RestHandlerBehaviorType_MOCK:
			resp, performed, err := executor.performMockLogic(b.HttpHandlerBehaviorUnion.HttpMockBehavior)
			if err != nil {
				return nil, err
			}
			if !performed {
				continue
			}
			return &HttpLogicResponse{
				Mock: resp,
			}, nil

		default:
			return nil, fmt.Errorf("unknown behavior type: %s", b.Type)
		}
	}

	return nil, errors.New("behavior was not set, check correctness of config")
}

func makeLogicExecutor(protocol config.RpcServiceType, clientParams HttpClientRequestParams) (iHttpLogicExecutor, bool) {
	switch protocol {
	case config.RpcServiceType_REST:
		return &restLogicExecutor{
			reqParams: clientParams,
		}, true
	case config.RpcServiceType_SOAP:
		return &soapLogicExecutor{
			reqParams: clientParams,
		}, true
	}
	return nil, false
}

type iHttpLogicExecutor interface {
	performStubLogic(behavior config.HttpStubBehavior) (*HttpResp, bool, error)
	performMockLogic(behavior config.HttpMockBehavior) (*HttpResp, bool, error)
}

type (
	restLogicExecutor struct {
		reqParams HttpClientRequestParams
	}

	soapLogicExecutor struct {
		reqParams HttpClientRequestParams
	}
)

// performMockLogic implements iHttpLogicExecutor.
func (r *restLogicExecutor) performMockLogic(behavior config.HttpMockBehavior) (*HttpResp, bool, error) {
	return nil, false, errors.New("python code executor needs implementation")
}

// performStubLogic implements iHttpLogicExecutor.
func (r *restLogicExecutor) performStubLogic(behavior config.HttpStubBehavior) (*HttpResp, bool, error) {
	reqParams := r.reqParams

	// check if all the request parameters fit the behavior logic
	if !maps.Equal(reqParams.Headers, behavior.Params.Headers) ||
		!maps.Equal(reqParams.Query, behavior.Params.Query) ||
		!maps.Equal(reqParams.Url, behavior.Params.Url) {
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
	parseObjWithType := func(bytesJson []byte, obj *any) (t int, err error) {
		if err := json.Unmarshal(bytesJson, obj); err != nil {
			return 0, err
		}
		if _, ok := (*obj).(map[string]any); ok {
			return jsonObjType_OBJECT, nil
		}
		if _, ok := (*obj).([]any); ok {
			return jsonObjType_LIST, nil
		}
		if (*obj) == nil {
			return jsonObjType_NULL, nil
		}
		return jsonObjType_PRIMITIVE, nil
	}
	var clientReqObj, behavReqObj any
	clientT, err := parseObjWithType(utils.S2B(reqParams.Body), &clientReqObj)
	if err != nil {
		return nil, false, err
	}
	behavT, err := parseObjWithType(utils.S2B(behavior.Params.Body), &behavReqObj)
	if err != nil {
		return nil, false, err
	}
	if clientT != behavT {
		return nil, false, err
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
	switch clientT {
	case jsonObjType_OBJECT:
		if anyMapsEqual(clientReqObj, behavReqObj) {
			return (*HttpResp)(lo.ToPtr(behavior.Response)), true, nil
		}
	case jsonObjType_LIST:
		// todo: include order strictness in comparison
		if reflect.DeepEqual(clientReqObj, behavReqObj) {
			return (*HttpResp)(lo.ToPtr(behavior.Response)), true, nil
		}
	case jsonObjType_NULL:
		if clientReqObj == nil && behavReqObj == nil {
			return (*HttpResp)(lo.ToPtr(behavior.Response)), true, nil
		}

	case jsonObjType_PRIMITIVE:
		if clientReqObj == behavReqObj {
			return (*HttpResp)(lo.ToPtr(behavior.Response)), true, nil
		}
	}

	return nil, false, nil
}

// performMockLogic implements iHttpLogicExecutor.
func (s *soapLogicExecutor) performMockLogic(behavior config.HttpMockBehavior) (*HttpResp, bool, error) {
	return nil, false, errors.New("python code executor needs implementation")
}

// performStubLogic implements iHttpLogicExecutor.
func (s *soapLogicExecutor) performStubLogic(behavior config.HttpStubBehavior) (*HttpResp, bool, error) {
	return nil, false, errors.New("soap stub logic needs implementation")
}

var (
	_ iHttpLogicExecutor = (*restLogicExecutor)(nil)
	_ iHttpLogicExecutor = (*soapLogicExecutor)(nil)
)
