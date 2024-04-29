package src

import (
	"makeLanguages/src/constants"
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/function"
	"makeLanguages/src/interprete"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
)

func MakeLanguage(syntax string, filePath string) {
	conf, ok := lexer.ReadLanguageConfiguraction(syntax)
	if !ok {
		return
	}

	input, ok := lexer.ReadFile(filePath)
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
		customErrors.InvalidSyntax(*parser_.CurrentToken, err.Error(), constants.STOP_EXECUTION)
		return
	}

	var languageContext_ = languageContext.NewContext(nil)
	languageContext_.Set("TRUE", interpreteStructs.VarType{
		Value:      booleans.NewBoolean(true),
		IsConstant: true,
	})

	languageContext_.Set("FALSE", interpreteStructs.VarType{
		Value:      booleans.NewBoolean(false),
		IsConstant: true,
	})

	functions := function.BuildFunctions(conf.Functions)

	for key, value := range functions {
		languageContext_.Set(key, value)
	}

	interprete_ := interprete.NewInterprete(ast, conf.Scope)
	interprete_.Run(languageContext_)
}
