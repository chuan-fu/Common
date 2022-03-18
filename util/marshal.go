package util

import (
	"encoding/json"
	"reflect"

	"github.com/pkg/errors"
)

const (
	NoneInt      = 0
	NoneStr      = ""
	NoneArrayStr = "[]"
	NoneObjStr   = "{}"
)

func IsArrayNone(data string) bool {
	if data == NoneStr || data == NoneArrayStr {
		return true
	}
	return false
}

func CommonUnmarshalArray(s string, v interface{}) error {
	if IsArrayNone(s) {
		return nil
	}
	if v == nil || reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.New("参数 v 不是指针类型！")
	}
	return json.Unmarshal(StringToBytes(s), v)
}

func MarshalObj(v interface{}) string {
	if v == nil {
		return NoneObjStr
	}
	data, _ := json.Marshal(v)
	return BytesToString(data)
}

func MarshalArray(v interface{}) string {
	if v == nil {
		return NoneArrayStr
	}
	data, _ := json.Marshal(v)
	return BytesToString(data)
}
