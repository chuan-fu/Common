package jsonx

import (
	"github.com/bytedance/sonic"
	"github.com/chuan-fu/Common/cdefs"
)

func IsArrayNone(data string) bool {
	if data == cdefs.NoneStr || data == cdefs.NoneArrayStr {
		return true
	}
	return false
}

func IsObjNone(data string) bool {
	if data == cdefs.NoneStr || data == cdefs.NoneObjStr {
		return true
	}
	return false
}

// v 应为指针
func Unmarshal(s string, v interface{}) error {
	if s == cdefs.NoneStr {
		return nil
	}
	return sonic.UnmarshalString(s, v)
}

// v 应为指针
func UnmarshalArray(s string, v interface{}) error {
	if IsArrayNone(s) {
		return nil
	}
	return sonic.UnmarshalString(s, v)
}

// v 应为指针
func UnmarshalObj(s string, v interface{}) error {
	if IsObjNone(s) {
		return nil
	}
	return sonic.UnmarshalString(s, v)
}

func Marshal(v interface{}) string {
	if v == nil {
		return cdefs.NoneStr
	}
	data, _ := sonic.MarshalString(v)
	return data
}

func MarshalObj(v interface{}) string {
	if v == nil {
		return cdefs.NoneObjStr
	}
	data, _ := sonic.MarshalString(v)
	return data
}

func MarshalArray(v interface{}) string {
	if v == nil {
		return cdefs.NoneArrayStr
	}
	data, _ := sonic.MarshalString(v)
	return data
}
