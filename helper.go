package govm

import (
	"fmt"
	"go/version"
	"os"
	"path/filepath"
	"strings"
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

func CompareVersion(v1, v2 string) int {
	return version.Compare(v1, v2)
}
func IsValidVersion(v string) bool {
	return version.IsValid(v)
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

func Println(a ...any) {
	if !Silence {
		fmt.Println(a...)
	}
}

func Printf(format string, a ...any) {
	if !Silence {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		fmt.Printf(format, a...)
	}
}
