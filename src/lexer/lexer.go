package lexer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/constants"

	CustomErrors "github.com/Gabrieltrinidad0101/Make-Language/src/customErrors"
)

type LanguageConfiguraction struct {
	Numbers           map[string]string   `json:"numbers"`
	Compares          map[string]string   `json:"compares"`
	ComparesContinues map[string][]string `json:"compares_continues"`
	LanguageSyntax    map[string]string   `json:"language_syntax"`
	CustomOperators   map[string]string   `json:"custom_operators"`
	Functions         map[string]string   `json:"functions"`
	Scope             string              `json:"scope"`
	ConstructorName   string              `json:"constructor_name"`
	Paren             bool                `json:"paren"`
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

func ReadLanguageConfiguraction(path string) (LanguageConfiguraction, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return LanguageConfiguraction{}, err
	}

	// Define a variable to hold the data
	var languageConfiguraction LanguageConfiguraction

	// Unmarshal the JSON data into the defined structure
	err = json.Unmarshal(file, &languageConfiguraction)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return LanguageConfiguraction{}, err
	}

	return languageConfiguraction, err
}

func ReadFile(path string) (*string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(file)
	return &text, nil
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

func (lexer *Lexer) Tokens() (*[]lexerStructs.Token, error) {
	lexer.advance()
	for lexer.current_char != nil {
		if *lexer.current_char == " " {
			lexer.advance()
			continue
		}

		isString, err := lexer.makeString()

		if err != nil {
			return nil, err
		}

		if isString {
			continue
		}

		isNumber := lexer.makeNumber()
		if isNumber {
			continue
		}

		operatorCustom := lexer.syntaxToken(lexer.languageConfiguraction.CustomOperators)

		if operatorCustom {
			continue
		}

		isCompare := lexer.syntaxToken(lexer.languageConfiguraction.Compares)

		if isCompare {
			continue
		}

		ok := lexer.syntaxToken(lexer.languageConfiguraction.LanguageSyntax)

		if ok {
			continue
		}

		ok, err = lexer.getIdentifier()

		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, CustomErrors.IllegalCharacter(lexerStructs.Token{
				Value: *lexer.current_char,
				IPositionBase: lexerStructs.PositionBase{
					PositionStart: lexer.Position,
					PositionEnd:   lexer.Position,
				},
			})
		}

	}
	*lexer.tokens = append(*lexer.tokens, lexerStructs.Token{
		Type_: constants.EOF,
	})
	return lexer.tokens, nil
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

func (lexer *Lexer) getIdentifier() (bool, error) {
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
			return false, CustomErrors.Show(
				lexerStructs.Token{
					IPositionBase: lexerStructs.PositionBase{
						PositionStart: position,
						PositionEnd:   position,
					},
				},
				"Variables Identifier cannot have spot",
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
		return false, nil
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
	return true, nil
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

func (lexer *Lexer) makeString() (bool, error) {
	languageSyntax := lexer.languageConfiguraction.LanguageSyntax
	if constants.TT_STRING != languageSyntax[*lexer.current_char] {
		return false, nil
	}

	stringValue := ""
	positionStart := lexer.PositionCopy()
	lexer.advance()
	positionEnd := lexer.PositionCopy()
	for constants.TT_STRING != languageSyntax[*lexer.current_char] {
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
			return false, CustomErrors.InvalidSyntax(token, "Is necesary to use \" to end a string ")
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

	return true, nil
}
