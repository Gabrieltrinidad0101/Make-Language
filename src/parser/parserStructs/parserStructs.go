package parserStructs

import "makeLanguages/src/lexer/lexerStructs"

type BinOP struct {
	LeftNode  interface{}
	Operation lexerStructs.Token
	RigthNode interface{}
	lexerStructs.PositionBase
}

type UnaryOP struct {
	Operation string
	RigthNode interface{}
	lexerStructs.PositionBase
}

type IfNode struct {
	Ifs   []*ConditionAndBody
	Else_ interface{}
}

type ConditionAndBody struct {
	Condition interface{}
	Body      interface{}
}

type VarAssignNode struct {
	Identifier string
	Node       interface{}
	IsConstant bool
}

type UpdateVariableNode struct {
	Identifier string
	Node       interface{}
}

type VarAccessNode struct {
	lexerStructs.PositionBase
	Identifier string
}

type ListNode struct {
	Nodes []interface{}
}

type WhileNode struct {
	Condition interface{}
	Body      interface{}
}

type FuncNode struct {
	Params *[]lexerStructs.Token
	Body   interface{}
	Name   string
}

type ClassNode struct {
	Methods    interface{}
	Properties string
	Name       string
}

type ForNode struct {
	Expr1     interface{}
	Expr2     interface{}
	Condition interface{}
	Body      interface{}
}

type CallFuncNode struct {
	Params *[]interface{}
	Name   string
}

type NewClassNode struct {
	Name string
}

type NullNode struct{}

type ContinueNode struct {
	lexerStructs.PositionBase
}

type BreakNode struct {
	lexerStructs.PositionBase
}
