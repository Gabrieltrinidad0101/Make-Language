package array

import (
	"makeLanguages/src/features/class"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer/lexerStructs"
	"slices"
)

type Array struct {
	Value *[]interface{}
	lexerStructs.IPositionBase
	class *class.BuildClass
}

func NewArray(value *[]interface{}) *Array {
	context := languageContext.NewContext(nil)
	newClass := class.NewBuildClass(context)
	newClass.AddProperty("length", numbers.NewNumbers(float64(len(*value)), nil))
	return &Array{
		class: newClass,
		Value: value,
	}
}

func (array *Array) GetClassContext() *languageContext.Context {
	return array.class.Context
}

func (array *Array) GetValue() interface{} {
	return array.Value
}

func (array *Array) PLUS(element interface{}) *Array {
	*array.Value = append(*array.Value, element)
	array.class.AddProperty("length", numbers.NewNumbers(float64(len(*array.Value)), nil))
	return array
}

func (array *Array) Concat(params *[]interpreteStructs.IBaseElement) interface{} {
	if len(*params) > 1 {
		panic("Concat")
	}

	value1 := (*params)[0].(interpreteStructs.IBaseElement)
	concatValue := slices.Concat[[]interface{}](*array.Value, value1.GetValue().([]interface{}))
	return NewArray(&concatValue)
}

func (array Array) Initial() {
	array.class.AddMethod("concat", array.Concat)
}
