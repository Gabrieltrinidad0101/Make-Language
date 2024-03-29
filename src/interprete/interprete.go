package interprete

import (
	"fmt"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/parser"
	"reflect"
)

type Interprete struct {
	ast         interface{}
	currentNode interface{}
	context     *languageContext.Context
}

func NewInterprete(ast interface{}, context *languageContext.Context) *Interprete {
	return &Interprete{
		ast:     ast,
		context: context,
	}
}

func (interprete *Interprete) Run() {
	value := interprete.call(interprete.ast)
	fmt.Println(value)
}

func (interprete *Interprete) getMethodName(node interface{}) string {
	method := reflect.TypeOf(node)
	if method.Kind() == reflect.Ptr {
		return method.Elem().Name()
	} else {
		return method.Name()
	}
}

func (interprete *Interprete) call(node interface{}) interface{} {
	methodName := interprete.getMethodName(node)
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

func (interprete Interprete) VarAssignNode(node interface{}) interface{} {
	varAssignNode := node.(parser.VarAssignNode)
	if _, ok := interprete.context.Get(varAssignNode.Identifier); ok && varAssignNode.IsConstant {
		panic("Const " + varAssignNode.Identifier)
	}
	result := interprete.call(varAssignNode.Node)
	interprete.context.Set(varAssignNode.Identifier, result)
	return parser.NullNode{}
}

func (interprete Interprete) VarAccessNode(node interface{}) interface{} {
	varAccessNode := node.(*parser.VarAccessNode)
	valueNode, ok := interprete.context.Get(varAccessNode.Identifier)
	if !ok {
		panic("Variable is undefined")
	}
	return interprete.call(valueNode)
}

func (interprete *Interprete) UnaryOP(node interface{}) *numbers.Number {
	unaryOP := node.(parser.UnaryOP)
	number := interprete.call(unaryOP.RigthNode).(*numbers.Number)

	if unaryOP.Operation == "MINUS" {
		number.Value *= -1
	}
	return number
}

func (interprete *Interprete) IfNode(node interface{}) interface{} {
	ifNode := node.(parser.IfNode)

	for _, if_ := range ifNode.Ifs {
		conditionInterface := interprete.call(if_.Condition)

		if interprete.getMethodName(conditionInterface) != "Boolean" {
			panic("Error if expression need to a condition")
		}

		condition := conditionInterface.(*booleans.Boolean)

		if condition.Value {
			node := interprete.call(if_.Body)
			return node
		}
	}

	if ifNode.Else_ != nil {
		node := interprete.call(ifNode.Else_)
		return node
	}

	return parser.NullNode{}
}

func (interprete *Interprete) WhileNode(node interface{}) interface{} {
	whileNode := node.(parser.WhileNode)

	for {
		boolean := interprete.call(whileNode.Condition).(*booleans.Boolean)
		if !boolean.Value {
			break
		}
		fmt.Print(interprete.call(whileNode.Body))
	}

	return parser.NullNode{}
}

func (interprete *Interprete) ListNode(node interface{}) interface{} {
	listNode := node.(parser.ListNode)
	for _, node := range listNode.Nodes {
		fmt.Println(interprete.call(node))
	}
	return 1
}

func (interprete *Interprete) Number(node interface{}) *numbers.Number {
	return node.(*numbers.Number)
}
