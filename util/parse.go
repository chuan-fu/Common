package util

import (
	"fmt"
	"strconv"
)

const (
	Int64Base  = 10
	NumBitSize = 64
)

func ToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, Int64Base, NumBitSize)
}

func ToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, NumBitSize)
}

func ToString(v interface{}) string {
	switch vt := v.(type) {
	case int:
		return strconv.Itoa(vt)
	case int64:
		return strconv.FormatInt(vt, Int64Base)
	case string:
		return vt
	default:
		return fmt.Sprintf("%v", v)
	}
}
