package class

import (
	"makeLanguages/src/features/function"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
)

type ClassBase interface {
	GetClassContext() *languageContext.Context
}

type BuildClass struct {
	Context *languageContext.Context
}

func NewBuildClass(Context *languageContext.Context) *BuildClass {
	return &BuildClass{
		Context,
	}
}

func (buildClass *BuildClass) AddMethod(name string, callBack func(params *[]interface{}) interface{}) {
	newMethod := function.NewBaseFunction(buildClass.Context, name, callBack)
	buildClass.Context.Set(name, interpreteStructs.VarType{
		Value:      newMethod,
		IsConstant: true,
	})
}

type Class struct {
	Context *languageContext.Context
	Name    string
}

func (class Class) GetClassContext() *languageContext.Context {
	return class.Context
}
