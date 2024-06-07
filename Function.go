package main

import (
	"fmt"
	"strconv"
	"strings"
)

// NewFunction creates a new Function instance.
func NewFunction(name *string, bodyNode *Node, argNames []string, Flag bool) *Value {
	baseFunc := NewBaseFunction(name)
	return &Value{Function: &Function{bodyNode, argNames, baseFunc, Flag}}
}

// Execute executes the function with the given arguments.
func (f *Function) Execute(args []*Value) *RTResult {
	res := NewRTResult()
	interpreter := NewInterpreter()
	execCtx := f.Base.GenerateNewContext()

	res.Register(f.Base.CheckAndPopulateArgs(f.ArgNames, args, execCtx, false))
	if res.ShouldReturn() {
		return res
	}

	if f.BodyNode == nil {
		res.Failure(NewRTError(
			f.PosStart(), f.PosEnd(),
			"Body of function '"+f.Base.Name+"' is not defined",
			f.Base.Context,
		))
	}

	value := res.Register(interpreter.visit(*f.BodyNode, execCtx))
	if res.ShouldReturn() && res.FuncReturnValue == nil {
		return res
	}

	var ReturnValue *Value
	if f.Flag {
		if !value.IsEmpty() {
			ReturnValue = value
		} else if res.FuncReturnValue != nil {
			ReturnValue = res.FuncReturnValue
		} else {
			ReturnValue = NewNull()
		}
	} else {
		if res.FuncReturnValue != nil {
			ReturnValue = res.FuncReturnValue
		} else {
			ReturnValue = NewNull()
		}
	}
	return res.Success(ReturnValue)
}

// Copy creates a copy of the function.
func (f *Function) Copy() *Value {
	return NewFunction(&f.Base.Name, f.BodyNode, f.ArgNames, f.Flag).SetContext(f.Base.Context).SetPos(f.PosStart(), f.PosEnd())
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

func (b *BaseFunction) CheckArgs(argNames []string, args []*Value, variadic bool) *RTResult {
	res := NewRTResult()

	if variadic {
		return res.Success(nil)
	}

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

func (b *BaseFunction) PopulateArgs(argNames []string, args []*Value, execCtx *Context, variadic bool) {
	if variadic {
		value := NewVariadicArray(args)
		execCtx.SymbolTable.Set(argNames[0], value, false)
	} else {
		for i, argName := range argNames {
			argValue := args[i]
			/*if argValue.Function != nil {
				argValue = argValue.Function.Execute(argValue.)
			} else if argValue.BuildInFunction != nil {

			} else if argValue.StdLibFunction != nil {

			}*/
			value := argValue.SetContext(execCtx)
			execCtx.SymbolTable.Set(argName, value, false)
		}
	}
}

func (b *BaseFunction) CheckAndPopulateArgs(argNames []string, args []*Value, execCtx *Context, variadic bool) *RTResult {
	res := NewRTResult()
	res.Register(b.CheckArgs(argNames, args, variadic))
	if res.Error != nil {
		return res
	}
	b.PopulateArgs(argNames, args, execCtx, variadic)
	return res.Success(nil)
}

// TODO rethink approach, maybe general declaration?

func NewBuildInFunction(name string) *Value {
	BuildInFn := &BuildInFunction{Base: NewBaseFunction(&name), Methods: make(map[string]Method)}
	BuildInFn.Methods["print"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executePrint}
	BuildInFn.Methods["println"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executePrintLn}
	BuildInFn.Methods["Input"] = Method{ArgsNames: nil, Fn: BuildInFn.executeInput}
	BuildInFn.Methods["isString"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executeIsString}
	BuildInFn.Methods["isNumber"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executeIsNumber}
	BuildInFn.Methods["isFunction"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.executeIsFunction}
	BuildInFn.Methods["isArray"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.ExecuteIsArray}
	BuildInFn.Methods["append"] = Method{ArgsNames: []string{"array", "value"}, Fn: BuildInFn.ExecuteAppend}
	BuildInFn.Methods["len"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.ExecuteLen}
	BuildInFn.Methods["pop"] = Method{ArgsNames: []string{"array", "index"}, Fn: BuildInFn.ExecutePop}
	BuildInFn.Methods["str"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.ExecuteStr}
	BuildInFn.Methods["num"] = Method{ArgsNames: []string{"value"}, Fn: BuildInFn.ExecuteNum}

	return &Value{BuildInFunction: BuildInFn}

}

func (b *BuildInFunction) Execute(args ...*Value) *RTResult {
	res := NewRTResult()
	execCtx := b.Base.GenerateNewContext()
	method, ok := b.Methods[b.Base.Name]
	if !ok {
		b.noVisitMethod()
	}

	if b.Base.Name == "print" || b.Base.Name == "println" {
		res.Register(b.Base.CheckAndPopulateArgs(method.ArgsNames, args, execCtx, true))
	} else {
		res.Register(b.Base.CheckAndPopulateArgs(method.ArgsNames, args, execCtx, false))
		if res.Error != nil {
			return res
		}
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

func (b *BuildInFunction) executeIsNumber(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if exists && value.Number != nil {
		return NewRTResult().Success(NewBoolean(One))
	} else {
		return NewRTResult().Success(NewBoolean(Zero))
	}
}

func (b *BuildInFunction) executeIsFunction(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if exists && value.Function != nil {
		return NewRTResult().Success(NewBoolean(One))
	} else {
		return NewRTResult().Success(NewBoolean(Zero))
	}
}

func (b *BuildInFunction) ExecuteIsArray(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if exists && value.Array != nil {
		return NewRTResult().Success(NewBoolean(One))
	} else {
		return NewRTResult().Success(NewBoolean(Zero))
	}
}

func (b *BuildInFunction) ExecuteAppend(execCtx *Context) *RTResult {
	res := NewRTResult()
	array, exists, _ := execCtx.SymbolTable.Get("array")
	value, _, _ := execCtx.SymbolTable.Get("value")

	if exists && array.Array == nil {
		return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), "First argument must be an array", execCtx))
	}

	array.Array.Elements = append(array.Array.Elements, value)
	return res.Success(NewNull())
}

func (b *BuildInFunction) ExecuteLen(execCtx *Context) *RTResult {
	value, exists, _ := execCtx.SymbolTable.Get("value")
	res := NewRTResult()
	if exists {
		result := value.Length()
		if result != nil {
			return res.Success(result)
		} else {
			return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("Cannot get length of %s", value.Type()), execCtx))
		}
	} else {
		return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("missing argument"), execCtx))
	}
}

func (b *BuildInFunction) ExecutePop(execCtx *Context) *RTResult {
	index, _, _ := execCtx.SymbolTable.Get("index")
	array, exists, _ := execCtx.SymbolTable.Get("array")

	if exists && array.Array != nil {
		if index.Number != nil {
			arr := array.Array.Elements
			idx := index.Number.ValueField.(int)
			if index.Number.ValueField.(int) < 0 || index.Number.ValueField.(int) >= len(arr) {
				return NewRTResult().Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), "Index out of bounds", execCtx))
			}

			element := arr[idx]

			array.Array.Elements = append(arr[:idx], arr[idx+1:]...)

			return NewRTResult().Success(element)
		} else {
			return NewRTResult().Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("Index must be an Number, got: %v", index.Type()), execCtx))
		}
	} else {
		return NewRTResult().Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("First argument must be an Array, got: %v", index.Type()), execCtx))
	}
}

func (b *BuildInFunction) ExecuteStr(execCtx *Context) *RTResult {
	res := NewRTResult()
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if !exists {
		return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), "Missing an argument for a conversion to string", execCtx))
	} else if value.Number != nil {
		switch number := value.Number.ValueField.(type) {
		case int:
			return res.Success(NewString(strconv.Itoa(number)))
		case float64:
			return res.Success(NewString(strconv.Itoa(int(number))))
		}
	} else if value.ByteArray != nil {
		return res.Success(value.ByteArray.ToString())
	}
	return res.Success(NewNull())
}

func (b *BuildInFunction) ExecuteNum(execCtx *Context) *RTResult {
	res := NewRTResult()
	value, exists, _ := execCtx.SymbolTable.Get("value")

	if !exists {
		return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), "Missing an argument for a conversion to Number (integer)", execCtx))
	} else if value.String != nil {
		if strings.Contains(value.String.ValueField, ".") {
			parsed, _ := strconv.ParseFloat(value.String.ValueField, 64)
			return res.Success(NewNumber(parsed))

		} else {
			parsed, _ := strconv.Atoi(value.String.ValueField)
			return res.Success(NewNumber(parsed))
		}
	} else if value.ByteArray != nil {
		result, err := value.ByteArray.ToNumber()
		if err != nil {
			return res.Failure(err)
		}
		return res.Success(result)
	} else if value.Number != nil {
		return res.Success(value)
	}
	return res.Failure(NewRTError(b.Base.PosStart(), b.Base.PosEnd(), fmt.Sprintf("Can not use given argument of type %s", value.Type()), b.Base.Context))
}
