package str

import (
	"makeLanguages/src/features/function"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/parser/parserStructs"
	"strings"
)

type BaseFunction struct {
	function.BaseFunction
	Context  *languageContext.Context
	callBack func(params *[]interface{}) interface{}
	Name     string
}

func (func_ BaseFunction) Execute(params *[]interface{}) (interface{}, bool) {
	func_.callBack(params)
	hasACustomExecute := true
	return parserStructs.NullNode{}, hasACustomExecute
}

type String_ struct {
	Context *languageContext.Context
	Value   string
}

func NewString(value string) *String_ {
	string_ := String_{
		Value: value,
	}
	string_.Initial()
	return &string_
}

func (string String_) baseFunction(name string, callBack func(params *[]interface{}) interface{}) function.IFunction {
	return BaseFunction{
		Context:  languageContext.NewContext(nil),
		callBack: callBack,
		Name:     name,
	}
}

func (string String_) Replace(params *[]interface{}) interface{} {
	if len(*params) > 2 {
		panic("Replace")
	}

	string1 := (*params)[0].(String_)
	string2 := (*params)[1].(String_)

	newString := strings.Replace(string.Value, string1.Value, string2.Value, 0)

	return NewString(newString)
}

func (string_ String_) Concat(params *[]interface{}) interface{} {
	if len(*params) > 1 {
		panic("Concat")
	}

	string1 := (*params)[0].(String_)

	newString := string1.Value + string_.Value
	return NewString(newString)
}

func (string_ String_) Initial() {
	string_.baseFunction("replace", string_.Replace)
	string_.baseFunction("concat", string_.Concat)
}
