package array

import (
	"makeLanguages/src/features/class"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
	"slices"
)

type Array struct {
	context *languageContext.Context
	Value   *[]interface{}
}

func NewArray(value *[]interface{}) *Array {
	context := languageContext.NewContext(nil)
	return &Array{
		context: context,
		Value:   value,
	}
}

func (array *Array) PLUS(element interface{}) *Array {
	*array.Value = append(*array.Value, element)
	return array
}

func (array *Array) Concat(params *[]interface{}) interface{} {
	if len(*params) > 1 {
		panic("Concat")
	}

	value1 := (*params)[0].(interpreteStructs.IBaseElement)
	concatValue := slices.Concat[[]interface{}](*array.Value, value1.GetValue().([]interface{}))
	return NewArray(&concatValue)
}

func (array Array) Initial() {
	newClass := class.NewBuildClass(array.context)
	newClass.AddMethod("concat", array.Concat)
}
