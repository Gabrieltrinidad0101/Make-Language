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

func NewParser(tokens *[]lexer.Token) *Parser {
	return &Parser{
		idx:          0,
		tokens:       tokens,
		currentToken: &lexer.Token{},
		len:          len(*tokens),
	}
}

func (parser *Parser) advance() {
	if parser.idx >= parser.len {
		return
	}
	*parser.currentToken = (*parser.tokens)[parser.idx]
	parser.idx++
}

func (parser *Parser) getToken(idx int) (*lexer.Token, bool) {
	if parser.idx >= parser.len {
		return nil, false
	}
	return &(*parser.tokens)[idx], true
}

func (parser *Parser) Parse() interface{} {
	return parser.expr()
}

func (parser *Parser) expr() interface{} {
	ast := parser.compare()
	return ast
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
	parser.advance()

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
		node := parser.expr()
		if !(parser.currentToken.Type_ == constants.TT_RPAREN) {
			panic("error )")
		}
		parser.advance()
		return node
	}
	panic("Error")
}

func (parser *Parser) binOP(callBack func() interface{}, ops ...string) interface{} {
	leftNode := callBack()
	for slices.Contains[[]string](ops, parser.currentToken.Type_) {
		operation := *parser.currentToken
		rigthNode := callBack()
		leftNode = BinOP{
			LeftNode:  leftNode,
			Operation: operation,
			RigthNode: rigthNode,
		}
	}

	return leftNode
}
