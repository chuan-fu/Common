package cdao

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/chuan-fu/Common/util"
)

func formatSec(dur time.Duration) int64 {
	if dur <= 0 {
		return -1
	}
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}

func formatMs(dur time.Duration) int64 {
	if dur <= 0 {
		return -1
	}
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

const (
	BitSize0   = 0
	BitSize8   = 8
	BitSize10  = 10
	BitSize16  = 16
	BitSize32  = 32
	BitSize64  = 64
	BitSize128 = 128
)

func setReflectValueByStr(value reflect.Value, val string) error {
	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, BitSize0, value)
	case reflect.Int8:
		return setIntField(val, BitSize8, value)
	case reflect.Int16:
		return setIntField(val, BitSize16, value)
	case reflect.Int32:
		return setIntField(val, BitSize32, value)
	case reflect.Int64:
		if _, ok := value.Interface().(time.Duration); ok {
			return setTimeDuration(val, value)
		}
		return setIntField(val, BitSize64, value)
	case reflect.Uint:
		return setUintField(val, BitSize0, value)
	case reflect.Uint8:
		return setUintField(val, BitSize8, value)
	case reflect.Uint16:
		return setUintField(val, BitSize16, value)
	case reflect.Uint32:
		return setUintField(val, BitSize32, value)
	case reflect.Uint64:
		return setUintField(val, BitSize64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, BitSize32, value)
	case reflect.Float64:
		return setFloatField(val, BitSize64, value)
	case reflect.Complex64:
		return setComplexField(val, BitSize64, value)
	case reflect.Complex128:
		return setComplexField(val, BitSize128, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return setTimeField(val, value)
		}
		return json.Unmarshal(util.StringToBytes(val), value.Addr().Interface())
	case reflect.Map, reflect.Array:
		return json.Unmarshal(util.StringToBytes(val), value.Addr().Interface())
	case reflect.Slice:
		if _, ok := value.Interface().([]byte); ok {
			return setByteSlice(val, value)
		}
		return json.Unmarshal(util.StringToBytes(val), value.Addr().Interface())
	default:
		return errors.New(fmt.Sprintf("setReflectValueByStr 不可转换类型:%d", value.Kind()))
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) (err error) {
	var intVal int64
	if intVal, err = strconv.ParseInt(val, BitSize10, bitSize); err != nil {
		return
	}
	field.SetInt(intVal)
	return
}

func setUintField(val string, bitSize int, field reflect.Value) (err error) {
	var intVal uint64
	if intVal, err = strconv.ParseUint(val, BitSize10, bitSize); err != nil {
		return
	}
	field.SetUint(intVal)
	return
}

func setBoolField(val string, field reflect.Value) (err error) {
	var boolVal bool
	if boolVal, err = strconv.ParseBool(val); err != nil {
		return
	}
	field.SetBool(boolVal)
	return
}

func setFloatField(val string, bitSize int, field reflect.Value) (err error) {
	var floatVal float64
	if floatVal, err = strconv.ParseFloat(val, bitSize); err != nil {
		return
	}
	field.SetFloat(floatVal)
	return
}

func setComplexField(val string, bitSize int, field reflect.Value) (err error) {
	var floatVal complex128
	if floatVal, err = strconv.ParseComplex(val, bitSize); err != nil {
		return
	}
	field.SetComplex(floatVal)
	return
}

func setTimeDuration(val string, value reflect.Value) error {
	t, err := time.ParseDuration(val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}

func setTimeField(val string, value reflect.Value) error {
	t, err := time.Parse(time.RFC3339Nano, val)
	if err != nil {
		return err
	}
	value.Set(reflect.ValueOf(t))
	return nil
}

func setByteSlice(val string, value reflect.Value) error {
	value.Set(reflect.ValueOf(util.StringToBytes(val)))
	return nil
}

func getTag(field reflect.Type, i int, tag string) string {
	if tag != "" {
		if d := field.Field(i).Tag.Get(tag); d != "" && d != "-" {
			return d
		}
	}
	return field.Field(i).Name
}