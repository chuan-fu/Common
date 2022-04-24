package util

import (
	"reflect"
)

// 返回非指针类型【取地址】
func Indirect(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return v
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	return rv.Interface()
}
