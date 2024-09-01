package test

import (
	"testing"

	"github.com/Gabrieltrinidad0101/Make-Language/src/api"
	"github.com/Gabrieltrinidad0101/Make-Language/src/customErrors"
	"github.com/Gabrieltrinidad0101/Make-Language/src/features/booleans"
	"github.com/Gabrieltrinidad0101/Make-Language/src/features/function"
	"github.com/Gabrieltrinidad0101/Make-Language/src/features/numbers"
	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete"
	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"
	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"
	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer"
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser"
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser/parserStructs"
	"github.com/stretchr/testify/assert"
)

func BaseInterprete(languageContext_ *languageContext.Context, filePath, confPath string) (*languageContext.Context, error) {
	input, err := lexer.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	conf, err := lexer.ReadLanguageConfiguraction(confPath)
	if err != nil {
		return nil, err
	}

	customErrors.New(*input)

	lexer_ := lexer.NewLexer(input, conf)
	tokens, err := lexer_.Tokens()

	if err != nil {
		return nil, err
	}

	parser_ := parser.NewParser(tokens, conf)
	ast, err := parser_.Parse()

	if err != nil {
		return nil, customErrors.InvalidSyntax(*parser_.CurrentToken, err.Error())
	}

	api_ := api.NewApi()
	interprete_ := interprete.NewInterprete(ast, conf.Scope, api_, conf)
	interprete_.Run(languageContext_)

	return languageContext_, nil
}
func getLanguageContext(confPath string) (*languageContext.Context, error) {
	conf, err := lexer.ReadLanguageConfiguraction(confPath)
	if err != nil {
		return nil, err
	}

	languageContext_ := languageContext.NewContext(nil)
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
	return languageContext_, nil
}

var call = 0

type Print struct {
	function.BaseFunction
	assert *assert.Assertions
}

func (p Print) GetContext() *languageContext.Context {
	return p.Context
}

func (p Print) CanChangeContextParent() bool {
	return true
}

func (print Print) Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error) {
	number := (*params)[0].(interpreteStructs.IBaseElement)
	if call == 0 {
		print.assert.Equal(float64(2), number.GetValue())
	} else if call == 1 {
		print.assert.Equal(float64(3), number.GetValue())
	} else if call == 2 {
		print.assert.Equal(float64(10), number.GetValue())
	} else if call == 3 {
		print.assert.Equal(float64(0), number.GetValue())
	} else if call == 4 {
		print.assert.Equal(float64(20), number.GetValue())
	} else if call == 5 {
		print.assert.Equal(float64(1), number.GetValue())
	} else if call == 6 {
		print.assert.Equal("hello world", number.GetValue())
	} else if call == 7 {
		print.assert.Equal("hell0 w0rld", number.GetValue())
	} else if call == 8 {
		print.assert.Equal("HELLO WORLD", number.GetValue())
	} else if call == 9 {
		print.assert.Equal("test1", number.GetValue())
	}

	call++
	return parserStructs.NullNode{}, true, nil
}

func TestVariablesAndIfs(t *testing.T) {
	assert := assert.New(t)
	context, err := getLanguageContext("./conf.json")
	assert.Nil(err)
	context.Set("print", &interpreteStructs.VarType{
		Value: Print{
			assert: assert,
			BaseFunction: function.BaseFunction{
				Context: &languageContext.Context{},
			},
		},
		IsConstant: true,
	})
	context, err = BaseInterprete(context, "./main.makeLanguage", "./conf.json")
	assert.Nil(err)

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
