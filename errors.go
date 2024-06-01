package main

import (
	"fmt"
	"path/filepath"
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
	pos := e.PosStart
	ErrorFilePath, _ := filepath.Abs(pos.Fn)

	result := fmt.Sprintf("%s%s: %s%s\n", bold, e.ErrorName, e.Details, reset)
	result += fmt.Sprintf("File %s, line %d\n", e.PosStart.Fn, e.PosStart.Ln+1)
	result += fmt.Sprintf(" %s:%v\n\n", ErrorFilePath, pos.Ln+1)
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

func NewExpectedCharError(posStart *Position, posEnd *Position, details string) *ExpectedCharError {
	return &ExpectedCharError{Error{posStart, posEnd, "Expected Character", details}}
}

func NewRTError(posStart *Position, posEnd *Position, details string, context *Context) *RuntimeError {
	return &RuntimeError{
		Error:   NewError(posStart, posEnd, "Runtime Error", details),
		Context: context,
	}
}

// stringWithArrows returns a string with arrows pointing to the error position.
/*func stringWithArrows(text string, posStart Position, posEnd Position) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	// Calculate indices
	//idxStart := max(0, strings.LastIndex(lines[posStart.Ln], "\n"))
	//idxEnd := min(len(text), strings.Index(lines[posEnd.Ln], "\n")+1)

	// Output line
	result.WriteString(lines[posStart.Ln])
	result.WriteString("\n")

	// Output arrows
	arrows := strings.Repeat(" ", posStart.Col) + strings.Repeat("^", max(0, posEnd.Idx-posStart.Idx+1))
	result.WriteString(arrows)
	return result.String()
}*/

func stringWithArrows(text string, posStart Position, posEnd Position) string {
	var result strings.Builder

	// Calculate indices
	idxStart := max(strings.LastIndex(text[:posStart.Idx], "\n"), 0)
	idxEnd := strings.Index(text[posStart.Idx:], "\n")
	if idxEnd < 0 {
		idxEnd = len(text)
	} else {
		idxEnd += posStart.Idx
	}

	// Generate each line
	lineCount := posEnd.Ln - posStart.Ln + 1
	for i := 0; i < lineCount; i++ {
		// Calculate line columns
		line := text[idxStart:idxEnd]
		colStart := posStart.Col
		if i != 0 {
			colStart = 0
		}
		colEnd := posEnd.Col
		if i != lineCount-1 {
			colEnd = len(line) - 1
		}

		// Append to result
		result.WriteString(line + "\n")

		if colEnd-colStart == 0 {
			result.WriteString(strings.Repeat(" ", colStart) + strings.Repeat("^", 1))
		} else {
			result.WriteString(strings.Repeat(" ", colStart) + strings.Repeat("^", colEnd-colStart))

		}

		// Re-calculate indices
		idxStart = idxEnd
		idxEnd = strings.Index(text[idxStart:], "\n")
		if idxEnd < 0 {
			idxEnd = len(text)
		} else {
			idxEnd += idxStart + 1
		}
	}

	return strings.ReplaceAll(result.String(), "\t", "")
}

func (e RuntimeError) generateTraceback() string {
	var result string
	pos := e.PosStart
	ctx := e.Context

	ErrorFilePath, _ := filepath.Abs(pos.Fn)
	for ctx != nil {
		result = fmt.Sprintf("File %s, line %d, in %s\n%s", pos.Fn, pos.Ln+1, ctx.DisplayName, result)
		result += fmt.Sprintf(" %s:%v\n", ErrorFilePath, pos.Ln+1)
		pos = ctx.ParentEntryPos
		ctx = ctx.Parent
	}

	return "Traceback (most recent call last):\n" + result
}

func (e RuntimeError) AsString() string {
	result := e.generateTraceback()
	result += fmt.Sprintf("%s%s: %s%s\n\n", bold, e.ErrorName, e.Details, reset)
	result += stringWithArrows(e.PosStart.Ftxt, *e.PosStart, *e.PosEnd)
	return result
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
