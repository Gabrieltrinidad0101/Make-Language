package token

type Token struct {
	Type_         string
	Value         interface{}
	PositionEnd   int
	PositionStart int
}
