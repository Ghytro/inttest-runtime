package internalErr

import (
	"errors"
	"fmt"
)

type errWithCode struct {
	err  error
	code ErrCode
}

func (err errWithCode) Error() string {
	return fmt.Sprintf("[%04d] %s", err.code, err.err.Error())
}

func (err errWithCode) Unwrap() error {
	return err.err
}

func WrapWithCode(err error, code ErrCode) error {
	var t *errWithCode
	if errors.As(err, &t) {
		return err
	}
	return &errWithCode{
		err:  err,
		code: code,
	}
}

func Explain(err error) (string, bool) {
	var (
		result string
		bRes   bool
	)
	ok := performIfInternal(err, func(err *errWithCode) {
		if err.code < 1 || int(err.code) > len(explains) {
			return
		}
		result = explains[err.code]
		bRes = true
	})
	return result, ok && bRes
}

func GetCode(err error) (ErrCode, bool) {
	var result ErrCode
	ok := performIfInternal(err, func(err *errWithCode) {
		result = err.code
	})
	return result, ok
}

func performIfInternal(err error, fn func(err *errWithCode)) bool {
	var t *errWithCode
	if errors.As(err, &t) {
		fn(t)
		return true
	}
	return false
}
