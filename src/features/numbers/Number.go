package numbers

import (
	"math"

	"github.com/Gabrieltrinidad0101/Make-Language/src/lexer/lexerStructs"

	"github.com/Gabrieltrinidad0101/Make-Language/src/features/booleans"
)

type Number struct {
	Value float64
	lexerStructs.IPositionBase
}

func NewNumbers(value float64, position lexerStructs.IPositionBase) *Number {
	return &Number{
		Value:         value,
		IPositionBase: position,
	}
}

func (number *Number) GetValue() interface{} {
	return number.Value
}

func (number *Number) PLUS(number_ *Number) *Number {
	value := number.Value + number_.Value
	return NewNumbers(value, nil)
}

func (number *Number) MINUS(number_ *Number) *Number {
	value := number.Value - number_.Value
	return NewNumbers(value, nil)
}

func (number *Number) PLUS1() *Number {
	value := number.Value + 1
	return NewNumbers(value, nil)
}

func (number *Number) MINUS1(number_ *Number) *Number {
	value := number.Value - 1
	return NewNumbers(value, nil)
}

func (number *Number) MUL(number_ *Number) *Number {
	value := number.Value * number_.Value
	return NewNumbers(value, nil)
}

func (number *Number) DIV(number_ *Number) *Number {
	value := number.Value / number_.Value
	return NewNumbers(value, nil)
}

func (number *Number) POW(number_ *Number) *Number {
	value := math.Pow(number.Value, number_.Value)
	return NewNumbers(value, nil)
}

func (number *Number) SQUARE_ROOT(number_ *Number) *Number {
	value := math.Pow(number.Value, 1.0/number_.Value)
	return NewNumbers(value, nil)
}

func (number *Number) GT(number_ *Number) *booleans.Boolean {
	value := number.Value > number_.Value
	return booleans.NewBoolean(value)
}

func (number *Number) GTE(number_ *Number) *booleans.Boolean {
	value := number.Value >= number_.Value
	return booleans.NewBoolean(value)
}

func (number *Number) LT(number_ *Number) *booleans.Boolean {
	value := number.Value < number_.Value
	return booleans.NewBoolean(value)
}

func (number *Number) LTE(number_ *Number) *booleans.Boolean {
	value := number.Value <= number_.Value
	return booleans.NewBoolean(value)
}

func (number *Number) EQE(number_ *Number) *booleans.Boolean {
	value := number.Value == number_.Value
	return booleans.NewBoolean(value)
}