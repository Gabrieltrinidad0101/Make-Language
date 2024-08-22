package parserStructs

import (
	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"
)

type BaseGetValue struct{}

func (unaryOp BaseGetValue) GetValue() interface{} {
	panic("Get Value No implement")
}

type BinOP struct {
	LeftNode  interpreteStructs.IBaseElement
	Operation lexerStructs.Token
	RigthNode interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type UnaryOP struct {
	Operation string
	RigthNode interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type IfNode struct {
	Ifs   []*ConditionAndBody
	Else_ interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type ConditionAndBody struct {
	Condition interpreteStructs.IBaseElement
	Body      interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type VarAssignNode struct {
	Identifier string
	Node       interpreteStructs.IBaseElement
	IsConstant bool
	lexerStructs.IPositionBase
	BaseGetValue
}

type UpdateVariableNode struct {
	Identifier string
	Node       interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type VarAccessNode struct {
	lexerStructs.IPositionBase
	Identifier string
	BaseGetValue
}

type ListNode struct {
	Nodes []interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type WhileNode struct {
	Condition interpreteStructs.IBaseElement
	Body      interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type FuncNode struct {
	Params *[]lexerStructs.Token
	Body   interpreteStructs.IBaseElement
	Name   string
	Class  bool
	lexerStructs.IPositionBase
	BaseGetValue
}

type ClassNode struct {
	Methods    interface{}
	Properties string
	Name       string
	lexerStructs.IPositionBase
	BaseGetValue
}

type ForNode struct {
	Expr1     interpreteStructs.IBaseElement
	Expr2     interpreteStructs.IBaseElement
	Condition interpreteStructs.IBaseElement
	Body      interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type CallObjectNode struct {
	Params *[]interpreteStructs.IBaseElement
	Name   string
	HasNew bool
	lexerStructs.IPositionBase
	BaseGetValue
}

type NullNode struct {
	lexerStructs.IPositionBase
	BaseGetValue
}

type ContinueNode struct {
	lexerStructs.IPositionBase
	BaseGetValue
}

type ReturnNode struct {
	lexerStructs.IPositionBase
	Value interpreteStructs.IBaseElement
	BaseGetValue
}

type BreakNode struct {
	lexerStructs.IPositionBase
	BaseGetValue
}

type StringNode struct {
	Value string
	lexerStructs.IPositionBase
	BaseGetValue
}

type ArrayAccess struct {
	Identifier string
	Node       interpreteStructs.IBaseElement
	lexerStructs.IPositionBase
	BaseGetValue
}

type ThisNode struct {
	lexerStructs.IPositionBase
	BaseGetValue
}
