package main

import (
	"makeLanguages/src/customErrors"
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

	customErrors.New(*input)

	lexer_ := lexer.NewLexer(input, conf)
	tokens, ok := lexer_.Tokens()

	if ok {
		return
	}

	parser_ := parser.NewParser(tokens)
	ast, err := parser_.Parse()

	if err != nil {
		customErrors.InvalidSyntax(*parser_.CurrentToken, err.Error())
		return
	}

	var languageContext_ = languageContext.NewContext(nil)
	languageContext_.Set("TRUE", true)
	languageContext_.Set("FALSE", false)

	interprete_ := interprete.NewInterprete(ast, &languageContext_)
	interprete_.Run()
}
