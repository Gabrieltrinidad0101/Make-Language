package function

import (
	"fmt"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
	lexerStructs "makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/parser/parserStructs"
)

type Print struct {
	BaseFunction
	Context *languageContext.Context
	Params  *[]lexerStructs.Token
	Body    interface{}
	lexerStructs.IPositionBase
	parserStructs.BaseGetValue
}

func (func_ Print) Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error) {
	for _, param := range *params {
		fmt.Print(param.GetValue())
	}
	fmt.Println()
	hasACustomExecute := true
	return parserStructs.NullNode{}, hasACustomExecute, nil
}
