package api

import "makeLanguages/src/interprete/interpreteStructs"

type operatorFunc func(interpreteStructs.IBaseElement, interpreteStructs.IBaseElement) interface{}

type Api struct {
	UnOperetor map[string]interface{}
	Operetor   map[string]operatorFunc
	Functions  map[string]interface{}
	Class      map[string]interface{}
}

func (api *Api) AddOperetor(tokenName string, callBack operatorFunc) {
	api.Operetor[tokenName] = callBack
}

func (api *Api) Call(tokenName string, value1 interpreteStructs.IBaseElement, value2 interpreteStructs.IBaseElement) (interface{}, bool) {
	callBack, ok := api.Operetor[tokenName]
	if !ok {
		return nil, ok
	}
	result := callBack(value1, value2)
	return result, false
}

func CustomFunctions() {

}

func CustomClass() {

}
