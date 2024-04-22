package interpreteStructs

import "makeLanguages/src/lexer/lexerStructs"

type VarType struct {
	Value      interface{}
	IsConstant bool
	Type       string
}

type IBaseElement interface {
	lexerStructs.IPositionBase
	GetValue() interface{}
}
