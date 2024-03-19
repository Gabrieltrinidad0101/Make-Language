package CustomErrors

import (
	"fmt"
	"os"
	"strings"
)

type CustomError struct{}

func New() *CustomError {
	return &CustomError{}
}

func (customError *CustomError) show(text string, positionStart int, positionEnd int, details string) {
	lines := strings.Split(text, "\n")
	errorText := ""
	for i, line := range lines {
		if i == 0 {
			padding := strings.Repeat(" ", len(line)-positionStart)
			errorSignal := strings.Repeat("^", max(positionEnd-positionStart, 1))
			errorText += fmt.Sprintf("%s\n %s\n", line, padding+errorSignal)
		}
	}

	fmt.Printf("%s:\n %s", details, errorText)
	os.Exit(0)
}

func (customErrors *CustomError) IllegalCharacter(text, character string, positionStart, positionEnd int) {
	customErrors.show(text, positionStart, positionEnd, fmt.Sprintf("Illegal Character: %s", character))
}
