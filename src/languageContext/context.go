package languageContext

import "makeLanguages/src/interprete/interpreteStructs"

type Variables *map[string]interpreteStructs.VarType

type Context struct {
	Parent    interface{}
	ScopeType string // open calibrase,function or global
	variables Variables
	IsClass   bool
}

func NewContext(parent *Context) *Context {
	context := &Context{
		Parent:    parent,
		variables: &map[string]interpreteStructs.VarType{},
	}
	return context
}

func (context *Context) Get(name string) (interpreteStructs.VarType, bool) {
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
	(*context.variables)[name] = varType
}

func (context *Context) SetClass(name string, varType interpreteStructs.VarType) {
	if !context.IsClass {
		currentContext := context
		if currentContext.Parent.(*Context) == nil {
			return
		}
		currentContext = currentContext.Parent.(*Context)
		currentContext.SetClass(name, varType)
		return
	}
	context.SetClass(name, varType)
}

func (context *Context) GetClassContext() (*Context, bool) {
	if !context.IsClass {
		if context.Parent.(*Context) == nil {
			return nil, false
		}
		currentContext := context.Parent.(*Context)
		return currentContext.GetClassContext()
	}
	return context, true
}
