package test

import (
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/function"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/interprete"
	interpreteStructs "makeLanguages/src/interprete/structs"
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
	languageContext_.Set("TRUE", interpreteStructs.VarType{
		Value:      true,
		IsConstant: true,
	})
	languageContext_.Set("FALSE", interpreteStructs.VarType{
		Value:      true,
		IsConstant: true,
	})

	for key, value := range functions {
		languageContext_.Set(key, value)
	}
	return &languageContext_
}

var call = 0

type Print struct {
	function.BaseFunction
	assert *assert.Assertions
}

func (print Print) Execute(params *[]interface{}) (interface{}, bool) {
	number := (*params)[0].(*numbers.Number)
	if call == 0 {
		print.assert.Equal(number.Value, float64(2))
	} else if call == 1 {
		print.assert.Equal(number.Value, float64(3))
	} else if call == 2 {
		print.assert.Equal(number.Value, float64(11))
	} else if call == 3 {
		print.assert.Equal(number.Value, float64(0))
	} else if call == 4 {
		print.assert.Equal(number.Value, float64(20))
	}

	call++
	return nil, true
}

func TestVariablesAndIfs(t *testing.T) {
	assert := assert.New(t)
	context := getLanguageContext("./conf.json")
	context.Set("print", interpreteStructs.VarType{
		Value: Print{
			assert: assert,
		},
		IsConstant: true,
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
	assert.Equal(a.Value.(*numbers.Number).Value, float64(1))
	assert.Equal(b.Value.(*numbers.Number).Value, float64(2))
	assert.Equal(c.Value.(*numbers.Number).Value, float64(2))
	assert.Equal(d.Value.(*numbers.Number).Value, float64(3))
}
