package jsonx

import (
	"encoding/json"

	"github.com/chuan-fu/Common/baseservice/stringx"
)

const (
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

func IsObjNone(data string) bool {
	if data == NoneStr || data == NoneObjStr {
		return true
	}
	return false
}

// v 应为指针
func UnmarshalArray(s string, v interface{}) error {
	if IsArrayNone(s) {
		return nil
	}
	return json.Unmarshal(stringx.StringToBytes(s), v)
}

// v 应为指针
func UnmarshalObj(s string, v interface{}) error {
	if IsObjNone(s) {
		return nil
	}
	return json.Unmarshal(stringx.StringToBytes(s), v)
}

func MarshalObj(v interface{}) string {
	if v == nil {
		return NoneObjStr
	}
	data, _ := json.Marshal(v)
	return stringx.BytesToString(data)
}

func MarshalArray(v interface{}) string {
	if v == nil {
		return NoneArrayStr
	}
	data, _ := json.Marshal(v)
	return stringx.BytesToString(data)
}
