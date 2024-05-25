package main

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strconv"
)

func ConvertBoolToInt(expr bool) Binary {
	if expr {
		return One
	} else {
		return Zero
	}
}

// SetContext sets the context
func (v *Value) SetContext(context *Context) *Value {
	if v.Number != nil {
		v.Number.Context = context
	} else if v.String != nil {
		v.String.Context = context
	} else if v.Function != nil {
		v.Function.Base.Context = context
	} else if v.BuildInFunction != nil {
		v.BuildInFunction.Base.Context = context
	} else if v.Array != nil {
		v.Array.Context = context
	} else if v.Null != nil {
		v.Null.Context = context
	} else if v.Boolean != nil {
		v.Boolean.Context = context
	}
	return v
}

// SetPos sets the position
func (v *Value) SetPos(posStart *Position, posEnd *Position) *Value {
	if v.Number != nil {
		v.Number.PositionStart = posStart
		v.Number.PositionEnd = posEnd
	} else if v.String != nil {
		v.String.PositionStart = posStart
		v.String.PositionEnd = posEnd
	} else if v.Function != nil {
		v.Function.Base.PositionStart = posStart
		v.Function.Base.PositionEnd = posEnd
	} else if v.Array != nil {
		v.Array.PositionStart = posStart
		v.Array.PositionEnd = posEnd
	} else if v.BuildInFunction != nil {
		v.BuildInFunction.Base.PositionStart = posStart
		v.BuildInFunction.Base.PositionEnd = posEnd
	} else if v.Null != nil {
		v.Null.PositionStart = posStart
		v.Null.PositionEnd = posEnd
	} else if v.Boolean != nil {
		v.Boolean.PositionStart = posStart
		v.Boolean.PositionEnd = posEnd
	}
	return v
}

// Value retrieves the Value of the type if available
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
	} else if v.Null != nil {
		return v.Null.String()
	} else if v.Boolean != nil {
		return v.Boolean.String()
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

// GetPosStart searches for the non-nil field in the Value struct and returns its StartPos if available.
func (v *Value) GetPosStart() *Position {
	if v.Number != nil {
		return v.Number.PosStart()
	} else if v.Function != nil {
		return v.Function.PosStart()
	} else if v.BuildInFunction != nil {
		return v.BuildInFunction.Base.PosStart()
	} else if v.String != nil {
		return v.String.PosStart()
	} else if v.Array != nil {
		return v.Array.PosStart()
	} else if v.Boolean != nil {
		return v.Boolean.PosStart()
	}
	return nil
}

// GetPosEnd searches for the non-nil field in the Value struct and returns its EndPos if available.
func (v *Value) GetPosEnd() *Position {
	if v.Number != nil {
		return v.Number.PosEnd()
	} else if v.Function != nil {
		return v.Function.PosEnd()
	} else if v.BuildInFunction != nil {
		return v.BuildInFunction.Base.PosEnd()
	} else if v.String != nil {
		return v.String.PosEnd()
	} else if v.Array != nil {
		return v.Array.PosEnd()
	} else if v.Boolean != nil {
		return v.Boolean.PosEnd()
	}
	return nil
}

func (v *Value) Type() string {
	if v.Number != nil {
		return "Number"
	} else if v.Function != nil {
		return "Function"
	} else if v.BuildInFunction != nil {
		return "BuildInFunction"
	} else if v.String != nil {
		return "String"
	} else if v.Array != nil {
		return "Array"
	} else if v.Boolean != nil {
		return "Boolean"
	}
	return ""
}

func interfaceToBytes(data interface{}) []byte {
	switch v := data.(type) {
	case string:
		return []byte(v)
	case int:
		dataString := strconv.Itoa(v)
		return []byte(dataString)
	}
	panic("unexpected type" + reflect.TypeOf(data).String())
}

func isString(data []byte) bool {
	for _, b := range data {
		if !isLetter(b) {
			return false
		}
	}
	return true
}

func isNumber(data []byte) bool {
	if isInt(data) {
		return true
	}
	return isFloat(data)
}

func isInt(data []byte) bool {
	var value int64
	err := binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
	if err == nil {
		return true
	} else {
		return false
	}
}

func isFloat(data []byte) bool {
	str := string(data)
	_, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return true
	} else {
		return false
	}
}
