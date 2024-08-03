package errorx

import (
	"errors"
	"fmt"
)

const (
	ErrorKind = "error"
	WarnKind  = "warn"
)

type KindError struct {
	Kind string
	Err  error
}

func (e KindError) Error() string {
	return e.Kind + ": " + e.Err.Error()
}

func Warn(s string) error {
	return KindError{Kind: WarnKind, Err: errors.New(s)}
}

func Warnf(s string, a ...any) error {
	return KindError{Kind: WarnKind, Err: fmt.Errorf(s, a...)}
}

func Error(s string) error {
	return KindError{Kind: ErrorKind, Err: errors.New(s)}
}

func Errorf(s string, a ...any) error {
	return KindError{Kind: ErrorKind, Err: fmt.Errorf(s, a...)}
}
