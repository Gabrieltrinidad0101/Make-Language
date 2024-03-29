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
	linesCanShow := min(3, token.Line-1)
	startText := token.Line - linesCanShow - 1
	endText := token.Line + min(3, len(lines)-token.Line-1)
	linesCut := lines[startText:endText]
	fmt.Printf("Line: %d\n", token.Line)
	fmt.Printf("Col: %d\n", token.Col)
	errorText := ""
	for i, line := range linesCut {
		lineNumber := i + token.Line - linesCanShow
		if i == linesCanShow {
			characterLength := len(fmt.Sprint(token.Value))
			padding := strings.Repeat(" ", token.Col+len(string(rune(lineNumber)))+1)
			errorSignal := strings.Repeat("^", characterLength)
			errorText += fmt.Sprintf("%d: %s\n%s\n", lineNumber, line, padding+errorSignal)
			continue
		}
		errorText += fmt.Sprintf("%d: %s\n", lineNumber, line)
	}

	fmt.Printf("%s\n\n", details)
	fmt.Printf("%s\n\n", errorText)
	os.Exit(0)
}

func IllegalCharacter(token token.Token) {
	show(token, fmt.Sprintf("Illegal Character: %s", token.Value))
}

func InvalidSyntax(token token.Token, details string) {
	fmt.Printf("Invalid Syntax: %s\n", token.Value)
	show(token, details)
}
