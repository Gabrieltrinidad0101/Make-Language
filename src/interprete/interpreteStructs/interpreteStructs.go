package interpreteStructs

import "makeLanguages/src/lexer/lexerStructs"

type VarType struct {
	Value      IBaseElement
	IsConstant bool
	Type       string
}

type IBaseElement interface {
	lexerStructs.IPositionBase
	GetValue() interface{}
}
