package main

import (
	"fmt"
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
		switch nVal := n.Value.(type) {
		case int:
			fmt.Println(nVal, other.Value)
			return NewNumber(nVal+other.Value.(int)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		case float64:
			return NewNumber(nVal+other.Value.(float64)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
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
	fmt.Println("FSD:", n.Value, other.Value, reflect.TypeOf(n.Value), reflect.TypeOf(other.Value))
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
		switch nVal := n.Value.(type) {
		case int:
			return NewNumber(nVal/other.Value.(int)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		case float64:
			return NewNumber(nVal/other.Value.(float64)).SetContext(n.Context).SetPos(n.PosStart, n.PosEnd), nil
		}
	}
	return nil, NewRTError(n.PosStart, n.PosEnd, "Invalid operation", n.Context)
}

func NewNumber(value interface{}) *Number {
	return &Number{Value: value}
}
