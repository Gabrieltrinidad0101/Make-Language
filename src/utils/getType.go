package utils

import "reflect"

func GetType(node interface{}) string {
	method := reflect.TypeOf(node)
	if method == nil {
		panic(node)
	}
	if method.Kind() == reflect.Ptr {
		return method.Elem().Name()
	} else {
		return method.Name()
	}
}
