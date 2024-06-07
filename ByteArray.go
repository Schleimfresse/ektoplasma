package main

import "strconv"

// NewByteArray is the constructor for ByteArray
func NewByteArray(value []byte) *Value {
	return &Value{ByteArray: &ByteArray{ValueField: value}}
}

func (b *ByteArray) Copy() *Value {
	return NewByteArray(append([]byte{}, b.ValueField...)).SetContext(b.Context).SetPos(b.PositionStart, b.PositionEnd)
}

func (b *ByteArray) Value() []byte {
	return b.ValueField
}

func (b *ByteArray) PosStart() *Position {
	return b.PositionStart
}

func (b *ByteArray) PosEnd() *Position {
	return b.PositionEnd
}

func (b *ByteArray) IllegalOperation(other interface{}) *RuntimeError {
	if other == nil {
		other = b
	}
	return NewRTError(b.PosStart(), b.PosEnd(), "Illegal operation", b.Context)
}

func (b *ByteArray) AddedTo(other *ByteArray) (*Value, *RuntimeError) {
	if other == nil {
		return nil, b.IllegalOperation(other)
	}
	value := NewByteArray(append(b.ValueField, other.ValueField...))
	value.SetContext(b.Context)
	return value, nil
}

func (b *ByteArray) MultipliedBy(other *Number) (*Value, *RuntimeError) {
	if ValueField, ok := other.ValueField.(int); ok {
		if ValueField <= 0 {
			value := NewByteArray([]byte{})
			value.SetContext(b.Context)
			return value, nil
		}

		result := make([]byte, 0, len(b.ValueField)*ValueField)
		for i := 0; i < ValueField; i++ {
			result = append(result, b.ValueField...)
		}
		value := NewByteArray(result)
		value.SetContext(b.Context)
		return value, nil
	}
	return nil, b.IllegalOperation(other)
}

func (b *ByteArray) GetComparisonEq(other *Value) (*Value, *RuntimeError) {
	if other != nil && other.ByteArray != nil {
		bVal := b.ValueField
		otherVal := other.ByteArray.ValueField
		value := NewBoolean(ConvertBoolToInt(string(bVal) == string(otherVal)))
		value.SetContext(b.Context)
		return value, nil
	}
	return nil, b.IllegalOperation(other)
}

// GetByte retrieves the byte at a specified index.
func (b *ByteArray) GetByte(index *Number) (*Value, *RuntimeError) {
	if idx, ok := index.ValueField.(int); ok {
		if idx < 0 || idx >= len(b.ValueField) {
			return nil, NewRTError(b.PosStart(), b.PosEnd(), "Index out of bounds", b.Context)
		}
		value := NewNumber(float64(b.ValueField[idx]))
		value.SetContext(b.Context)
		return value, nil
	}
	return nil, b.IllegalOperation(index)
}

// Slice retrieves a sub-array of the byte array.
func (b *ByteArray) Slice(startIndex, endIndex *Number) (*Value, *RuntimeError) {
	if start, ok := startIndex.ValueField.(int); ok {
		if end, ok := endIndex.ValueField.(int); ok {
			if start < 0 || end > len(b.ValueField) || start >= end {
				return nil, NewRTError(b.PosStart(), b.PosEnd(), "Invalid slice indices", b.Context)
			}
			value := NewByteArray(b.ValueField[start:end])
			value.SetContext(b.Context)
			return value, nil
		}
	}
	return nil, b.IllegalOperation(startIndex)
}

// Length returns the length of the byte array.
func (b *ByteArray) Length() *Value {
	value := NewNumber(float64(len(b.ValueField)))
	value.SetContext(b.Context)
	return value
}

// ToString converts the byte array to a string.
func (b *ByteArray) ToString() *Value {
	value := NewString(string(b.ValueField))
	value.SetContext(b.Context)
	return value
}

// ToNumber converts the byte array to a number if possible.
func (b *ByteArray) ToNumber() (*Value, *RuntimeError) {
	strValue := string(b.ValueField)
	numValue, err := strconv.ParseFloat(strValue, 64)
	if err != nil {
		return nil, NewRTError(b.PosStart(), b.PosEnd(), "Invalid byte array for conversion to number", b.Context)
	}
	value := NewNumber(numValue)
	value.SetContext(b.Context)
	return value, nil
}
