package customErrors

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"
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

func Show(token lexerStructs.IPositionBase, details string) error {
	lines := strings.Split(CustomError.text, "\n")
	positionStart := token.GetPositionStart()
	positionEnd := token.GetPositionEnd()
	linesCanShow := min(3, positionStart.Line-1)
	startText := positionStart.Line - linesCanShow - 1
	endText := positionStart.Line + min(3, len(lines)-(positionStart.Line-1)) - 1
	linesCut := lines[startText:endText]
	err := fmt.Sprintf("Line: %d\n", positionStart.Line)
	err += fmt.Sprintf("Col: %d\n", positionStart.Col)
	errorText := ""
	for i, line := range linesCut {
		lineNumber := i + positionStart.Line - linesCanShow
		if i == linesCanShow {
			characterLength := max(1, positionEnd.Col-positionStart.Col)
			padding := strings.Repeat(" ", positionStart.Col+len(strconv.Itoa(lineNumber))+1)
			errorSignal := strings.Repeat("^", characterLength)
			errorText += fmt.Sprintf("%d: %s\n%s\n", lineNumber, line, padding+errorSignal)
			continue
		}
		errorText += fmt.Sprintf("%d: %s\n", lineNumber, line)
	}

	err += fmt.Sprintf("%s\n\n", details)
	err += fmt.Sprintf("%s\n\n", errorText)
	return fmt.Errorf(err)
}

func IllegalCharacter(token lexerStructs.Token) error {
	return Show(token.IPositionBase, fmt.Sprintf("Illegal Character: %s", token.Value))
}

func InvalidSyntax(token lexerStructs.Token, details string) error {
	fmt.Printf("Invalid Syntax: %s\n", token.Value)
	return Show(token.IPositionBase, details)
}

func RunTimeError(token lexerStructs.IPositionBase, details string) error {
	fmt.Println("Run Time Error")
	return Show(token, details)
}
