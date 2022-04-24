package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

func ToIntI(v interface{}) (int, error) {
	if v = Indirect(v); v == nil {
		return 0, errors.New("ToIntI: unable to ToIntI(nil)")
	}
	switch s := v.(type) {
	case int:
		return s, nil
	case int64:
		return int(s), nil
	case int32:
		return int(s), nil
	case int16:
		return int(s), nil
	case int8:
		return int(s), nil
	case uint:
		return int(s), nil
	case uint64:
		return int(s), nil
	case uint32:
		return int(s), nil
	case uint16:
		return int(s), nil
	case uint8:
		return int(s), nil
	case float64:
		return int(s), nil
	case float32:
		return int(s), nil
	case string:
		return ToInt(s)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("ToIntI: unable to cast %#v of type %T to int", v, v)
	}
}

func ToInt64I(v interface{}) (int64, error) {
	if v = Indirect(v); v == nil {
		return 0, errors.New("ToInt64I: unable to ToInt64I(nil)")
	}
	switch s := v.(type) {
	case int64:
		return s, nil
	case int:
		return int64(s), nil
	case int32:
		return int64(s), nil
	case int16:
		return int64(s), nil
	case int8:
		return int64(s), nil
	case uint:
		return int64(s), nil
	case uint64:
		return int64(s), nil
	case uint32:
		return int64(s), nil
	case uint16:
		return int64(s), nil
	case uint8:
		return int64(s), nil
	case float64:
		return int64(s), nil
	case float32:
		return int64(s), nil
	case string:
		return ToInt64(s)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("ToInt64I: unable to cast %#v of type %T to int64", v, v)
	}
}

func ToFloat64I(v interface{}) (float64, error) {
	if v = Indirect(v); v == nil {
		return 0, errors.New("ToFloat64I: unable to ToFloat64I(nil)")
	}
	switch s := v.(type) {
	case float64:
		return s, nil
	case float32:
		return float64(s), nil
	case int:
		return float64(s), nil
	case int64:
		return float64(s), nil
	case int32:
		return float64(s), nil
	case int16:
		return float64(s), nil
	case int8:
		return float64(s), nil
	case uint:
		return float64(s), nil
	case uint64:
		return float64(s), nil
	case uint32:
		return float64(s), nil
	case uint16:
		return float64(s), nil
	case uint8:
		return float64(s), nil
	case string:
		return ToFloat(s)
	case bool:
		if s {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("ToFloat64I: unable to cast %#v of type %T to float64", v, v)
	}
}

func ToString(v interface{}) string {
	if v = Indirect(v); v == nil {
		return ""
	}
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
