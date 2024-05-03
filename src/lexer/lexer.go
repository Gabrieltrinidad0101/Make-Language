package lexer

import (
	"encoding/json"
	"fmt"
	"makeLanguages/src/constants"
	CustomErrors "makeLanguages/src/customErrors"
	"makeLanguages/src/lexer/lexerStructs"
	"os"
	"strconv"
	"strings"
)

type LanguageConfiguraction struct {
	Numbers           map[string]string   `json:"numbers"`
	Compares          map[string]string   `json:"compares"`
	ComparesContinues map[string][]string `json:"compares_continues"`
	LanguageSyntax    map[string]string   `json:"language_syntax"`
	OperatorCustom    map[string]string   `json:"operator_custom"`
	Functions         map[string]string   `json:"functions"`
	Scope             string              `json:"scope"`
}

type Simbols struct {
	Text      string `json:"text"`
	TokenName string `json:"token_name"`
}

const LETTERS = "qwertyuiopasdfghjklñzxcvbnmQWERTYUIOPASDFGHJKLÑZXCVBNM_ "

type Lexer struct {
	lexerStructs.Position
	idx                    int
	current_char           *string
	tokens                 *[]lexerStructs.Token
	text                   *string
	len                    int
	languageConfiguraction LanguageConfiguraction
	characterMaxLength     int
}

func ReadLanguageConfiguraction(path string) (LanguageConfiguraction, bool) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return LanguageConfiguraction{}, false
	}

	// Define a variable to hold the data
	var languageConfiguraction LanguageConfiguraction

	// Unmarshal the JSON data into the defined structure
	err = json.Unmarshal(file, &languageConfiguraction)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return LanguageConfiguraction{}, false
	}

	return languageConfiguraction, true
}

func ReadFile(path string) (*string, bool) {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, false
	}
	text := string(file)
	return &text, true
}

func getMaxLengthCharacter(characters map[string]string) int {
	maxLength := 0
	for key := range characters {
		length := len(key)
		if length > maxLength {
			maxLength = length
		}
	}

	return maxLength
}

func (lexer *Lexer) advance() bool {
	lexer.idx++
	if lexer.idx >= lexer.len {
		lexer.current_char = nil
		return false
	}

	lexer.Col++
	*lexer.current_char = string((*lexer.text)[lexer.idx])
	if *lexer.current_char == "\n" {
		lexer.Line++
		lexer.Col = 0
	}

	return true
}

func NewLexer(text *string, languageConfiguraction LanguageConfiguraction) *Lexer {
	currentChar := ""
	return &Lexer{
		text: text,
		Position: lexerStructs.Position{
			Line: 1,
		},
		current_char:           &currentChar,
		len:                    len(*text),
		tokens:                 &[]lexerStructs.Token{},
		idx:                    -1,
		languageConfiguraction: languageConfiguraction,
		characterMaxLength:     getMaxLengthCharacter(languageConfiguraction.LanguageSyntax),
	}
}

func (lexer *Lexer) Tokens() (*[]lexerStructs.Token, bool) {
	lexer.advance()
	for lexer.current_char != nil {
		if *lexer.current_char == " " {
			lexer.advance()
			continue
		}

		isString := lexer.makeString()

		if isString {
			continue
		}

		isNumber := lexer.makeNumber()
		if isNumber {
			continue
		}

		operatorCustom := lexer.syntaxToken(lexer.languageConfiguraction.OperatorCustom)

		if operatorCustom {
			continue
		}

		isCompare := lexer.syntaxToken(lexer.languageConfiguraction.Compares)

		if isCompare {
			continue
		}

		ok := lexer.syntaxToken(lexer.languageConfiguraction.LanguageSyntax)

		if !ok {
			ok = lexer.getIdentifier()
		}

		if !ok {
			CustomErrors.IllegalCharacter(lexerStructs.Token{
				Value: *lexer.current_char,
				IPositionBase: lexerStructs.PositionBase{
					PositionStart: lexer.Position,
					PositionEnd:   lexer.Position,
				},
			}, constants.STOP_EXECUTION)
			return nil, true
		}
	}
	*lexer.tokens = append(*lexer.tokens, lexerStructs.Token{
		Type_: constants.EOF,
	})
	return lexer.tokens, false
}

func (lexer *Lexer) syntaxToken(syntax map[string]string) bool {
	var type_ *string = nil
	var simbolText string = ""

	positionCopy := lexer.PositionCopy()
	for i := lexer.idx; i < lexer.len && i-lexer.idx < lexer.characterMaxLength; i++ {
		simbolText += string((*lexer.text)[i])

		simbol, ok := syntax[simbolText]

		if !ok {
			continue
		}
		type_ = &simbol
		break
	}

	if type_ == nil {
		return false
	}

	for i := 0; i < len(simbolText); i++ {
		lexer.advance()
	}

	token := lexerStructs.Token{
		Type_: *type_,
		Value: simbolText,
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: positionCopy,
			PositionEnd:   lexer.PositionCopy(),
		},
	}

	*lexer.tokens = append(*lexer.tokens, token)

	return true
}

func (lexer *Lexer) getIdentifier() bool {
	identifier := ""
	positionStart := lexer.PositionCopy()
	positionEnd := positionStart
	hasString := false
	for {
		lastToken := lexerStructs.Token{}
		if len(*lexer.tokens) > 0 {
			lastToken = (*lexer.tokens)[len(*lexer.tokens)-1]
		}

		if lastToken.Type_ == constants.TT_VAR && *lexer.current_char == "." {
			position := lexerStructs.Position{
				Line: lexer.Line,
				Col:  lexer.idx + 1,
			}
			CustomErrors.Show(
				lexerStructs.Token{
					IPositionBase: lexerStructs.PositionBase{
						PositionStart: position,
						PositionEnd:   position,
					},
				},
				"Variables Identifier cannot have spot",
				constants.STOP_EXECUTION,
			)
		}

		canNotHasNumber := hasString && !strings.Contains("1234567890", *lexer.current_char)

		if !hasString {
			hasString = strings.Contains(LETTERS, *lexer.current_char)
		}
		isNotLegalCharacter := !strings.Contains(LETTERS, *lexer.current_char)
		if isNotLegalCharacter && canNotHasNumber {
			break
		}
		identifier += *lexer.current_char
		positionEnd = lexer.PositionCopy()
		ok := lexer.advance()
		if !ok {
			break
		}
	}

	if identifier == "" {
		return false
	}

	token := lexerStructs.Token{
		Type_: constants.TT_IDENTIFIER,
		Value: strings.Trim(identifier, " "),
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: positionStart,
			PositionEnd:   positionEnd,
		},
	}

	*lexer.tokens = append(*lexer.tokens, token)
	return true
}

func (lexer *Lexer) makeNumber() bool {
	numberString, ok := lexer.languageConfiguraction.Numbers[*lexer.current_char]
	if !ok {
		return false
	}
	dotNumber := 0
	positionStart := lexer.PositionCopy()
	for lexer.current_char != nil && lexer.advance() {
		numberNext, ok := lexer.languageConfiguraction.Numbers[*lexer.current_char]
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

	*lexer.tokens = append(*lexer.tokens, lexerStructs.Token{
		Type_: "number",
		Value: number,
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: positionStart,
			PositionEnd:   lexer.PositionCopy(),
		},
	})
	return true
}

func (lexer *Lexer) makeString() bool {

	if constants.TT_STRING != lexer.languageConfiguraction.LanguageSyntax[*lexer.current_char] {
		return false
	}

	stringValue := ""
	positionStart := lexer.PositionCopy()
	lexer.advance()
	positionEnd := lexer.PositionCopy()
	for constants.TT_STRING != lexer.languageConfiguraction.LanguageSyntax[*lexer.current_char] {
		stringValue += *lexer.current_char
		ok := lexer.advance()
		positionEnd = lexer.PositionCopy()
		if !ok || *lexer.current_char == "\n" {
			token := lexerStructs.Token{
				Value: "\"",
				IPositionBase: lexerStructs.PositionBase{
					PositionStart: positionStart,
					PositionEnd:   positionEnd,
				},
			}
			CustomErrors.InvalidSyntax(token, "Is necesary to use \" to end a string ", constants.STOP_EXECUTION)
		}
	}

	lexer.advance()

	stringToken := lexerStructs.Token{
		Value: stringValue,
		Type_: constants.TT_STRING,
		IPositionBase: lexerStructs.PositionBase{
			PositionStart: positionStart,
			PositionEnd:   positionEnd,
		},
	}

	*lexer.tokens = append(*lexer.tokens, stringToken)

	return true
}
