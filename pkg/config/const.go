package config

import (
	"fmt"
	"net/http"
	"slices"
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

var allBrokerTypes = []BrokerType{
	"redis",
	"kafka",
}

type enumeratedConst interface {
	HttpMethod | HttpStatus | BrokerType
}

func validateEnumConst[T enumeratedConst](c T) error {
	var (
		enumStringName string
		collection     any
		formattedValue string
	)

	switch t := any(c).(type) {
	case HttpMethod:
		enumStringName = "http method"
		collection = allHttpMethods
		formattedValue = fmt.Sprintf("%q", t)

	case HttpStatus:
		enumStringName = "http status"
		collection = allHttpStatus
		formattedValue = fmt.Sprintf("%d", t)

	case BrokerType:
		enumStringName = "broker type"
		collection = allBrokerTypes
		formattedValue = fmt.Sprintf("%q", t)
	}

	if !slices.Contains(collection.([]T), c) {
		return fmt.Errorf(
			"incorrect value for %s: %s",
			enumStringName,
			formattedValue,
		)
	}

	return nil
}
