package lexer

import (
	"fmt"
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
	Type_ string
	Value interface{}
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

var LanguageSyntax = map[string]string{
	"+": "PLUS",
	"-": "MINUS",
	"*": "MUL",
	"/": "DIV",
	"(": "LPAREN",
	")": "RPAREN",
	"^": "POW",
	"~": "SQUARE_ROOT",
}

func (lexer *Lexer) Tokens() *[]Token {
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

		type_, ok := LanguageSyntax[*lexer.current_char]
		if !ok {
			fmt.Print("Error")
			break
		}
		*lexer.tokens = append(*lexer.tokens, Token{
			Type_: type_,
			Value: nil,
		})
		lexer.advance()
	}
	return lexer.tokens
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
		Type_: "number",
		Value: number,
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
