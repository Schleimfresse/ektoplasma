package main

import (
	"math"
	"reflect"
)

func NewNumber(value interface{}) *Value {
	return &Value{Number: &Number{ValueField: value}}
}

func (n *Number) IllegalOperation(other *Number) *RuntimeError {
	if other == nil {
		other = n
	}
	return NewRTError(n.PosStart(), n.PosEnd(), "Illegal operation", n.Context)
}

func (n *Number) Copy() *Value {
	return NewNumber(n.ValueField).SetContext(n.Context).SetPos(n.PosStart(), n.PosEnd())
}

func (n *Number) Type() reflect.Type {
	return reflect.TypeOf(n)
}

func (n *Number) Value() interface{} {
	return n.ValueField
}

func (n *Number) PosStart() *Position {
	return n.PositionStart
}

func (n *Number) PosEnd() *Position {
	return n.PositionEnd
}

// AddedTo performs addition with another number.
func (n *Number) AddedTo(other *Number) (*Value, *RuntimeError) {
	if n.ValueField != nil && other.ValueField != nil {
		switch nVal := n.ValueField.(type) {
		case int:
			if otherIsInt, ok := other.ValueField.(int); ok {
				value := NewNumber(float64(nVal) + float64(otherIsInt))
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			} else if otherIsFloat, ok := other.ValueField.(float64); ok {
				value := NewNumber(float64(nVal) + otherIsFloat)
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			}
		case float64:
			if otherIsInt, ok := other.ValueField.(int); ok {
				value := NewNumber(nVal + float64(otherIsInt))
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			} else if otherIsFloat, ok := other.ValueField.(float64); ok {
				value := NewNumber(nVal + otherIsFloat)
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			}
		}
	}
	return nil, n.IllegalOperation(other)
}

// SubtractedBy performs subtraction with another number.
func (n *Number) SubtractedBy(other *Number) (*Value, *RuntimeError) {
	if n.ValueField != nil && other.ValueField != nil {
		switch nVal := n.ValueField.(type) {
		case int:
			value := NewNumber(nVal - other.ValueField.(int))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		case float64:
			value := NewNumber(nVal - other.ValueField.(float64))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	}
	return nil, n.IllegalOperation(other)
}

// MultipliedBy performs multiplication with another number.
func (n *Number) MultipliedBy(other *Number) (*Value, *RuntimeError) {
	if n.ValueField != nil && other.ValueField != nil {
		switch nVal := n.ValueField.(type) {
		case int:
			value := NewNumber(nVal * other.ValueField.(int))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		case float64:
			value := NewNumber(nVal * other.ValueField.(float64))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	}
	return nil, n.IllegalOperation(other)
}

// DividedBy performs division with another number.
func (n *Number) DividedBy(other *Number) (*Value, *RuntimeError) {
	if n.ValueField != nil && other.ValueField != nil {
		if other.ValueField == 0 {
			return nil, NewRTError(other.PosStart(), other.PosEnd(), "Division by zero", other.Context)
		}
		switch nVal := n.ValueField.(type) {
		case int:
			value := NewNumber(nVal / other.ValueField.(int))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		case float64:
			value := NewNumber(nVal / other.ValueField.(float64))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) PowedBy(other *Number) (*Value, *RuntimeError) {
	if n.ValueField != nil && other.ValueField != nil {
		var nVal, otherVal float64
		switch val := n.ValueField.(type) {
		case float64:
			nVal = val
		case int:
			nVal = float64(val)
		default:
			// Handle unsupported types here
			return nil, n.IllegalOperation(other)
		}

		switch val := other.ValueField.(type) {
		case float64:
			otherVal = val
		case int:
			otherVal = float64(val)
		default:
			// Handle unsupported types here
			return nil, NewRTError(other.PosStart(), other.PosEnd(), "Invalid operation", other.Context)
		}

		value := NewNumber(math.Pow(nVal, otherVal))
		NewNumber(math.Pow(nVal, otherVal)).SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
		return value, nil
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) GetComparisonEq(other *Number) (*Value, *RuntimeError) {
	if other != nil {
		nVal := toInt(n.ValueField)
		otherVal := toInt(other.ValueField)
		value := NewBoolean(ConvertBoolToInt(nVal == otherVal))
		value.SetContext(n.Context)
		return value, nil
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) GetComparisonNe(other *Number) (*Value, *RuntimeError) {
	if other != nil {
		value := NewBoolean(ConvertBoolToInt(n.ValueField != other.ValueField))
		value.SetContext(n.Context)
		return value, nil
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) GetComparisonLt(other *Number) (*Value, *RuntimeError) {
	if other != nil {
		switch nVal := n.ValueField.(type) {
		case int:
			switch otherVal := other.ValueField.(type) {
			case int:
				value := NewBoolean(ConvertBoolToInt(n.ValueField.(int) < otherVal))
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			case float64:
				value := NewBoolean(ConvertBoolToInt(n.ValueField.(float64) < otherVal))
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			}
		case float64:
			switch otherVal := other.ValueField.(type) {
			case int:
				value := NewBoolean(ConvertBoolToInt(int(nVal) < otherVal))
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			case float64:
				value := NewBoolean(ConvertBoolToInt(nVal < otherVal))
				value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
				return value, nil
			}
		}
	}

	return nil, n.IllegalOperation(other)
}

func (n *Number) GetComparisonGt(other *Number) (*Value, *RuntimeError) {
	switch nVal := n.ValueField.(type) {
	case int:
		if otherIsInt, ok := other.ValueField.(int); ok {
			value := NewBoolean(ConvertBoolToInt(nVal > otherIsInt))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		} else if otherIsFloat, ok := other.ValueField.(float64); ok {
			value := NewBoolean(ConvertBoolToInt(float64(nVal) > otherIsFloat))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	case float64:
		if otherIsInt, ok := other.ValueField.(int); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(int) > otherIsInt))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		} else if otherIsFloat, ok := other.ValueField.(float64); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(float64) > otherIsFloat))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) GetComparisonLte(other *Number) (*Value, *RuntimeError) {
	switch nVal := n.ValueField.(type) {
	case int:
		if otherIsInt, ok := other.ValueField.(int); ok {
			value := NewBoolean(ConvertBoolToInt(nVal <= otherIsInt))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		} else if otherIsFloat, ok := other.ValueField.(float64); ok {
			value := NewBoolean(ConvertBoolToInt(float64(nVal) <= otherIsFloat))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	case float64:
		if otherIsInt, ok := other.ValueField.(int); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(int) <= otherIsInt))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		} else if otherIsFloat, ok := other.ValueField.(float64); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(float64) <= otherIsFloat))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) GetComparisonGte(other *Number) (*Value, *RuntimeError) {
	switch nVal := n.ValueField.(type) {
	case int:
		if otherIsInt, ok := other.ValueField.(int); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(int) >= otherIsInt))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		} else if otherIsFloat, ok := other.ValueField.(float64); ok {
			value := NewBoolean(ConvertBoolToInt(float64(nVal) >= otherIsFloat))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	case float64:
		if otherIsInt, ok := other.ValueField.(int); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(int) >= otherIsInt))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		} else if otherIsFloat, ok := other.ValueField.(float64); ok {
			value := NewBoolean(ConvertBoolToInt(n.ValueField.(float64) >= otherIsFloat))
			value.SetContext(n.Context).SetPos(n.PositionStart, n.PositionEnd)
			return value, nil
		}
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) AndedBy(other *Number) (*Value, *RuntimeError) {
	if other != nil {
		value := NewBoolean(ConvertBoolToInt(n.ValueField != 0 && other.ValueField != 0))
		value.SetContext(n.Context)
		return value, nil
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) OredBy(other *Number) (*Value, *RuntimeError) {
	if other != nil {
		value := NewBoolean(ConvertBoolToInt(n.ValueField != 0 || other.ValueField != 0))
		value.SetContext(n.Context)
		return value, nil
	}
	return nil, n.IllegalOperation(other)
}

func (n *Number) Notted() (*Value, *RuntimeError) {
	value := NewBoolean(ConvertBoolToInt(n.ValueField != 0))
	value.SetContext(n.Context)
	return value, nil
}
