package parser

import (
	"fmt"
	"makeLanguages/src/constants"
	"makeLanguages/src/features/numbers"
	"makeLanguages/src/lexer/structs"
	lexerStructs "makeLanguages/src/lexer/structs"
	"slices"
)

type Parser struct {
	tokens       *[]lexerStructs.Token
	idx          int
	CurrentToken *lexerStructs.Token
	len          int
}

type BinOP struct {
	LeftNode  interface{}
	Operation lexerStructs.Token
	RigthNode interface{}
}

type UnaryOP struct {
	Operation string
	RigthNode interface{}
}

type IfNode struct {
	Ifs   []*ConditionAndBody
	Else_ interface{}
}

type ConditionAndBody struct {
	Condition interface{}
	Body      interface{}
}

type VarAssignNode struct {
	Identifier string
	Node       interface{}
	IsConstant bool
}

type UpdateVariableNode struct {
	Identifier string
	Node       interface{}
}

type VarAccessNode struct {
	Identifier string
}

type ListNode struct {
	Nodes []interface{}
}

type WhileNode struct {
	Condition interface{}
	Body      interface{}
}

type FuncNode struct {
	Params *[]structs.Token
	Body   interface{}
	Name   string
}

type ForNode struct {
	Expr1     interface{}
	Expr2     interface{}
	Condition interface{}
	Body      interface{}
}

type CallFuncNode struct {
	Params *[]interface{}
	Name   string
}

type NullNode struct{}

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
		if ok && token.Type_ != type_ {
			return nil, fmt.Errorf("Expect: %s", type_)
		}
		lastToken = token
	}
	parser.advances(len(tokensType))
	return lastToken, nil
}

func (parser *Parser) binOP(callBack func() (interface{}, error), ops ...string) (interface{}, error) {
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
		leftNode = BinOP{
			LeftNode:  leftNode,
			Operation: operation,
			RigthNode: rigthNode,
		}
	}

	return leftNode, nil
}

func (parser *Parser) Parse() (interface{}, error) {
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

func (parser *Parser) expr() (interface{}, error) {
	return parser.statements(constants.EOF)
}

func (parser *Parser) statements(tokenEnd string) (interface{}, error) {
	for parser.CurrentToken.Type_ == constants.TT_NEWLINE {
		parser.advance()
	}
	listNodes := ListNode{}
	ast, err := parser.statement()
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

		ast, err := parser.statement()
		if err != nil {
			return nil, err
		}
		listNodes.Nodes = append(listNodes.Nodes, ast)
	}
	return listNodes, nil
}

func (parser *Parser) statement() (interface{}, error) {

	variableAndConst, err := parser.variableAndConst()
	if variableAndConst != nil || err != nil {
		return variableAndConst, err
	}

	updateVariable, err := parser.updateVariable()
	if updateVariable != nil || err != nil {
		return updateVariable, err
	}

	while, err := parser.while()
	if while != nil || err != nil {
		return while, err
	}

	for_, err := parser.for_()
	if for_ != nil || err != nil {
		return for_, err
	}

	return parser.compare()
}

func (parser *Parser) variableAndConst() (interface{}, error) {
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

	return VarAssignNode{
		Identifier: identifier.(string),
		Node:       node,
		IsConstant: constError == nil,
	}, nil
}

func (parser *Parser) updateVariable() (interface{}, error) {
	value := parser.CurrentToken.Value
	_, err := parser.verifyNextToken(constants.TT_IDENTIFIER, constants.TT_EQ)
	if err != nil {
		return nil, nil
	}
	expr, err := parser.compare()

	if err != nil {
		return nil, err
	}

	return UpdateVariableNode{
		Identifier: value.(string),
		Node:       expr,
	}, nil
}

func (parser *Parser) while() (interface{}, error) {
	_, err := parser.verifyNextToken(constants.TT_WHILE)
	if err != nil {
		return nil, nil
	}
	conditionAndBodyBase, err := parser.conditionAndBodyBase()

	if err != nil {
		return nil, err
	}
	return WhileNode{
		Condition: conditionAndBodyBase.Condition,
		Body:      conditionAndBodyBase.Body,
	}, nil
}

func (parser *Parser) for_() (interface{}, error) {
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

	return ForNode{
		Expr1:     expr1,
		Condition: condition,
		Expr2:     expr2,
		Body:      body,
	}, nil
}

func (parser *Parser) AndOr() (interface{}, error) {
	return parser.binOP(parser.compare, constants.TT_AND, constants.TT_AND)
}

func (parser *Parser) compare() (interface{}, error) {
	return parser.binOP(parser.plus, constants.TT_GT, constants.TT_GTE, constants.TT_GT, constants.TT_LT, constants.TT_LTE, constants.TT_EQE)
}

func (parser *Parser) plus() (interface{}, error) {
	return parser.binOP(parser.factor, constants.TT_PLUS, constants.TT_MINUS)
}

func (parser *Parser) factor() (interface{}, error) {
	return parser.binOP(parser.pow, constants.TT_MUL, constants.TT_DIV)
}

func (parser *Parser) pow() (interface{}, error) {
	return parser.binOP(parser.term, constants.TT_POW, constants.TT_SQUARE_ROOT)
}

func (parser *Parser) term() (interface{}, error) {
	nodeType := parser.CurrentToken.Type_

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

		unaryOP := UnaryOP{
			Operation: nodeType,
			RigthNode: rigthNode,
		}
		return unaryOP, nil
	}

	if nodeType == "number" {
		value := parser.CurrentToken.Value.(float64)
		number := numbers.NewNumbers(value)
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

	if varAccess, err := parser.varAccess(); err != nil || varAccess != nil {
		return varAccess, nil
	}

	if funcNode, err := parser.func_(); funcNode != nil || err != nil {
		return funcNode, err
	}

	return nil, fmt.Errorf("")
}

func (parser *Parser) if_() (interface{}, error) {
	ifs := []*ConditionAndBody{}
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

	for (*parser.CurrentToken).Type_ == constants.TT_ELIF {
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

	return IfNode{
		Ifs:   ifs,
		Else_: elseNode,
	}, nil
}

func (parser *Parser) func_() (*FuncNode, error) {
	identoifierToken, err := parser.verifyNextToken(constants.TT_FUNC, constants.TT_IDENTIFIER)
	if err != nil {
		return nil, err
	}
	params, err := parser.params()
	if err != nil {
		return nil, err
	}
	body, err := parser.BodyBase()
	if err != nil {
		return nil, err
	}
	return &FuncNode{
		params,
		body,
		identoifierToken.Value.(string),
	}, nil
}

func (parser *Parser) callFunc() (*CallFuncNode, error) {
	funcName := parser.CurrentToken.Value
	_, err := parser.verifyNextToken(constants.TT_IDENTIFIER, constants.TT_LPAREN)
	if err != nil {
		return nil, nil
	}
	params, err := parser.args()
	if err != nil {
		return nil, err
	}
	return &CallFuncNode{
		Params: params,
		Name:   funcName.(string),
	}, nil
}

func (parser *Parser) conditionAndBodyBase() (*ConditionAndBody, error) {
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

	return &ConditionAndBody{
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

func (parser *Parser) varAccess() (*VarAccessNode, error) {
	if parser.CurrentToken.Type_ != constants.TT_IDENTIFIER {
		return nil, nil
	}

	varAccessNode := &VarAccessNode{
		Identifier: parser.CurrentToken.Value.(string),
	}
	parser.advance()
	return varAccessNode, nil
}

func (parser *Parser) params() (*[]lexerStructs.Token, error) {
	_, err := parser.verifyNextToken(constants.TT_LPAREN)
	params := &[]lexerStructs.Token{}
	if err != nil {
		return nil, err
	}
	for {
		*params = append(*params, *parser.CurrentToken)
		_, err = parser.verifyNextToken(constants.TT_IDENTIFIER)
		if err != nil {
			return nil, err
		}

		_, err = parser.verifyNextToken(constants.TT_COMMA)
		if err != nil {
			_, err = parser.verifyNextToken(constants.TT_RPAREN)
			if err == nil {
				break
			}
			return nil, err
		}

	}

	return params, nil
}

func (parser *Parser) args() (*[]interface{}, error) {
	params := []interface{}{}
	for {
		param, err := parser.compare()
		if err != nil {
			return nil, err
		}
		params = append(params, param)
		_, err = parser.verifyNextToken(constants.TT_COMMA)
		if err != nil {
			_, err = parser.verifyNextToken(constants.TT_RPAREN)
			if err == nil {
				break
			}
			return nil, err
		}
	}

	return &params, nil
}
