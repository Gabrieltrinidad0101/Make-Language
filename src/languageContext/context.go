package languageContext

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
	return value, ok
}

func (context *Context) Set(name string, value interface{}) {
	(*context.variables)[name] = value
}
