package interprete

import (
	"fmt"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/parser"
	"reflect"
)

type Interprete struct {
	ast         interface{}
	currentNode interface{}
}

func NewInterprete(ast interface{}) *Interprete {
	return &Interprete{
		ast: ast,
	}
}

func (interprete *Interprete) Run() {
	value := interprete.call(interprete.ast)
	fmt.Println(value)
}

func getMethodName(node interface{}) string {
	method := reflect.TypeOf(node)
	if method.Kind() == reflect.Ptr {
		return method.Elem().Name()
	} else {
		return method.Name()
	}
}

func (interprete *Interprete) call(node interface{}) interface{} {
	methodName := getMethodName(node)
	return interprete.callMethod(interprete, methodName, node)
}

func (interprete *Interprete) callMethod(object interface{}, methodName string, values ...interface{}) interface{} {
	method := reflect.ValueOf(object).MethodByName(methodName)
	var params []reflect.Value
	for _, value := range values {
		params = append(params, reflect.ValueOf(value))
	}

	returnValue := method.Call(params)
	return returnValue[0].Interface()
}

func (interprete *Interprete) BinOP(node interface{}) interface{} {
	binOP := node.(parser.BinOP)
	nodeLeft := interprete.call(binOP.LeftNode)
	nodeRigth := interprete.call(binOP.RigthNode)
	newNode := interprete.callMethod(nodeLeft, binOP.Operation.Type_, nodeRigth)
	return newNode
}

func (interprete *Interprete) UnaryOP(node interface{}) *numbers.Number {
	unaryOP := node.(parser.UnaryOP)
	number := interprete.call(unaryOP.RigthNode).(*numbers.Number)

	if unaryOP.Operation == "MINUS" {
		number.Value *= -1
	}
	return number
}

func (interprete *Interprete) Number(node interface{}) *numbers.Number {
	return node.(*numbers.Number)
}
