package govm

import (
	"fmt"
	"github.com/Open-Source-CQUT/gover"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unsafe"
)

// Bytes2string convert a byte slice to a string without allocating new memory.
func Bytes2string(bytes []byte) string {
	return unsafe.String(unsafe.SliceData(bytes), len(bytes))
}

// String2bytes convert a string to a byte slice without allocating new memory.
func String2bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func MaxVersion(vs []string) string {
	if len(vs) == 0 {
		return ""
	}
	return slices.MaxFunc(vs, func(v1, v2 string) int {
		return CompareVersion(v1, v2)
	})
}

func CompareVersion(v1, v2 string) int {
	return gover.Compare(v1, v2)
}

func IsValidVersion(v string) bool {
	return gover.IsValid(v)
}

func CheckVersion(v string) (string, bool) {
	if !strings.HasPrefix(v, "go") {
		v = "go" + v
	}
	return v, gover.IsValid(v)
}

func OpenFile(filename string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(filename)
	if len(dir) == 1 && dir != "." && !os.IsPathSeparator(dir[0]) || len(dir) > 1 {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filename, flag, perm)
}

var Silence bool

func Tipf(format string, a ...any) {
	if !Silence {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		fmt.Printf(format, a...)
	}
}
