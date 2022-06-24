package stringx

import (
	"strings"
	"unicode/utf8"
	"unsafe"
)

// 删除所有前缀
func TrimPrefixAll(s, prefix string) string {
	for strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}
	return s
}

// 删除所有后缀
func TrimSuffixAll(s, suffix string) string {
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
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

// 统计
func CountByte(s string, b byte) (n int) {
	for k := range s {
		if s[k] == b {
			n++
		}
	}
	return
}

// 清除s字符串中所有b
func TrimByte(s string, b byte) string {
	n := CountByte(s, b)
	if n == 0 {
		return s
	}
	var bs strings.Builder
	bs.Grow(len(s) - n)
	for k := range s {
		if s[k] != b {
			bs.WriteByte(s[k])
		}
	}
	return bs.String()
}

func ReplaceByte(s string, o, n byte) string {
	if o == n {
		return s
	}
	if strings.IndexByte(s, o) < 0 { // 不存在
		return s
	}
	var bs strings.Builder
	bs.Grow(len(s))
	for k := range s {
		if s[k] != o {
			bs.WriteByte(s[k])
		} else {
			bs.WriteByte(n)
		}
	}
	return bs.String()
}

func RuneLen(s string) int {
	return utf8.RuneCountInString(s)
}
