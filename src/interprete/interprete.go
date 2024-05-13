package interprete

import (
	"fmt"
	"makeLanguages/src/api"
	"makeLanguages/src/constants"
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/array"
	"makeLanguages/src/features/booleans"
	"makeLanguages/src/features/class"
	"makeLanguages/src/features/function"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/interprete/interpreteStructs"
	"makeLanguages/src/languageContext"
	"makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/parser/parserStructs"
	"reflect"
)

type Interprete struct {
	ast         interpreteStructs.IBaseElement
	currentNode interface{}
	scope       string
	api         *api.Api
}

func NewInterprete(ast interpreteStructs.IBaseElement, scope string, api *api.Api) *Interprete {
	return &Interprete{
		ast:   ast,
		scope: scope,
		api:   api,
	}
}

func (interprete *Interprete) Run(context *languageContext.Context) interface{} {
	value := interprete.call(interprete.ast, context)
	return value
}

func (interprete *Interprete) getMethodName(node interface{}) string {
	method := reflect.TypeOf(node)
	if method == nil {
		panic(node)
	}
	if method.Kind() == reflect.Ptr {
		return method.Elem().Name()
	} else {
		return method.Name()
	}
}

func (interprete *Interprete) stopExecute(node interface{}) string {
	if interprete.getMethodName(node) == "ContinueNode" ||
		interprete.getMethodName(node) == "BreakNode" ||
		interprete.getMethodName(node) == "ReturnNode" {
		return interprete.getMethodName(node)
	}

	return ""
}

func (interprete *Interprete) call(node interpreteStructs.IBaseElement, context *languageContext.Context) interpreteStructs.IBaseElement {
	methodName := interprete.getMethodName(node)
	return interprete.callMethod(interprete, methodName, node, context)
}

func (interprete *Interprete) callMethod(object interface{}, methodName string, values ...interface{}) interpreteStructs.IBaseElement {
	method := reflect.ValueOf(object).MethodByName(methodName)
	var params []reflect.Value
	for _, value := range values {
		params = append(params, reflect.ValueOf(value))
	}
	if !method.IsValid() {
		customErrors.RunTimeError(object.(lexerStructs.IPositionBase), fmt.Sprintf("Error tring to access the method %s", methodName), constants.STOP_EXECUTION)
	}

	returnValue := method.Call(params)
	interface_ := returnValue[0].Interface()
	return interface_.(interpreteStructs.IBaseElement)
}

func (interprete *Interprete) callMethodByOp(object interpreteStructs.IBaseElement, op lexerStructs.Token, value interpreteStructs.IBaseElement) interface{} {
	method := reflect.ValueOf(object).MethodByName(op.Type_)
	result, ok := interprete.api.Call(op.Type_, object, value)
	if ok {
		return result
	}
	params := []reflect.Value{reflect.ValueOf(value)}
	if !method.IsValid() {
		customErrors.RunTimeError(op, fmt.Sprintf("Error tring to access the method %s", op.Type_), constants.STOP_EXECUTION)
	}

	returnValue := method.Call(params)
	return returnValue[0].Interface()
}

func (interprete *Interprete) ClassNode(node interface{}, context *languageContext.Context) interface{} {
	classNode := node.(*parserStructs.ClassNode)
	newContext := languageContext.NewContext(context)
	newContext.IsClass = true
	listNode := classNode.Methods.(parserStructs.ListNode)

	for _, func_ := range listNode.Nodes {
		funcNode := func_.(parserStructs.FuncNode)
		interprete.FuncNode(funcNode, newContext)
	}

	class_ := class.Class{
		Context: newContext,
		Name:    classNode.Name,
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
	if binOP.Operation.Type_ == constants.TT_SPOT {
		return interprete.methodAccess(binOP, context)
	}
	nodeLeft := interprete.call(binOP.LeftNode, context)
	nodeRigth := interprete.call(binOP.RigthNode, context)
	newNode := interprete.callMethodByOp(nodeLeft, binOP.Operation, nodeRigth)
	return newNode
}

func (interprete Interprete) NullNode(node interface{}, context *languageContext.Context) interface{} {
	return node.(parserStructs.NullNode)
}

func (interprete *Interprete) methodAccess(node parserStructs.BinOP, context *languageContext.Context) interface{} {
	for node.Operation.Type_ == constants.TT_SPOT {
		classNode := interprete.call(node.LeftNode, context).(class.ClassBase)
		if interprete.getMethodName(node.RigthNode) == "BinOP" {
			subNode := node.RigthNode.(parserStructs.BinOP)
			interprete.call(subNode.LeftNode, classNode.GetClassContext())
			node = subNode
			continue
		}
		return interprete.call(node.RigthNode, classNode.GetClassContext())
	}
	return parserStructs.NullNode{}
}

func (interprete Interprete) VarAssignNode(node interface{}, context *languageContext.Context) interface{} {
	varAssignNode := node.(parserStructs.VarAssignNode)
	if _, ok := context.Get(varAssignNode.Identifier); ok && varAssignNode.IsConstant {
		customErrors.RunTimeError(varAssignNode.IPositionBase, "The "+varAssignNode.Identifier+" is a const variable", constants.STOP_EXECUTION)
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
	if !ok && !context.IsClass {
		panic("The variable no exist" + updateVariableNode.Identifier)
	}
	if varType.IsConstant {
		customErrors.RunTimeError(updateVariableNode.IPositionBase, "The "+updateVariableNode.Identifier+" is a const variable", constants.STOP_EXECUTION)
	}

	result := interprete.call(updateVariableNode.Node, context)

	varType.Value = result

	ok = context.Update(updateVariableNode.Identifier, varType)
	if !ok {
		panic("The variable no exist" + updateVariableNode.Identifier)
	}
	return parserStructs.NullNode{}
}

func (interprete Interprete) VarAccessNode(node interface{}, context *languageContext.Context) interface{} {
	varAccessNode := node.(*parserStructs.VarAccessNode)
	varType, ok := context.Get(varAccessNode.Identifier)
	if !ok {
		customErrors.RunTimeError(varAccessNode.IPositionBase, "Variable is undefined "+varAccessNode.Identifier, constants.STOP_EXECUTION)
	}
	return interprete.call(varType.Value, context)
}

func (interprete *Interprete) FuncNode(node interface{}, context *languageContext.Context) interface{} {
	funcNode := node.(parserStructs.FuncNode)
	newContext := context
	if interprete.scope != "GLOBAL" {
		newContext = languageContext.NewContext(context)
	}

	for _, param := range *funcNode.Params {
		newContext.Set(param.Value.(string), interpreteStructs.VarType{
			Value: parserStructs.NullNode{},
		})
	}

	func_ := function.Function{
		Body:    funcNode.Body,
		Context: newContext,
		Params:  funcNode.Params,
	}

	context.Set(funcNode.Name, interpreteStructs.VarType{
		Value:      func_,
		IsConstant: true,
	})
	return parserStructs.NullNode{}
}

func (interprete *Interprete) CallObjectNode(node interface{}, context *languageContext.Context) interface{} {
	callFuncNode := node.(*parserStructs.CallObjectNode)
	varType, ok := context.Get(callFuncNode.Name)
	if !ok || (callFuncNode.HasNew && interprete.getMethodName(varType.Value) != "Class") {
		customErrors.RunTimeError(callFuncNode.IPositionBase, fmt.Sprintf("The %s is undefined", callFuncNode.Name), constants.STOP_EXECUTION)
	}

	if interprete.getMethodName(varType.Value) != "Class" {
		funcNode := varType.Value.(function.IFunction)

		var params []interpreteStructs.IBaseElement

		for _, param := range *callFuncNode.Params {
			params = append(params, interprete.call(param, context))
		}

		funcNodeBody, hasACustomExecute, err := funcNode.Execute(&params)
		if err != nil {
			customErrors.RunTimeError(callFuncNode.IPositionBase, err.Error(), constants.STOP_EXECUTION)
		}
		if hasACustomExecute {
			return funcNodeBody
		}
		node := interprete.call(funcNode.GetBody(), funcNode.GetContext())

		isReturn := interprete.stopExecute(node)

		if isReturn == "ReturnNode" {
			return_ := node.(*parserStructs.ReturnNode)
			return interprete.call(return_.Value, context)
		}

		return parserStructs.NullNode{}
	}

	return varType.Value
}

func (interprete *Interprete) String_(node interface{}, context *languageContext.Context) interface{} {
	return node
}

func (interprete *Interprete) Boolean(node interface{}, context *languageContext.Context) interface{} {
	return node
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

func (interprete *Interprete) createNewContext(context *languageContext.Context) *languageContext.Context {
	newContext := context
	if interprete.scope == "CURLY_BRACE" {
		newContext = languageContext.NewContext(context)
	}
	return newContext
}

func (interprete *Interprete) IfNode(node interface{}, context *languageContext.Context) interface{} {
	ifNode := node.(parserStructs.IfNode)

	context = interprete.createNewContext(context)

	for _, if_ := range ifNode.Ifs {
		conditionInterface := interprete.call(if_.Condition, context)

		if interprete.getMethodName(conditionInterface) != "Boolean" {
			customErrors.RunTimeError(
				conditionInterface.(lexerStructs.IPositionBase),
				fmt.Sprintf("The return value is %s", interprete.getMethodName(conditionInterface)),
				constants.SHOW_ERROR)
			customErrors.RunTimeError(if_.Condition.(lexerStructs.IPositionBase), "Error if expression need to a condition", constants.STOP_EXECUTION)
		}

		condition := conditionInterface.GetValue().(bool)

		if condition {
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
	context = interprete.createNewContext(context)
	for {
		boolean := interprete.call(whileNode.Condition, context).GetValue().(*booleans.Boolean)
		if !boolean.Value {
			break
		}
		node := interprete.call(whileNode.Body, context)
		stop := interprete.stopExecute(node)

		if stop == "CONTINUE" {
			continue
		}

		if stop == "BreakNode" {
			break
		}
	}

	return parserStructs.NullNode{}
}

func (interprete *Interprete) ForNode(node interface{}, context *languageContext.Context) interface{} {
	forNode := node.(parserStructs.ForNode)
	context = interprete.createNewContext(context)

	for {
		interprete.call(forNode.Expr1, context)
		condition := interprete.call(forNode.Condition, context)
		coditionNode := condition.GetValue().(*booleans.Boolean)
		if !coditionNode.Value {
			break
		}
		bodyNode := interprete.call(forNode.Body, context)
		stop := interprete.stopExecute(bodyNode)
		if stop == "CONTINUE" {
			continue
		}
		if stop == "BreakNode" {
			break
		}
		interprete.call(forNode.Expr2, context)
	}

	return parserStructs.NullNode{}
}

func (interprete *Interprete) ThisNode(node interface{}, context *languageContext.Context) interface{} {
	thisNode := node.(parserStructs.ThisNode)
	classContext, ok := context.GetClassContext()
	if !ok {
		customErrors.RunTimeError(thisNode, "This need to be inside of class", constants.STOP_EXECUTION)
	}

	return class.Class{
		Context: classContext,
	}
}

func (interprete *Interprete) Array(node interface{}, context *languageContext.Context) interface{} {
	return node
}

func (interprete *Interprete) ArrayAccess(node interface{}, context *languageContext.Context) interface{} {
	arrayAccess := node.(parserStructs.ArrayAccess)
	varType, ok := context.Get(arrayAccess.Identifier)
	if !ok {
		customErrors.RunTimeError(arrayAccess.IPositionBase, "Variable is undefined "+arrayAccess.Identifier, constants.STOP_EXECUTION)
	}
	index := interprete.call(arrayAccess.Node, context)
	if interprete.getMethodName(index) != "Number" {
		customErrors.RunTimeError(arrayAccess.Node.(lexerStructs.IPositionBase), "The index is not a number ", constants.STOP_EXECUTION)
	}
	array_ := varType.Value.GetValue().(*array.Array)
	element := (*array_.Value)[int(index.(*numbers.Number).Value)]
	return element
}

func (interprete *Interprete) ListNode(node interface{}, context *languageContext.Context) interface{} {
	listNode := node.(parserStructs.ListNode)
	values := []interface{}{}
	for _, node := range listNode.Nodes {
		result := interprete.call(node, context)
		if interprete.stopExecute(result) != "" {
			return result
		}
		values = append(values, result)
	}
	return array.NewArray(&values)
}

func (interprete *Interprete) ContinueNode(node interface{}, context *languageContext.Context) interface{} {
	return node.(*parserStructs.ContinueNode)
}

func (interprete *Interprete) BreakNode(node interface{}, context *languageContext.Context) interface{} {
	return node.(*parserStructs.BreakNode)
}

func (interprete *Interprete) ReturnNode(node interface{}, context *languageContext.Context) interface{} {
	return node.(*parserStructs.ReturnNode)
}

func (interprete *Interprete) Number(node interface{}, context *languageContext.Context) *numbers.Number {
	return node.(*numbers.Number)
}

func (interprete *Interprete) Class(node interface{}, context *languageContext.Context) class.Class {
	return node.(class.Class)
}
