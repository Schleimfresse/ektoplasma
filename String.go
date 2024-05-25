package main

func NewString(value string) *Value {
	return &Value{String: &String{ValueField: value}}
}

func (s *String) IllegalOperation(other interface{}) *RuntimeError {
	if other == nil {
		other = s
	}
	return NewRTError(s.PosStart(), s.PosEnd(), "Illegal operation", s.Context)
}

func (s *String) Copy() *Value {
	return NewString(s.ValueField).SetContext(s.Context).SetPos(s.PositionStart, s.PositionEnd)
}

func (s *String) IsTrue() bool {
	return len(s.ValueField) > 0
}

func (s *String) Value() string {
	return s.ValueField
}

func (s *String) PosStart() *Position {
	return s.PositionStart
}

func (s *String) PosEnd() *Position {
	return s.PositionEnd
}

// MultipliedBy performs multiplication with another number.
func (s *String) MultipliedBy(other *Number) (*Value, *RuntimeError) {
	if ValueField, ok := other.ValueField.(int); ok {
		if ValueField <= 0 {
			value := NewString("")
			value.SetContext(s.Context)
			return value, nil
		}

		var result string
		for i := 0; i < ValueField; i++ {
			result += s.ValueField
		}
		value := NewString(result)
		value.SetContext(s.Context)
		return value, nil
	}
	return nil, s.IllegalOperation(other)
}

// AddedTo performs addition with another number.
func (s *String) AddedTo(other *String) (*Value, *RuntimeError) {
	value := NewString(s.ValueField + other.ValueField)
	value.SetContext(s.Context)
	return value, nil
}
