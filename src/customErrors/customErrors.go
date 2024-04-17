package customErrors

import (
	"fmt"
	"makeLanguages/src/lexer/lexerStructs"
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

func Show(token lexerStructs.IPositionBase, details string) {
	lines := strings.Split(CustomError.text, "\n")
	positionStart := token.GetPositionStart()
	positionEnd := token.GetPositionEnd()
	linesCanShow := min(3, positionStart.Line-1)
	startText := positionStart.Line - linesCanShow - 1
	endText := positionStart.Line + min(3, len(lines)-(positionStart.Line-1)) - 1
	linesCut := lines[startText:endText]
	fmt.Printf("Line: %d\n", positionStart.Line)
	fmt.Printf("Col: %d\n", positionStart.Col)
	errorText := ""
	for i, line := range linesCut {
		lineNumber := i + positionStart.Line - linesCanShow
		if i == linesCanShow {
			characterLength := max(1, positionEnd.Col-positionStart.Col)
			padding := strings.Repeat(" ", positionStart.Col+len(string(rune(lineNumber)))+2)
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

func IllegalCharacter(token lexerStructs.Token) {
	Show(token.IPositionBase, fmt.Sprintf("Illegal Character: %s", token.Value))
}

func InvalidSyntax(token lexerStructs.Token, details string) {
	fmt.Printf("Invalid Syntax: %s\n", token.Value)
	Show(token.IPositionBase, details)
}

func RunTimeError(token lexerStructs.IPositionBase, details string) {
	fmt.Println("Run Time Error")
	Show(token, details)
}
