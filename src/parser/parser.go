package parser

import (
	"fmt"
	"makeLanguages/src/constants"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/token"
	"slices"
)

type Parser struct {
	tokens       *[]token.Token
	idx          int
	CurrentToken *token.Token
	len          int
}

type BinOP struct {
	LeftNode  interface{}
	Operation token.Token
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
	IsConstant bool
}

type VarAccessNode struct {
	Identifier string
}

type ListNode struct {
	Nodes []interface{}
}

type NullNode struct{}

func NewParser(tokens *[]token.Token) *Parser {
	return &Parser{
		idx:          -1,
		tokens:       tokens,
		CurrentToken: &token.Token{},
		len:          len(*tokens),
	}
}

func (parser *Parser) advance() bool {
	parser.idx++
	if parser.idx >= parser.len {
		return false
	}
	*parser.CurrentToken = (*parser.tokens)[parser.idx]
	return true
}

func (parser *Parser) getToken(idx int) (*token.Token, bool) {
	if parser.idx >= parser.len {
		return nil, false
	}
	return &(*parser.tokens)[idx], true
}

func (parser *Parser) verifyNextToken(tokenType string) error {
	if parser.CurrentToken.Type_ != tokenType {
		return fmt.Errorf("Expect: %s", tokenType)
	}
	parser.advance()
	return nil
}

func (parser *Parser) binOP(callBack func() (interface{}, error), ops ...string) (interface{}, error) {
	leftNode, err := callBack()
	if err != nil {
		return nil, err
	}
	for slices.Contains[[]string](ops, parser.CurrentToken.Type_) {
		operation := *parser.CurrentToken
		parser.advance()
		rigthNode, err := callBack()
		if err != nil {
			return nil, err
		}
		leftNode = BinOP{
			LeftNode:  leftNode,
			Operation: operation,
			RigthNode: rigthNode,
		}
	}

	return leftNode, nil
}

func (parser *Parser) Parse() (interface{}, error) {
	parser.advance()
	ast, err := parser.expr()

	if err != nil {
		return nil, err
	}

	if parser.CurrentToken.Type_ != constants.EOF {
		return nil, fmt.Errorf("Expect +,-,*,/")
	}

	return ast, nil
}

func (parser *Parser) expr() (interface{}, error) {
	return parser.statements(constants.EOF)
}

func (parser *Parser) statements(tokenEnd string) (interface{}, error) {
	for parser.CurrentToken.Type_ == constants.TT_NEWLINE {
		parser.advance()
	}
	listNodes := ListNode{}
	ast, err := parser.statement()
	if err != nil {
		return nil, err
	}
	listNodes.Nodes = append(listNodes.Nodes, ast)

	for {
		thereIsANewLine := false
		for parser.CurrentToken.Type_ == constants.TT_NEWLINE {
			parser.advance()
			thereIsANewLine = true
		}

		if !thereIsANewLine {
			break
		}

		if parser.CurrentToken.Type_ == tokenEnd {
			parser.advance()
			break
		}

		ast, err := parser.statement()
		if err != nil {
			return nil, err
		}
		listNodes.Nodes = append(listNodes.Nodes, ast)
	}
	return listNodes, nil
}

func (parser *Parser) statement() (interface{}, error) {

	variableAndConst, err := parser.variableAndConst()
	if variableAndConst != nil || err != nil {
		return variableAndConst, err
	}

	return parser.compare()
}

func (parser *Parser) variableAndConst() (interface{}, error) {
	if err := parser.verifyNextToken(constants.TT_VAR); err == nil {
		identifier := parser.CurrentToken.Value
		err := parser.verifyNextToken(constants.TT_IDENTIFIER)
		if err != nil {
			return nil, err
		}
		err = parser.verifyNextToken(constants.TT_EQ)
		if err != nil {
			return nil, err
		}
		node, err := parser.compare()

		if err != nil {
			return nil, err
		}

		return VarAssignNode{
			Identifier: identifier.(string),
			Node:       node,
		}, nil
	}
	return nil, nil
}

func (parser *Parser) AndOr() (interface{}, error) {
	return parser.binOP(parser.compare, constants.TT_AND, constants.TT_AND)
}

func (parser *Parser) compare() (interface{}, error) {
	return parser.binOP(parser.plus, constants.TT_GT, constants.TT_GTE, constants.TT_GT, constants.TT_LT, constants.TT_LTE, constants.TT_EQE)
}

func (parser *Parser) plus() (interface{}, error) {
	return parser.binOP(parser.factor, constants.TT_PLUS, constants.TT_MINUS)
}

func (parser *Parser) factor() (interface{}, error) {
	return parser.binOP(parser.pow, constants.TT_MUL, constants.TT_DIV)
}

func (parser *Parser) pow() (interface{}, error) {
	return parser.binOP(parser.term, constants.TT_POW, constants.TT_SQUARE_ROOT)
}

func (parser *Parser) term() (interface{}, error) {
	nodeType := parser.CurrentToken.Type_

	if nodeType == constants.TT_PLUS || nodeType == constants.TT_MINUS {
		token, ok := parser.getToken(parser.idx + 1)
		if ok && token.Type_ == constants.TT_PLUS || token.Type_ == constants.TT_MINUS {
			return nil, fmt.Errorf("Error is necesery a ( between - and + simbols")
		}
		parser.advance()
		rigthNode, err := parser.term()

		if err != nil {
			return nil, err
		}

		unaryOP := UnaryOP{
			Operation: nodeType,
			RigthNode: rigthNode,
		}
		return unaryOP, nil
	}

	if nodeType == "number" {
		value := parser.CurrentToken.Value.(float64)
		number := numbers.NewNumbers(value)
		parser.advance()
		return number, nil
	}

	if parser.CurrentToken.Type_ == constants.TT_LPAREN {
		parser.advance()
		node, err := parser.statement()
		if err != nil {
			return nil, err
		}
		if !(parser.CurrentToken.Type_ == constants.TT_RPAREN) {
			return nil, fmt.Errorf("Expect )")
		}
		parser.advance()
		return node, nil
	}

	if ifNode, err := parser.if_(); ifNode != nil || err != nil {
		return ifNode, err
	}

	if varAccess, ok := parser.varAccess(); ok {
		return varAccess, nil
	}

	return nil, fmt.Errorf("")
}

func (parser *Parser) if_() (interface{}, error) {
	ifs := []*ConditionAndBody{}
	var elseNode interface{}
	if parser.CurrentToken.Type_ != constants.TT_IF {
		return nil, nil
	}

	parser.advance()

	node, err := parser.conditionAndBodyBase()

	if err != nil {
		return nil, err
	}

	ifs = append(ifs, node)

	for (*parser.CurrentToken).Type_ == constants.TT_ELIF {
		node, err := parser.conditionAndBodyBase()

		if err != nil {
			return nil, err
		}

		ifs = append(ifs, node)
	}
	if parser.CurrentToken.Type_ == constants.TT_ELSE {
		parser.advance()
		elseNode, err = parser.BodyBase()
	}

	return IfNode{
		Ifs:   ifs,
		Else_: elseNode,
	}, nil
}

func (parser *Parser) conditionAndBodyBase() (*ConditionAndBody, error) {
	err := parser.verifyNextToken(constants.TT_LPAREN)
	if err != nil {
		return nil, err
	}

	condition, err := parser.AndOr()
	if err != nil {
		return nil, err
	}

	err = parser.verifyNextToken(constants.TT_RPAREN)
	if err != nil {
		return nil, err
	}

	body, err := parser.BodyBase()

	if err != nil {
		return nil, err
	}

	return &ConditionAndBody{
		Condition: condition,
		Body:      body,
	}, err
}

func (parser *Parser) BodyBase() (interface{}, error) {
	if err := parser.verifyNextToken(constants.TT_START_BODY); err == nil {
		parser.verifyNextToken(constants.TT_NEWLINE)
		return parser.statements(constants.TT_END_BODY)
	}
	return parser.statement()
}

func (parser *Parser) varAccess() (*VarAccessNode, bool) {
	if parser.CurrentToken.Type_ != constants.TT_IDENTIFIER {
		return nil, false
	}

	varAccessNode := &VarAccessNode{
		Identifier: parser.CurrentToken.Value.(string),
	}
	parser.advance()
	return varAccessNode, true
}
