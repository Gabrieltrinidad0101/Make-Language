package lexer

import (
	"encoding/json"
	"fmt"
	"makeLanguages/src/constants"
	CustomErrors "makeLanguages/src/customErrors"
	"os"
	"strconv"
	"strings"
)

type LanguageConfiguraction struct {
	Numbers           map[string]string   `json:"numbers"`
	Compares          map[string]string   `json:"compares"`
	ComparesContinues map[string][]string `json:"compares_continues"`
	LanguageSyntax    map[string]string   `json:"language_syntax"`
}

type Token struct {
	Type_         string
	Value         interface{}
	PositionEnd   int
	PositionStart int
}

type Simbols struct {
	Text      string `json:"text"`
	TokenName string `json:"token_name"`
}

const LETTERS = "qwertyuiopasdfghjklñzxcvbnmQWERTYUIOPASDFGHJKLÑZXCVBNM_"
const ASCII = "123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

type Lexer struct {
	col                    int
	line                   int
	idx                    int
	current_char           *string
	tokens                 *[]Token
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

func NewLexer(text *string, languageConfiguraction LanguageConfiguraction) *Lexer {
	currentChar := ""
	return &Lexer{
		text:                   text,
		current_char:           &currentChar,
		len:                    len(*text),
		tokens:                 &[]Token{},
		idx:                    -1,
		languageConfiguraction: languageConfiguraction,
		characterMaxLength:     getMaxLengthCharacter(languageConfiguraction.LanguageSyntax),
	}
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

		ok := lexer.syntaxToken()

		if !ok {
			ok = lexer.getIdentifier()
		}

		if !ok {
			customErrors := CustomErrors.New()
			customErrors.IllegalCharacter(*lexer.text, *lexer.current_char, lexer.idx, lexer.idx)
			return nil, true
		}
	}
	*lexer.tokens = append(*lexer.tokens, Token{
		Type_: constants.EOF,
	})
	return lexer.tokens, false
}

func (lexer *Lexer) syntaxToken() bool {
	simbols, ok := lexer.languageConfiguraction.LanguageSyntax[*lexer.current_char]

	if ok {
		token := &Token{
			Type_: simbols,
			Value: nil,
		}
		*lexer.tokens = append(*lexer.tokens, *token)
		lexer.advance()
		return true
	}

	var type_ *string = nil
	var simbolText string = ""

main:
	for i := 0; i < lexer.len; i++ {
		if i >= lexer.characterMaxLength || lexer.idx+i >= lexer.len {
			return false
		}

		simbolText += string((*lexer.text)[lexer.idx+i])

		for simbol, value := range lexer.languageConfiguraction.LanguageSyntax {

			if simbol != simbolText {
				continue
			}

			type_ = &value
			break main
		}
	}

	if type_ == nil {
		return false
	}

	for i := 0; i < len(simbolText); {
		i++
		lexer.advance()
	}

	token := Token{
		Type_: *type_,
		Value: nil,
	}

	*lexer.tokens = append(*lexer.tokens, token)

	return true
}

func (lexer *Lexer) getIdentifier() bool {
	if !strings.Contains(LETTERS, *lexer.current_char) {
		return false
	}
	identifier := ""
	for {
		if !strings.Contains(LETTERS, *lexer.current_char) {
			break
		}
		identifier += *lexer.current_char
		ok := lexer.advance()
		if !ok {
			break
		}
	}

	token := Token{
		Type_: constants.TT_IDENTIFIER,
		Value: identifier,
	}

	*lexer.tokens = append(*lexer.tokens, token)
	lexer.advance()
	return true
}

func (lexer *Lexer) advance() bool {
	lexer.idx++
	if lexer.idx >= lexer.len {
		lexer.current_char = nil
		return false
	}

	*lexer.current_char = string((*lexer.text)[lexer.idx])
	if *lexer.current_char == "\n" {
		lexer.line++
		lexer.col = 0
	}

	return true
}

func (lexer *Lexer) makeNumber() bool {
	numberString, ok := lexer.languageConfiguraction.Numbers[*lexer.current_char]
	if !ok {
		return false
	}
	dotNumber := 0
	positionStart := lexer.idx
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

	*lexer.tokens = append(*lexer.tokens, Token{
		Type_:         "number",
		Value:         number,
		PositionStart: positionStart,
		PositionEnd:   lexer.idx,
	})
	return true
}

func (lexer *Lexer) makeCompares() bool {
	compare, ok := lexer.languageConfiguraction.Compares[*lexer.current_char]
	if !ok {
		return false
	}
	lexer.advance()
	continues := lexer.languageConfiguraction.ComparesContinues[compare]
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
