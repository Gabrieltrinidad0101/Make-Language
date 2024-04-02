package interprete

import (
	"fmt"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/function"
	"makeLanguages/src/features/numbers"
	interpreteStructs "makeLanguages/src/interprete/structs"
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
	context.Set(varAssignNode.Identifier, interpreteStructs.VarType{
		Value:      result,
		IsConstant: varAssignNode.IsConstant,
	})
	return parser.NullNode{}
}

func (interprete Interprete) UpdateVariableNode(node interface{}, context *languageContext.Context) interface{} {
	updateVariableNode := node.(parser.UpdateVariableNode)
	varType, ok := context.Get(updateVariableNode.Identifier)
	if !ok {
		panic("The variable no exist" + updateVariableNode.Identifier)
	}
	varTypeNode := varType.(interpreteStructs.VarType)
	if varTypeNode.IsConstant {
		panic("Const")
	}

	result := interprete.call(updateVariableNode.Node, context)

	varTypeNode.Value = result

	ok = context.Update(updateVariableNode.Identifier, varTypeNode)
	if !ok {
		panic("The variable no exist" + updateVariableNode.Identifier)
	}
	return parser.NullNode{}
}

func (interprete Interprete) VarAccessNode(node interface{}, context *languageContext.Context) interface{} {
	varAccessNode := node.(*parser.VarAccessNode)
	varType, ok := context.Get(varAccessNode.Identifier)
	if !ok {
		panic("Variable is undefined " + varAccessNode.Identifier)
	}
	valueNode, ok := varType.(interpreteStructs.VarType)
	if !ok {
		fmt.Print("")
	}
	return interprete.call(valueNode.Value, context)
}

func (interprete *Interprete) FuncNode(node interface{}, context *languageContext.Context) interface{} {
	funcNode := node.(*parser.FuncNode)
	newContext := languageContext.NewContext(context)

	for _, param := range *funcNode.Params {
		newContext.Set(param.Value.(string), interpreteStructs.VarType{
			Value: parser.NullNode{},
		})
	}

	func_ := function.Function{
		Body:    funcNode.Body,
		Context: &newContext,
		Params:  funcNode.Params,
	}

	context.Set(funcNode.Name, interpreteStructs.VarType{
		Value:      func_,
		IsConstant: true,
	})
	return parser.NullNode{}
}

func (interprete *Interprete) CallFuncNode(node interface{}, context *languageContext.Context) interface{} {
	callFuncNode := node.(*parser.CallFuncNode)
	func_, ok := context.Get(callFuncNode.Name)
	if !ok {
		panic(callFuncNode.Name)
	}
	varType := func_.(interpreteStructs.VarType)

	funcNode := varType.Value.(function.Function)

	var params []interface{}

	for _, param := range *callFuncNode.Params {
		params = append(params, interprete.call(param, context))
	}

	funcNodeBody, hasACustomExecute := funcNode.Execute(&params)
	if hasACustomExecute {
		return funcNodeBody
	}
	interprete.call(funcNode.GetBody(), funcNode.GetContext())
	return parser.NullNode{}
}

func (interprete *Interprete) UnaryOP(node interface{}, context *languageContext.Context) *numbers.Number {
	unaryOP := node.(parser.UnaryOP)
	number := interprete.call(unaryOP.RigthNode, context).(*numbers.Number)

	if unaryOP.Operation == "MINUS" {
		number.Value *= -1
	}
	if unaryOP.Operation == "PLUS1" {
		number.Value += 1
	}
	if unaryOP.Operation == "MINUS1" {
		number.Value -= 1
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

func (interprete *Interprete) ForNode(node interface{}, context *languageContext.Context) interface{} {
	forNode := node.(parser.ForNode)

	for {
		interprete.call(forNode.Expr1, context)
		condition := interprete.call(forNode.Condition, context)
		coditionNode := condition.(*booleans.Boolean)
		if !coditionNode.Value {
			break
		}
		interprete.call(forNode.Expr2, context)
		interprete.call(forNode.Body, context)
	}

	return parser.NullNode{}
}

func (interprete *Interprete) ListNode(node interface{}, context *languageContext.Context) interface{} {
	listNode := node.(parser.ListNode)
	for _, node := range listNode.Nodes {
		interprete.call(node, context)
	}
	return 1
}

func (interprete *Interprete) Number(node interface{}, context *languageContext.Context) *numbers.Number {
	return node.(*numbers.Number)
}
