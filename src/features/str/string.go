package str

import (
	"fmt"
	"makeLanguages/src/features"
	"makeLanguages/src/features/class"
	"makeLanguages/src/languageContext"
	"strings"
)

type String_ struct {
	Context *languageContext.Context
	Value   string
}

func NewString(value string) *String_ {
	string_ := String_{
		Value:   value,
		Context: languageContext.NewContext(nil),
	}
	string_.Initial()
	return &string_
}

func (class String_) GetClassContext() *languageContext.Context {
	return class.Context
}

func (string String_) Replace(params *[]interface{}) interface{} {
	if len(*params) > 2 {
		panic("Replace")
	}

	string1 := (*params)[0].(*String_)
	string2 := (*params)[1].(*String_)

	newString := strings.ReplaceAll(string.Value, string1.Value, string2.Value)

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

func (string_ String_) PLUS(node features.Type) *String_ {
	return NewString(string_.Value + fmt.Sprint(node.GetValue()))
}

func (string_ String_) Initial() {
	newClass := class.NewBuildClass(string_.Context)
	newClass.AddMethod("replace", string_.Replace)
	newClass.AddMethod("concat", string_.Concat)
}
