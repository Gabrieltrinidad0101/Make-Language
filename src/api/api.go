package api

import (
	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/function"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/class"
)

type operatorFunc func(interpreteStructs.IBaseElement, interpreteStructs.IBaseElement) interface{}

type Methods map[string]func(params *[]interpreteStructs.IBaseElement) interface{}

type CustomClassValues struct {
	Methods    Methods
	Properties map[string]interpreteStructs.IBaseElement
}

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
	baseFunction := function.NewBaseFunction(languageContext.NewContext(nil), name, function_, true)
	api.Functions[name] = baseFunction
}

func (api *Api) AddClass(name string, customClassValues CustomClassValues) {
	customClass := class.NewBuildClass(languageContext.NewContext(nil))
	for key, value := range customClassValues.Methods {
		customClass.AddMethod(key, value)
	}
	for key, value := range customClassValues.Properties {
		customClass.AddProperty(key, value)
	}
	api.Class[name] = class.Class{
		Context: customClass.Context,
		Name:    name,
	}
}
