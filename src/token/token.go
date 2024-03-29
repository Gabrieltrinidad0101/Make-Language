package token

type Token struct {
	Position
	Type_ string
	Value interface{}
}

type Position struct {
	Line int
	Col  int
}

func (position *Position) Copy() Position {
	return *position
}
