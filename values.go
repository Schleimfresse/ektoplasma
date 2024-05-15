package main

import (
	"log"
	"math"
	"reflect"
)

// SetPos sets the position of the number.
func (n *Number) SetPos(posStart *Position, posEnd *Position) *Number {
	n.PosStart = posStart
	n.PosEnd = posEnd
	return n
}

// SetContext sets the context of the number.
func (n *Number) SetContext(context *Context) *Number {
	n.Context = context
	return n
}

// AddedTo performs addition with another number.
func (n *Number) AddedTo(other *Number) (*Number, *RuntimeError) {
	if n.Value != nil && other.Value != nil {
		log.Println(n.Value, other.Value, reflect.TypeOf(n.Value), reflect.TypeOf(other.Value))
		switch nVal := n.Value.(type) {
		case int:
			if otherIsInt, ok := other.Value.(int); ok {
				return NewNumber(float64(nVal)+float64(otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			} else if otherIsFloat, ok := other.Value.(float64); ok {
				return NewNumber(float64(nVal)+otherIsFloat).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			}
		case float64:
			if otherIsInt, ok := other.Value.(int); ok {
				return NewNumber(nVal+float64(otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			} else if otherIsFloat, ok := other.Value.(float64); ok {
				return NewNumber(nVal+otherIsFloat).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			}
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

// SubbedBy performs subtraction with another number.
func (n *Number) SubtractedBy(other *Number) (*Number, *RuntimeError) {
	if n.Value != nil && other.Value != nil {
		switch nVal := n.Value.(type) {
		case int:
			return NewNumber(nVal-other.Value.(int)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		case float64:
			return NewNumber(nVal-other.Value.(float64)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

// MultedBy performs multiplication with another number.
func (n *Number) MultipliedBy(other *Number) (*Number, *RuntimeError) {
	if n.Value != nil && other.Value != nil {
		switch nVal := n.Value.(type) {
		case int:
			return NewNumber(nVal*other.Value.(int)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		case float64:
			return NewNumber(nVal*other.Value.(float64)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

// DivedBy performs division with another number.
func (n *Number) DividedBy(other *Number) (*Number, *RuntimeError) {
	if n.Value != nil && other.Value != nil {
		if other.Value == 0 {
			return nil, NewRTError(other.PosStart, other.PosEnd, "Division by zero", other.Context)
		}
		switch nVal := n.Value.(type) {
		case int:
			return NewNumber(nVal/other.Value.(int)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		case float64:
			return NewNumber(nVal/other.Value.(float64)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) PowedBy(other *Number) (*Number, *RuntimeError) {
	if n.Value != nil && other.Value != nil {
		var nVal, otherVal float64
		switch val := n.Value.(type) {
		case float64:
			nVal = val
		case int:
			nVal = float64(val)
		default:
			// Handle unsupported types here
			return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
		}

		switch val := other.Value.(type) {
		case float64:
			otherVal = val
		case int:
			otherVal = float64(val)
		default:
			// Handle unsupported types here
			return nil, NewRTError(other.PosStart, other.PosEnd, "Invalid operation", other.Context)
		}

		return NewNumber(math.Pow(nVal, otherVal)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) GetComparisonEq(other *Number) (*Number, *RuntimeError) {
	if other != nil {
		return NewNumber(ConvertBoolToInt(n.Value == other.Value)).SetContext(n.Context), nil
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) GetComparisonNe(other *Number) (*Number, *RuntimeError) {
	if other != nil {
		return NewNumber(ConvertBoolToInt(n.Value != other.Value)).SetContext(n.Context), nil
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) GetComparisonLt(other *Number) (*Number, *RuntimeError) {
	if other != nil {
		switch nVal := n.Value.(type) {
		case int:
			if otherIsInt, ok := other.Value.(int); ok {
				return NewNumber(ConvertBoolToInt(n.Value.(int) < otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			} else if otherIsFloat, ok := other.Value.(float64); ok {
				return NewNumber(ConvertBoolToInt(float64(nVal) < otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			}
		case float64:
			if otherIsInt, ok := other.Value.(int); ok {
				return NewNumber(ConvertBoolToInt(n.Value.(int) < otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			} else if otherIsFloat, ok := other.Value.(float64); ok {
				return NewNumber(ConvertBoolToInt(nVal < otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
			}
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) GetComparisonGt(other *Number) (*Number, *RuntimeError) {
	switch nVal := n.Value.(type) {
	case int:
		if otherIsInt, ok := other.Value.(int); ok {
			return NewNumber(ConvertBoolToInt(n.Value.(int) > otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		} else if otherIsFloat, ok := other.Value.(float64); ok {
			return NewNumber(ConvertBoolToInt(float64(nVal) > otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	case float64:
		if otherIsInt, ok := other.Value.(int); ok {
			return NewNumber(ConvertBoolToInt(n.Value.(int) > otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		} else if otherIsFloat, ok := other.Value.(float64); ok {
			return NewNumber(ConvertBoolToInt(nVal > otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) GetComparisonLte(other *Number) (*Number, *RuntimeError) {
	switch nVal := n.Value.(type) {
	case int:
		if otherIsInt, ok := other.Value.(int); ok {
			return NewNumber(ConvertBoolToInt(n.Value.(int) <= otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		} else if otherIsFloat, ok := other.Value.(float64); ok {
			return NewNumber(ConvertBoolToInt(float64(nVal) <= otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	case float64:
		if otherIsInt, ok := other.Value.(int); ok {
			return NewNumber(ConvertBoolToInt(n.Value.(int) <= otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		} else if otherIsFloat, ok := other.Value.(float64); ok {
			return NewNumber(ConvertBoolToInt(nVal <= otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) GetComparisonGte(other *Number) (*Number, *RuntimeError) {
	switch nVal := n.Value.(type) {
	case int:
		if otherIsInt, ok := other.Value.(int); ok {
			return NewNumber(ConvertBoolToInt(n.Value.(int) >= otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		} else if otherIsFloat, ok := other.Value.(float64); ok {
			return NewNumber(ConvertBoolToInt(float64(nVal) >= otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	case float64:
		if otherIsInt, ok := other.Value.(int); ok {
			return NewNumber(ConvertBoolToInt(n.Value.(int) >= otherIsInt)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		} else if otherIsFloat, ok := other.Value.(float64); ok {
			return NewNumber(ConvertBoolToInt(nVal >= otherIsFloat)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) AndedBy(other *Number) (*Number, *RuntimeError) {
	if other != nil {
		return NewNumber(ConvertBoolToInt(n.Value != 0 && other.Value != 0)).SetContext(n.Context), nil
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) OredBy(other *Number) (*Number, *RuntimeError) {
	if other != nil {
		return NewNumber(ConvertBoolToInt(n.Value != 0 || other.Value != 0)).SetContext(n.Context), nil
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func (n *Number) Notted() (*Number, *RuntimeError) {
	return NewNumber(ConvertBoolToInt(n.Value == 0)).SetContext(n.Context), nil
}

func (n *Number) IsTrue() bool {
	return n.Value != 0
}

func NewNumber(value interface{}) *Number {
	return &Number{Value: value}
}

func ConvertBoolToInt(expr bool) interface{} {
	if expr {
		return 1
	} else {
		return 0
	}
}
