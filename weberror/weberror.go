package weberror

import (
	"fmt"
	"strings"
)

const (
	NotFoundErrorMsg = "没有对应记录"
)

type BaseWebError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (b *BaseWebError) Error() string {
	return b.Msg
}

func NewBaseWebError(code int, msg string) BaseWebError {
	return BaseWebError{Code: code, Msg: msg}
}

func NewNotFoundError(msgs ...string) (b BaseWebError) {
	b.Code = 300
	b.Msg = NotFoundErrorMsg
	if len(msgs) > 0 {
		b.Msg = fmt.Sprintf("%s:%s", NotFoundErrorMsg, strings.Join(msgs, " "))
	}
	return
}
