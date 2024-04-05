package function

import (
	"fmt"
	"makeLanguages/src/languageContext"
	lexerStructs "makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/parser/parserStructs"
)

type Print struct {
	BaseFunction
	Context *languageContext.Context
	Params  *[]lexerStructs.Token
	Body    interface{}
}

func (func_ Print) Execute(params *[]interface{}) (interface{}, bool) {
	fmt.Println(*params...)
	hasACustomExecute := true
	return parserStructs.NullNode{}, hasACustomExecute
}
