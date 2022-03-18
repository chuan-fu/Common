package util

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
		d, _ := json.Marshal(v)
		return BytesToString(d)
	case reflect.String:
		return rv.String()
	default:
		return fmt.Sprintf("%v", rv.Interface())
	}
}
