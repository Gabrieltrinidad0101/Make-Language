package main

import (
	"makeLanguages/src"
	"makeLanguages/src/api"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/interprete/interpreteStructs"
)

func lessOrGreaterOne(value1 interpreteStructs.IBaseElement, value2 interpreteStructs.IBaseElement) interface{} {
	number1 := value1.GetValue().(float64)
	number2 := value2.GetValue().(float64)
	boolean := number1+1 == number2 || number1-1 == number2
	return booleans.NewBoolean(boolean)
}

func main() {
	api_ := api.NewApi()
	api_.AddOperetor("LESS_OR_GREATER_ONE", lessOrGreaterOne)
	src.MakeLanguage("./conf.json", "./main.mkL", api_)
}
