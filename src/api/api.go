package api

import "makeLanguages/src/interprete/interpreteStructs"

type operatorFunc func(interpreteStructs.IBaseElement, interpreteStructs.IBaseElement) interface{}

type Api struct {
	UnOperetor map[string]interface{}
	Operetor   map[string]operatorFunc
	Functions  map[string]interface{}
	Class      map[string]interface{}
}

func NewApi() *Api {
	return &Api{
		UnOperetor: make(map[string]interface{}),
		Operetor:   make(map[string]operatorFunc),
		Functions:  make(map[string]interface{}),
		Class:      make(map[string]interface{}),
	}
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
	return result, true
}

func CustomFunctions() {

}

func CustomClass() {

}
