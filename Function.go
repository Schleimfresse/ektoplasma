package main

import (
	"fmt"
	"log"
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
	log.Println(f.Base)
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
		value := *argValue.SetContext(execCtx)
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
	log.Println("METHOD:", method)
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
		fmt.Println(value.Value())
	}
	null, _ := GlobalSymbolTable.Get("null")
	return NewRTResult().Success(&null)
}
