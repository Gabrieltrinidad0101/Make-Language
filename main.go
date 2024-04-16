package main

import (
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/function"
	"makeLanguages/src/interprete"
	"makeLanguages/src/interprete/structs"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
)

func main() {
	input, ok := lexer.ReadFile("./main.mkL")
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
	languageContext_.Set("TRUE", structs.VarType{
		Value:      booleans.NewBoolean(true),
		IsConstant: true,
	})
	languageContext_.Set("FALSE", structs.VarType{
		Value:      booleans.NewBoolean(false),
		IsConstant: true,
	})

	functions := function.BuildFunctions(conf.Functions)

	for key, value := range functions {
		languageContext_.Set(key, value)
	}

	interprete_ := interprete.NewInterprete(ast)
	interprete_.Run(languageContext_)
}
