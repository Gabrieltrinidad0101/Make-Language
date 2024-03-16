package numbers

type Number struct {
	Value int
}

func NewNumbers(value int) *Number {
	return &Number{
		value,
	}
}

func (number *Number) PLUS(number_ *Number) *Number {
	value := number.Value + number_.Value
	return NewNumbers(value)
}

func (number *Number) MINUS(number_ *Number) *Number {
	value := number.Value - number_.Value
	return NewNumbers(value)
}

func (number *Number) MUL(number_ *Number) *Number {
	value := number.Value * number_.Value
	return NewNumbers(value)
}

func (number *Number) DIV(number_ *Number) *Number {
	value := number.Value / number_.Value
	return NewNumbers(value)
}
