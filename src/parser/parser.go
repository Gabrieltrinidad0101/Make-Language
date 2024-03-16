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
	ast := parser.binOP(parser.factor, "PLUS", "MINUS")
	return ast
}

func (parser *Parser) factor() interface{} {
	return parser.binOP(parser.term, "MUL", "DIV")
}

func (parser *Parser) term() interface{} {
	parser.advance()
	if parser.currentToken.Type_ == "number" {
		value := parser.currentToken.Value.(int)
		number := numbers.NewNumbers(value)
		parser.advance()
		return number
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
