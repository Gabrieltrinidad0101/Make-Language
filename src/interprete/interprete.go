package interprete

import (
	"fmt"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/function"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/languageContext"
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

func (interprete *Interprete) Run(context *languageContext.Context) interface{} {
	value := interprete.call(interprete.ast, context)
	return value
}

func (interprete *Interprete) getMethodName(node interface{}) string {
	method := reflect.TypeOf(node)
	if method.Kind() == reflect.Ptr {
		return method.Elem().Name()
	} else {
		return method.Name()
	}
}

func (interprete *Interprete) call(node interface{}, context *languageContext.Context) interface{} {
	methodName := interprete.getMethodName(node)
	return interprete.callMethod(interprete, methodName, node, context)
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

func (interprete *Interprete) BinOP(node interface{}, context *languageContext.Context) interface{} {
	binOP := node.(parser.BinOP)
	nodeLeft := interprete.call(binOP.LeftNode, context)
	nodeRigth := interprete.call(binOP.RigthNode, context)
	newNode := interprete.callMethod(nodeLeft, binOP.Operation.Type_, nodeRigth)
	return newNode
}

func (interprete Interprete) VarAssignNode(node interface{}, context *languageContext.Context) interface{} {
	varAssignNode := node.(parser.VarAssignNode)
	if _, ok := context.Get(varAssignNode.Identifier); ok && varAssignNode.IsConstant {
		panic("Const " + varAssignNode.Identifier)
	}
	result := interprete.call(varAssignNode.Node, context)
	context.Set(varAssignNode.Identifier, result)
	return parser.NullNode{}
}

func (interprete Interprete) VarAccessNode(node interface{}, context *languageContext.Context) interface{} {
	varAccessNode := node.(*parser.VarAccessNode)
	valueNode, ok := context.Get(varAccessNode.Identifier)
	if !ok {
		panic("Variable is undefined")
	}
	return interprete.call(valueNode, context)
}

func (interprete *Interprete) CallFuncNode(node interface{}, context *languageContext.Context) {
	callFuncNode := node.(*parser.CallFuncNode)
	func_, ok := context.Get(callFuncNode.Name)
	if !ok {
		panic(callFuncNode.Name)
	}
	func_.(function.Function).SetParams(callFuncNode.Params)
}

func (interprete *Interprete) UnaryOP(node interface{}, context *languageContext.Context) *numbers.Number {
	unaryOP := node.(parser.UnaryOP)
	number := interprete.call(unaryOP.RigthNode, context).(*numbers.Number)

	if unaryOP.Operation == "MINUS" {
		number.Value *= -1
	}
	return number
}

func (interprete *Interprete) IfNode(node interface{}, context *languageContext.Context) interface{} {
	ifNode := node.(parser.IfNode)

	for _, if_ := range ifNode.Ifs {
		conditionInterface := interprete.call(if_.Condition, context)

		if interprete.getMethodName(conditionInterface) != "Boolean" {
			panic("Error if expression need to a condition")
		}

		condition := conditionInterface.(*booleans.Boolean)

		if condition.Value {
			node := interprete.call(if_.Body, context)
			return node
		}
	}

	if ifNode.Else_ != nil {
		node := interprete.call(ifNode.Else_, context)
		return node
	}

	return parser.NullNode{}
}

func (interprete *Interprete) WhileNode(node interface{}, context *languageContext.Context) interface{} {
	whileNode := node.(parser.WhileNode)

	for {
		boolean := interprete.call(whileNode.Condition, context).(*booleans.Boolean)
		if !boolean.Value {
			break
		}
		fmt.Print(interprete.call(whileNode.Body, context))
	}

	return parser.NullNode{}
}

func (interprete *Interprete) FuncNode(node interface{}, context *languageContext.Context) interface{} {
	funcNode := node.(*parser.FuncNode)
	newContext := languageContext.NewContext(context)

	for _, param := range *funcNode.Params {
		newContext.Set(param.Value.(string), parser.NullNode{})
	}
	context.Set(funcNode.Name, function.Function{
		Body:    funcNode.Body,
		Context: newContext,
		Params:  funcNode.Params,
	})
	return parser.NullNode{}
}

func (interprete *Interprete) ListNode(node interface{}, context *languageContext.Context) interface{} {
	listNode := node.(parser.ListNode)
	for _, node := range listNode.Nodes {
		fmt.Println(interprete.call(node, context))
	}
	return 1
}

func (interprete *Interprete) Number(node interface{}) *numbers.Number {
	return node.(*numbers.Number)
}
