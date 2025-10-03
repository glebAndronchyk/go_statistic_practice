package shared

type Mode struct {
	value      float32
	occurences int

	Value      float32
	Occurences int
}

func NewMode(value float32, occurences int) Mode {
	return Mode{value: value, occurences: occurences, Value: value, Occurences: occurences}
}
