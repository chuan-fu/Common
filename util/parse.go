package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
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
	// 部分类型特殊处理
	// 和cdao/internal.go:45里的反序列化相匹配，勿轻易修改
	switch vt := v.(type) {
	case []byte:
		return BytesToString(vt)
	case time.Time:
		return vt.Format(time.RFC3339Nano)
	case time.Duration:
		return timeDurationToString(vt)
	}

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

func timeDurationToString(vt time.Duration) string {
	if vt%time.Hour == 0 {
		return fmt.Sprintf("%dh", vt/time.Hour)
	}
	if vt%time.Minute == 0 {
		return fmt.Sprintf("%dm", vt/time.Minute)
	}
	if vt%time.Second == 0 {
		return fmt.Sprintf("%ds", vt/time.Second)
	}
	if vt%time.Millisecond == 0 {
		return fmt.Sprintf("%dms", vt/time.Millisecond)
	}
	if vt%time.Microsecond == 0 {
		return fmt.Sprintf("%dµs", vt/time.Microsecond)
	}
	return fmt.Sprintf("%dns", vt/time.Nanosecond)
}
