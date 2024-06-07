package main

import (
	"fmt"
	"strings"
)

func NewArray(elements []*Value) *Value {
	return &Value{
		Array: &Array{
			Elements: elements,
		},
	}
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
	newArray.Array.Elements = append(newArray.Array.Elements, other)
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
func (a *Array) GetIndex(index *Number) (*Value, *RuntimeError) {
	if index.ValueField.(int) < 0 || index.ValueField.(int) >= len(a.Elements) {
		return nil, NewRTError(index.PosStart(), index.PosEnd(), fmt.Sprintf("Element at index %v could not be retrieved from array, index is out of bounds", index.ValueField.(int)), a.Context)
	}

	return a.Elements[index.ValueField.(int)], nil
}

// Error for illegal operation
func (a *Array) IllegalOperation(other *Array) *RuntimeError {
	if other == nil {
		other = a
	}
	return NewRTError(a.PosStart(), a.PosEnd(), "Illegal operation", a.Context)
}

// Length returns the length of the byte array.
func (a *Array) Length() *Value {
	value := NewNumber(float64(len(a.Elements)))
	value.SetContext(a.Context)
	return value
}

// String representation of Array
func (a *Array) String() string {
	elementStrings := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		switch {
		case element.Number != nil:
			elementStrings[i] = fmt.Sprintf("%v", element.Number.ValueField)
		case element.String != nil:
			elementStrings[i] = fmt.Sprintf("%q", element.String.ValueField)
		case element.Array != nil:
			elementStrings[i] = element.Array.String() // Recursively call String for nested arrays
		case element.Function != nil:
			elementStrings[i] = element.Function.String()
		case element.BuildInFunction != nil:
			elementStrings[i] = element.BuildInFunction.Base.Name
		case element.Boolean != nil:
			elementStrings[i] = element.Boolean.String()
		default:
			elementStrings[i] = "<null>"
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(elementStrings, ", "))
}

func NewVariadicArray(elements []*Value) *Value {
	var array []*Value
	for e := range elements {
		array = append(array, elements[e])
	}
	return &Value{VariadicArray: &VariadicArray{array, array[0].GetPosStart(), array[len(array)-1].GetPosEnd(), array[0].GetContext()}}
}
