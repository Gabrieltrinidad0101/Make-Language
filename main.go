package main

import (
	"fmt"
	"os"

	"github.com/Gabrieltrinidad0101/Make-Language/src/api"
	"github.com/Gabrieltrinidad0101/Make-Language/src/utils"

	"github.com/Gabrieltrinidad0101/Make-Language/src/parser/parserStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/booleans"

	src "github.com/Gabrieltrinidad0101/Make-Language/src"
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

func makeFile(params *[]interpreteStructs.IBaseElement) interface{} {
	utils.ValidateTypes(params, "String_")
	_, err := os.Create((*params)[0].GetValue().(string))
	if err != nil {
		fmt.Println(err)
	}

	return parserStructs.NullNode{}
}

func main() {
	makeLanguage := src.NewMakeLanguage("./conf.json", "./main.mkl")
	makeLanguage.AddOperetor("<1>", lessOrGreaterOne)
	makeLanguage.AddFunction("printLn2", printLn2)

	methods := api.Methods{}
	methods["create"] = makeFile
	makeLanguage.AddClass("File", api.CustomClassValues{
		Methods: methods,
	})
	makeLanguage.Run()
}
