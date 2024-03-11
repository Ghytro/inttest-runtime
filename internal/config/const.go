package config

import (
	"net/http"
)

// preallocated slices with constants

var allHttpMethods = []HttpMethod{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

var allHttpStatus = []HttpStatus{
	http.StatusContinue,
	http.StatusSwitchingProtocols,
	http.StatusProcessing,
	http.StatusEarlyHints,

	http.StatusOK,
	http.StatusCreated,
	http.StatusAccepted,
	http.StatusNonAuthoritativeInfo,
	http.StatusNoContent,
	http.StatusResetContent,
	http.StatusPartialContent,
	http.StatusMultiStatus,
	http.StatusAlreadyReported,
	http.StatusIMUsed,

	http.StatusMultipleChoices,
	http.StatusMovedPermanently,
	http.StatusFound,
	http.StatusSeeOther,
	http.StatusNotModified,
	http.StatusUseProxy,

	http.StatusTemporaryRedirect,
	http.StatusPermanentRedirect,

	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusPaymentRequired,
	http.StatusForbidden,
	http.StatusNotFound,
	http.StatusMethodNotAllowed,
	http.StatusNotAcceptable,
	http.StatusProxyAuthRequired,
	http.StatusRequestTimeout,
	http.StatusConflict,
	http.StatusGone,
	http.StatusLengthRequired,
	http.StatusPreconditionFailed,
	http.StatusRequestEntityTooLarge,
	http.StatusRequestURITooLong,
	http.StatusUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed,
	http.StatusTeapot,
	http.StatusMisdirectedRequest,
	http.StatusUnprocessableEntity,
	http.StatusLocked,
	http.StatusFailedDependency,
	http.StatusTooEarly,
	http.StatusUpgradeRequired,
	http.StatusPreconditionRequired,
	http.StatusTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons,

	http.StatusInternalServerError,
	http.StatusNotImplemented,
	http.StatusBadGateway,
	http.StatusServiceUnavailable,
	http.StatusGatewayTimeout,
	http.StatusHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates,
	http.StatusInsufficientStorage,
	http.StatusLoopDetected,
	http.StatusNotExtended,
	http.StatusNetworkAuthenticationRequired,
}

const (
	BrokerType_REDIS_PUBSUB BrokerType = "redis-pubsub"
)

var allBrokerTypes = []BrokerType{
	BrokerType_REDIS_PUBSUB,
}

const (
	RedisTopicGeneratorType_CONST RedisTopicGeneratorType = "constant"
	RedisTopicGeneratorType_PROG  RedisTopicGeneratorType = "programmable"
)

var allRedisTopicGeneratorTypes = []RedisTopicGeneratorType{
	RedisTopicGeneratorType_CONST,
	RedisTopicGeneratorType_PROG,
}

const (
	RestHandlerBehaviorType_STUB RestHandlerBehaviorType = "stub"
	RestHandlerBehaviorType_MOCK RestHandlerBehaviorType = "mock"
)

var allRestHandlerBehaviorTypes = []RestHandlerBehaviorType{
	RestHandlerBehaviorType_STUB,
	RestHandlerBehaviorType_MOCK,
}

const (
	RpcServiceType_REST RpcServiceType = "rest"
	RpcServiceType_SOAP RpcServiceType = "soap"
)

var allRpcServiceTypes = []RpcServiceType{
	RpcServiceType_REST,
	RpcServiceType_SOAP,
}
