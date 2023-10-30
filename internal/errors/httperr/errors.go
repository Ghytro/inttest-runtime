package httperr

import (
	"errors"
	"fmt"
)

type HTTPCodeError struct {
	err  error
	code int
}

func (err *HTTPCodeError) Error() string {
	return fmt.Sprintf("%s (%d)", err.err.Error(), err.code)
}

func (err *HTTPCodeError) Unwrap() error {
	return err.err
}

func WithHTTPCode(err error, code int) error {
	var errT *HTTPCodeError
	if errors.As(err, &errT) {
		return errT
	}
	return &HTTPCodeError{
		err:  err,
		code: code,
	}
}

func GetHTTPCode(err error) (int, bool) {
	var _err *HTTPCodeError
	if errors.As(err, &_err) {
		return _err.code, true
	}
	return 0, false
}
