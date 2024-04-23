package parserStructs

import (
	"makeLanguages/src/lexer/lexerStructs"
)

type baseGetValue struct{}

func (unaryOp baseGetValue) GetValue() interface{} {
	panic("Get Value No implement")
}

type BinOP struct {
	LeftNode  interface{}
	Operation lexerStructs.Token
	RigthNode interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type UnaryOP struct {
	Operation string
	RigthNode interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type IfNode struct {
	Ifs   []*ConditionAndBody
	Else_ interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type ConditionAndBody struct {
	Condition interface{}
	Body      interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type VarAssignNode struct {
	Identifier string
	Node       interface{}
	IsConstant bool
	lexerStructs.IPositionBase
	baseGetValue
}

type UpdateVariableNode struct {
	Identifier string
	Node       interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type VarAccessNode struct {
	lexerStructs.IPositionBase
	Identifier string
	baseGetValue
}

type ListNode struct {
	Nodes []interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type WhileNode struct {
	Condition interface{}
	Body      interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type FuncNode struct {
	Params *[]lexerStructs.Token
	Body   interface{}
	Name   string
	lexerStructs.IPositionBase
	baseGetValue
}

type ClassNode struct {
	Methods    interface{}
	Properties string
	Name       string
	lexerStructs.IPositionBase
	baseGetValue
}

type ForNode struct {
	Expr1     interface{}
	Expr2     interface{}
	Condition interface{}
	Body      interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type CallObjectNode struct {
	Params *[]interface{}
	Name   string
	HasNew bool
	lexerStructs.IPositionBase
	baseGetValue
}

type NullNode struct{}

type ContinueNode struct {
	lexerStructs.IPositionBase
	baseGetValue
}

type ReturnNode struct {
	lexerStructs.IPositionBase
	Value interface{}
	baseGetValue
}

type BreakNode struct {
	lexerStructs.IPositionBase
	baseGetValue
}

type StringNode struct {
	Value string
	lexerStructs.IPositionBase
	baseGetValue
}

type ArrayAccess struct {
	Identifier string
	Node       interface{}
	lexerStructs.IPositionBase
	baseGetValue
}

type ThisNode struct {
	lexerStructs.IPositionBase
	baseGetValue
}
