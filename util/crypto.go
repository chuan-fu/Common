package util

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5 加密
func Md5(encodeString string) string {
	h := md5.New()
	h.Write(StringToBytes(encodeString))
	return hex.EncodeToString(h.Sum(nil))
}
