package main

func NewNull() *Value {
	return &Value{Null: &Null{}}
}

func (n *Null) Copy() *Value {
	return NewNull().SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
}

func (n *Null) PosStart() *Position {
	return n.PositionStart
}

func (n *Null) PosEnd() *Position {
	return n.PositionEnd
}

func (n *Null) String() string {
	return "<null>"
}

func (n *Null) GetComparisonEq(other *Value) (*Value, *RuntimeError) {
	if other.Null == n {
		return NewBoolean(One), nil
	} else {
		return NewBoolean(Zero), nil
	}
}

func (n *Null) GetComparisonNe(other *Value) (*Value, *RuntimeError) {
	if other.Null != n {
		return NewBoolean(One), nil
	} else {
		return NewBoolean(Zero), nil
	}
}

func (n *Null) IllegalOperation(other interface{}) *RuntimeError {
	if other == nil {
		other = n
	}
	return NewRTError(n.PosStart(), n.PosEnd(), "Illegal operation", n.Context)
}
