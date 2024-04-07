package parserStructs

import (
	"makeLanguages/src/lexer/lexerStructs"
)

type BinOP struct {
	LeftNode  interface{}
	Operation lexerStructs.Token
	RigthNode interface{}
	lexerStructs.IPositionBase
}

type UnaryOP struct {
	Operation string
	RigthNode interface{}
	lexerStructs.IPositionBase
}

type IfNode struct {
	Ifs   []*ConditionAndBody
	Else_ interface{}
	lexerStructs.IPositionBase
}

type ConditionAndBody struct {
	Condition interface{}
	Body      interface{}
	lexerStructs.IPositionBase
}

type VarAssignNode struct {
	Identifier string
	Node       interface{}
	IsConstant bool
	lexerStructs.IPositionBase
}

type UpdateVariableNode struct {
	Identifier string
	Node       interface{}
	lexerStructs.IPositionBase
}

type VarAccessNode struct {
	lexerStructs.IPositionBase
	Identifier string
}

type ListNode struct {
	Nodes []interface{}
	lexerStructs.IPositionBase
}

type WhileNode struct {
	Condition interface{}
	Body      interface{}
	lexerStructs.IPositionBase
}

type FuncNode struct {
	Params *[]lexerStructs.Token
	Body   interface{}
	Name   string
	lexerStructs.IPositionBase
}

type ClassNode struct {
	Methods    interface{}
	Properties string
	Name       string
	lexerStructs.IPositionBase
}

type ForNode struct {
	Expr1     interface{}
	Expr2     interface{}
	Condition interface{}
	Body      interface{}
	lexerStructs.IPositionBase
}

type CallObjectNode struct {
	Params *[]interface{}
	Name   string
	HasNew bool
	lexerStructs.IPositionBase
}

type NullNode struct{}

type ContinueNode struct {
	lexerStructs.IPositionBase
}

type ReturnNode struct {
	lexerStructs.IPositionBase
	Value interface{}
}

type BreakNode struct {
	lexerStructs.IPositionBase
}

type ClassAccessNode struct {
	Name   string
	Method interface{}
	lexerStructs.IPositionBase
}

type StringNode struct {
	Value string
	lexerStructs.IPositionBase
}
