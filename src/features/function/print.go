package function

import (
	"fmt"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/parser/parserStructs"
)

type Print struct {
	BaseFunction
}

func (func_ Print) Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error) {
	for _, param := range *params {
		fmt.Print(param.GetValue())
	}
	fmt.Println()
	hasACustomExecute := true
	return parserStructs.NullNode{}, hasACustomExecute, nil
}
