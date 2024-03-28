package main

import (
	"makeLanguages/src/interprete"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
)

func main() {

	input, ok := lexer.ReadFile("./main.makeLanguage")
	if !ok {
		return
	}

	conf, ok := lexer.ReadLanguageConfiguraction("./conf.json")
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

	var languageContext_ = languageContext.NewContext(nil)
	languageContext_.Set("TRUE", true)
	languageContext_.Set("FALSE", false)

	interprete_ := interprete.NewInterprete(ast, &languageContext_)
	interprete_.Run()
}
