package main

import (
	"makeLanguages/src/interprete"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
)

func main() {
	input := "if ( 2 == 1) { 1 + 2 } else { 4 }"
	lexer_ := lexer.NewLexer(&input)
	tokens, ok := lexer_.Tokens()

	if ok {
		return
	}

	parser_ := parser.NewParser(tokens)
	ast := parser_.Parse()

	interprete_ := interprete.NewInterprete(ast)
	interprete_.Run()
}
