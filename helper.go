package govm

import (
	"unsafe"
)

// convert a byte slice to a string without allocating new memory.
func bytes2string(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// convert a string to a byte slice without allocating new memory.
func string2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
