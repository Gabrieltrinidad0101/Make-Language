package main

import (
	"makeLanguages/src/interprete"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
)

func main() {
	input := "10 ^ 2 == 2 "
	lexer_ := lexer.NewLexer(&input)
	tokens := lexer_.Tokens()

	parser_ := parser.NewParser(tokens)
	ast := parser_.Parse()

	interprete_ := interprete.NewInterprete(ast)
	interprete_.Run()
}
