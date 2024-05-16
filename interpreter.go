package main

import (
	"fmt"
	"reflect"
)

func (i *Interpreter) visit(node Node, context *Context) *RTResult {
	//log.Println("NODE HEADER:", reflect.TypeOf(node))
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
	return NewRTResult().Success(
		NewNumber(node.Value).SetContext(context).SetPos(node.PosStart(), node.PosEnd()),
	)
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

	left := leftRTValue.Value.(*Number)
	right := rightRTValue.Value.(*Number)

	var result *Number
	var err *RuntimeError

	// get the operation type and use the left and the right node from the operation symbol as values
	switch node.OpTok.Type {
	case TT_PLUS:
		result, err = left.AddedTo(right)
	case TT_MINUS:
		result, err = left.SubtractedBy(right)
	case TT_MUL:
		result, err = left.MultipliedBy(right)
	case TT_DIV:
		result, err = left.DividedBy(right)
	case TT_POW:
		result, err = left.PowedBy(right)
	case TT_EE:
		result, err = left.GetComparisonEq(right)
	case TT_NE:
		result, err = left.GetComparisonNe(right)
	case TT_LT:
		result, err = left.GetComparisonLt(right)
	case TT_GT:
		result, err = left.GetComparisonGt(right)
	case TT_LTE:
		result, err = left.GetComparisonLte(right)
	case TT_GTE:
		result, err = left.GetComparisonGte(right)
	case TT_KEYWORD:
		if node.OpTok.Value == "AND" {
			result, err = left.AndedBy(right)
		} else if node.OpTok.Value == "OR" {
			result, err = left.OredBy(right)
		}
	default:
		return res.Failure(NewRTError(node.OpTok.PosStart, node.OpTok.PosEnd, "Invalid operation", context))
	}
	if err != nil {
		return res.Failure(err)
	}
	return res.Success(result.SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visitUnaryOpNode(node UnaryOpNode, context *Context) *RTResult {
	res := NewRTResult()
	numValue := res.Register(i.visit(node.Node, context))
	if res.Error != nil {
		return res
	}

	var result *Number
	var err *RuntimeError

	num, ok := numValue.(*Number)
	if !ok {
		return res.Failure(NewRTError(node.Node.PosStart(), node.Node.PosEnd(), "Expected a number", context))
	}

	// else if for some reason required, when not expressions like +1 won't work because the context is not set
	if node.OpTok.Type == TT_MINUS {
		result, err = num.MultipliedBy(NewNumber(-1))
		if err != nil {
			return res.Failure(err)
		}
	} else if node.OpTok.Type == TT_PLUS {
		result, err = num.MultipliedBy(NewNumber(1))
		if err != nil {
			return res.Failure(err)
		}
	} else if node.OpTok.Matches(TT_KEYWORD, "NOT") {
		result, err = num.Notted()
	}

	if err != nil {
		return res.Failure(err)
	} else {
		return res.Success(result.SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
	}
}

// visitVarAccessNode visits a VarAccessNode and retrieves its value from the symbol table.
func (i *Interpreter) visitVarAccessNode(node VarAccessNode, context *Context) *RTResult {
	res := NewRTResult()
	varName := node.VarNameTok.Value

	value, exists := context.SymbolTable.Get(varName.(string))
	f := value.(*Number)
	if !exists {
		return res.Failure(NewRTError(
			node.PosStart(), node.PosEnd(),
			fmt.Sprintf("'%s' is not defined", varName),
			context))
	}

	value = f.SetPos(node.PosStart(), node.PosEnd())
	return res.Success(value)
}

// visitVarAssignNode visits a VarAssignNode and assigns a value to the variable in the symbol table.
func (i *Interpreter) visitVarAssignNode(node VarAssignNode, context *Context) *RTResult {
	res := NewRTResult()
	varName := node.VarNameTok.Value

	value := res.Register(i.visit(node.ValueNode, context))
	if res.Error != nil {
		return res
	}

	context.SymbolTable.Set(varName.(string), value)
	return res.Success(value)
}

func (i *Interpreter) visitIfNode(node IfNode, context *Context) *RTResult {
	res := NewRTResult()

	for _, ifcase := range node.Cases {

		value := res.Register(i.visit(ifcase.Condition, context))
		if res.Error != nil {
			return res
		}

		conditionValue := value.(*Number)
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
		stepValue = res.Register(i.visit(node.StepValueNode, context)).(*Number)
		if res.Error != nil {
			return res
		}
	} else {
		stepValue = NewNumber(1)
	}

	var iVal int
	var endValue int

	if IsInt(start.(*Number).Value) {
		iVal = start.(*Number).Value.(int)
	} else {
		startNum := start.(*Number)
		return res.Failure(NewRTError(startNum.PosStart, startNum.PosEnd, fmt.Sprintf("Can not use type %s as type int", reflect.TypeOf(startNum.Value)), startNum.Context))
	}
	if IsInt(end.(*Number).Value) {
		endValue = end.(*Number).Value.(int)
	} else {
		endNum := end.(*Number)
		return res.Failure(NewRTError(endNum.PosStart, endNum.PosEnd, fmt.Sprintf("Can not use type %s as type int", reflect.TypeOf(endNum.Value)), endNum.Context))
	}

	var condition func() bool
	if stepValue.Value.(int) >= 0 {
		condition = func() bool { return iVal < endValue }
	} else {
		condition = func() bool { return iVal > endValue }
	}

	for condition() {
		context.SymbolTable.Set(node.VarNameTok.Value.(string), NewNumber(iVal))

		iVal += stepValue.Value.(int)

		res.Register(i.visit(node.BodyNode, context))
		if res.Error != nil {
			return res
		}
	}

	return res.Success(nil)
}

func (i *Interpreter) visitWhileNode(node WhileNode, context *Context) *RTResult {
	res := NewRTResult()

	for {
		condition := res.Register(i.visit(node.ConditionNode, context)).(*Number)
		if res.Error != nil {
			return res
		}

		if !condition.IsTrue() {
			break
		}

		res.Register(i.visit(node.BodyNode, context))
		if res.Error != nil {
			return res
		}
	}

	return res.Success(nil)
}

// NewRTResult creates a new RTResult instance.
func NewRTResult() *RTResult {
	return &RTResult{}
}

// Register registers the result of a runtime operation.
func (r *RTResult) Register(res *RTResult) interface{} {
	if res.Error != nil {
		r.Error = res.Error
	}
	return res.Value
}

// Success indicates a successful runtime operation.
func (r *RTResult) Success(value interface{}) *RTResult {
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
	return &Context{
		DisplayName:    displayName,
		Parent:         parent,
		ParentEntryPos: parentEntryPos,
	}
}

// NewSymbolTable creates a new SymbolTable instance.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]interface{}),
		parent:  nil,
	}
}

// Get retrieves the value associated with the name from the symbol table.
func (st *SymbolTable) Get(name string) (interface{}, bool) {
	value, exists := st.symbols[name]
	if !exists && st.parent != nil {
		return st.parent.Get(name)
	}
	return value, exists
}

// Set sets the value associated with the name in the symbol table.
func (st *SymbolTable) Set(name string, value interface{}) {
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
