package interprete

import (
	"fmt"
	"reflect"

	"github.com/Gabrieltrinidad0101/Make-Language/src/api"
	"github.com/Gabrieltrinidad0101/Make-Language/src/parser/parserStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer"

	"github.com/Gabrieltrinidad0101/Make-Language/src/languageContext"

	"github.com/Gabrieltrinidad0101/Make-Language/src/interprete/interpreteStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/numbers"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/function"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/class"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/array"

	"github.com/Gabrieltrinidad0101/Make-Language/src/customErrors"

	"github.com/Gabrieltrinidad0101/Make-Language/src/constants"
)

type Interprete struct {
	ast         interpreteStructs.IBaseElement
	currentNode interface{}
	api         *api.Api
	conf        lexer.LanguageConfiguraction
}

func NewInterprete(ast interpreteStructs.IBaseElement, scope string, api *api.Api, conf lexer.LanguageConfiguraction) *Interprete {
	return &Interprete{
		ast:  ast,
		api:  api,
		conf: conf,
	}
}

func (interprete *Interprete) Run(context *languageContext.Context) (interface{}, error) {
	value, err := interprete.call(interprete.ast, context)
	return value, err
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

func (interprete *Interprete) call(node interpreteStructs.IBaseElement, context *languageContext.Context) (interpreteStructs.IBaseElement, error) {
	methodName := interprete.getMethodName(node)
	return interprete.callMethod(interprete, methodName, node, context)
}

func (interprete *Interprete) callMethod(object interface{}, methodName string, values ...interface{}) (interpreteStructs.IBaseElement, error) {
	method := reflect.ValueOf(object).MethodByName(methodName)
	var params []reflect.Value
	for _, value := range values {
		params = append(params, reflect.ValueOf(value))
	}
	if !method.IsValid() {
		typeName := interprete.getMethodName(object)
		return nil, customErrors.RunTimeError(object.(lexerStructs.IPositionBase), fmt.Sprintf("Error %s to access the method %s", typeName, methodName))
	}

	returnValue := method.Call(params)
	interface_ := returnValue[0].Interface()
	if interface_ == nil {
		err := returnValue[1].Interface().(error)
		return nil, err
	}
	return interface_.(interpreteStructs.IBaseElement), nil
}

func (interprete *Interprete) callMethodByOp(object interpreteStructs.IBaseElement, op lexerStructs.Token, value interpreteStructs.IBaseElement) (interface{}, error) {
	method := reflect.ValueOf(object).MethodByName(op.Type_)
	result, ok := interprete.api.Call(op.Type_, object, value)
	if ok {
		return result, nil
	}
	params := []reflect.Value{reflect.ValueOf(value)}
	if !method.IsValid() {
		typeName := interprete.getMethodName(object)
		return nil, customErrors.RunTimeError(op, fmt.Sprintf("Error %s does not have the method %s", typeName, op.Type_))
	}

	returnValue := method.Call(params)
	if returnValue[1].Interface() == nil {
		return returnValue[0].Interface(), nil
	}
	return returnValue[0].Interface(), returnValue[1].Interface().(error)
}

func (interprete *Interprete) ClassNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	classNode := node.(*parserStructs.ClassNode)
	newContext := languageContext.NewContext(context)
	newContext.IsClass = true

	method := interprete.getMethodName(classNode.Methods) != "NullNode"
	if method {
		listNode := classNode.Methods.(parserStructs.ListNode)
		for _, func_ := range listNode.Nodes {
			funcNode := func_.(parserStructs.FuncNode)
			interprete.FuncNode(funcNode, newContext)
		}
	}

	class_ := class.Class{
		Context: newContext,
		Name:    classNode.Name,
	}

	context.Set(classNode.Name, &interpreteStructs.VarType{
		Value:      class_,
		IsConstant: true,
		Type:       constants.TT_CLASS,
	})

	return parserStructs.NullNode{}, nil
}

func (interprete *Interprete) BinOP(node interface{}, context *languageContext.Context) (interface{}, error) {
	binOP := node.(parserStructs.BinOP)
	if binOP.Operation.Type_ == constants.TT_SPOT {
		return interprete.methodAccess(binOP, context)
	}

	nodeLeft, err := interprete.call(binOP.LeftNode, context)
	if err != nil {
		return nil, err
	}
	nodeRigth, err := interprete.call(binOP.RigthNode, context)
	if err != nil {
		return nil, err
	}

	newNode, err := interprete.callMethodByOp(nodeLeft, binOP.Operation, nodeRigth)
	return newNode, err
}

func (interprete Interprete) NullNode(node interface{}, context *languageContext.Context) interface{} {
	return node.(parserStructs.NullNode)
}

func (interprete *Interprete) methodAccess(node parserStructs.BinOP, context *languageContext.Context) (interface{}, error) {
	for node.Operation.Type_ == constants.TT_SPOT {
		class_, err := interprete.call(node.LeftNode, context)
		if err != nil {
			return nil, err
		}

		if interprete.getMethodName(class_) == "NullNode" {
			return nil, customErrors.RunTimeError(node.RigthNode, "Invalid method access")
		}

		classNode := class_.(class.ClassBase)
		if interprete.getMethodName(node.RigthNode) == "BinOP" {
			subNode := node.RigthNode.(parserStructs.BinOP)
			interprete.call(subNode.LeftNode, classNode.GetClassContext())
			node = subNode
			continue
		}
		if interprete.getMethodName(node.RigthNode) == "UpdateVariableNode" {
			updateVariableNode := node.RigthNode.(parserStructs.UpdateVariableNode)
			updateVariableNode.SetValueContext = context
			return interprete.call(updateVariableNode, classNode.GetClassContext())
		}
		return interprete.call(node.RigthNode, classNode.GetClassContext())
	}
	return parserStructs.NullNode{}, nil
}

func (interprete Interprete) VarAssignNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	varAssignNode := node.(parserStructs.VarAssignNode)
	if _, ok := context.Get(varAssignNode.Identifier); ok && varAssignNode.IsConstant {
		return nil, customErrors.RunTimeError(varAssignNode.IPositionBase, "The "+varAssignNode.Identifier+" is a const variable")
	}
	result, err := interprete.call(varAssignNode.Node, context)
	if err != nil {
		return nil, err
	}
	context.Set(varAssignNode.Identifier, &interpreteStructs.VarType{
		Value:      result,
		IsConstant: varAssignNode.IsConstant,
	})
	return parserStructs.NullNode{}, nil
}

func (interprete Interprete) UpdateVariableNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	updateVariableNode := node.(parserStructs.UpdateVariableNode)
	varType, ok := context.Get(updateVariableNode.Identifier)
	if !ok && !context.IsClass {
		err := customErrors.RunTimeError(updateVariableNode.Node, "The "+updateVariableNode.Identifier+" no exist")
		return nil, err
	}
	if varType.IsConstant && !context.IsClass {
		err := customErrors.RunTimeError(updateVariableNode.IPositionBase, "The "+updateVariableNode.Identifier+" is a const variable")
		if err != nil {
			return nil, err
		}
	}

	setValueContext := updateVariableNode.SetValueContext
	if setValueContext == nil {
		setValueContext = context
	}

	result, err := interprete.call(updateVariableNode.Node, setValueContext)

	if varType.OnUpdateVariable != nil {
		varType.OnUpdateVariable(result.GetValue())
	}

	if err != nil {
		return nil, err
	}
	varType.Value = result
	ok = context.Update(updateVariableNode.Identifier, varType)
	if !ok {
		panic("The variable no exist" + updateVariableNode.Identifier)
	}
	return parserStructs.NullNode{}, nil
}

func (interprete Interprete) VarAccessNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	varAccessNode := node.(*parserStructs.VarAccessNode)
	varType, ok := context.Get(varAccessNode.Identifier)
	if !ok {
		return nil, customErrors.RunTimeError(varAccessNode.IPositionBase, "Variable is undefined "+varAccessNode.Identifier)
	}
	return interprete.call(varType.Value, context)
}

func (interprete *Interprete) FuncNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	funcNode := node.(parserStructs.FuncNode)
	newContext := context
	if interprete.conf.Scope != "GLOBAL" {
		newContext = languageContext.NewContext(context)
	}

	for _, param := range *funcNode.Params {
		newContext.Set(param.Value.(string), &interpreteStructs.VarType{
			Value: parserStructs.NullNode{},
		})
	}

	func_ := function.Function{
		Body:    funcNode.Body,
		Context: newContext,
		Params:  funcNode.Params,
	}

	context.Set(funcNode.Name, &interpreteStructs.VarType{
		Value:      func_,
		IsConstant: true,
	})
	return parserStructs.NullNode{}, nil
}

func (interprete *Interprete) CallObjectNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	callFuncNode := node.(*parserStructs.CallObjectNode)
	varType, ok := context.Get(callFuncNode.Name)
	if !ok || (callFuncNode.HasNew && interprete.getMethodName(varType.Value) != "Class") {
		return nil, customErrors.RunTimeError(callFuncNode.IPositionBase, fmt.Sprintf("The %s is undefined", callFuncNode.Name))
	}

	if interprete.getMethodName(varType.Value) != "Class" {
		funcNode := varType.Value.(function.IFunction)
		var params []interpreteStructs.IBaseElement

		if funcNode.CanChangeContextParent() {
			funcNode.GetContext().Parent = context
		}

		for _, param := range *callFuncNode.Params {
			node, err := interprete.call(param, funcNode.GetContext())
			if err != nil {
				return nil, err
			}
			params = append(params, node)
		}

		funcNodeBody, hasACustomExecute, err := funcNode.Execute(&params)
		if err != nil {
			return nil, customErrors.RunTimeError(callFuncNode.IPositionBase, err.Error())
		}
		if hasACustomExecute {
			return funcNodeBody, nil
		}
		node, err := interprete.call(funcNode.GetBody(), funcNode.GetContext())
		if err != nil {
			return nil, err
		}
		isReturn := interprete.stopExecute(node)

		if isReturn == "ReturnNode" {
			return_ := node.(*parserStructs.ReturnNode)
			return interprete.call(return_.Value, funcNode.GetContext())
		}

		return parserStructs.NullNode{}, nil
	}

	class := varType.Value.(class.Class)

	if interprete.conf.ConstructorName == "CLASS_NAME" {
		callFuncNode.HasNew = false
		_, ok := class.Context.GetBase(callFuncNode.Name)
		if ok {
			interprete.CallObjectNode(callFuncNode, class.Context)
		}
	}
	return varType.Value, nil
}

func (interprete *Interprete) String_(node interface{}, context *languageContext.Context) interface{} {
	return node
}

func (interprete *Interprete) Boolean(node interface{}, context *languageContext.Context) interface{} {
	return node
}

func (interprete *Interprete) UnaryOP(node interface{}, context *languageContext.Context) (*numbers.Number, error) {
	unaryOP := node.(*parserStructs.UnaryOP)
	number_, err := interprete.call(unaryOP.RigthNode, context)
	if err != nil {
		return nil, err
	}

	number := number_.(*numbers.Number)

	if unaryOP.Operation == "MINUS" {
		number.Value *= -1
	}
	if unaryOP.Operation == "PLUS1" {
		number.Value += 1
	}
	if unaryOP.Operation == "MINUS1" {
		number.Value -= 1
	}
	return number, nil
}

func (interprete *Interprete) createNewContext(context *languageContext.Context) *languageContext.Context {
	newContext := context
	if interprete.conf.Scope == "CURLY_BRACE" {
		newContext = languageContext.NewContext(context)
	}
	return newContext
}

func (interprete *Interprete) IfNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	ifNode := node.(parserStructs.IfNode)

	context = interprete.createNewContext(context)

	for _, if_ := range ifNode.Ifs {
		conditionInterface, err := interprete.call(if_.Condition, context)
		if err != nil {
			return nil, err
		}
		if interprete.getMethodName(conditionInterface) != "Boolean" {
			err1 := customErrors.RunTimeError(conditionInterface.(lexerStructs.IPositionBase), fmt.Sprintf("The return value is %s", interprete.getMethodName(conditionInterface)))
			err2 := customErrors.RunTimeError(if_.Condition.(lexerStructs.IPositionBase), "Error if expression need to a condition")
			return nil, fmt.Errorf("%v \n\n %v", err1.Error(), err2.Error())

		}

		condition := conditionInterface.GetValue().(bool)

		if condition {
			node, err := interprete.call(if_.Body, context)
			if err != nil {
				return nil, err
			}
			return node, nil
		}
	}

	if ifNode.Else_ != nil {
		node, err := interprete.call(ifNode.Else_, context)
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	return parserStructs.NullNode{}, nil
}

func (interprete *Interprete) WhileNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	whileNode := node.(parserStructs.WhileNode)
	context = interprete.createNewContext(context)
	for {
		boolean_, err := interprete.call(whileNode.Condition, context)
		boolean := boolean_.GetValue().(bool)
		if err != nil {
			return nil, err
		}
		if !boolean {
			break
		}
		node, err := interprete.call(whileNode.Body, context)
		if err != nil {
			return nil, err
		}
		stop := interprete.stopExecute(node)

		if stop == "CONTINUE" {
			continue
		}

		if stop == "BreakNode" {
			break
		}
	}

	return parserStructs.NullNode{}, nil
}

func (interprete *Interprete) ForNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	forNode := node.(parserStructs.ForNode)
	context = interprete.createNewContext(context)

	for {
		interprete.call(forNode.Expr1, context)
		condition, err := interprete.call(forNode.Condition, context)
		if err != nil {
			return nil, err
		}
		coditionNode := condition.GetValue().(bool)
		if !coditionNode {
			break
		}
		bodyNode, err := interprete.call(forNode.Body, context)
		if err != nil {
			return nil, err
		}
		stop := interprete.stopExecute(bodyNode)
		if stop == "CONTINUE" {
			continue
		}
		if stop == "BreakNode" {
			break
		}
		interprete.call(forNode.Expr2, context)
	}

	return parserStructs.NullNode{}, nil
}

func (interprete *Interprete) ThisNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	thisNode := node.(parserStructs.ThisNode)
	classContext, ok := context.GetClassContext()
	if !ok {
		err := customErrors.RunTimeError(thisNode, "This need to be inside of class")
		return nil, err
	}

	return class.Class{
		Context: classContext,
	}, nil
}

func (interprete *Interprete) Array(node interface{}, context *languageContext.Context) interface{} {
	return node
}

func (interprete *Interprete) ArrayAccess(node interface{}, context *languageContext.Context) (interface{}, error) {
	arrayAccess := node.(parserStructs.ArrayAccess)
	varType, ok := context.Get(arrayAccess.Identifier)
	if !ok {
		err := customErrors.RunTimeError(arrayAccess.IPositionBase, "Variable is undefined "+arrayAccess.Identifier)
		return nil, err
	}
	index, err := interprete.call(arrayAccess.Node, context)
	if err != nil {
		return nil, err
	}
	if interprete.getMethodName(index) != "Number" {
		err := customErrors.RunTimeError(arrayAccess.Node.(lexerStructs.IPositionBase), "The index is not a number ")
		if err != nil {
			return nil, err
		}
	}
	array_ := varType.Value.GetValue().(*array.Array)
	element := (*array_.Value)[int(index.(*numbers.Number).Value)]
	return element, nil
}

func (interprete *Interprete) ListNode(node interface{}, context *languageContext.Context) (interface{}, error) {
	listNode := node.(parserStructs.ListNode)
	values := []interface{}{}
	for _, node := range listNode.Nodes {
		result, err := interprete.call(node, context)
		if err != nil {
			return nil, err
		}
		if interprete.stopExecute(result) != "" {
			return result, nil
		}
		values = append(values, result)
	}
	return array.NewArray(&values), nil
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
