package parser

import (
	"makeLanguages/src/lexer"
	"slices"
)

type Parser struct {
	tokens       *[]lexer.Token
	idx          int
	currentToken *lexer.Token
}

type Number struct {
	value int64
}

type BinOP struct {
	leftNode  Number
	operation lexer.Token
	rigthNode     Number
}

func NewParser() *Parser {
	return &Parser{
		idx: 0,
	}
}

func (parser *Parser) advance() {
	*parser.currentToken = (*parser.tokens)[parser.idx]
	parser.idx++
}

func (parser *Parser) parse() {

}

func (parser *Parser) expr() {

}

func (parser *Parser) term() interface{} {
	if parser.currentToken.Type_ == "number" {
		return Number{
			value: parser.currentToken.Value.(int64),
		}
	}
}

func (parser *Parser) binOP(ops... string) {
	var leftNode interface{} 
	leftNode = parser.term()

	for slices.Contains[[]string](ops,parser.currentToken.Type_) {
		parser.advance()
		operation := parser.currentToken
		rigthNode := parser.term()
		leftNode = interface{}{
			leftNode: leftNode,
			operation: operation,
			rigthNode: rigthNode
		}

	}
	
}
