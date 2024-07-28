package govm

import (
	"fmt"
	"os"
	"path/filepath"
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

// compare two version strings.
// eg.
//
//	compareVersion(1.0.0,1.0.1) = -1
//	compareVersion(1.2.0,1.0.1) = 1
//	compareVersion(1.1.0,1.1.0) = 0
func compareVersion(version1 string, version2 string) int {
	var i, j int
	for i < len(version1) || j < len(version2) {
		var a, b int
		for ; i < len(version1) && version1[i] != '.'; i++ {
			a = a*10 + int(version1[i]-'0')
		}
		for ; j < len(version2) && version2[j] != '.'; j++ {
			b = b*10 + int(version2[j]-'0')
		}
		if a > b {
			return 1
		} else if a < b {
			return -1
		}

		i++
		j++
	}
	return 0
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
		fmt.Printf(format+"\n", a...)
	}
}
