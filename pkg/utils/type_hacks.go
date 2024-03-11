package utils

import (
	"fmt"

	"github.com/samber/lo"
)

func SliceTypeAssert[T any, S ~[]E, E any](s S) []T {
	return any(s).([]T)
}

func TypeInstAny[T any]() any {
	return any(lo.Empty[T]())
}

func ForEachType[S ~[]T, T, TConv any](s S, fn func(item TConv, index int) error) error {
	for i := range s {
		casted, ok := any(s[i]).(TConv)
		if !ok {
			return fmt.Errorf(
				"runtime error: instance of %T is non-castable to %T",
				lo.Empty[T],
				lo.Empty[TConv],
			)
		}
		if err := fn(casted, i); err != nil {
			return err
		}
	}
	return nil
}

func ForEachPtr[S ~[]T, T any](s S, fn func(item *T, index int) error) error {
	for i := range s {
		if err := fn(&s[i], i); err != nil {
			return err
		}
	}
	return nil
}
