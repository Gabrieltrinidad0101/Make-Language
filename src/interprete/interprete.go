package interprete

import (
	"makeLanguages/src/constants"
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/class"
	"makeLanguages/src/features/function"
	"makeLanguages/src/features/numbers"
	interpreteStructs "makeLanguages/src/interprete/structs"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/parser/parserStructs"
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

func (interprete *Interprete) stopExecute(node interface{}) string {
	if interprete.getMethodName(node) == "ContinueNode" ||
		interprete.getMethodName(node) == "BREAK" ||
		interprete.getMethodName(node) == "RETURN" {
		return interprete.getMethodName(node)
	}

	return ""

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

func (interprete *Interprete) ClassNode(node interface{}, context *languageContext.Context) interface{} {
	classNode := node.(*parserStructs.ClassNode)
	newContext := languageContext.NewContext(context)

	class_ := class.Class{
		Methods: classNode.Methods,
		Context: newContext,
	}

	context.Set(classNode.Name, interpreteStructs.VarType{
		Value:      class_,
		IsConstant: true,
		Type:       constants.TT_CLASS,
	})

	return parserStructs.NullNode{}
}

func (interprete *Interprete) BinOP(node interface{}, context *languageContext.Context) interface{} {
	binOP := node.(parserStructs.BinOP)
	nodeLeft := interprete.call(binOP.LeftNode, context)
	nodeRigth := interprete.call(binOP.RigthNode, context)
	newNode := interprete.callMethod(nodeLeft, binOP.Operation.Type_, nodeRigth)
	return newNode
}

func (interprete Interprete) VarAssignNode(node interface{}, context *languageContext.Context) interface{} {
	varAssignNode := node.(parserStructs.VarAssignNode)
	if _, ok := context.Get(varAssignNode.Identifier); ok && varAssignNode.IsConstant {
		panic("Const " + varAssignNode.Identifier)
	}
	result := interprete.call(varAssignNode.Node, context)
	context.Set(varAssignNode.Identifier, interpreteStructs.VarType{
		Value:      result,
		IsConstant: varAssignNode.IsConstant,
	})
	return parserStructs.NullNode{}
}

func (interprete Interprete) UpdateVariableNode(node interface{}, context *languageContext.Context) interface{} {
	updateVariableNode := node.(parserStructs.UpdateVariableNode)
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
	return parserStructs.NullNode{}
}

func (interprete Interprete) VarAccessNode(node interface{}, context *languageContext.Context) interface{} {
	varAccessNode := node.(*parserStructs.VarAccessNode)
	varType, ok := context.Get(varAccessNode.Identifier)
	if !ok {
		customErrors.RunTimeError(&varAccessNode.PositionBase, "Variable is undefined "+varAccessNode.Identifier)
	}
	valueNode, ok := varType.(interpreteStructs.VarType)
	return interprete.call(valueNode.Value, context)
}

func (interprete *Interprete) FuncNode(node interface{}, context *languageContext.Context) interface{} {
	funcNode := node.(parserStructs.FuncNode)
	newContext := languageContext.NewContext(context)

	for _, param := range *funcNode.Params {
		newContext.Set(param.Value.(string), interpreteStructs.VarType{
			Value: parserStructs.NullNode{},
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
	return parserStructs.NullNode{}
}

func (interprete *Interprete) CallFuncNode(node interface{}, context *languageContext.Context) interface{} {
	callFuncNode := node.(*parserStructs.CallFuncNode)
	func_, ok := context.Get(callFuncNode.Name)
	if !ok {
		panic(callFuncNode.Name)
	}
	varType := func_.(interpreteStructs.VarType)

	funcNode := varType.Value.(function.IFunction)

	var params []interface{}

	for _, param := range *callFuncNode.Params {
		params = append(params, interprete.call(param, context))
	}

	funcNodeBody, hasACustomExecute := funcNode.Execute(&params)
	if hasACustomExecute {
		return funcNodeBody
	}
	interprete.call(funcNode.GetBody(), funcNode.GetContext())
	return parserStructs.NullNode{}
}

func (interprete *Interprete) UnaryOP(node interface{}, context *languageContext.Context) *numbers.Number {
	unaryOP := node.(*parserStructs.UnaryOP)
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
	ifNode := node.(parserStructs.IfNode)

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

	return parserStructs.NullNode{}
}

func (interprete *Interprete) WhileNode(node interface{}, context *languageContext.Context) interface{} {
	whileNode := node.(parserStructs.WhileNode)

	for {
		boolean := interprete.call(whileNode.Condition, context).(*booleans.Boolean)
		if !boolean.Value {
			break
		}
		node := interprete.call(whileNode.Body, context)
		if interprete.getMethodName(node) == "CONTINUE" {
			continue
		}
	}

	return parserStructs.NullNode{}
}

func (interprete *Interprete) ForNode(node interface{}, context *languageContext.Context) interface{} {
	forNode := node.(parserStructs.ForNode)

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

	return parserStructs.NullNode{}
}

func (interprete *Interprete) ListNode(node interface{}, context *languageContext.Context) interface{} {
	listNode := node.(parserStructs.ListNode)
	for _, node := range listNode.Nodes {
		result := interprete.call(node, context)
		if interprete.stopExecute(result) != "" {
			return result
		}
	}
	return parserStructs.NullNode{}
}

func (interprete *Interprete) ContinueNode(node interface{}, context *languageContext.Context) interface{} {
	return node.(*parserStructs.ContinueNode)
}

func (interprete *Interprete) Number(node interface{}, context *languageContext.Context) *numbers.Number {
	return node.(*numbers.Number)
}
