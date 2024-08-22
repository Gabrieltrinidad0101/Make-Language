package utils

import (
	"fmt"
	"reflect"

	"github.com/Gabrieltrinidad0101/Make-Language/src/constants"

	"github.com/Gabrieltrinidad0101/Make-Language/src/customErrors"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"
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
