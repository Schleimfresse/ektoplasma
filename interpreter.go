package main

import (
	"fmt"
	"reflect"
)

// PosStart returns the start position of the number node.
func (n *NumberNode) PosStart() *Position {
	return n.Position
}

// PosEnd returns the end position of the number node.
func (n *NumberNode) PosEnd() *Position {
	return n.Position
}

// PosStart returns the start position of the binary operation node.
func (b *BinOpNode) PosStart() *Position {
	return b.Position
}

// PosEnd returns the end position of the binary operation node.
func (b *BinOpNode) PosEnd() *Position {
	return b.Position
}

// PosStart returns the start position of the unary operation node.
func (u *UnaryOpNode) PosStart() *Position {
	return u.Position
}

// PosEnd returns the end position of the unary operation node.
func (u *UnaryOpNode) PosEnd() *Position {
	return u.Position
}

func (i *Interpreter) visit(node Node, context *Context) *RTResult {
	switch n := node.(type) {
	case *UnaryOpNode:
		return i.visit_UnaryOpNode(*n, context)
	case *BinOpNode:
		return i.visit_BinOpNode(*n, context)
	case *NumberNode:
		return i.visit_NumberNode(*n, context)
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

func (i *Interpreter) visit_NumberNode(node NumberNode, context *Context) *RTResult {
	fmt.Println(node.Tok.Value, reflect.TypeOf(node.Tok.Value))
	return NewRTResult().Success(
		NewNumber(node.Tok.Value).SetContext(context).SetPos(node.PosStart(), node.PosEnd()),
	)
}

func (i *Interpreter) visit_BinOpNode(node BinOpNode, context *Context) *RTResult {
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
	/*fmt.Println("VALS:", rightRTValue.Value, leftRTValue.Value, reflect.TypeOf(leftRTValue.Value))
	left := NewNumber(leftRTValue.Value)
	fmt.Println("f", node.LeftNode, reflect.TypeOf(node.LeftNode))
	fmt.Println("LEFT 2", reflect.TypeOf(node))
	fmt.Println("left n VAL:", left.Value)
	right := NewNumber(rightRTValue.Value)*/

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
	default:
		return res.Failure(NewRTError(node.OpTok.PosStart, node.OpTok.PosEnd, "Invalid operation", context))
	}
	fmt.Println("RESULT", result, right.Value, left.Value)
	if err != nil {
		return res.Failure(err)
	}

	return res.Success(result.SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
}

func (i *Interpreter) visit_UnaryOpNode(node UnaryOpNode, context *Context) *RTResult {
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

	if node.OpTok.Type == TT_MINUS {
		result, err = num.MultipliedBy(NewNumber(-1))
		if err != nil {
			return res.Failure(err)
		}
	}

	return res.Success(result.SetContext(context).SetPos(node.PosStart(), node.PosEnd()))
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
