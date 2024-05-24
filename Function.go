package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"
)

// NewFunction creates a new Function instance.
func NewFunction(name *string, bodyNode *Node, argNames []string) *Value {
	baseFunc := NewBaseFunction(name)
	return &Value{Function: &Function{bodyNode, argNames, baseFunc}}
}

// Execute executes the function with the given arguments.
func (f *Function) Execute(args []*Value) *RTResult {
	res := NewRTResult()
	interpreter := NewInterpreter()
	execCtx := f.Base.GenerateNewContext()

	res.Register(f.Base.CheckAndPopulateArgs(f.ArgNames, args, execCtx))
	if res.Error != nil {
		return res
	}

	if f.BodyNode == nil {
		res.Failure(NewRTError(
			f.PosStart(), f.PosEnd(),
			"Body of function '"+f.Base.Name+"' is not defined",
			f.Base.Context,
		))
	}

	Value := res.Register(interpreter.visit(*f.BodyNode, execCtx))
	if res.Error != nil {
		return res
	}
	return res.Success(Value)
}

// Copy creates a copy of the function.
func (f *Function) Copy() *Value {
	return NewFunction(&f.Base.Name, f.BodyNode, f.ArgNames).SetContext(f.Base.Context).SetPos(f.PosStart(), f.PosEnd())
}

// String returns the string representation of the function.
func (f *Function) String() string {
	return "<function " + f.Base.Name + ">"
}

func (f *Function) PosStart() *Position {
	return f.Base.PositionStart
}

func (f *Function) PosEnd() *Position {
	return f.Base.PositionEnd
}

func (f *Function) IllegalOperation(other *Value) *RuntimeError {
	if other.Function == nil {
		other.Function = f
	}
	return NewRTError(f.PosStart(), f.PosEnd(), "Illegal operation", f.Base.Context)
}

// NewBaseFunction creates a new BaseFunction instance.
func NewBaseFunction(name *string) *BaseFunction {
	if name == nil {
		var a = "<anonymous>"
		name = &a
	}

	return &BaseFunction{*name, nil, nil, nil}
}

func (b *BaseFunction) PosStart() *Position {
	return b.PositionStart
}

func (b *BaseFunction) PosEnd() *Position {
	return b.PositionEnd
}

func (b *BaseFunction) GenerateNewContext() *Context {
	newContext := NewContext(b.Name, b.Context, b.PosStart())
	log.Println(newContext.Parent, b.Context)
	newContext.SymbolTable = NewSymbolTable(newContext.Parent.SymbolTable)
	return newContext
}

func (b *BaseFunction) CheckArgs(argNames []string, args []*Value) *RTResult {
	res := NewRTResult()

	if len(args) > len(argNames) {
		return res.Failure(NewRTError(
			b.PosStart(), b.PosEnd(),
			"Too many args passed into '"+b.Name+"'",
			b.Context,
		))
	}

	if len(args) < len(argNames) {
		return res.Failure(NewRTError(
			b.PosStart(), b.PosEnd(),
			"Too few args passed into '"+b.Name+"'",
			b.Context,
		))
	}

	return res.Success(nil)
}

func (b *BaseFunction) PopulateArgs(argNames []string, args []*Value, execCtx *Context) {
	for i, argName := range argNames {
		argValue := args[i]
		value := argValue.SetContext(execCtx)
		execCtx.SymbolTable.Set(argName, value)
	}
}

func (b *BaseFunction) CheckAndPopulateArgs(argNames []string, args []*Value, execCtx *Context) *RTResult {
	res := NewRTResult()
	res.Register(b.CheckArgs(argNames, args))
	if res.Error != nil {
		return res
	}
	b.PopulateArgs(argNames, args, execCtx)
	return res.Success(nil)
}

func NewBuildInFunction(name string) *Value {
	BuildInFn := &BuildInFunction{Base: NewBaseFunction(&name), Methods: make(map[string]Method)}
	BuildInFn.Methods["Print"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executePrint}
	BuildInFn.Methods["Input"] = Method{ArgsNames: nil, Fn: BuildInFn.executeInput}
	BuildInFn.Methods["isString"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executeIsString}
	BuildInFn.Methods["isNumber"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executeIsNumber}
	BuildInFn.Methods["isFunction"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executeIsFunction}
	BuildInFn.Methods["isArray"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.ExecuteIsArray}
	BuildInFn.Methods["append"] = Method{ArgsNames: []string{"array", "value"}, Fn: BuildInFn.ExecuteAppend}

	return &Value{BuildInFunction: BuildInFn}

}

func (b *BuildInFunction) Execute(args []*Value) *RTResult {
	res := NewRTResult()
	execCtx := b.Base.GenerateNewContext()

	method, ok := b.Methods[b.Base.Name]
	if !ok {
		b.noVisitMethod()
	}

	res.Register(b.Base.CheckAndPopulateArgs(method.ArgsNames, args, execCtx))
	if res.Error != nil {
		return res
	}
	returnValue := res.Register(method.Fn(execCtx))
	if res.Error != nil {
		return res
	}

	return res.Success(returnValue)
}

// Default method for handling unknown functions
func (b *BuildInFunction) noVisitMethod() {
	panic("No execute" + b.Base.Name + " method defined")
}

func (b *BuildInFunction) String() string {
	return "<build-in function " + b.Base.Name + ">"
}

func (b *BuildInFunction) Copy() *Value {
	return NewBuildInFunction(b.Base.Name).SetContext(b.Base.Context).SetPos(b.Base.PosStart(), b.Base.PosEnd())
}

func (b *BuildInFunction) executePrint(execCtx *Context) *RTResult {
	value, exists := execCtx.SymbolTable.Get("value")

	if exists {
		_, err := syscall.Write(syscall.Stdout, interfaceToBytes(value.Value()))
		if err != nil {
			return nil
		}
	}
	null, _ := GlobalSymbolTable.Get("null")
	return NewRTResult().Success(null)
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
	value, exists := execCtx.SymbolTable.Get("value")

	if exists && value.String != nil {
		t, _ := GlobalSymbolTable.Get("true")
		return NewRTResult().Success(t)
	} else {
		f, _ := GlobalSymbolTable.Get("false")
		return NewRTResult().Success(f)
	}
}

func (b *BuildInFunction) executeIsNumber(execCtx *Context) *RTResult {
	value, exists := execCtx.SymbolTable.Get("value")

	if exists && value.Number != nil {
		t, _ := GlobalSymbolTable.Get("true")
		return NewRTResult().Success(t)
	} else {
		f, _ := GlobalSymbolTable.Get("false")
		return NewRTResult().Success(f)
	}
}

func (b *BuildInFunction) executeIsFunction(execCtx *Context) *RTResult {
	value, exists := execCtx.SymbolTable.Get("value")

	if exists && value.Function != nil {
		t, _ := GlobalSymbolTable.Get("true")
		return NewRTResult().Success(t)
	} else {
		f, _ := GlobalSymbolTable.Get("false")
		return NewRTResult().Success(f)
	}
}

func (b *BuildInFunction) ExecuteIsArray(execCtx *Context) *RTResult {
	value, exists := execCtx.SymbolTable.Get("value")

	if exists && value.Array != nil {
		t, _ := GlobalSymbolTable.Get("true")
		return NewRTResult().Success(t)
	} else {
		f, _ := GlobalSymbolTable.Get("false")
		return NewRTResult().Success(f)
	}
}

func (b *BuildInFunction) ExecuteAppend(execCtx *Context) *RTResult {
	array, _ := execCtx.SymbolTable.Get("array")
	value, _ := execCtx.SymbolTable.Get("value")

	if array.Array == nil {
		return NewRTResult().Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), "First argument must be an array", execCtx))
	}

	array.Array.Elements = append(array.Array.Elements, value)
	return NewRTResult().Success(NewNull())
}
