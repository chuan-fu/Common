package errorx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chuan-fu/Common/cdefs"
)

const (
	wrapP  = "%s: %v"
	errorP = "[code: %d]: %v"
)

type ErrorX interface {
	Error() string
	Code() int64
}

type errorX struct {
	code int64
	error
}

func New(s string) error {
	return &errorX{error: errors.New(s)}
}

func Newf(format string, args ...interface{}) error {
	return &errorX{error: errors.New(fmt.Sprintf(format, args...))}
}

func News(s ...string) error {
	switch len(s) {
	case 0:
		return nil
	case 1:
		return &errorX{error: errors.New(s[0])}
	default:
		return &errorX{error: errors.New(strings.Join(s, cdefs.HalfWidthComma))}
	}
}

func WithCode(code int64, err error) error {
	return &errorX{code: code, error: err}
}

func Wrap(err error, key string) error {
	return &errorX{error: &withMessage{err, key}}
}

func Wrapf(err error, format string, args ...interface{}) error {
	return &errorX{error: &withMessage{err, fmt.Sprintf(format, args...)}}
}

func (b *errorX) Error() string {
	if b.code == 0 {
		return b.error.Error()
	}
	return fmt.Sprintf(errorP, b.code, b.error)
}

func (b *errorX) Code() int64 {
	return b.code
}

func ToErrorX(err error) ErrorX {
	if e, ok := err.(*errorX); ok {
		return e
	}
	return &errorX{error: err}
}

type withMessage struct {
	err error
	key string
}

func (w *withMessage) Error() string {
	return fmt.Sprintf(wrapP, w.key, w.err)
}
