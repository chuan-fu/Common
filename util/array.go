package util

import (
	"bytes"
	"strconv"
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

func IsInStrArray(key string, list []string) bool {
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

func ConvertToIntArray(arr []string) ([]int, bool) {
	result := make([]int, 0)
	for _, i := range arr {
		res, err := strconv.Atoi(i)
		if err != nil {
			return result, false
		}
		result = append(result, res)
	}
	return result, true
}

func Int64Join(list []int64, sep string) string {
	if len(list) == 0 {
		return ""
	}
	b := bytes.Buffer{}
	for k, v := range list {
		if k > 0 {
			b.WriteString(sep)
		}
		b.WriteString(ToString(v))
	}
	return b.String()
}
