package main

import (
	"makeLanguages/src/customErrors"
	"makeLanguages/src/interprete"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
	"testing"
)

func TestLanguage(t *testing.T) {
	code := "var a = 1 + 1\n" +
		"var b = a * 3\n" +
		"var c = if(a > b) 1 else 4\n" +
		"var d = if(a < b) 2 else 5\n" +
		"a\n" +
		"b\n" +
		"c\n" +
		"d\n"

	conf, ok := lexer.ReadLanguageConfiguraction("./conf.json")
	if !ok {
		return
	}

	customErrors.New(code)

	lexer_ := lexer.NewLexer(&code, conf)
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

	interprete_ := interprete.NewInterprete(ast)
	res := interprete_.Run(&languageContext_)
	nodeList := res.(parser.ListNode)

	require

}
