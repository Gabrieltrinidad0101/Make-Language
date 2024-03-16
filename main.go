package main

import (
	"makeLanguages/src/lexer"
)

func main() {

	input := "() 123 + 321321 * 131231"
	lexer := lexer.NewLexer(&input)
	lexer.Tokens()

}
