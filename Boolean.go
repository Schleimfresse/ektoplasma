package main

func NewBoolean(value Binary) *Value {
	return &Value{Boolean: &Boolean{Binary: value}}
}

func (b *Boolean) Copy() *Value {
	return NewBoolean(b.Binary).SetContext(b.Context).SetPos(b.PositionStart, b.PositionEnd)
}

func (b *Boolean) PosStart() *Position {
	return b.PositionStart
}

func (b *Boolean) PosEnd() *Position {
	return b.PositionEnd
}

func (b *Boolean) IsTrue() bool {
	return false
}

func (b *Boolean) String() string {
	if b.Binary == One {
		return "true"
	} else {
		return "false"
	}
}

func (b *Boolean) GetComparisonEq(other *Boolean) (*Value, *RuntimeError) {
	if other != nil {
		var value *Value
		if b.Binary == other.Binary {
			value = NewBoolean(One)
		} else {
			value = NewBoolean(Zero)
		}
		value.SetContext(b.Context)
		return value, nil
	}
	return nil, b.IllegalOperation(other)
}

func (b *Boolean) GetComparisonNe(other *Boolean) (*Value, *RuntimeError) {
	if other != nil {
		var value *Value
		if b.Binary != other.Binary {
			value = NewBoolean(One)
		} else {
			value = NewBoolean(Zero)
		}
		value.SetContext(b.Context)
		return value, nil
	}
	return nil, b.IllegalOperation(other)
}

func (b *Boolean) IllegalOperation(other interface{}) *RuntimeError {
	if other == nil {
		other = b
	}
	return NewRTError(b.PosStart(), b.PosEnd(), "Illegal operation", b.Context)
}
