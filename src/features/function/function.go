package function

import (
	interpreteStructs "makeLanguages/src/interprete/structs"
	"makeLanguages/src/languageContext"
	lexerStructs "makeLanguages/src/lexer/lexerStructs"
)

type IFunction interface {
	Execute(params *[]interface{}) (interface{}, bool)
	GetParams() *[]lexerStructs.Token
	GetBody() interface{}
	GetContext() *languageContext.Context
}

type BaseFunction struct {
	Context  *languageContext.Context
	callBack func(params *[]interface{}) interface{}
	Name     string
}

func NewBaseFunction(Context *languageContext.Context, Name string, callBack func(params *[]interface{}) interface{}) *BaseFunction {
	return &BaseFunction{
		Context,
		callBack,
		Name,
	}
}

func (func_ BaseFunction) GetParams() *[]lexerStructs.Token {
	panic("internal error hasACustomExecute need to be true ")
}
func (func_ BaseFunction) GetBody() interface{} {
	panic("internal error hasACustomExecute need to be true ")
}
func (func_ BaseFunction) GetContext() *languageContext.Context {
	panic("internal error hasACustomExecute need to be true ")
}

func (func_ BaseFunction) Execute(params *[]interface{}) (interface{}, bool) {
	value := func_.callBack(params)
	hasACustomExecute := true
	return value, hasACustomExecute
}

type Function struct {
	Context *languageContext.Context
	Params  *[]lexerStructs.Token
	Body    interface{}
}

func (func_ Function) Execute(params *[]interface{}) (interface{}, bool) {
	i := 0
	for _, token := range *func_.Params {
		func_.Context.Set(token.Value.(string),
			interpreteStructs.VarType{
				Value: (*params)[i],
			})
		i++
	}
	hasACustomExecute := false
	return nil, hasACustomExecute
}
func (func_ Function) GetParams() *[]lexerStructs.Token {
	return func_.Params
}
func (func_ Function) GetBody() interface{} {
	return func_.Body
}
func (func_ Function) GetContext() *languageContext.Context {
	return func_.Context
}
