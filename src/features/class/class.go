package class

import (
	"strings"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"
	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"
	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser/parserStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/function"
)

type ClassBase interface {
	GetClassContext() *languageContext.Context
	lexerStructs.IPositionBase
}

type BuildClass struct {
	Context *languageContext.Context
}

func NewBuildClass(Context *languageContext.Context) *BuildClass {
	return &BuildClass{
		Context,
	}
}

func (buildClass *BuildClass) AddMethod(name string, callBack func(params *[]interpreteStructs.IBaseElement) interface{}) {
	newMethod := function.NewBaseFunction(buildClass.Context, name, callBack, false)
	buildClass.Context.Set(name, &interpreteStructs.VarType{
		Value:      newMethod,
		IsConstant: true,
	})
}

func (buildClass *BuildClass) AddProperty(name string, varType *interpreteStructs.VarType) {
	buildClass.Context.Set(name, varType)
}

type Class struct {
	Context *languageContext.Context
	Name    string
	lexerStructs.IPositionBase
	parserStructs.BaseGetValue
}

func (class Class) GetClassContext() *languageContext.Context {
	return class.Context
}

func (class Class) GetValue() interface{} {
	result := class.Name + "\n"
	result += class.Context.GetString(strings.Repeat(" ", len(class.Name)-1))
	return result
}
