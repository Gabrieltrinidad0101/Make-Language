package customErrors

import (
	"fmt"
	"makeLanguages/src/token"
	"os"
	"strings"
)

type customError struct {
	text string
}

var CustomError *customError

func New(text string) *customError {
	if CustomError != nil {
		return CustomError
	}
	CustomError = &customError{
		text,
	}
	return CustomError
}

func show(token token.Token, details string) {
	lines := strings.Split(CustomError.text, "\n")
	fmt.Printf("Line: %d\n", token.Line)
	fmt.Printf("Col: %d\n", token.Col)
	errorText := ""
	for i, line := range lines {
		if i == 0 {
			characterLength := max(token.PositionEnd-token.PositionStart, 1)
			padding := strings.Repeat(" ", len(string(line))-characterLength-1)
			errorSignal := strings.Repeat("^", characterLength)
			errorText += fmt.Sprintf("%s\n %s\n", line, padding+errorSignal)
			continue
		}
		errorText += string(line)
	}

	fmt.Printf("%s\n", details)
	fmt.Printf("%s\n", errorText)
	os.Exit(0)
}

func IllegalCharacter(token token.Token) {
	show(token, fmt.Sprintf("Illegal Character: %s", token.Value))
}

func InvalidSyntax(token token.Token, details string) {
	fmt.Printf("Invalid Syntax: %s\n", token.Value)
	show(token, details)
}
