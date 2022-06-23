package util

import "reflect"

func Type(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func IsPtr(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Ptr
}
