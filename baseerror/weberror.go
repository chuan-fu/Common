package baseerror

import (
	"fmt"
	"strings"
)

type BaseError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (b *BaseError) Error() string {
	return b.Msg
}

func NewBaseError(code int, msg string) error {
	return &BaseError{Code: code, Msg: msg}
}

func NewBaseCodeError(code int) error {
	return &BaseError{Code: code, Msg: func() string {
		if m, ok := errCodeMsg[code]; ok {
			return m
		}
		return "未知错误"
	}()}
}

func NewNotFoundError(msgs ...string) error {
	return &BaseError{Code: NotFoundErrorCode, Msg: func() string {
		if len(msgs) == 0 {
			return NotFoundErrorMsg
		}
		return fmt.Sprintf("%s:%s", NotFoundErrorMsg, strings.Join(msgs, " "))
	}()}
}
