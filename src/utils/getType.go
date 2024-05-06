package utils

import (
	"fmt"
	"makeLanguages/src/constants"
	"makeLanguages/src/customErrors"
	"makeLanguages/src/interprete/interpreteStructs"
	"reflect"
)

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

func ValidateTypes(params *[]interpreteStructs.IBaseElement, types ...string) {
	for i, type_ := range types {
		paramType := GetType((*params)[i])
		if paramType != type_ {
			customErrors.RunTimeError((*params)[i], fmt.Sprintf("Invalid params got %s expect %s", paramType, type_), constants.STOP_EXECUTION)
		}
	}
}
