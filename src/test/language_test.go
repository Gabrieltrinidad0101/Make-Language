package test

import (
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/function"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/interprete"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer"
	"makeLanguages/src/parser"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BaseInterprete(languageContext_ *languageContext.Context, filePath, confPath string) *languageContext.Context {
	input, ok := lexer.ReadFile(filePath)
	if !ok {
		return nil
	}

	conf, ok := lexer.ReadLanguageConfiguraction(confPath)
	if !ok {
		return nil
	}

	customErrors.New(*input)

	lexer_ := lexer.NewLexer(input, conf)
	tokens, ok := lexer_.Tokens()

	if ok {
		return nil
	}

	parser_ := parser.NewParser(tokens)
	ast, err := parser_.Parse()

	if err != nil {
		customErrors.InvalidSyntax(*parser_.CurrentToken, err.Error())
		return nil
	}

	interprete_ := interprete.NewInterprete(ast)
	interprete_.Run(languageContext_)

	return languageContext_
}
func getLanguageContext(confPath string) *languageContext.Context {
	conf, ok := lexer.ReadLanguageConfiguraction(confPath)
	if !ok {
		return nil
	}

	functions := function.BuildFunctions(conf.Functions)

	languageContext_ := languageContext.NewContext(nil)
	languageContext_.Set("TRUE", true)
	languageContext_.Set("FALSE", false)

	for key, value := range functions {
		languageContext_.Set(key, value)
	}
	return &languageContext_
}

type Print struct {
	assert *assert.Assertions
	call   int
}

func (print Print) Execute(params *[]interface{}) (interface{}, bool) {
	number := (*params)[0].(numbers.Number)
	if print.call == 0 {
		print.assert.Equal(number.Value, 2)
	} else {
		print.assert.Equal(number.Value, 3)
	}
	return nil, true
}

func TestVariablesAndIfs(t *testing.T) {
	assert := assert.New(t)
	context := getLanguageContext("./conf.json")
	context.Set("print", Print{
		assert: assert,
	})
	context = BaseInterprete(context, "./main.makeLanguage", "./conf.json")

	a, ok := context.Get("a")
	assert.True(ok)
	b, ok := context.Get("b")
	assert.True(ok)
	c, ok := context.Get("c")
	assert.True(ok)
	d, ok := context.Get("d")
	assert.True(ok)
	assert.Equal(a.(numbers.Number), 1)
	assert.Equal(b.(numbers.Number), 2)
	assert.Equal(c.(numbers.Number), 2)
	assert.Equal(d.(numbers.Number), 3)
}
