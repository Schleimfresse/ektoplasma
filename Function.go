package main

import (
	"log"
	"reflect"
)

// NewFunction creates a new Function instance.
func NewFunction(name *string, bodyNode *Node, argNames []string) *Value {
	if name == nil {
		var a = "<anonymous>"
		name = &a
	}

	return &Value{Function: &Function{*name, bodyNode, argNames, nil, nil, nil}}
}

// SetPos sets the position of the function.
func (f *Function) SetPos(posStart *Position, posEnd *Position) *Function {
	f.PositionStart = posStart
	f.PositionEnd = posEnd
	return f
}

// SetContext sets the context of the function.
func (f *Function) SetContext(context *Context) *Function {
	f.Context = context
	return f
}

// Execute executes the function with the given arguments.
func (f *Function) Execute(args []*Value) *RTResult {
	res := NewRTResult()
	interpreter := NewInterpreter()
	newContext := NewContext(f.Name, f.Context, f.PosStart())
	newContext.SymbolTable = NewSymbolTable(newContext.Parent.SymbolTable)

	if len(args) > len(f.ArgNames) {
		return res.Failure(NewRTError(
			f.PosStart(), f.PosEnd(),
			"Too many args passed into '"+f.Name+"'",
			f.Context,
		))
	}

	if len(args) < len(f.ArgNames) {
		return res.Failure(NewRTError(
			f.PosStart(), f.PosEnd(),
			"Too few args passed into '"+f.Name+"'",
			f.Context,
		))
	}

	for i, argName := range f.ArgNames {
		argValue := args[i]
		num := argValue.Number.SetContext(newContext)
		newContext.SymbolTable.Set(argName, Value{Number: num})
	}

	if f.BodyNode == nil {
		res.Failure(NewRTError(
			f.PosStart(), f.PosEnd(),
			"Body of function '"+f.Name+"' is not defined",
			f.Context,
		))
	}

	Value := res.Register(interpreter.visit(*f.BodyNode, newContext))
	if res.Error != nil {
		return res
	}

	log.Println("VALUE IN EXECUTE CALL:", reflect.TypeOf(Value), Value)
	return res.Success(Value)
}

// Copy creates a copy of the function.
func (f *Function) Copy() *Function {
	copy := NewFunction(&f.Name, f.BodyNode, f.ArgNames).Function
	copy.SetContext(f.Context)
	copy.SetPos(f.PosStart(), f.PosEnd())
	return copy
}

// String returns the string representation of the function.
func (f *Function) String() string {
	return "<function " + f.Name + ">"
}

func (f *Function) PosStart() *Position {
	return f.PositionStart
}

func (f *Function) PosEnd() *Position {
	return f.PositionEnd
}

func (f *Function) IllegalOperation(other *Value) *RuntimeError {
	if other.Function == nil {
		other.Function = f
	}
	return NewRTError(f.PosStart(), f.PosEnd(), "Illegal operation", f.Context)
}

// TODO ILLEGAL operation err allgemein machen EP. 8, sowie Value type
