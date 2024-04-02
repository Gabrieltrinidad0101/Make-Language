package languageContext

import (
	interpreteStructs "makeLanguages/src/interprete/structs"
)

type Variables *map[string]interface{}

type Context struct {
	Parent    interface{}
	ScopeType string // open calibrase,function or global
	variables Variables
}

func NewContext(parent *Context) Context {
	context := &Context{
		Parent:    parent,
		variables: &map[string]interface{}{},
	}
	return *context
}

func (context *Context) Get(name string) (interface{}, bool) {
	value, ok := (*context.variables)[name]
	if !ok {
		currentContext := context
		if currentContext.Parent.(*Context) == nil {
			return value, ok
		}
		currentContext = currentContext.Parent.(*Context)
		return currentContext.Get(name)
	}
	return value, ok
}

func (context *Context) Update(name string, varType interpreteStructs.VarType) bool {
	_, ok := (*context.variables)[name]
	if !ok {
		currentContext := context
		if currentContext.Parent.(*Context) == nil {
			return false
		}
		currentContext = currentContext.Parent.(*Context)
		currentContext.Set(name, varType)
	}
	context.Set(name, varType)
	return true
}

func (context *Context) Set(name string, varType interpreteStructs.VarType) {
	(*context.variables)[name] = varType.Value
}
