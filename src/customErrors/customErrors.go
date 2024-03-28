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
	errorText := ""
	for i, line := range lines {
		if i == 0 {
			padding := strings.Repeat(" ", len(line)-token.PositionStart)
			errorSignal := strings.Repeat("^", max(token.PositionEnd-token.PositionStart, 1))
			errorText += fmt.Sprintf("%s\n %s\n", line, padding+errorSignal)
		}
	}

	fmt.Printf("%s:\n %s", details, errorText)
	os.Exit(0)
}

func IllegalCharacter(token token.Token) {
	show(token, fmt.Sprintf("Illegal Character: %s", token.Value))
}

func InvalidSyntax(token token.Token, details string) {
	show(token, fmt.Sprintf("Invalid: %s\n%s", token.Value, details))
}
