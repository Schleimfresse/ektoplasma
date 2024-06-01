package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"
)

func (b *BuildInFunction) executePrint(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if exists {
		log.Print(interfaceToBytes(value.Value()))
		_, err := syscall.Write(syscall.Stdout, interfaceToBytes(value.Value()))
		if err != nil {
			return nil
		}
	}

	return NewRTResult().Success(NewEmptyValue())
}

func (b *BuildInFunction) executePrintLn(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if exists {
		_, err := syscall.Write(syscall.Stdout, interfaceToBytes(value.Value()))
		if err != nil {
			return nil
		}

		_, err = syscall.Write(syscall.Stdout, []byte{10})
		if err != nil {
			return nil
		}

	}

	return NewRTResult().Success(NewEmptyValue())
}

func (b *BuildInFunction) executeInput(execCtx *Context) *RTResult {
	res := NewRTResult()
	buf := make([]byte, 1024)
	n, err := syscall.Read(syscall.Stdin, buf)
	if err != nil {
		return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("Error reading input: %s", err), execCtx))
	}
	inputStr := strings.TrimSpace(string(buf[:n]))

	// Try parsing as float64
	var value float64
	var errFloat error
	if value, errFloat = strconv.ParseFloat(inputStr, 64); errFloat == nil {
		return res.Success(NewNumber(value))
	}

	// Try parsing as int
	var intValue int64
	var errInt error
	if intValue, errInt = strconv.ParseInt(inputStr, 10, 64); errInt == nil {
		return res.Success(NewNumber(float64(intValue)))
	}

	// If parsing an int or float is not successful, return a string
	return res.Success(NewString(inputStr))
}

func (b *BuildInFunction) executeIsString(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if exists && value.String != nil {
		return NewRTResult().Success(NewBoolean(One))
	} else {
		return NewRTResult().Success(NewBoolean(Zero))
	}
}
