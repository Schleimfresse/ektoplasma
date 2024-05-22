package main

import (
	"fmt"
	"log"
	"reflect"
)

func NewInterpreter() Interpreter {
	return Interpreter{}
}

func (i *Interpreter) visit(node Node, context *Context) *RTResult {
	switch n := node.(type) {
	case *UnaryOpNode:
		return i.visitUnaryOpNode(*n, context)
	case *BinOpNode:
		return i.visitBinOpNode(*n, context)
	case *NumberNode:
		return i.visitNumberNode(*n, context)
	case *VarAccessNode:
		return i.visitVarAccessNode(*n, context)
	case *VarAssignNode:
		return i.visitVarAssignNode(*n, context)
	case *IfNode:
		return i.visitIfNode(*n, context)
	case *ForNode:
		return i.visitForNode(*n, context)
	case *WhileNode:
		return i.visitWhileNode(*n, context)
	case *CallNode:
		return i.visitCallNode(*n, context)
	case *FuncDefNode:
		return i.visitFuncDefNode(*n, context)
	case *StringNode:
		return i.visitStringNode(*n, context)
	case *ArrayNode:
		return i.visitArrayNode(*n, context)
	default:
		// Handle unknown node types
		return NewRTResult().Failure(NewRTError(node.PosStart(), node.PosEnd(), fmt.Sprintf("No visit method defined for node type %T", node), context))
	}
}

// NoVisitMethod raises an exception for non-existing visit methods.
func (i *Interpreter) NoVisitMethod(node Node, context *Context) *RTResult {
	errMsg := fmt.Sprintf("No visit_%T method defined", reflect.TypeOf(node))
	return NewRTResult().Failure(NewRTError(node.PosStart(), node.PosEnd(), errMsg, context))
}

func (i *Interpreter) visitNumberNode(node NumberNode, context *Context) *RTResult {
	value := NewNumber(node.Value)
	value.SetContext(context).SetPos(node.PosStart(), node.PosEnd())
	return NewRTResult().Success(value)
}

func (i *Interpreter) visitStringNode(node StringNode, context *Context) *RTResult {
	if TokenValue, ok := node.Value.(string); ok {
		value := NewString(TokenValue)
		value.SetContext(context).SetPos(node.PosStart(), node.PosEnd())
		return NewRTResult().Success(value)
	}
	return NewRTResult().Failure(NewRTError(node.PosStart(), node.PosEnd(), fmt.Sprintf("%s is not of type string", reflect.TypeOf(node.Value)), context))
}

func (i *Interpreter) visitArrayNode(node ArrayNode, context *Context) *RTResult {
	res := NewRTResult()
	var elements []Value

	for idx := 0; idx < len(node.ElementNodes); idx++ {
		elements = append(elements, *res.Register(i.visit(node.ElementNodes[idx], context)))
		if res.Error != nil {
			return res
		}
	}

	return res.Success(NewArray(elements).SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visitBinOpNode(node BinOpNode, context *Context) *RTResult {
	res := NewRTResult()

	leftRTValue := i.visit(node.LeftNode, context)
	res.Register(leftRTValue)
	if res.Error != nil {
		return res
	}

	rightRTValue := i.visit(node.RightNode, context)
	res.Register(rightRTValue)
	if res.Error != nil {
		return res
	}

	left := leftRTValue.Value
	right := rightRTValue.Value

	var result *Value
	var err *RuntimeError

	// get the operation type and use the left and the right.Number node from the operation symbol as values
	switch node.OpTok.Type {
	case TT_PLUS:
		if left.String != nil && right.String != nil {
			result, err = left.String.AddedTo(right.String)
		} else if left.Number != nil && right.Number != nil {
			result, err = left.Number.AddedTo(right.Number)
		} else if left.Array != nil && right.Number != nil {
			result, err = left.Array.AddedTo(right)
		} else {
			return res.Failure(NewRTError(node.OpTok.PosStart, node.OpTok.PosEnd, "Invalid types for addition", context))
		}
	case TT_MINUS:
		if left.Array != nil && right.Number != nil {
			result, err = left.Array.SubtractedBy(right)
		} else {
			result, err = left.Number.SubtractedBy(right.Number)
		}
	case TT_MUL:
		if left.String != nil && right.Number != nil {
			result, err = left.String.MultipliedBy(right.Number)
		} else if left.Number != nil && right.Number != nil {
			result, err = left.Number.MultipliedBy(right.Number)
		} else if left.Array != nil && right.Array != nil {
			result, err = left.Array.MultipliedBy(right)
		} else {
			return res.Failure(NewRTError(node.OpTok.PosStart, node.OpTok.PosEnd, "Invalid types for multiplication", context))
		}
	case TT_DIV:
		if left.Array != nil && right.Number != nil {
			result, err = left.Array.DividedBy(right)
		} else {
			result, err = left.Number.DividedBy(right.Number)
		}
	case TT_POW:
		result, err = left.Number.PowedBy(right.Number)
	case TT_EE:
		result, err = left.Number.GetComparisonEq(right.Number)
	case TT_NE:
		result, err = left.Number.GetComparisonNe(right.Number)
	case TT_LT:
		result, err = left.Number.GetComparisonLt(right.Number)
	case TT_GT:
		result, err = left.Number.GetComparisonGt(right.Number)
	case TT_LTE:
		result, err = left.Number.GetComparisonLte(right.Number)
	case TT_GTE:
		result, err = left.Number.GetComparisonGte(right.Number)
	case TT_KEYWORD:
		if node.OpTok.Value == "AND" {
			result, err = left.Number.AndedBy(right.Number)
		} else if node.OpTok.Value == "OR" {
			result, err = left.Number.OredBy(right.Number)
		}
	default:
		return res.Failure(NewRTError(node.OpTok.PosStart, node.OpTok.PosEnd, "Invalid operation", context))
	}
	if err != nil {
		return res.Failure(err)
	}

	result.SetContext(context).SetPos(node.PosStart(), node.PosEnd())

	return res.Success(result)
}

func (i *Interpreter) visitUnaryOpNode(node UnaryOpNode, context *Context) *RTResult {
	res := NewRTResult()
	numValue := res.Register(i.visit(node.Node, context))
	if res.Error != nil {
		return res
	}

	var result *Value
	var err *RuntimeError

	num := numValue.Number
	if num == nil {
		return res.Failure(NewRTError(node.Node.PosStart(), node.Node.PosEnd(), "Expected a number", context))
	}

	// else if for some reason required, when not expressions like +1 won't work because the context is not set
	if node.OpTok.Type == TT_MINUS {
		result, err = num.MultipliedBy(NewNumber(-1).Number)
		if err != nil {
			return res.Failure(err)
		}
	} else if node.OpTok.Type == TT_PLUS {
		result, err = num.MultipliedBy(NewNumber(1).Number)
		if err != nil {
			return res.Failure(err)
		}
	} else if node.OpTok.Matches(TT_KEYWORD, "NOT") {
		result, err = num.Notted()
	}

	if err != nil {
		return res.Failure(err)
	} else {
		result.SetContext(context).SetPos(node.PosStart(), node.PosEnd())
		return res.Success(result)
	}
}

// visitVarAccessNode visits a VarAccessNode and retrieves its value from the symbol table.
func (i *Interpreter) visitVarAccessNode(node VarAccessNode, context *Context) *RTResult {
	res := NewRTResult()
	varName := node.VarNameTok.Value

	value, exists := context.SymbolTable.Get(varName.(string))
	if !exists {
		return res.Failure(NewRTError(
			node.PosStart(), node.PosEnd(),
			fmt.Sprintf("'%s' is not defined", varName),
			context))
	}

	if value.Number != nil {
		value = *value.Number.Copy().SetPos(value.Number.PosStart(), value.Number.PosEnd()).SetContext(context)
	} else if value.Function != nil {
		value = *value.Function.Copy().SetPos(value.Function.PosStart(), value.Function.PosEnd()).SetContext(context)
	}

	return res.Success(&value)
}

// visitVarAssignNode visits a VarAssignNode and assigns a value to the variable in the symbol table.
func (i *Interpreter) visitVarAssignNode(node VarAssignNode, context *Context) *RTResult {
	res := NewRTResult()
	varName := node.VarNameTok.Value

	value := res.Register(i.visit(node.ValueNode, context))
	if res.Error != nil {
		return res
	}

	context.SymbolTable.Set(varName.(string), *value)
	return res.Success(value)
}

func (i *Interpreter) visitIfNode(node IfNode, context *Context) *RTResult {
	res := NewRTResult()

	for _, ifcase := range node.Cases {

		value := res.Register(i.visit(ifcase.Condition, context))
		if res.Error != nil {
			return res
		}

		conditionValue := value.Number
		if conditionValue.IsTrue() {
			exprValue := res.Register(i.visit(ifcase.Expr, context))
			if res.Error != nil {
				return res
			}
			return res.Success(exprValue)
		}
	}

	if node.ElseCase != nil {
		elseValue := res.Register(i.visit(node.ElseCase, context))
		if res.Error != nil {
			return res
		}
		return res.Success(elseValue)
	}

	return res.Success(nil)
}

func (i *Interpreter) visitForNode(node ForNode, context *Context) *RTResult {
	res := NewRTResult()
	var elements []Value

	start := res.Register(i.visit(node.StartValueNode, context))
	if res.Error != nil {
		return res
	}

	end := res.Register(i.visit(node.EndValueNode, context))
	if res.Error != nil {
		return res
	}

	var stepValue *Number
	if node.StepValueNode != nil {
		stepValue = res.Register(i.visit(node.StepValueNode, context)).Number
		if res.Error != nil {
			return res
		}
	} else {
		stepValue = NewNumber(1).Number
	}

	var iVal int
	var endValue int

	if IsInt(start.Number.ValueField) {
		iVal = start.Number.ValueField.(int)
	} else {
		startNum := start.Number
		return res.Failure(NewRTError(startNum.PosStart(), startNum.PosEnd(), fmt.Sprintf("Can not use type %s as type int", reflect.TypeOf(startNum.ValueField)), startNum.Context))
	}
	if IsInt(end.Number.ValueField) {
		endValue = end.Number.ValueField.(int)
	} else {
		endNum := end.Number
		return res.Failure(NewRTError(endNum.PosStart(), endNum.PosEnd(), fmt.Sprintf("Can not use type %s as type int", reflect.TypeOf(endNum.ValueField)), endNum.Context))
	}

	var condition func() bool
	if stepValue.ValueField.(int) >= 0 {
		condition = func() bool { return iVal < endValue }
	} else {
		condition = func() bool { return iVal > endValue }
	}

	for condition() {
		context.SymbolTable.Set(node.VarNameTok.Value.(string), *NewNumber(iVal))

		iVal += stepValue.ValueField.(int)

		elements = append(elements, *res.Register(i.visit(node.BodyNode, context)))
		if res.Error != nil {
			return res
		}
	}

	return res.Success(NewArray(elements).SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visitWhileNode(node WhileNode, context *Context) *RTResult {
	res := NewRTResult()
	var elements []Value

	for {
		condition := res.Register(i.visit(node.ConditionNode, context))
		if res.Error != nil {
			return res
		}

		if !condition.Number.IsTrue() {
			break
		}

		elements = append(elements, *res.Register(i.visit(node.BodyNode, context)))
		if res.Error != nil {
			return res
		}
	}

	return res.Success(NewArray(elements).SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visitFuncDefNode(node FuncDefNode, context *Context) *RTResult {
	res := NewRTResult()

	var funcName *string
	if node.VarNameTok != nil {
		value := node.VarNameTok.Value.(string)
		funcName = &value
	} else {
		funcName = nil
	}

	argNames := make([]string, len(node.ArgNameToks))
	for idx, argName := range node.ArgNameToks {
		argNames[idx] = argName.Value.(string)
	}

	value := NewFunction(funcName, &node.BodyNode, argNames)
	value.SetContext(context).SetPos(node.PosStart(), node.PosEnd())

	if node.VarNameTok != nil {
		context.SymbolTable.Set(*funcName, *value)
	}

	return res.Success(value)
}

func (i *Interpreter) visitCallNode(node CallNode, context *Context) *RTResult {
	res := NewRTResult()
	args := make([]*Value, len(node.ArgNodes))
	Call := res.Register(i.visit(node.NodeToCall, context))
	log.Println(node.NodeToCall, context, Call)

	if res.Error != nil {
		return res
	}

	valueToCall := Call.SetPos(node.PosStart(), node.PosEnd())

	for idx, argNode := range node.ArgNodes {
		args[idx] = res.Register(i.visit(argNode, context))
		if res.Error != nil {
			return res
		}
	}
	var returnValue *Value
	if valueToCall.Function != nil {
		returnValue = res.Register(valueToCall.Function.Execute(args))
	} else if valueToCall.BuildInFunction != nil {
		returnValue = res.Register(valueToCall.BuildInFunction.Execute(args))
	}

	if res.Error != nil {
		return res
	}

	returnValue = returnValue.Copy().SetPos(node.PosStart(), node.PosEnd()).SetContext(context)
	return res.Success(returnValue)
}

// NewRTResult creates a new RTResult instance.
func NewRTResult() *RTResult {
	return &RTResult{}
}

// Register registers the result of a runtime operation.
func (r *RTResult) Register(res *RTResult) *Value {
	if res.Error != nil {
		r.Error = res.Error
	}
	return res.Value
}

// Success indicates a successful runtime operation.
func (r *RTResult) Success(value *Value) *RTResult {
	r.Value = value
	return r
}

// Failure indicates a failed runtime operation.
func (r *RTResult) Failure(error *RuntimeError) *RTResult {
	r.Error = error
	return r
}

// NewContext creates a new context with the given display name, parent, and parent entry position.
func NewContext(displayName string, parent *Context, parentEntryPos *Position) *Context {
	if parent == nil {
		parent = &Context{}
	}
	return &Context{
		DisplayName:    displayName,
		Parent:         parent,
		ParentEntryPos: parentEntryPos,
		SymbolTable:    &SymbolTable{},
	}
}

// NewSymbolTable creates a new SymbolTable instance.
func NewSymbolTable(symboltable *SymbolTable) *SymbolTable {
	if symboltable == nil {

		return &SymbolTable{
			symbols: make(map[string]Value),
			parent:  nil,
		}
	} else {
		return &SymbolTable{symbols: symboltable.symbols, parent: symboltable.parent}
	}
}

// Get retrieves the value associated with the name from the symbol table.
func (st *SymbolTable) Get(name string) (Value, bool) {
	value, exists := st.symbols[name]
	if !exists && st.parent != nil {
		return st.parent.Get(name)
	}
	return value, exists
}

// Set sets the value associated with the name in the symbol table.
func (st *SymbolTable) Set(name string, value Value) {
	st.symbols[name] = value
}

// Remove removes the entry associated with the name from the symbol table.
func (st *SymbolTable) Remove(name string) {
	delete(st.symbols, name)
}

func IsInt(value interface{}) bool {
	switch value.(type) {
	case int:
		return true
	default:
		return false
	}
}
