package utils

import (
	"unsafe"
)

type Pair[T, E any] struct {
	First  T
	Second E
}

func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
