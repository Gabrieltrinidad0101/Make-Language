package array

import "makeLanguages/src/languageContext"

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
