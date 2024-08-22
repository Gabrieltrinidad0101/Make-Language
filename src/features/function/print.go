package function

import (
	"fmt"

	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser/parserStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"
)

type Print struct {
	BaseFunction
}

func NewPrint(context *languageContext.Context) *Print {
	return &Print{
		BaseFunction: BaseFunction{
			Context: languageContext.NewContext(context),
		},
	}
}
func (func_ Print) GetContext() *languageContext.Context {
	return func_.Context
}

func (func_ Print) CanChangeContextParent() bool {
	return true
}

func (func_ Print) Execute(params *[]interpreteStructs.IBaseElement) (interface{}, bool, error) {
	for _, param := range *params {
		fmt.Println(param.GetValue())
	}
	hasACustomExecute := true
	return parserStructs.NullNode{}, hasACustomExecute, nil
}
