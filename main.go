package main

import (
	"makeLanguages/src/interprete"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
	"strings"
)

var languageContext_ = languageContext.NewContext(nil)

func main() {

	input, ok := lexer.ReadFile("./main.makeLanguage")
	if !ok {
		return
	}

	lines := strings.Split(*input, "\n")

	conf, ok := lexer.ReadLanguageConfiguraction("./conf.json")
	if !ok {
		return
	}

	for _, line := range lines {

		lexer_ := lexer.NewLexer(&line, conf)
		tokens, ok := lexer_.Tokens()

		if ok {
			return
		}

		parser_ := parser.NewParser(tokens)
		ast := parser_.Parse()

		interprete_ := interprete.NewInterprete(ast, &languageContext_)
		interprete_.Run()
	}
}
