package main

import (
	"fmt"
	"strings"
)

func NewArray(value []Value) *Value {
	return &Value{Array: &Array{Elements: value}}
}

func (a *Array) Copy() *Value {
	return NewArray(a.Elements).SetContext(a.Context).SetPos(a.PosStart(), a.PosEnd())
}

func (a *Array) PosStart() *Position {
	return a.PositionStart
}

func (a *Array) PosEnd() *Position {
	return a.PositionEnd
}

// Add element to Array
func (a *Array) AddedTo(other *Value) (*Value, *RuntimeError) {
	newArray := a.Copy()
	newArray.Array.Elements = append(newArray.Array.Elements, *other)
	return newArray, nil
}

// Remove element by index from Array
func (a *Array) SubtractedBy(other *Value) (*Value, *RuntimeError) {
	newArray := a.Copy()
	if other.Number.ValueField.(int) < 0 || other.Number.ValueField.(int) >= len(newArray.Array.Elements) {
		return nil, NewRTError(other.Number.PosStart(), other.Number.PosEnd(), fmt.Sprintf("Element at index %v could not be removed from array, index is out of bounds", other.Value()), a.Context)
	}
	newArray.Array.Elements = append(newArray.Array.Elements[:other.Number.ValueField.(int)], newArray.Array.Elements[other.Number.ValueField.(int)+1:]...)
	return newArray, nil
}

// Extend Array with another Array
func (a *Array) MultipliedBy(other *Value) (*Value, *RuntimeError) {
	newArray := a.Copy()
	newArray.Array.Elements = append(newArray.Array.Elements, other.Array.Elements...)
	return newArray, nil
}

// Retrieve element by index from Array
func (a *Array) DividedBy(other *Value) (*Value, *RuntimeError) {

	if other.Number.ValueField.(int) < 0 || other.Number.ValueField.(int) >= len(a.Elements) {
		return nil, NewRTError(other.Number.PosStart(), other.Number.PosEnd(), fmt.Sprintf("Element at index %v could not be retrieved from array, index is out of bounds", other.Value()), a.Context)
	}
	return &a.Elements[other.Number.ValueField.(int)], nil
}

// Error for illegal operation
func (a *Array) IllegalOperation(other *Array) *RuntimeError {
	if other == nil {
		other = a
	}
	return NewRTError(a.PosStart(), a.PosEnd(), "Illegal operation", a.Context)
}

// String representation of Array
func (a *Array) String() string {
	elementStrings := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		elementStrings[i] = fmt.Sprintf("%v", element.Value())
	}
	return fmt.Sprintf("[%s]", strings.Join(elementStrings, ", "))
}
