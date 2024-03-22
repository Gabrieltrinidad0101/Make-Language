package main

import (
	"makeLanguages/src/interprete"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
)

func main() {
	conf, ok := lexer.ReadLanguageConfiguraction("./conf.json")
	if !ok {
		return
	}

	input, ok := lexer.ReadFile("./main.makeLanguage")
	if !ok {
		return
	}

	lexer_ := lexer.NewLexer(input, conf)
	tokens, ok := lexer_.Tokens()

	if ok {
		return
	}

	parser_ := parser.NewParser(tokens)
	ast := parser_.Parse()

	interprete_ := interprete.NewInterprete(ast)
	interprete_.Run()
}
