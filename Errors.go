package main

import (
	"fmt"
	"strings"
)

const (
	bold  = "\033[1m"
	reset = "\033[0m"
)

func NewError(posStart *Position, posEnd *Position, errorName string, details string) *Error {
	return &Error{posStart, posEnd, errorName, details}
}

// AsString converts the error to a string format.
func (e *Error) AsString() string {
	result := fmt.Sprintf("%s%s: %s%s\n", bold, e.ErrorName, e.Details, reset)
	result += fmt.Sprintf("File %s, line %d\n\n", e.PosStart.Fn, e.PosStart.Ln+1)
	result += stringWithArrows(e.PosStart.Ftxt, *e.PosStart, *e.PosEnd)
	return result
}

// NewIllegalCharError creates a new IllegalCharError instance.
func NewIllegalCharError(posStart *Position, posEnd *Position, details string) *IllegalCharError {
	return &IllegalCharError{Error{posStart, posEnd, "Illegal Character", details}}
}

// NewInvalidSyntaxError creates a new InvalidSyntaxError instance.
func NewInvalidSyntaxError(posStart *Position, posEnd *Position, details string) *InvalidSyntaxError {
	return &InvalidSyntaxError{Error{posStart, posEnd, "Invalid Syntax", details}}
}

func NewRTError(posStart *Position, posEnd *Position, details string, context *Context) *RuntimeError {
	return &RuntimeError{
		Error:   NewError(posStart, posEnd, "Runtime Error", details),
		Context: context,
	}
}

// stringWithArrows returns a string with arrows pointing to the error position.
func stringWithArrows(text string, posStart Position, posEnd Position) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	// Calculate indices
	idxStart := max(0, strings.LastIndex(lines[posStart.Ln], "\n"))
	idxEnd := min(len(text), strings.Index(lines[posEnd.Ln], "\n")+1)

	// Output line
	result.WriteString(lines[posStart.Ln])
	result.WriteString("\n")

	// Output arrows
	arrows := strings.Repeat(" ", posStart.Col) + strings.Repeat("^", max(0, idxEnd-idxStart+1))
	result.WriteString(arrows)
	return result.String()
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
