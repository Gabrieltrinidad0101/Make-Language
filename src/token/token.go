package token

type Token struct {
	Position
	Type_ string
	Value interface{}
}

type Position struct {
	PositionEnd   int
	PositionStart int
	Line          int
	Col           int
}
