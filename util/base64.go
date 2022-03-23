package util

import (
	"encoding/base64"
	"net/url"
	"strings"
)

// base64编码
func Encode(s string) string {
	if s == "" {
		return s
	}
	return base64.StdEncoding.EncodeToString(StringToBytes(s))
}

func EncodeV2(s string) string {
	result := base64.StdEncoding.EncodeToString([]byte(s))
	result = strings.Replace(strings.Replace(strings.Replace(result, "=", "", -1), "/", "_", -1), "+", "-", -1)
	return result
}

// base64解码
func Decode(s string) (string, error) {
	if s == "" {
		return s, nil
	}
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return BytesToString(b), nil
}

func DecodeV2(s string) (string, error) {
	remainder := len(s) % 4
	// base64编码需要为4的倍数，如果不是4的倍数，则填充"="号
	if remainder > 0 {
		padlen := 4 - remainder
		s = s + strings.Repeat("=", padlen)
	}
	// 将原字符串中的"_","-"分别用"/"和"+"替换
	s = strings.Replace(strings.Replace(s, "_", "/", -1), "-", "+", -1)
	result, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return BytesToString(result), nil
}

// url编码
func URLEncode(s string) string {
	if s == "" {
		return s
	}
	return url.QueryEscape(s)
}

// url解码
func URLDecode(s string) (string, error) {
	if s == "" {
		return s, nil
	}
	b, err := url.QueryUnescape(s)
	if err != nil {
		return "", err
	}
	return b, nil
}
