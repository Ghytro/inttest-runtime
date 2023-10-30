package config

import (
	"encoding/json"
	"fmt"
	"inttest-runtime/pkg/utils"
	"math"
)

type Config struct {
	RestServices []RestService `json:"rest_services"`
	GrpcServices []GrpcService `json:"grpc_services"`
	Brokers      []Broker      `json:"brokers"`
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

func (p *Port) UnmarshalJSON(data []byte) error {
	return validatorUnmarshal(data, p)
}

var _ interface {
	json.Unmarshaler
	validator
} = (*Port)(nil)

type RestService struct {
	Port
	Handlers []RestHandler `json:"handlers"`
}

type RestHandler struct {
	ApiPrefix string     `json:"api_prefix"`
	Routes    RestRoutes `json:"routes"`
}

type RestRoutes []RestRoute

type RestRoute struct {
	Route    string              `json:"route"`
	Method   HttpMethod          `json:"method"`
	Headers  HttpHeaders         `json:"headers"`
	Behavior RestHandlerBehavior `json:"behavior"`
}

type RestHandlerBehavior []RestHandlerBehaviorItem

type RestHandlerBehaviorItem struct {
	Params   RestHandlerParams   `json:"parameters"`
	Response RestHandlerResponse `json:"response"`
}

type RestHandlerResponse struct {
	Status  HttpStatus  `json:"status"`
	Headers HttpHeaders `json:"headers"`
	Payload string      `json:"payload"`
}

type GrpcService struct {
	Port
	Protobuf []PbPayload `json:"protobuf"`
}

type PbPayload string

type Broker struct {
	Port
	Type   BrokerType    `json:"type"`
	Topics []BrokerTopic `json:"topics"`
}

type BrokerTopic struct {
	Name     string         `json:"name"`
	Behavior BrokerBehavior `json:"behavior"`
}

type BrokerBehavior string

type HttpMethod string

func (m *HttpMethod) UnmarshalJSON(data []byte) error {
	return validatorUnmarshal(data, m)
}

func (m HttpMethod) Validate() error {
	return validateEnumConst(m)
}

var _ interface {
	json.Unmarshaler
	validator
} = (*HttpMethod)(nil)

type HttpStatus int

func (s *HttpStatus) UnmarshalJSON(data []byte) error {
	return validatorUnmarshal(data, s)
}

func (s HttpStatus) Validate() error {
	return validateEnumConst(s)
}

var _ interface {
	json.Unmarshaler
	validator
} = (*HttpStatus)(nil)

type BrokerType string

func (t BrokerType) Validate() error {
	return validateEnumConst(t)
}

type validator interface {
	Validate() error
}

func validatorUnmarshal[T validator](data []byte, receiver *T) error {
	var temp T
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	if err := temp.Validate(); err != nil {
		return err
	}
	*receiver = temp
	return nil
}

type kvStore interface {
	Get(k string) (string, bool)
	PGet(k string) *string
	MustGet(k string) string
	Contains(k string) bool
	IsNil(k string) bool
	Set(k, v string)
	Unset(k string)
}

type HttpHeaders map[string]string

// Get implements kvStore.
func (h HttpHeaders) Get(header string) (string, bool) {
	result, ok := h[header]
	return result, ok
}

// PGet implements kvStore.
func (h HttpHeaders) PGet(header string) *string {
	result, ok := h[header]
	if !ok {
		return nil
	}
	return utils.ToPtr(result)
}

// MustGet implements kvStore.
func (h HttpHeaders) MustGet(header string) string {
	return h[header]
}

// Contains implements kvStore.
func (h HttpHeaders) Contains(header string) bool {
	_, ok := h[header]
	return ok
}

// IsNIl implements kvStore.
func (h HttpHeaders) IsNil(header string) bool {
	return false
}

// Set implements kvStore.
func (h HttpHeaders) Set(header, val string) {
	h[header] = val
}

// Unset implements kvStore.
func (h HttpHeaders) Unset(header string) {
	delete(h, header)
}

var _ kvStore = (*HttpHeaders)(nil)

type RestHandlerParams map[string]*string

// Contains implements kvStore.
func (p RestHandlerParams) Contains(param string) bool {
	_, ok := p[param]
	return ok
}

// Get implements kvStore.
func (p RestHandlerParams) Get(param string) (string, bool) {
	v, ok := p[param]
	if !ok {
		return "", false
	}
	return *v, true
}

// IsNil implements kvStore.
func (p RestHandlerParams) IsNil(param string) bool {
	v, ok := p[param]
	if !ok {
		return false
	}
	return v == nil
}

// MustGet implements kvStore.
func (p RestHandlerParams) MustGet(param string) string {
	v, ok := p[param]
	if !ok {
		return ""
	}
	return *v
}

// PGet implements kvStore.
func (p RestHandlerParams) PGet(param string) *string {
	return p[param]
}

// Set implements kvStore.
func (p RestHandlerParams) Set(param string, v string) {
	p[param] = utils.ToPtr(v)
}

// Unset implements kvStore.
func (p RestHandlerParams) Unset(param string) {
	delete(p, param)
}

var _ kvStore = (*RestHandlerParams)(nil)
