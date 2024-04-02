package function

import (
	"fmt"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/token"
)

type Print struct {
	BaseFunction
	Context *languageContext.Context
	Params  *[]token.Token
	Body    interface{}
}

func (func_ Print) Execute(params *[]interface{}) (interface{}, bool) {
	fmt.Println(*params...)
	hasACustomExecute := true
	return nil, hasACustomExecute
}
