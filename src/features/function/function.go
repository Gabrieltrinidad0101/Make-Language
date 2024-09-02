package function

import (
	"fmt"

	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"
	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser/parserStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"
)

type IFunction interface {
	Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error)
	GetParams() *[]lexerStructs.Token
	GetBody() interpreteStructs.IBaseElement
	GetContext() *languageContext.Context
	CanChangeContextParent() bool
}

type BaseFunction struct {
	Context  *languageContext.Context
	callBack func(params *[]interpreteStructs.IBaseElement) interface{}
	Name     string
	lexerStructs.IPositionBase
	parserStructs.BaseGetValue
	canChangeContextParent bool
}

func NewBaseFunction(Context *languageContext.Context, Name string, callBack func(params *[]interpreteStructs.IBaseElement) interface{}, canChangeContextParent bool) *BaseFunction {
	return &BaseFunction{
		Context:                Context,
		callBack:               callBack,
		Name:                   Name,
		canChangeContextParent: canChangeContextParent,
	}
}

func (func_ BaseFunction) GetParams() *[]lexerStructs.Token {
	panic("internal error hasACustomExecute need to be true ")
}
func (func_ BaseFunction) GetBody() interpreteStructs.IBaseElement {
	panic("internal error hasACustomExecute need to be true ")
}

func (func_ BaseFunction) GetContext() *languageContext.Context {
	if func_.Context != nil {
		return func_.Context
	}
	panic("internal error hasACustomExecute need to be true ")
}

func (func_ BaseFunction) CanChangeContextParent() bool {
	return func_.canChangeContextParent
}

func (func_ BaseFunction) Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error) {
	value := func_.callBack(params)
	hasACustomExecute := true
	return value, hasACustomExecute, nil
}

type Function struct {
	Context *languageContext.Context
	Params  *[]lexerStructs.Token
	Body    interpreteStructs.IBaseElement
	parserStructs.BaseGetValue
	lexerStructs.IPositionBase
}

func (func_ Function) Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error) {
	i := 0
	if len(*func_.Params) != len(*params) {
		return nil, false, fmt.Errorf("Invalid params expect %d, got %d", len(*func_.Params), len(*params))
	}
	for _, token := range *func_.Params {
		func_.Context.Set(token.Value.(string),
			&interpreteStructs.VarType{
				Value: (*params)[i],
			})
		i++
	}
	hasACustomExecute := false
	return nil, hasACustomExecute, nil
}
func (func_ Function) GetParams() *[]lexerStructs.Token {
	return func_.Params
}
func (func_ Function) GetBody() interpreteStructs.IBaseElement {
	return func_.Body
}

func (func_ Function) GetContext() *languageContext.Context {
	return func_.Context
}

func (func_ Function) CanChangeContextParent() bool {
	return false
}
