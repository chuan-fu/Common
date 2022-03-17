package util

import "reflect"

func Type(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func IsPtrStruct(v interface{}) bool {
	rt := reflect.TypeOf(v)
	return rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct
}
