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
		v.Function.Context = context
		return v
	}
	return nil
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
		v.Function.PositionStart = posStart
		v.Function.PositionEnd = posEnd
		return v
	}
	return nil
}
