package parser

import (
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/lexer"
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

func (parser *Parser) Parse() interface{} {
	return parser.expr()
}

func (parser *Parser) expr() interface{} {
	ast := parser.compare()
	return ast
}

func (parser *Parser) compare() interface{} {
	return parser.binOP(parser.plus, "GT", "GTE", "LT", "LTE", "EQE")
}

func (parser *Parser) plus() interface{} {
	ast := parser.binOP(parser.factor, "PLUS", "MINUS")
	return ast
}

func (parser *Parser) factor() interface{} {
	return parser.binOP(parser.pow, "MUL", "DIV")
}

func (parser *Parser) pow() interface{} {
	return parser.binOP(parser.term, "POW", "SQUARE_ROOT")
}

func (parser *Parser) term() interface{} {
	parser.advance()
	if parser.currentToken.Type_ == "number" {
		value := parser.currentToken.Value.(float64)
		number := numbers.NewNumbers(value)
		parser.advance()
		return number
	}
	if parser.currentToken.Type_ == "LPAREN" {
		node := parser.expr()
		if !(parser.currentToken.Type_ == "RPAREN") {
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
