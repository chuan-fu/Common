package util

import "reflect"

func Type(v interface{}) string {
	return reflect.TypeOf(v).String()
}
