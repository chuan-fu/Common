package util

import (
	"encoding/json"
	"fmt"
	"unsafe"
)

func IsInArray(s int64, array []int64) bool {
	for _, v := range array {
		if s == v {
			return true
		}
	}
	return false
}

func StrInArray(list []string, key string) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Fmt(i interface{}) {
	d, _ := json.Marshal(i)
	fmt.Println(BytesToString(d))
}