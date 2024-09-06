package interpreteStructs

import "github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"

type VarType struct {
	Value            IBaseElement
	IsConstant       bool
	Type             string
	OnUpdateVariable func(value interface{})
}

type IBaseElement interface {
	lexerStructs.IPositionBase
	GetValue() interface{}
}
