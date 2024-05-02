package booleans

import (
	"makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/parser/parserStructs"
)

type Boolean struct {
	Value bool
	lexerStructs.IPositionBase
	parserStructs.BaseGetValue
}

func NewBoolean(value_ bool) *Boolean {
	return &Boolean{
		Value: value_,
	}
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
