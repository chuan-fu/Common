package util

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/chuan-fu/Common/baseservice/stringx"
)

// Md5 加密
func Md5(s string) string {
	h := md5.New()
	h.Write(stringx.StringToBytes(s))
	return hex.EncodeToString(h.Sum(nil))
}
