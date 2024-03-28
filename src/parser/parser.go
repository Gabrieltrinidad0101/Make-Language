package parser

import (
	"makeLanguages/src/constants"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/lexer"
	"reflect"
	"slices"
)

type Parser struct {
	tokens       *[]lexer.Token
	idx          int
	currentToken *lexer.Token
	len          int
}

type BinOP struct {
	LeftNode  interface{}
	Operation lexer.Token
	RigthNode interface{}
}

type UnaryOP struct {
	Operation string
	RigthNode interface{}
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
}

type VarAccessNode struct {
	Identifier string
}

type ListNode struct {
	Nodes []interface{}
}

type NullNode struct{}

func NewParser(tokens *[]lexer.Token) *Parser {
	return &Parser{
		idx:          -1,
		tokens:       tokens,
		currentToken: &lexer.Token{},
		len:          len(*tokens),
	}
}

func (parser *Parser) advance() bool {
	parser.idx++
	if parser.idx >= parser.len {
		return false
	}
	*parser.currentToken = (*parser.tokens)[parser.idx]
	return true
}

func (parser *Parser) getToken(idx int) (*lexer.Token, bool) {
	if parser.idx >= parser.len {
		return nil, false
	}
	return &(*parser.tokens)[idx], true
}

func (parser *Parser) verifyNextToken(tokenType string) bool {
	if parser.currentToken.Type_ != tokenType {
		panic("Expect " + tokenType)
	}

	parser.advance()

	return true
}

func (parser *Parser) binOP(callBack func() interface{}, ops ...string) interface{} {
	leftNode := callBack()
	for slices.Contains[[]string](ops, parser.currentToken.Type_) {
		operation := *parser.currentToken
		parser.advance()
		rigthNode := callBack()
		leftNode = BinOP{
			LeftNode:  leftNode,
			Operation: operation,
			RigthNode: rigthNode,
		}
	}

	return leftNode
}

func (parser *Parser) Parse() interface{} {
	parser.advance()
	return parser.expr()
}

func (parser *Parser) expr() interface{} {
	return parser.statements(constants.EOF)
}

func (parser *Parser) statements(tokenEnd string) interface{} {
	listNodes := ListNode{}
	ast := parser.statement()
	listNodes.Nodes = append(listNodes.Nodes, ast)

	for {
		thereIsANewLine := false
		for parser.currentToken.Type_ == constants.TT_NEWLINE {
			parser.advance()
			thereIsANewLine = true
		}

		if !thereIsANewLine {
			break
		}

		if parser.currentToken.Type_ == tokenEnd {
			parser.advance()
			break
		}

		ast := parser.statement()
		listNodes.Nodes = append(listNodes.Nodes, ast)
	}
	return listNodes
}

func (parser *Parser) statement() interface{} {

	variableAndConst := parser.variableAndConst()
	if variableAndConst != nil {
		return variableAndConst
	}

	ast := parser.compare()
	return ast
}

func (parser *Parser) variableAndConst() interface{} {
	if parser.currentToken.Type_ == constants.TT_VAR {
		parser.advance()
		parser.verifyNextToken(constants.TT_IDENTIFIER)
		identifier := parser.currentToken.Value
		parser.advance()
		parser.verifyNextToken(constants.TT_EQ)
		parser.advance()
		node := parser.compare()

		return VarAssignNode{
			Identifier: identifier.(string),
			Node:       node,
		}
	}
	return nil
}

func (parser *Parser) AndOr() interface{} {
	return parser.binOP(parser.compare, constants.TT_AND, constants.TT_AND)
}

func (parser *Parser) compare() interface{} {
	return parser.binOP(parser.plus, constants.TT_GT, constants.TT_GTE, constants.TT_GT, constants.TT_LT, constants.TT_LTE, constants.TT_EQE)
}

func (parser *Parser) plus() interface{} {
	ast := parser.binOP(parser.factor, constants.TT_PLUS, constants.TT_MINUS)
	return ast
}

func (parser *Parser) factor() interface{} {
	return parser.binOP(parser.pow, constants.TT_MUL, constants.TT_DIV)
}

func (parser *Parser) pow() interface{} {
	return parser.binOP(parser.term, constants.TT_POW, constants.TT_SQUARE_ROOT)
}

func (parser *Parser) term() interface{} {
	nodeType := parser.currentToken.Type_

	if nodeType == constants.TT_PLUS || nodeType == constants.TT_MINUS {
		token, ok := parser.getToken(parser.idx + 1)

		rigthNode := parser.term()

		if ok && reflect.TypeOf(rigthNode).Name() == "UnaryOP" && token.Type_ != constants.TT_LPAREN {
			panic("Error is necesery a ( between - and + simbols")
		}

		unaryOP := UnaryOP{
			Operation: nodeType,
			RigthNode: rigthNode,
		}
		return unaryOP
	}

	if nodeType == "number" {
		value := parser.currentToken.Value.(float64)
		number := numbers.NewNumbers(value)
		parser.advance()
		return number
	}

	if parser.currentToken.Type_ == constants.TT_LPAREN {
		parser.advance()
		node := parser.statement()
		if !(parser.currentToken.Type_ == constants.TT_RPAREN) {
			panic("error )")
		}
		parser.advance()
		return node
	}

	if ifNode, ok := parser.if_(); ok {
		return ifNode
	}

	if varAccess, ok := parser.varAccess(); ok {
		return varAccess
	}

	panic("Error")
}

func (parser *Parser) if_() (interface{}, bool) {
	ifs := []*ConditionAndBody{}
	var elseNode interface{}
	if parser.currentToken.Type_ != constants.TT_IF {
		return nil, false
	}

	parser.advance()

	node, ok := parser.conditionAndBodyBase()

	if !ok {
		return nil, false
	}

	ifs = append(ifs, node)

	for (*parser.currentToken).Type_ == constants.TT_ELIF {
		node, ok := parser.conditionAndBodyBase()

		if !ok {
			return nil, false
		}

		ifs = append(ifs, node)
	}
	if parser.currentToken.Type_ == constants.TT_ELSE {
		parser.advance()
		elseNode = parser.BodyBase()
	}

	return IfNode{
		Ifs:   ifs,
		Else_: elseNode,
	}, true
}

func (parser *Parser) conditionAndBodyBase() (*ConditionAndBody, bool) {

	ok := parser.verifyNextToken(constants.TT_LPAREN)
	if !ok {
		return nil, false
	}

	condition := parser.AndOr()

	ok = parser.verifyNextToken(constants.TT_RPAREN)
	if !ok {
		return nil, false
	}

	body := parser.BodyBase()

	return &ConditionAndBody{
		Condition: condition,
		Body:      body,
	}, true
}

func (parser *Parser) BodyBase() interface{} {
	var body interface{}
	if ok := parser.verifyNextToken(constants.TT_START_BODY); ok {
		parser.verifyNextToken(constants.TT_NEWLINE)
		body = parser.statements(constants.TT_END_BODY)
	} else {
		body = parser.statement()
	}
	return body
}

func (parser *Parser) varAccess() (*VarAccessNode, bool) {
	if parser.currentToken.Type_ != constants.TT_IDENTIFIER {
		return nil, false
	}

	varAccessNode := &VarAccessNode{
		Identifier: parser.currentToken.Value.(string),
	}
	parser.advance()
	return varAccessNode, true
}
