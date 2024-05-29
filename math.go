package main

import (
	"fmt"
	"math"
)

func (b *BuildInFunction) Cos(execCtx *Context) *RTResult {
	res := NewRTResult()
	value, exists := execCtx.SymbolTable.Get("x")
	if exists && value.Number != nil {
		switch v := value.Number.ValueField.(type) {
		case float64:
			return res.Success(NewNumber(math.Cos(v)))
		case int:
			return res.Success(NewNumber(math.Cos(float64(v))))
		}
	} else {
		return NewRTResult().Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("argument must be an Number, got: %v", value.Type()), execCtx))
	}
	return res.Success(NewEmptyValue())
}

func (b *BuildInFunction) Sin(execCtx *Context) *RTResult {
	res := NewRTResult()
	value, exists := execCtx.SymbolTable.Get("x")
	if exists && value.Number != nil {
		switch v := value.Number.ValueField.(type) {
		case float64:
			return res.Success(NewNumber(math.Sin(v)))
		case int:
			return res.Success(NewNumber(math.Sin(float64(v))))
		}
	} else {
		return NewRTResult().Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("argument must be an Number, got: %v", value.Type()), execCtx))
	}
	return res.Success(NewEmptyValue())
}
