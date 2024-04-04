package lexerStructs

type Token struct {
	PositionBase
	Type_ string
	Value interface{}
}

type IPositionBase interface {
	GetPositionStart() Position
	GetPositionEnd() Position
}

type PositionBase struct {
	PositionStart Position
	PositionEnd   Position
}

func (NodeBase *PositionBase) GetPositionStart() Position {
	return NodeBase.PositionStart
}

func (NodeBase *PositionBase) GetPositionEnd() Position {
	return NodeBase.PositionEnd
}

type Position struct {
	Line int
	Col  int
}

func (position *Position) PositionCopy() Position {
	return *position
}
