package utils

import (
	"cmp"
	"slices"
)

func ToPtr[T any](val T) *T {
	result := new(T)
	*result = val
	return result
}

func IsUniq[S ~[]E, E cmp.Ordered](s S) bool {
	cp := slices.Clone(s)
	slices.Sort(cp)
	return len(slices.Compact(cp)) == len(s)
}
