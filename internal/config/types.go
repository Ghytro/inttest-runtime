package config

import (
	"fmt"
	"math"

	"github.com/samber/lo"
)

type Config struct {
	Port
	RpcServices []RpcService `json:"rpc_services"`
	Brokers     []Broker     `json:"brokers"`
}

type Port struct {
	Port int `json:"port"`
}

func (p Port) GetPort() int {
	return p.Port
}

func (p Port) Validate() error {
	if p.Port <= 0 || p.Port > math.MaxUint16 {
		return fmt.Errorf("incorrect port value: %d", p.Port)
	}
	return nil
}

type RpcService struct {
	RpcServiceCommon
	RpcServiceUnion
}

type RpcServiceCommon struct {
	Type RpcServiceType `json:"type"`
	ID   string         `json:"id"`
	Port
}

type RpcServiceType string

func (t RpcServiceType) Validate() error {
	if !lo.Contains(allRpcServiceTypes, t) {
		return fmt.Errorf("некорректное значение типа rpc-сервиса: %s", string(t))
	}
	return nil
}

type RpcServiceUnion struct {
	HttpService
}

type HttpService struct {
	Routes []HttpRoute `json:"routes"`
}

type HttpRoute struct {
	Route    ParametrizedRestRoute     `json:"route"`
	Method   HttpMethod                `json:"method"`
	Behavior []HttpHandlerBehaviorItem `json:"behavior"`
}

type ParametrizedRestRoute string

func (r ParametrizedRestRoute) String() string {
	return string(r)
}

type RestRouteParam struct {
	Name string
	Pos  int
}

type HttpMethod string

func (m HttpMethod) Validate() error {
	if !lo.Contains(allHttpMethods, m) {
		return fmt.Errorf("некорректное значение http-метода (%s)", string(m))
	}
	return nil
}

type HttpHandlerBehaviorItem struct {
	Type HttpHandlerBehaviorType
	HttpHandlerBehaviorUnion
}

type HttpHandlerBehaviorUnion struct {
	HttpStubBehavior
	HttpMockBehavior
}

type HttpHandlerBehaviorType string

func (bt HttpHandlerBehaviorType) Validate() error {
	if !lo.Contains(allRestHandlerBehaviorTypes, bt) {
		return fmt.Errorf("некорректное значение типа поведения rest-хендлера: %s", string(bt))
	}
	return nil
}

type HttpStubBehavior struct {
	Params   HttpStubBehaviorParams   `json:"parameters"`
	Response HttpStubBehaviorResponse `json:"response"`
}

type HttpStubBehaviorParams struct {
	Headers map[string]string `json:"headers"`
	Query   map[string]string `json:"query"`
	Body    string            `json:"body"`
	Url     map[string]string `json:"url"`
}

type HttpStubBehaviorResponse struct {
	Status  HttpStatus        `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type HttpStatus int

func (s HttpStatus) Validate() error {
	if !lo.Contains(allHttpStatus, s) {
		return fmt.Errorf("некорректное значение HTTP статуса (%d)", s)
	}
	return nil
}

type HttpMockBehavior struct {
	// скорее всего можно будет взять из пакета
	// тип вроде code-snippet
	// (или все таки отделить мух от котлет?)
	Impl []string `json:"impl"`
}

type Broker struct {
	ID   string     `json:"id"`
	Type BrokerType `json:"type"`
	Port
	BrokerBehaviorUnion
}

type BrokerType string

func (t BrokerType) Validate() error {
	if !lo.Contains(allBrokerTypes, t) {
		return fmt.Errorf("некорректное значение типа брокера: %s", string(t))
	}
	return nil
}

type BrokerBehaviorUnion struct {
	BrokerBehaviorRedis
}

type BrokerBehaviorRedis struct {
	Behavior []BrokerBehaviorRedisItem `json:"behavior"`
}

type BrokerBehaviorRedisItem struct {
	Topic      string                `json:"topic"`
	Generators []RedisTopicGenerator `json:"generators"`
}

type RedisTopicGenerator struct {
	Interval string                  `json:"interval"`
	Type     RedisTopicGeneratorType `json:"type"`
	RedisTopicGeneratorUnion
}

type RedisTopicGeneratorUnion struct {
	Const *RedisTopicGeneratorConst
	Prog  *RedisTopicGeneratorProg
}

type RedisTopicGeneratorConst struct {
	Payload string `json:"payload"`
}

type RedisTopicGeneratorProg struct {
	Behavior []string `json:"behavior"`
}

type RedisTopicGeneratorType string
