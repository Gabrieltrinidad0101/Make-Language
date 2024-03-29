package lexer

import (
	"encoding/json"
	"fmt"
	"makeLanguages/src/constants"
	CustomErrors "makeLanguages/src/customErrors"
	"makeLanguages/src/token"
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

type Simbols struct {
	Text      string `json:"text"`
	TokenName string `json:"token_name"`
}

const LETTERS = "qwertyuiopasdfghjklñzxcvbnmQWERTYUIOPASDFGHJKLÑZXCVBNM_"
const ASCII = "123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

type Lexer struct {
	token.Position
	idx                    int
	current_char           *string
	tokens                 *[]token.Token
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
		Position: token.Position{
			Line: 1,
		},
		current_char:           &currentChar,
		len:                    len(*text),
		tokens:                 &[]token.Token{},
		idx:                    -1,
		languageConfiguraction: languageConfiguraction,
		characterMaxLength:     getMaxLengthCharacter(languageConfiguraction.LanguageSyntax),
	}
}

func (lexer *Lexer) Tokens() (*[]token.Token, bool) {
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

		isCompare := lexer.syntaxToken(lexer.languageConfiguraction.Compares)

		if isCompare {
			continue
		}

		ok := lexer.syntaxToken(lexer.languageConfiguraction.LanguageSyntax)

		if !ok {
			ok = lexer.getIdentifier()
		}

		if !ok {
			CustomErrors.IllegalCharacter(token.Token{
				Value:    *lexer.current_char,
				Position: lexer.Position,
			})
			return nil, true
		}
	}
	*lexer.tokens = append(*lexer.tokens, token.Token{
		Type_: constants.EOF,
	})
	return lexer.tokens, false
}

func (lexer *Lexer) syntaxToken(syntax map[string]string) bool {
	var type_ *string = nil
	var simbolText string = ""

	positionCopy := lexer.Position.Copy()
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

	token := token.Token{
		Type_:    *type_,
		Value:    simbolText,
		Position: positionCopy,
	}

	*lexer.tokens = append(*lexer.tokens, token)

	return true
}

func (lexer *Lexer) getIdentifier() bool {
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

	if identifier == "" {
		return false
	}

	token := token.Token{
		Type_: constants.TT_IDENTIFIER,
		Value: identifier,
	}

	*lexer.tokens = append(*lexer.tokens, token)
	lexer.advance()
	return true
}

func (lexer *Lexer) makeNumber() bool {
	numberString, ok := lexer.languageConfiguraction.Numbers[*lexer.current_char]
	if !ok {
		return false
	}
	dotNumber := 0
	positionStart := lexer.Position.Copy()
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

	*lexer.tokens = append(*lexer.tokens, token.Token{
		Type_:    "number",
		Value:    number,
		Position: positionStart,
	})
	return true
}
