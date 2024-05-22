package main

func ConvertBoolToInt(expr bool) interface{} {
	if expr {
		return 1
	} else {
		return 0
	}
}

func (v *Value) SetContext(context *Context) *Value {
	if v.Number != nil {
		v.Number.Context = context
		return v
	} else if v.String != nil {
		v.String.Context = context
		return v
	} else if v.Function != nil {
		v.Function.Base.Context = context
		return v
	} else if v.Array != nil {
		v.Array.Context = context
	}
	return v
}

func (v *Value) SetPos(posStart *Position, posEnd *Position) *Value {
	if v.Number != nil {
		v.Number.PositionStart = posStart
		v.Number.PositionEnd = posEnd
		return v
	} else if v.String != nil {
		v.String.PositionStart = posStart
		v.String.PositionEnd = posEnd
		return v
	} else if v.Function != nil {
		v.Function.Base.PositionStart = posStart
		v.Function.Base.PositionEnd = posEnd
		return v
	} else if v.Array != nil {
		v.Array.PositionStart = posStart
		v.Array.PositionEnd = posEnd
	} else if v.BuildInFunction != nil {
		v.BuildInFunction.Base.PositionStart = posStart
		v.BuildInFunction.Base.PositionEnd = posEnd
	}
	return v
}

func (v *Value) Value() interface{} {
	if v.Number != nil {
		return v.Number.Value()
	} else if v.String != nil {
		return v.String.Value()
	} else if v.Function != nil {
		return v.Function.String()
	} else if v.Array != nil {
		return v.Array.Elements
	} else if v.BuildInFunction != nil {
		return v.BuildInFunction.String()
	}
	return v
}

func (v *Value) Copy() *Value {
	if v.Number != nil {
		return v.Number.Copy()
	} else if v.String != nil {
		return v.String.Copy()
	} else if v.Function != nil {
		return v.Function.Copy()
	} else if v.Array != nil {
		return v.Array.Copy()
	} else if v.BuildInFunction != nil {
		return v.BuildInFunction.Copy()
	}
	return v
}
