package utils

import "github.com/samber/lo"

func SliceTypeAssert[T any, S ~[]E, E any](s S) []T {
	return any(s).([]T)
}

func TypeInstAny[T any]() any {
	return any(lo.Empty[T]())
}
