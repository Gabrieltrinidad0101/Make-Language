package main

import (
	"fmt"
	"makeLanguages/src"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/parser/parserStructs"
	"makeLanguages/src/utils"
)

func lessOrGreaterOne(value1 interpreteStructs.IBaseElement, value2 interpreteStructs.IBaseElement) interface{} {
	params := &[]interpreteStructs.IBaseElement{
		value1,
		value2,
	}
	utils.ValidateTypes(params, "Number", "Number")
	number1 := value1.GetValue().(float64)
	number2 := value2.GetValue().(float64)
	boolean := number1+1 == number2 || number1-1 == number2
	return booleans.NewBoolean(boolean)
}

func printLn2(params *[]interpreteStructs.IBaseElement) interface{} {
	fmt.Println((*params)[0].GetValue())
	fmt.Println()
	return parserStructs.NullNode{}
}

func main() {
	makeLanguage := src.NewMakeLanguage("./conf.json", "./main.mkL")
	makeLanguage.AddOperetor("<1>", lessOrGreaterOne)
	makeLanguage.AddFunction("printLn2", printLn2)
	makeLanguage.Run()

}
