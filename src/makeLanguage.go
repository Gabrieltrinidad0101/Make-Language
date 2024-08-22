package src

import (
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer"

	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/function"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/booleans"

	"github.com/Gabrieltrinidad0101/Make-Language/src/customErrors"

	"github.com/Gabrieltrinidad0101/Make-Language/src/constants"

	"github.com/Gabrieltrinidad0101/Make-Language/src/api"
)

type MakeLanguage struct {
	*api.Api
	syntax   string
	filePath string
}

func NewMakeLanguage(syntax, filePath string) *MakeLanguage {
	return &MakeLanguage{
		syntax:   syntax,
		filePath: filePath,
		Api:      api.NewApi(),
	}
}

func (m MakeLanguage) Run() {
	conf, ok := lexer.ReadLanguageConfiguraction(m.syntax)
	conf.CustomOperators = m.CustomOperetor
	if !ok {
		return
	}

	input, ok := lexer.ReadFile(m.filePath)
	if !ok {
		return
	}

	customErrors.New(*input)

	lexer_ := lexer.NewLexer(input, conf)
	tokens, ok := lexer_.Tokens()

	if ok {
		return
	}

	parser_ := parser.NewParser(tokens, conf)
	ast, err := parser_.Parse()

	if err != nil {
		customErrors.InvalidSyntax(*parser_.CurrentToken, err.Error(), constants.STOP_EXECUTION)
		return
	}

	var languageContext_ = languageContext.NewContext(nil)
	languageContext_.Set("TRUE", &interpreteStructs.VarType{
		Value:      booleans.NewBoolean(true),
		IsConstant: true,
	})

	languageContext_.Set("FALSE", &interpreteStructs.VarType{
		Value:      booleans.NewBoolean(false),
		IsConstant: true,
	})

	functions := function.BuildFunctions(languageContext_, conf.Functions)

	for key, value := range functions {
		languageContext_.Set(key, &value)
	}

	for key, value := range m.Api.Functions {
		languageContext_.Set(key, &interpreteStructs.VarType{
			Value:      value,
			IsConstant: true,
		})
	}

	for key, value := range m.Api.Class {
		languageContext_.Set(key, &interpreteStructs.VarType{
			Value:      value,
			IsConstant: true,
		})
	}

	interprete_ := interprete.NewInterprete(ast, conf.Scope, m.Api, conf)
	interprete_.Run(languageContext_)
}
