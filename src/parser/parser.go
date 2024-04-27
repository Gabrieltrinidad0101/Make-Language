package parser

import (
	"fmt"
	"makeLanguages/src/constants"
	"makeLanguages/src/customErrors"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/features/str"
	"makeLanguages/src/interprete/interpreteStructs"
	lexerStructs "makeLanguages/src/lexer/lexerStructs"
	"makeLanguages/src/parser/parserStructs"
	"slices"
)

type Parser struct {
	tokens       *[]lexerStructs.Token
	idx          int
	CurrentToken *lexerStructs.Token
	len          int
	scoopClass   bool
}

func NewParser(tokens *[]lexerStructs.Token) *Parser {
	return &Parser{
		idx:          -1,
		tokens:       tokens,
		CurrentToken: &lexerStructs.Token{},
		len:          len(*tokens),
	}
}

func (parser *Parser) advance() bool {
	return parser.advances(1)
}

func (parser *Parser) advances(number int) bool {
	for number > 0 {
		parser.idx++
		if parser.idx >= parser.len {
			return false
		}
		*parser.CurrentToken = (*parser.tokens)[parser.idx]
		number--
	}
	return true
}

func (parser *Parser) getToken(idx int) (*lexerStructs.Token, bool) {
	if parser.idx >= parser.len {
		return nil, false
	}
	return &(*parser.tokens)[idx], true
}

func (parser *Parser) verifyNextToken(tokensType ...string) (*lexerStructs.Token, error) {
	var lastToken *lexerStructs.Token = nil
	for i, type_ := range tokensType {
		token, ok := parser.getToken(parser.idx + i)
		if !ok || token.Type_ != type_ {
			return nil, fmt.Errorf("Expect: %s", type_)
		}
		lastToken = token
	}
	parser.advances(len(tokensType))
	return lastToken, nil
}

func (parser *Parser) binOP(callBack func() (interpreteStructs.IBaseElement, error), ops ...string) (interpreteStructs.IBaseElement, error) {
	leftNode, err := callBack()
	if err != nil {
		return nil, err
	}
	for slices.Contains[[]string](ops, parser.CurrentToken.Type_) {
		operation := *parser.CurrentToken
		parser.advance()
		rigthNode, err := callBack()
		if err != nil {
			return nil, err
		}
		leftNode = parserStructs.BinOP{
			LeftNode:  leftNode,
			Operation: operation,
			RigthNode: rigthNode,
		}
	}

	return leftNode, nil
}

func (parser *Parser) Parse() (interpreteStructs.IBaseElement, error) {
	parser.advance()
	ast, err := parser.expr()

	if err != nil {
		return nil, err
	}

	if parser.CurrentToken.Type_ != constants.EOF {
		return nil, fmt.Errorf("Expect +,-,*,/")
	}

	return ast, nil
}

func (parser *Parser) expr() (interpreteStructs.IBaseElement, error) {
	return parser.statements(constants.EOF)
}

func (parser *Parser) statements(tokenEnd string) (interpreteStructs.IBaseElement, error) {
	return parser.statementsBase(tokenEnd, parser.statement)
}

func (parser *Parser) statementsBase(tokenEnd string, callBack func() (interpreteStructs.IBaseElement, error)) (interpreteStructs.IBaseElement, error) {
	for parser.CurrentToken.Type_ == constants.TT_NEWLINE {
		parser.advance()
	}
	listNodes := parserStructs.ListNode{}
	ast, err := callBack()
	if err != nil {
		return nil, err
	}
	listNodes.Nodes = append(listNodes.Nodes, ast)

	for {
		thereIsANewLine := false
		for parser.CurrentToken.Type_ == constants.TT_NEWLINE {
			parser.advance()
			thereIsANewLine = true
		}

		if !thereIsANewLine {
			break
		}

		if parser.CurrentToken.Type_ == tokenEnd {
			parser.advance()
			break
		}

		ast, err := callBack()
		if err != nil {
			return nil, err
		}
		listNodes.Nodes = append(listNodes.Nodes, ast)
	}
	return listNodes, nil
}

func (parser *Parser) statement() (interpreteStructs.IBaseElement, error) {

	variableAndConst, err := parser.variableAndConst()
	if variableAndConst != nil || err != nil {
		return variableAndConst, err
	}

	updateVariable, err := parser.updateVariable()
	if updateVariable != nil || err != nil {
		return updateVariable, err
	}

	continue_ := parser.continue_()
	if continue_ != nil {
		return continue_, nil
	}

	break_ := parser.break_()
	if break_ != nil {
		return break_, nil
	}

	while, err := parser.while()
	if while != nil || err != nil {
		return while, err
	}

	for_, err := parser.for_()
	if for_ != nil || err != nil {
		return for_, err
	}

	class, err := parser.class()
	if class != nil || err != nil {
		return class, err
	}

	return_, err := parser.return_()
	if return_ != nil || err != nil {
		return return_, err
	}

	thisTop, err := parser.thisTop()
	if thisTop != nil || err != nil {
		return thisTop, err
	}

	return parser.compare()
}

func (parser *Parser) thisTop() (interpreteStructs.IBaseElement, error) {
	this, err := parser.this()
	if this == nil && err == nil {
		return nil, nil
	}

	spot, err := parser.verifyNextToken(constants.TT_SPOT)

	if err != nil {
		return this, nil
	}

	node, err := parser.binOP(parser.thisTopNext, constants.TT_POW)

	if err != nil {
		return nil, err
	}
	return parserStructs.BinOP{
		LeftNode:  this,
		Operation: *spot,
		RigthNode: node,
	}, nil
}

func (parser *Parser) thisTopNext() (interpreteStructs.IBaseElement, error) {
	updateVariable, err := parser.updateVariable()
	if updateVariable != nil && err == nil {
		return updateVariable, err
	}
	return parser.term()
}

func (parser *Parser) continue_() *parserStructs.ContinueNode {
	continue_, err := parser.verifyNextToken(constants.TT_CONTINUE)
	if err != nil {
		return nil
	}
	return &parserStructs.ContinueNode{
		IPositionBase: continue_.IPositionBase,
	}
}

func (parser *Parser) return_() (*parserStructs.ReturnNode, error) {
	continue_, err := parser.verifyNextToken(constants.TT_RETURN)
	if err != nil {
		return nil, nil
	}

	value, err := parser.compare()

	if err != nil {
		return nil, err
	}

	return &parserStructs.ReturnNode{
		IPositionBase: continue_.IPositionBase,
		Value:         value,
	}, nil
}

func (parser *Parser) break_() *parserStructs.BreakNode {
	break_, err := parser.verifyNextToken(constants.TT_BREAK)
	if err != nil {
		return nil
	}
	return &parserStructs.BreakNode{
		IPositionBase: break_.IPositionBase,
	}
}

func (parser *Parser) class() (*parserStructs.ClassNode, error) {
	identoifierToken, err := parser.verifyNextToken(constants.TT_CLASS, constants.TT_IDENTIFIER)
	if err != nil {
		return nil, nil
	}
	parser.scoopClass = true

	_, err = parser.verifyNextToken(constants.TT_START_BODY)
	if err != nil {
		return nil, err
	}

	methodes, err := parser.statementsBase(constants.TT_END_BODY, parser.func_)
	parser.scoopClass = false
	if err != nil {
		return nil, err
	}
	return &parserStructs.ClassNode{
		Methods: methodes,
		Name:    identoifierToken.Value.(string),
	}, nil
}

func (parser *Parser) variableAndConst() (interpreteStructs.IBaseElement, error) {
	_, constError := parser.verifyNextToken(constants.TT_CONST)
	_, varError := parser.verifyNextToken(constants.TT_VAR)
	if constError != nil && varError != nil {
		return nil, nil
	}
	identifier := parser.CurrentToken.Value
	_, err := parser.verifyNextToken(constants.TT_IDENTIFIER, constants.TT_EQ)
	if err != nil {
		return nil, err
	}
	node, err := parser.compare()

	if err != nil {
		return nil, err
	}

	return parserStructs.VarAssignNode{
		Identifier: identifier.(string),
		Node:       node,
		IsConstant: constError == nil,
	}, nil
}

func (parser *Parser) updateVariable() (interpreteStructs.IBaseElement, error) {
	value := parser.CurrentToken.Value
	_, err := parser.verifyNextToken(constants.TT_IDENTIFIER, constants.TT_EQ)
	if err != nil {
		return nil, nil
	}
	expr, err := parser.compare()

	if err != nil {
		return nil, err
	}

	return parserStructs.UpdateVariableNode{
		Identifier: value.(string),
		Node:       expr,
	}, nil
}

func (parser *Parser) while() (interpreteStructs.IBaseElement, error) {
	_, err := parser.verifyNextToken(constants.TT_WHILE)
	if err != nil {
		return nil, nil
	}
	conditionAndBodyBase, err := parser.conditionAndBodyBase()

	if err != nil {
		return nil, err
	}
	return parserStructs.WhileNode{
		Condition: conditionAndBodyBase.Condition,
		Body:      conditionAndBodyBase.Body,
	}, nil
}

func (parser *Parser) for_() (interpreteStructs.IBaseElement, error) {
	_, err := parser.verifyNextToken(constants.TT_FOR)
	if err != nil {
		return nil, nil
	}

	_, err = parser.verifyNextToken(constants.TT_LPAREN)
	if err != nil {
		return nil, nil
	}

	expr1, err := parser.statement()
	if err != nil {
		return nil, nil
	}

	_, err = parser.verifyNextToken(constants.TT_SEMICOLON)
	if err != nil {
		return nil, nil
	}

	condition, err := parser.compare()
	if err != nil {
		return nil, nil
	}

	_, err = parser.verifyNextToken(constants.TT_SEMICOLON)
	if err != nil {
		return nil, nil
	}

	expr2, err := parser.statement()
	if err != nil {
		return nil, nil
	}

	_, err = parser.verifyNextToken(constants.TT_RPAREN)
	if err != nil {
		return nil, nil
	}

	body, err := parser.BodyBase()
	if err != nil {
		return nil, nil
	}

	return parserStructs.ForNode{
		Expr1:     expr1,
		Condition: condition,
		Expr2:     expr2,
		Body:      body,
	}, nil
}

func (parser *Parser) AndOr() (interpreteStructs.IBaseElement, error) {
	return parser.binOP(parser.compare, constants.TT_AND, constants.TT_AND)
}

func (parser *Parser) compare() (interpreteStructs.IBaseElement, error) {
	return parser.binOP(parser.plus, constants.TT_GT, constants.TT_GTE, constants.TT_GT, constants.TT_LT, constants.TT_LTE, constants.TT_EQE)
}

func (parser *Parser) plus() (interpreteStructs.IBaseElement, error) {
	return parser.binOP(parser.factor, constants.TT_PLUS, constants.TT_MINUS)
}

func (parser *Parser) factor() (interpreteStructs.IBaseElement, error) {
	return parser.binOP(parser.pow, constants.TT_MUL, constants.TT_DIV)
}

func (parser *Parser) pow() (interpreteStructs.IBaseElement, error) {
	return parser.binOP(parser.spot, constants.TT_POW, constants.TT_SQUARE_ROOT)
}

func (parser *Parser) spot() (interpreteStructs.IBaseElement, error) {
	return parser.binOP(parser.term, constants.TT_SPOT)
}

func (parser *Parser) term() (interpreteStructs.IBaseElement, error) {
	currentNode := *parser.CurrentToken
	nodeType := currentNode.Type_
	if nodeType == constants.TT_PLUS || nodeType == constants.TT_MINUS || nodeType == constants.TT_PLUS1 || nodeType == constants.TT_MINUS1 {
		token, ok := parser.getToken(parser.idx + 1)
		if ok && token.Type_ == constants.TT_PLUS || token.Type_ == constants.TT_MINUS {
			return nil, fmt.Errorf("Error is necesery a ( between - and + simbols")
		}
		parser.advance()
		rigthNode, err := parser.term()

		if err != nil {
			return nil, err
		}

		unaryOP := parserStructs.UnaryOP{
			Operation: nodeType,
			RigthNode: rigthNode,
		}
		return &unaryOP, nil
	}
	if nodeType == "number" {
		value := parser.CurrentToken.Value.(float64)
		number := numbers.NewNumbers(value, parser.CurrentToken.IPositionBase)
		parser.advance()

		return number, nil
	}
	if parser.CurrentToken.Type_ == constants.TT_LPAREN {
		parser.advance()
		node, err := parser.statement()
		if err != nil {
			return nil, err
		}
		if !(parser.CurrentToken.Type_ == constants.TT_RPAREN) {
			return nil, fmt.Errorf("Expect )")
		}
		parser.advance()
		return node, nil
	}
	if ifNode, err := parser.if_(); ifNode != nil || err != nil {
		return ifNode, err
	}
	if callFuncNode, err := parser.callFunc(); callFuncNode != nil || err != nil {
		return callFuncNode, err
	}
	if funcNode, err := parser.func_(); funcNode != nil || err != nil {
		return funcNode, err
	}
	if arrayAccess, err := parser.arrayAccess(); err != nil || arrayAccess != nil {
		return arrayAccess, err
	}
	if varAccess, err := parser.varAccess(); err != nil || varAccess != nil {
		return varAccess, err
	}
	if string_, err := parser.string_(); err != nil || string_ != nil {
		return string_, nil
	}
	if array, err := parser.array(); err != nil || array != nil {
		return array, nil
	}
	if this, err := parser.this(); err != nil || this != nil {
		return this, nil
	}

	return nil, fmt.Errorf("")
}

func (parser *Parser) this() (interpreteStructs.IBaseElement, error) {
	token, err := parser.verifyNextToken(constants.TT_THIS)
	if err != nil {
		return nil, nil
	}

	if !parser.scoopClass {
		customErrors.InvalidSyntax(*token, "This can only be inside of class", constants.STOP_EXECUTION)
		return nil, nil
	}

	return parserStructs.ThisNode{
		IPositionBase: token.IPositionBase,
	}, nil
}

func (parser *Parser) arrayAccess() (interpreteStructs.IBaseElement, error) {
	identifier := *parser.CurrentToken
	if _, err := parser.verifyNextToken(constants.TT_IDENTIFIER, constants.TT_LSQUAREBRACKET); err != nil {
		return nil, nil
	}

	node, err := parser.term()

	if err != nil {
		return nil, nil
	}

	if _, err := parser.verifyNextToken(constants.TT_RSQUAREBRACKET); err != nil {
		return nil, err
	}

	return parserStructs.ArrayAccess{
		Identifier: identifier.Value.(string),
		Node:       node,
	}, nil
}

func (parser *Parser) array() (interpreteStructs.IBaseElement, error) {
	_, err := parser.verifyNextToken(constants.TT_LSQUAREBRACKET)
	if err != nil {
		return nil, nil
	}
	listNode := parserStructs.ListNode{}
	for {
		if _, err := parser.verifyNextToken(constants.TT_RSQUAREBRACKET); err == nil {
			break
		}
		node, err := parser.term()
		if err != nil {
			return nil, err
		}
		listNode.Nodes = append(listNode.Nodes, node)
		if parser.CurrentToken.Type_ != constants.TT_RSQUAREBRACKET {
			_, err := parser.verifyNextToken(constants.TT_COMMA)
			if err != nil {
				return nil, err
			}
		}
	}

	return listNode, nil
}

func (parser *Parser) string_() (interpreteStructs.IBaseElement, error) {
	stringToken, err := parser.verifyNextToken(constants.TT_STRING)
	if err != nil {
		return nil, nil
	}
	return str.NewString(stringToken.Value.(string), stringToken.IPositionBase), nil
}

func (parser *Parser) if_() (interpreteStructs.IBaseElement, error) {
	ifs := []*parserStructs.ConditionAndBody{}
	var elseNode interface{}
	if parser.CurrentToken.Type_ != constants.TT_IF {
		return nil, nil
	}

	parser.advance()

	node, err := parser.conditionAndBodyBase()

	if err != nil {
		return nil, err
	}

	ifs = append(ifs, node)

	for {
		_, err := parser.verifyNextToken(constants.TT_ELIF)
		if err != nil {
			break
		}
		node, err := parser.conditionAndBodyBase()

		if err != nil {
			return nil, err
		}

		ifs = append(ifs, node)
	}
	if parser.CurrentToken.Type_ == constants.TT_ELSE {
		parser.advance()
		elseNode, err = parser.BodyBase()
	}

	return parserStructs.IfNode{
		Ifs:   ifs,
		Else_: elseNode,
	}, nil
}

func (parser *Parser) func_() (interpreteStructs.IBaseElement, error) {
	func_ := parser.CurrentToken
	identoifierToken, err := parser.verifyNextToken(constants.TT_FUNC, constants.TT_IDENTIFIER)
	if err != nil {
		return nil, nil
	}
	params, position, err := parser.params()
	if err != nil {
		return nil, err
	}
	body, err := parser.BodyBase()
	if err != nil {
		return nil, err
	}
	return parserStructs.FuncNode{
		Params: params,
		Body:   body,
		Name:   identoifierToken.Value.(string),
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: func_.GetPositionStart(),
			PositionEnd:   position.PositionCopy(),
		},
	}, nil
}

func (parser *Parser) callFunc() (*parserStructs.CallObjectNode, error) {
	newNode, _ := parser.verifyNextToken(constants.TT_NEW)
	func_ := *parser.CurrentToken
	_, err := parser.verifyNextToken(constants.TT_IDENTIFIER, constants.TT_LPAREN)
	if err != nil {
		return nil, nil
	}
	params, positionEnd, err := parser.args()
	if err != nil {
		return nil, err
	}
	fistToken := newNode
	if fistToken == nil {
		fistToken = &func_
	}

	return &parserStructs.CallObjectNode{
		Params: params,
		Name:   func_.Value.(string),
		HasNew: newNode != nil,
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: fistToken.GetPositionStart(),
			PositionEnd:   positionEnd.PositionCopy(),
		},
	}, nil
}

func (parser *Parser) conditionAndBodyBase() (*parserStructs.ConditionAndBody, error) {
	_, err := parser.verifyNextToken(constants.TT_LPAREN)
	if err != nil {
		return nil, err
	}

	condition, err := parser.AndOr()
	if err != nil {
		return nil, err
	}

	_, err = parser.verifyNextToken(constants.TT_RPAREN)
	if err != nil {
		return nil, err
	}

	body, err := parser.BodyBase()

	if err != nil {
		return nil, err
	}

	return &parserStructs.ConditionAndBody{
		Condition: condition,
		Body:      body,
	}, err
}

func (parser *Parser) BodyBase() (interface{}, error) {
	if _, err := parser.verifyNextToken(constants.TT_START_BODY); err == nil {
		parser.verifyNextToken(constants.TT_NEWLINE)
		return parser.statements(constants.TT_END_BODY)
	}
	return parser.statement()
}

func (parser *Parser) varAccess() (*parserStructs.VarAccessNode, error) {
	if parser.CurrentToken.Type_ != constants.TT_IDENTIFIER {
		return nil, nil
	}

	varAccessNode := &parserStructs.VarAccessNode{
		Identifier: parser.CurrentToken.Value.(string),
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: parser.CurrentToken.GetPositionStart(),
			PositionEnd:   parser.CurrentToken.GetPositionEnd(),
		},
	}
	parser.advance()
	return varAccessNode, nil
}

func (parser *Parser) params() (*[]lexerStructs.Token, *lexerStructs.Position, error) {
	_, err := parser.verifyNextToken(constants.TT_LPAREN)
	params := &[]lexerStructs.Token{}
	var lastToken *lexerStructs.Token
	if err != nil {
		return nil, nil, err
	}
	for {
		identifier, err := parser.verifyNextToken(constants.TT_IDENTIFIER)
		if err != nil {
			lastToken, err = parser.verifyNextToken(constants.TT_RPAREN)
			if err == nil {
				position := lastToken.GetPositionEnd()
				return &[]lexerStructs.Token{}, &position, nil
			}
			return nil, nil, err
		}
		*params = append(*params, *identifier)

		_, err = parser.verifyNextToken(constants.TT_COMMA)
		if err != nil {
			lastToken, err = parser.verifyNextToken(constants.TT_RPAREN)
			if err == nil {
				break
			}
			return nil, nil, err
		}

	}
	position := lastToken.GetPositionEnd()
	return params, &position, nil
}

func (parser *Parser) args() (*[]interface{}, *lexerStructs.Position, error) {
	params := []interface{}{}
	var lastToken *lexerStructs.Token
	var err error
	for {
		lastToken, err = parser.verifyNextToken(constants.TT_RPAREN)
		if err == nil {
			break
		}
		param, err := parser.compare()
		if err != nil {
			return nil, nil, err
		}
		params = append(params, param)
		_, err = parser.verifyNextToken(constants.TT_COMMA)
		if err != nil {
			lastToken, err = parser.verifyNextToken(constants.TT_RPAREN)
			if err == nil {
				break
			}
			return nil, nil, err
		}
	}

	position := lastToken.GetPositionEnd()
	return &params, &position, nil
}
