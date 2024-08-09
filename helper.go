package govm

import (
	"fmt"
	"github.com/Open-Source-CQUT/gover"
	"os"
	"path"
	"path/filepath"
	"runtime"
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

func CheckVersion(v string) (string, bool) {
	if !strings.HasPrefix(v, "go") {
		v = "go" + v
	}
	return v, gover.IsValid(v)
}

func OpenFile(filename string, flag int, perm os.FileMode) (*os.File, error) {
	dir := filepath.Dir(filename)
	if len(dir) == 1 && dir != "." && !os.IsPathSeparator(dir[0]) || len(dir) > 1 {
		err := os.MkdirAll(dir, perm)
		if err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filename, flag, perm)
}

func ErrPrintln(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func Println(a ...any) {
	fmt.Fprintln(os.Stdout, a...)
}

func Printf(format string, a ...any) {
	fmt.Fprintf(os.Stdout, format, a...)
}

var Silence bool

func Tipf(format string, a ...any) {
	if !Silence {
		if !strings.HasSuffix(format, "\n") {
			format = format + "\n"
		}
		Printf(format, a...)
	}
}

func Warnf(format string, a ...any) {
	Tipf("warn: "+format, a...)
}

// UserHomeDir returns the home directory for the current user.
func UserHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS != "windows" {
		// for sudo
		sudoUser, exist := os.LookupEnv("SUDO_USER")
		if exist {
			return filepath.Join("/home", sudoUser), nil
		}
	}

	return homeDir, nil
}

// ToUnixPath converts Windows path to Unix-style path.
func ToUnixPath(p string) string {
	volumeName := filepath.VolumeName(p)
	if volumeName == "" {
		return p
	}
	subPath := strings.ReplaceAll(strings.TrimPrefix(p, volumeName), "\\", "/")
	// volume must be lowercase
	volume := path.Join("/", strings.ToLower(strings.TrimSuffix(volumeName, ":")))
	return path.Join(volume, subPath)
}
