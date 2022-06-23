package cast

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/chuan-fu/Common/baseservice/jsonx"

	"github.com/chuan-fu/Common/baseservice/stringx"
	"github.com/pkg/errors"
)

const (
	Int64Base  = 10
	NumBitSize = 64
)

func ToInt(s string) (int, error) {
	if s == "" {
		return 0, errors.New(`ToInt: unable to ToInt("")`)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrap(err, "ToInt")
	}
	return i, nil
}

func ToInt64(s string) (int64, error) {
	if s == "" {
		return 0, errors.New(`ToInt64: unable to ToInt64("")`)
	}
	i, err := strconv.ParseInt(s, Int64Base, NumBitSize)
	if err != nil {
		return 0, errors.Wrap(err, "ToInt64")
	}
	return i, nil
}

func ToFloat(s string) (float64, error) {
	if s == "" {
		return 0, errors.New(`ToFloat: unable to ToFloat("")`)
	}
	i, err := strconv.ParseFloat(s, NumBitSize)
	if err != nil {
		return 0, errors.Wrap(err, "ToFloat")
	}
	return i, nil
}

func ToIntI(v interface{}) (int, error) {
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
	switch vt := v.(type) {
	case []byte:
		return stringx.BytesToString(vt)
	case string:
		return vt
	case int64:
		return strconv.FormatInt(vt, Int64Base)
	case int:
		return strconv.Itoa(vt)
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
		return jsonx.Marshal(v)
	default:
		return fmt.Sprintf("%v", rv.Interface())
	}
}

// 返回非指针类型【取地址】
func Indirect(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	return rv.Interface()
}
