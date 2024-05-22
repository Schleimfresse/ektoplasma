package main

import "fmt"

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

func NewIfCaseNode(condition, expr Node) *IfCaseNode {
	return &IfCaseNode{condition, expr}
}

func NewIfNode(cases []*IfCaseNode, elseCase *ParseResult) *IfNode {
	var posEnd *Position

	if elseCase != nil {
		posEnd = elseCase.Node.PosEnd()
	} else {
		lastCondition := cases[len(cases)-1].Condition
		posEnd = lastCondition.PosEnd()
	}
	return &IfNode{cases, elseCase.Node.(*NumberNode), cases[0].Expr.PosStart(), posEnd}
}

func NewForNode(varNameTok *Token, startValueNode, endValueNode, stepValueNode, bodyNode Node) *ForNode {
	return &ForNode{
		VarNameTok:     varNameTok,
		StartValueNode: startValueNode,
		EndValueNode:   endValueNode,
		StepValueNode:  stepValueNode,
		BodyNode:       bodyNode,
		PositionStart:  varNameTok.PosStart,
		PositionEnd:    bodyNode.PosEnd(),
	}
}

func NewWhileNode(conditionNode, bodyNode Node) *WhileNode {
	return &WhileNode{
		ConditionNode: conditionNode,
		BodyNode:      bodyNode,
		PositionStart: conditionNode.PosStart(),
		PositionEnd:   bodyNode.PosEnd(),
	}
}

func NewFuncDefNode(varNameTok *Token, argNameToks []*Token, bodyNode Node) *FuncDefNode {
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
	return ""
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
