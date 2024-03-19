package lexer

import (
	"fmt"
	CustomErrors "makeLanguages/src/customErrors"
	"strconv"
)

type Lexer struct {
	col          int
	line         int
	idx          int
	current_char *string
	tokens       *[]Token
	text         *string
	len          int
}

func NewLexer(text *string) *Lexer {
	currentChar := ""
	return &Lexer{
		text:         text,
		current_char: &currentChar,
		len:          len(*text),
		tokens:       &[]Token{},
	}
}

type Token struct {
	Type_         string
	Value         interface{}
	PositionEnd   int
	PositionStart int
}

type Simbols struct {
	Text      string
	TokenName string
}

var numbers = map[string]string{
	"1": "1",
	"2": "2",
	"3": "3",
	"4": "4",
	"5": "5",
	"6": "6",
	"7": "7",
	"8": "8",
	"9": "9",
	"0": "0",
}

var compares = map[string]string{
	"=": "EQ",
	">": "GT",
	"<": "LT",
	"!": "NEQ",
}

var comparesContinues = map[string][]string{
	"EQ":  {"="},
	"GT":  {"="},
	"LT":  {"="},
	"NEQ": {"="},
}

var LanguageSyntax = map[string][]Simbols{
	"+": {
		{
			Text:      "+",
			TokenName: "PLUS",
		},
	},
	"-": {
		{
			Text:      "-",
			TokenName: "MINUS",
		},
	},
	"*": {
		{
			Text:      "*",
			TokenName: "MUL",
		},
	},
	"/": {
		{
			Text:      "/",
			TokenName: "DIV",
		},
	},
	"(": {
		{
			Text:      "(",
			TokenName: "LPAREN",
		},
	},
	")": {
		{
			Text:      ")",
			TokenName: "RPAREN",
		},
	},
	"^": {
		{
			Text:      "^",
			TokenName: "POW",
		},
	},
	"~": {
		{
			Text:      "~",
			TokenName: "SQUARE_ROOT",
		}},
	"if": {
		{
			Text:      "if",
			TokenName: "IF",
		},
	},
	"{": {
		{
			Text:      "{",
			TokenName: "START",
		},
	},
	"}": {
		{
			Text:      "}",
			TokenName: "END",
		},
	},
}

func (lexer *Lexer) Tokens() (*[]Token, bool) {
	lexer.advance()
	for lexer.current_char != nil {
		if *lexer.current_char == " " {
			lexer.advance()
			continue
		}

		isNumber := lexer.makeNumber()
		if isNumber {
			continue
		}

		isCompare := lexer.makeCompares()
		if isCompare {
			continue
		}

		token, ok := lexer.getToken()

		if !ok {
			customErrors := CustomErrors.New()
			customErrors.IllegalCharacter(*lexer.text, *lexer.current_char, lexer.idx, lexer.idx)
			return nil, true
		}

		*lexer.tokens = append(*lexer.tokens, *token)
		lexer.advance()
	}
	return lexer.tokens, false
}

func (lexer *Lexer) getToken() (*Token, bool) {
	simbols, ok := LanguageSyntax[*lexer.current_char]

	if !ok {
		return nil, false
	}

	if len(simbols) == 1 && len(simbols[0].Text) == 1 {
		return &Token{
			Type_: simbols[0].TokenName,
			Value: nil,
		}, true
	}

	var simbolText string = ""
	for i := 0; simbolText != " " || simbolText != "\r" || simbolText != "\n" || simbolText != "\t"; i++ {
		simbolText += string((*lexer.text)[lexer.idx+i])
	}

	var type_ *string = nil
	for _, simbol := range simbols {
		if simbolText != simbol.Text {
			continue
		}

		length := len(simbol.Text)
		for i := 0; i < length; i++ {
			lexer.advance()
		}

		*type_ = simbol.TokenName
		break
	}

	if type_ == nil {
		return nil, true
	}

	return &Token{
		Type_: *type_,
		Value: nil,
	}, false
}

func (lexer *Lexer) advance() bool {
	if lexer.idx > lexer.len-1 {
		lexer.current_char = nil
		return false
	}
	*lexer.current_char = string((*lexer.text)[lexer.idx])
	lexer.idx++
	if *lexer.current_char == "\n" {
		lexer.line++
		lexer.col = 0
	}

	return true
}

func (lexer *Lexer) makeNumber() bool {
	numberString, ok := numbers[*lexer.current_char]
	if !ok {
		return false
	}
	dotNumber := 0
	positionStart := lexer.idx
	for lexer.current_char != nil && lexer.advance() {
		numberNext, ok := numbers[*lexer.current_char]
		if !ok {
			if *lexer.current_char == "." && dotNumber < 1 {
				dotNumber++
				continue
			}
			break
		}
		numberString += numberNext
	}

	number, err := strconv.ParseFloat(numberString, 64)

	if err != err {
		panic(fmt.Sprintf("Internal error analize the number %s", numberString))
	}

	*lexer.tokens = append(*lexer.tokens, Token{
		Type_:         "number",
		Value:         number,
		PositionStart: positionStart,
		PositionEnd:   lexer.idx,
	})
	return true
}

func (lexer *Lexer) makeCompares() bool {
	compare, ok := compares[*lexer.current_char]
	if !ok {
		return false
	}
	lexer.advance()
	continues := comparesContinues[compare]
	if continues[0] != *lexer.current_char {
		*lexer.tokens = append(*lexer.tokens,
			Token{
				Type_: compare,
				Value: nil,
			},
		)
		return true
	}
	for _, character := range continues {
		if *lexer.current_char != character {
			panic("Error compare")
		}
		lexer.advance()
	}

	*lexer.tokens = append(*lexer.tokens,
		Token{
			Type_: compare + "E",
			Value: nil,
		},
	)

	return true
}
