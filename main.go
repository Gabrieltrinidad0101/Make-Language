package main

import (
	"makeLanguages/src"
	"makeLanguages/src/api"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/interprete/interpreteStructs"
)

func add(value1 interpreteStructs.IBaseElement, value2 interpreteStructs.IBaseElement) interface{} {
	return numbers.NewNumbers(value1.GetValue().(float64)*value2.GetValue().(float64), nil)
}

func main() {
	api_ := api.NewApi()
	api_.AddOperetor("PLUS", add)
	src.MakeLanguage("./conf.json", "./main.mkL", api_)
}
