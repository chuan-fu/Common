package util

import (
	"encoding/base64"
	"net/url"
)

// base64编码
func Encode(s string) string {
	if s == "" {
		return s
	}
	return base64.StdEncoding.EncodeToString(StringToBytes(s))
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
