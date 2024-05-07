package api

import (
	"makeLanguages/src/features/class"
	"makeLanguages/src/features/function"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
)

type operatorFunc func(interpreteStructs.IBaseElement, interpreteStructs.IBaseElement) interface{}

type Api struct {
	unOperetor     map[string]interface{}
	operetor       map[string]operatorFunc
	CustomOperetor map[string]string
	Functions      map[string]*function.BaseFunction
	Class          map[string]class.Class
}

func NewApi() *Api {
	return &Api{
		unOperetor:     make(map[string]interface{}),
		operetor:       make(map[string]operatorFunc),
		CustomOperetor: map[string]string{},
		Functions:      make(map[string]*function.BaseFunction),
		Class:          map[string]class.Class{},
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

func (api *Api) AddClass(name string, methods map[string]func(params *[]interpreteStructs.IBaseElement) interface{}) {
	customClass := class.NewBuildClass(languageContext.NewContext(nil))
	for key, value := range methods {
		customClass.AddMethod(key, value)
	}
	api.Class[name] = class.Class{
		Context: customClass.Context,
		Name:    name,
	}
}
