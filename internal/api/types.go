package api

import (
	"inttest-runtime/internal/errors/httperr"
	"inttest-runtime/internal/errors/internalErr"
)

type errResponse struct {
	err      error               `json:"-"`
	httpCode int                 `json:"-"`
	Code     internalErr.ErrCode `json:"code"`
	Msg      string              `json:"message"`
	Reason   errReason           `json:"reason"`
	Location string              `json:"location"`
}

func (err errResponse) Error() string {
	return err.err.Error()
}

type errResponseOption func(resp *errResponse)

func ErrRespWithCode(code internalErr.ErrCode) errResponseOption {
	return func(resp *errResponse) {
		resp.Code = code
	}
}

func ErrRespWithHttpCode(code int) errResponseOption {
	return func(resp *errResponse) {
		resp.httpCode = code
	}
}

func NewErrResponse(err error, reason errReason, location string, opts ...errResponseOption) *errResponse {
	result := errResponse{
		Msg:      err.Error(),
		Reason:   reason,
		Location: location,
	}
	for _, o := range opts {
		o(&result)
	}
	if code, ok := internalErr.GetCode(err); ok {
		result.Code = code
	}
	if code, ok := httperr.GetHTTPCode(err); ok {
		result.httpCode = code
	}
	return &result
}

type errReason string

const (
	ErrReasonBadRequest          errReason = "bad_request"
	ErrReasonInternalServerError errReason = "internal_server_error"
)
