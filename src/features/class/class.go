package class

import (
	"makeLanguages/src/features/function"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/parser/parserStructs"
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
	newMethod := function.NewBaseFunction(buildClass.Context, name, callBack)
	buildClass.Context.Set(name, &interpreteStructs.VarType{
		Value:      newMethod,
		IsConstant: true,
	})
}

func (buildClass *BuildClass) AddProperty(name string, property interpreteStructs.IBaseElement) {
	buildClass.Context.Set(name, &interpreteStructs.VarType{
		Value:      property,
		IsConstant: true,
	})
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
