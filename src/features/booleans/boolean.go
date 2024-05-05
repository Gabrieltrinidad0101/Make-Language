package booleans

import (
	"makeLanguages/src/lexer/lexerStructs"
)

type Boolean struct {
	Value bool
	lexerStructs.IPositionBase
}

func NewBoolean(value_ bool) *Boolean {
	return &Boolean{
		Value: value_,
	}
}

func (boolean *Boolean) GetValue() interface{} {
	return boolean.Value
}

func (boolean *Boolean) AND(boolean_ *Boolean) *Boolean {
	if boolean.Value && boolean_.Value {
		return NewBoolean(true)
	}
	return NewBoolean(false)
}

func (boolean *Boolean) OR(boolean_ *Boolean) *Boolean {
	if boolean.Value || boolean_.Value {
		return NewBoolean(true)
	}
	return NewBoolean(false)
}
