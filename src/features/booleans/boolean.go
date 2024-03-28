package booleans

type Boolean struct {
	value bool
}

func NewBoolean(value_ bool) *Boolean {
	return &Boolean{
		value: value_,
	}
}

func (boolean *Boolean) AND(boolean_ *Boolean) *Boolean {
	if boolean.value && boolean_.value {
		return NewBoolean(true)
	}
	return NewBoolean(false)
}

func (boolean *Boolean) OR(boolean_ *Boolean) *Boolean {
	if boolean.value || boolean_.value {
		return NewBoolean(true)
	}
	return NewBoolean(false)
}
