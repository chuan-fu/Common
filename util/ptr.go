package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/chuan-fu/Common/cdefs"

	"github.com/chuan-fu/Common/baseservice/stringx"
)

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

func SetReflectValueByStr(value reflect.Value, val string) error {
	// 部分类型特殊处理
	switch value.Interface().(type) {
	case time.Duration:
		return setTimeDuration(val, value)
	case time.Time:
		return setTimeField(val, value)
	case []byte:
		return setByteSlice(val, value)
	}

	switch value.Kind() {
	case reflect.Int:
		return setIntField(val, cdefs.BitSize0, value)
	case reflect.Int8:
		return setIntField(val, cdefs.BitSize8, value)
	case reflect.Int16:
		return setIntField(val, cdefs.BitSize16, value)
	case reflect.Int32:
		return setIntField(val, cdefs.BitSize32, value)
	case reflect.Int64:
		return setIntField(val, cdefs.BitSize64, value)
	case reflect.Uint:
		return setUintField(val, cdefs.BitSize0, value)
	case reflect.Uint8:
		return setUintField(val, cdefs.BitSize8, value)
	case reflect.Uint16:
		return setUintField(val, cdefs.BitSize16, value)
	case reflect.Uint32:
		return setUintField(val, cdefs.BitSize32, value)
	case reflect.Uint64:
		return setUintField(val, cdefs.BitSize64, value)
	case reflect.Bool:
		return setBoolField(val, value)
	case reflect.Float32:
		return setFloatField(val, cdefs.BitSize32, value)
	case reflect.Float64:
		return setFloatField(val, cdefs.BitSize64, value)
	case reflect.Complex64:
		return setComplexField(val, cdefs.BitSize64, value)
	case reflect.Complex128:
		return setComplexField(val, cdefs.BitSize128, value)
	case reflect.String:
		value.SetString(val)
	case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
		return json.Unmarshal(stringx.StringToBytes(val), value.Addr().Interface())
	default:
		return fmt.Errorf("setReflectValueByStr: unable to set type %d", value.Kind())
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) (err error) {
	var intVal int64
	if intVal, err = strconv.ParseInt(val, cdefs.BitSize10, bitSize); err != nil {
		return
	}
	field.SetInt(intVal)
	return
}

func setUintField(val string, bitSize int, field reflect.Value) (err error) {
	var intVal uint64
	if intVal, err = strconv.ParseUint(val, cdefs.BitSize10, bitSize); err != nil {
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
	value.Set(reflect.ValueOf(stringx.StringToBytes(val)))
	return nil
}
