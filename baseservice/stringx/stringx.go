package stringx

import (
	"strings"
	"unsafe"
)

// 删除所有前缀
func TrimPrefixAll(s, r string) string {
	for strings.HasPrefix(s, r) {
		s = strings.TrimPrefix(s, r)
	}
	return s
}

// 删除所有后缀
func TrimSuffixAll(s, r string) string {
	for strings.HasSuffix(s, r) {
		s = strings.TrimSuffix(s, r)
	}
	return s
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
