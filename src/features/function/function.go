package function

import (
	"makeLanguages/src/languageContext"
	"makeLanguages/src/token"
)

type IFunction interface {
	Execute(params *[]interface{}) (interface{}, bool)
	GetParams() *[]token.Token
	GetBody() interface{}
	GetContext() *languageContext.Context
}

type BaseFunction struct{}

func (func_ BaseFunction) GetParams() *[]token.Token {
	panic("internal error hasACustomExecute need to be false ")
}
func (func_ BaseFunction) GetBody() interface{} {
	panic("internal error hasACustomExecute need to be false ")
}
func (func_ BaseFunction) GetContext() *languageContext.Context {
	panic("internal error hasACustomExecute need to be false ")
}

type Function struct {
	Context *languageContext.Context
	Params  *[]token.Token
	Body    interface{}
}

func (func_ Function) Execute(params *[]interface{}) (interface{}, bool) {
	i := 0
	for _, token := range *func_.Params {
		func_.Context.Set(token.Value.(string), (*params)[i])
		i++
	}
	hasACustomExecute := false
	return nil, hasACustomExecute
}

func (func_ Function) GetParams() *[]token.Token {
	return func_.Params
}
func (func_ Function) GetBody() interface{} {
	return func_.Body
}
func (func_ Function) GetContext() *languageContext.Context {
	return func_.Context
}
