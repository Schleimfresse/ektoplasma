package main

import (
	"fmt"
)

// NewNumberNode creates a new NumberNode instance.
func NewNumberNode(tok *Token) *NumberNode {
	return &NumberNode{tok, tok.Value, tok.PosStart, tok.PosEnd}
}

// String returns the string representation of the NumberNode.
func (n *NumberNode) String() string {
	return fmt.Sprintf("%v", n.Tok)
}

// NewStringNode creates a new StringNode instance.
func NewStringNode(tok *Token) *StringNode {
	return &StringNode{tok, tok.Value, tok.PosStart, tok.PosEnd}
}

// NewBinOpNode creates a new BinOpNode instance.
func NewBinOpNode(left Node, opTok *Token, right Node) *BinOpNode {
	return &BinOpNode{left, opTok, right, left.PosStart(), right.PosEnd()}
}

// String returns the string representation of the BinOpNode.
func (b *BinOpNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", b.LeftNode, b.OpTok, b.RightNode)
}

// NewUnaryOpNode creates a new UnaryOpNode instance.
func NewUnaryOpNode(opTok *Token, node Node) *UnaryOpNode {
	return &UnaryOpNode{opTok, node, nil}
}

// String returns the string representation of the UnaryOpNode.
func (u *UnaryOpNode) String() string {
	return fmt.Sprintf("(%v, %v)", u.OpTok, u.Node)
}

// NewVarAccessNode creates a new VarAccessNode instance.
func NewVarAccessNode(varNameTok *Token) *VarAccessNode {
	return &VarAccessNode{varNameTok, varNameTok.PosStart, varNameTok.PosEnd}
}

// String returns the string representation of the UnaryOpNode.
func (v *VarAccessNode) String() string {
	return fmt.Sprintf("(%v)", v.VarNameTok)
}

// NewVarAssignNode creates a new VarAssignNode instance.
func NewVarAssignNode(varNameTok *Token, valueNode Node) *VarAssignNode {
	return &VarAssignNode{varNameTok, valueNode, varNameTok.PosStart, varNameTok.PosEnd}
}

func NewIfCaseNode(condition, expr Node, flag bool) *IfCaseNode {
	return &IfCaseNode{condition, expr, flag}
}

func NewIfNode(cases []*IfCaseNode, elseCase *ElseCaseNode) *IfNode {
	var posEnd *Position
	var posStart *Position

	if elseCase != nil {
		posEnd = elseCase.PosEnd()
	} else {
		posEnd = cases[len(cases)-1].Condition.PosEnd()
	}

	// in case we handle an else, cases will be nil
	if cases == nil {
		posStart = elseCase.PosStart()
	} else {
		posStart = cases[0].PosStart()
	}

	return &IfNode{cases, elseCase, posStart, posEnd}
}

func NewElseCaseNode(statement Node, flag bool) *ElseCaseNode {
	return &ElseCaseNode{statement, flag}
}

func NewForNode(varNameTok *Token, startValueNode, endValueNode, stepValueNode, bodyNode Node, Flag bool) *ForNode {
	return &ForNode{
		VarNameTok:     varNameTok,
		StartValueNode: startValueNode,
		EndValueNode:   endValueNode,
		StepValueNode:  stepValueNode,
		BodyNode:       bodyNode,
		PositionStart:  varNameTok.PosStart,
		PositionEnd:    bodyNode.PosEnd(),
		Flag:           Flag,
	}
}

func NewWhileNode(conditionNode, bodyNode Node, Flag bool) *WhileNode {
	return &WhileNode{
		ConditionNode: conditionNode,
		BodyNode:      bodyNode,
		PositionStart: conditionNode.PosStart(),
		PositionEnd:   bodyNode.PosEnd(),
		Flag:          Flag,
	}
}

func NewFuncDefNode(varNameTok *Token, argNameToks []*Token, bodyNode Node, Flag bool) *FuncDefNode {
	var posStart, posEnd *Position

	if varNameTok != nil {
		posStart = varNameTok.PosStart
	} else if len(argNameToks) > 0 {
		posStart = argNameToks[0].PosStart
	} else {
		posStart = bodyNode.PosStart()
	}

	posEnd = bodyNode.PosEnd()

	return &FuncDefNode{
		VarNameTok:    varNameTok,
		ArgNameToks:   argNameToks,
		BodyNode:      bodyNode,
		PositionStart: posStart,
		PositionEnd:   posEnd,
		Flag:          Flag,
	}
}

func NewCallNode(nodeToCall Node, argNodes []Node) *CallNode {
	var posStart, posEnd *Position

	posStart = nodeToCall.PosStart()

	if len(argNodes) > 0 {
		posEnd = argNodes[len(argNodes)-1].PosEnd()
	} else {
		posEnd = nodeToCall.PosEnd()
	}

	return &CallNode{
		NodeToCall:    nodeToCall,
		ArgNodes:      argNodes,
		PositionStart: posStart,
		PositionEnd:   posEnd,
	}
}

func NewArrayNode(ElementNodes []Node, PosStart *Position, PosEnd *Position) *ArrayNode {
	return &ArrayNode{ElementNodes, PosStart, PosEnd}
}

func NewReturnNode(NodeToReturn Node, PosStart *Position, PosEnd *Position) *ReturnNode {
	return &ReturnNode{NodeToReturn, PosStart, PosEnd}
}

func NewContinueNode(PosStart *Position, PosEnd *Position) *ContinueNode {
	return &ContinueNode{PosStart, PosEnd}
}

func NewBreakNode(PosStart *Position, PosEnd *Position) *BreakNode {
	return &BreakNode{PosStart, PosEnd}
}

func NewImportNode(functionName *Token, modulName *Token, posStart *Position, posEnd *Position) *ImportNode {
	return &ImportNode{functionName, modulName, posStart, posEnd}
}

// String returns the string representation of the array node.
func (i *ImportNode) String() string {
	return fmt.Sprintf("(%v, %v)", i.FunctionName, i.ModuleName)
}

// PosStart returns the start position of the array node.
func (i *ImportNode) PosStart() *Position {
	return i.PositionStart
}

// PosEnd returns the end position of the array node.
func (i *ImportNode) PosEnd() *Position {
	return i.PositionEnd
}

// String returns the string representation of the array node.
func (r *ReturnNode) String() string {
	return fmt.Sprintf("(%v)", r.NodeToReturn)
}

// PosStart returns the start position of the array node.
func (r *ReturnNode) PosStart() *Position {
	return r.PositionStart
}

// PosEnd returns the end position of the array node.
func (r *ReturnNode) PosEnd() *Position {
	return r.PositionEnd
}

// String returns the string representation of the array node.
func (c *ContinueNode) String() string {
	return "<continue>"
}

// PosStart returns the start position of the array node.
func (c *ContinueNode) PosStart() *Position {
	return c.PositionStart
}

// PosEnd returns the end position of the array node.
func (c *ContinueNode) PosEnd() *Position {
	return c.PositionEnd
}

// String returns the string representation of the array node.
func (b *BreakNode) String() string {
	return "<break>"
}

// PosStart returns the start position of the array node.
func (b *BreakNode) PosStart() *Position {
	return b.PositionStart
}

// PosEnd returns the end position of the array node.
func (b *BreakNode) PosEnd() *Position {
	return b.PositionEnd
}

// String returns the string representation of the array node.
func (a *ArrayNode) String() string {
	return fmt.Sprintf("(%v)", a.ElementNodes)
}

// PosStart returns the start position of the array node.
func (a *ArrayNode) PosStart() *Position {
	return a.PositionStart
}

// PosEnd returns the end position of the array node.
func (a *ArrayNode) PosEnd() *Position {
	return a.PositionEnd
}

// String returns the string representation of the UnaryOpNode.
func (c *CallNode) String() string {
	return fmt.Sprintf("(%v, %v)", c.ArgNodes, c.NodeToCall)
}

// PosStart returns the start position of the number node.
func (c *CallNode) PosStart() *Position {
	return c.PositionStart
}

// PosEnd returns the end position of the number node.
func (c *CallNode) PosEnd() *Position {
	return c.PositionEnd
}

// String returns the string representation of the UnaryOpNode.
func (f *FuncDefNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", f.VarNameTok, f.ArgNameToks, f.BodyNode)
}

// PosStart returns the start position of the number node.
func (f *FuncDefNode) PosStart() *Position {
	return f.PositionStart
}

// PosEnd returns the end position of the number node.
func (f *FuncDefNode) PosEnd() *Position {
	return f.PositionEnd
}

// String returns the string representation of the UnaryOpNode.
func (w *WhileNode) String() string {
	return fmt.Sprintf("(%v, %v)", w.ConditionNode, w.BodyNode)
}

// PosStart returns the start position of the number node.
func (w *WhileNode) PosStart() *Position {
	return w.PositionStart
}

// PosEnd returns the end position of the number node.
func (w *WhileNode) PosEnd() *Position {
	return w.PositionEnd
}

// String returns the string representation of the UnaryOpNode.
func (f *ForNode) String() string {
	return fmt.Sprintf("(%v, %v)", f.BodyNode, f.EndValueNode)
}

// PosStart returns the start position of the number node.
func (f *ForNode) PosStart() *Position {
	return f.PositionStart
}

// PosEnd returns the end position of the number node.
func (f *ForNode) PosEnd() *Position {
	return f.PositionEnd
}

// String returns the string representation of the UnaryOpNode.
func (i *IfNode) String() string {
	return fmt.Sprintf("(cases: %v, elsecase: %v)", i.Cases, i.ElseCase)
}

// PosStart returns the start position of the number node.
func (i *IfNode) PosStart() *Position {
	return i.PositionStart
}

// PosEnd returns the end position of the number node.
func (i *IfNode) PosEnd() *Position {
	return i.PositionEnd
}

// String returns the string representation of the UnaryOpNode.
func (v *VarAssignNode) String() string {
	return fmt.Sprintf("(%v, %v)", v.VarNameTok, v.ValueNode)
}

// PosStart returns the start position of the number node.
func (n *NumberNode) PosStart() *Position {
	return n.Tok.PosStart
}

// PosEnd returns the end position of the number node.
func (n *NumberNode) PosEnd() *Position {
	return n.Tok.PosEnd
}

func (s *StringNode) String() string {
	return fmt.Sprintf("%v", s.Tok)
}

// PosStart returns the start position of the number node.
func (s *StringNode) PosStart() *Position {
	return s.Tok.PosStart
}

// PosEnd returns the end position of the number node.
func (s *StringNode) PosEnd() *Position {
	return s.Tok.PosEnd
}

// PosStart returns the start position of the binary operation node.
func (b *BinOpNode) PosStart() *Position {
	return b.LeftNode.PosStart()
}

// PosEnd returns the end position of the binary operation node.
func (b *BinOpNode) PosEnd() *Position {
	return b.RightNode.PosEnd()
}

// PosStart returns the start position of the unary operation node.
func (u *UnaryOpNode) PosStart() *Position {
	return u.Position
}

// PosEnd returns the end position of the unary operation node.
func (u *UnaryOpNode) PosEnd() *Position {
	return u.Position
}

func (v *VarAssignNode) PosStart() *Position {
	return v.VarNameTok.PosStart
}

func (v *VarAssignNode) PosEnd() *Position {
	return v.VarNameTok.PosEnd
}

func (v *VarAccessNode) PosStart() *Position {
	return v.VarNameTok.PosStart
}

func (v *VarAccessNode) PosEnd() *Position {
	return v.VarNameTok.PosEnd
}

// PosStart returns the start position of the number node.
func (i *IfCaseNode) PosStart() *Position {
	return i.Condition.PosStart()
}

// PosEnd returns the end position of the number node.
func (i *IfCaseNode) PosEnd() *Position {
	return i.Condition.PosEnd()
}

func (i *IfCaseNode) String() string {
	return fmt.Sprintf("%v, %v", i.Expr, i.Condition)
}

// PosStart returns the start position of the number node.
func (e *ElseCaseNode) PosStart() *Position {
	return e.Expr.PosStart()
}

// PosEnd returns the end position of the number node.
func (e *ElseCaseNode) PosEnd() *Position {
	return e.Expr.PosEnd()
}

func (e *ElseCaseNode) String() string {
	return fmt.Sprintf("%v, %v", e.Expr, e.Flag)
}
