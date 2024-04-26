package str

import (
	"fmt"
	"makeLanguages/src/features"
	"makeLanguages/src/features/class"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/utils"
	"strings"
)

type String_ struct {
	lexerStructs.IPositionBase
	Context *languageContext.Context
	Value   string
}

func NewString(value string, position lexerStructs.IPositionBase) *String_ {
	string_ := String_{
		Value:         value,
		Context:       languageContext.NewContext(nil),
		IPositionBase: position,
	}
	string_.Initial()
	return &string_
}

func (string_ String_) GetValue() interface{} {
	return string_.Value
}

func (string_ String_) GetClassContext() *languageContext.Context {
	return string_.Context
}

func (string String_) Replace(params *[]interface{}) interface{} {
	if len(*params) > 2 {
		panic("Replace")
	}

	string1 := (*params)[0].(*String_)
	string2 := (*params)[1].(*String_)

	newString := strings.ReplaceAll(string.Value, string1.Value, string2.Value)

	return NewString(newString, nil)
}

func (string_ String_) Upper(params *[]interface{}) interface{} {
	if len(*params) > 0 {
		panic("Upper")
	}

	newString := strings.ToUpper(string_.Value)
	return NewString(newString, nil)
}

func (string_ String_) PLUS(node features.Type) *String_ {
	return NewString(string_.Value+fmt.Sprint(node.GetValue()), nil)
}

func (string_ String_) MUL(node features.Type) *String_ {
	if utils.GetType(node.GetValue()) != "float64" {
		panic("Mul string")
	}

	stringRepeat := strings.Repeat(string_.Value, int(node.GetValue().(float64)))
	return NewString(stringRepeat, nil)
}

func (string_ String_) Initial() {
	newClass := class.NewBuildClass(string_.Context)
	newClass.AddMethod("replace", string_.Replace)
	newClass.AddMethod("upper", string_.Upper)
}
