package api

import (
	"makeLanguages/src/features/function"
	"makeLanguages/src/interprete/interpreteStructs"
)

type operatorFunc func(interpreteStructs.IBaseElement, interpreteStructs.IBaseElement) interface{}

type Api struct {
	unOperetor     map[string]interface{}
	operetor       map[string]operatorFunc
	CustomOperetor map[string]string
	Functions      map[string]*function.BaseFunction
	class          map[string]interface{}
}

func NewApi() *Api {
	return &Api{
		unOperetor:     make(map[string]interface{}),
		operetor:       make(map[string]operatorFunc),
		CustomOperetor: map[string]string{},
		Functions:      make(map[string]*function.BaseFunction),
		class:          make(map[string]interface{}),
	}
}

func (api *Api) AddOperetor(simbol string, callBack operatorFunc) {
	api.CustomOperetor[simbol] = simbol
	api.operetor[simbol] = callBack
}

func (api *Api) Call(tokenName string, value1 interpreteStructs.IBaseElement, value2 interpreteStructs.IBaseElement) (interface{}, bool) {
	callBack, ok := api.operetor[tokenName]
	if !ok {
		return nil, ok
	}
	result := callBack(value1, value2)
	return result, true
}

func (api *Api) AddFunction(name string, function_ func(params *[]interpreteStructs.IBaseElement) interface{}) {
	baseFunction := function.NewBaseFunction(nil, name, function_)
	api.Functions[name] = baseFunction
}

func CustomClass() {

}
