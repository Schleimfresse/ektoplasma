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

// NewUnaryOpNode creates a new UnaryOpNode instance.
func NewUnaryOpNode(opTok *Token, node Node) *UnaryOpNode {
	return &UnaryOpNode{opTok, node, nil}
}

// NewVarAccessNode creates a new VarAccessNode instance.
func NewVarAccessNode(varNameTok *Token) *VarAccessNode {
	return &VarAccessNode{varNameTok, varNameTok.PosStart, varNameTok.PosEnd}
}

func NewIndexNode(varAccessNode *VarAccessNode, index *NumberNode) *IndexNode {
	return &IndexNode{varAccessNode, index, varAccessNode.PosStart(), index.PosEnd()}
}

// NewVarAssignNode creates a new VarAssignNode instance.
func NewVarAssignNode(varNameTok *Token, valueNode *Node, isConst bool, declaration bool) *VarAssignNode {
	return &VarAssignNode{varNameTok, valueNode, isConst, declaration, varNameTok.PosStart, varNameTok.PosEnd}
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

func NewImportNode(importNames []*Token, packageName []*Token, posStart *Position, posEnd *Position) *ImportNode {
	return &ImportNode{importNames, packageName, posStart, posEnd}
}

func NewPackageMethod(packageTok *Token, methodName string, callNode Node) *PackageMethod {
	return &PackageMethod{
		PackageName:   packageTok.Value.(string),
		MethodName:    methodName,
		CallNode:      &callNode,
		PositionStart: packageTok.PosStart,
		PositionEnd:   callNode.PosEnd(),
	}
}

// String returns the string representation of the array node.
func (i *ImportNode) String() string {
	var importNames []Token
	var packageNames []Token
	for _, name := range i.ImportNames {
		importNames = append(importNames, *name)
	}
	for _, name := range i.PackageNames {
		packageNames = append(packageNames, *name)
	}
	return fmt.Sprintf("(%v, %v)", importNames, packageNames)
}

func (i *ImportNode) PosStart() *Position {
	return i.PositionStart
}

func (i *ImportNode) PosEnd() *Position {
	return i.PositionEnd
}

func (r *ReturnNode) String() string {
	return fmt.Sprintf("(%v)", r.NodeToReturn)
}

func (r *ReturnNode) PosStart() *Position {
	return r.PositionStart
}

func (r *ReturnNode) PosEnd() *Position {
	return r.PositionEnd
}

func (c *ContinueNode) String() string {
	return "<continue>"
}

func (c *ContinueNode) PosStart() *Position {
	return c.PositionStart
}

func (c *ContinueNode) PosEnd() *Position {
	return c.PositionEnd
}

func (b *BreakNode) String() string {
	return "<break>"
}

func (b *BreakNode) PosStart() *Position {
	return b.PositionStart
}

func (b *BreakNode) PosEnd() *Position {
	return b.PositionEnd
}

func (a *ArrayNode) String() string {
	return fmt.Sprintf("(%v)", a.ElementNodes)
}

func (a *ArrayNode) PosStart() *Position {
	return a.PositionStart
}

func (a *ArrayNode) PosEnd() *Position {
	return a.PositionEnd
}

func (c *CallNode) String() string {
	return fmt.Sprintf("(%v, %v)", c.ArgNodes, c.NodeToCall)
}

func (c *CallNode) PosStart() *Position {
	return c.PositionStart
}

func (c *CallNode) PosEnd() *Position {
	return c.PositionEnd
}

func (f *FuncDefNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", f.VarNameTok, f.ArgNameToks, f.BodyNode)
}

func (f *FuncDefNode) PosStart() *Position {
	return f.PositionStart
}

func (f *FuncDefNode) PosEnd() *Position {
	return f.PositionEnd
}

func (w *WhileNode) String() string {
	return fmt.Sprintf("(%v, %v)", w.ConditionNode, w.BodyNode)
}

func (w *WhileNode) PosStart() *Position {
	return w.PositionStart
}

func (w *WhileNode) PosEnd() *Position {
	return w.PositionEnd
}

func (f *ForNode) String() string {
	return fmt.Sprintf("(%v, %v)", f.BodyNode, f.EndValueNode)
}

func (f *ForNode) PosStart() *Position {
	return f.PositionStart
}

func (f *ForNode) PosEnd() *Position {
	return f.PositionEnd
}

func (i *IfNode) String() string {
	return fmt.Sprintf("(cases: %v, elsecase: %v)", i.Cases, i.ElseCase)
}

func (i *IfNode) PosStart() *Position {
	return i.PositionStart
}

func (i *IfNode) PosEnd() *Position {
	return i.PositionEnd
}

func (v *VarAssignNode) PosStart() *Position {
	return v.VarNameTok.PosStart
}

func (v *VarAssignNode) PosEnd() *Position {
	return v.VarNameTok.PosEnd
}

func (v *VarAssignNode) String() string {
	return fmt.Sprintf("(%v, %v)", v.VarNameTok, v.ValueNode)
}

func (n *NumberNode) PosStart() *Position {
	return n.Tok.PosStart
}

func (n *NumberNode) PosEnd() *Position {
	return n.Tok.PosEnd
}

func (b *BinOpNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", b.LeftNode, b.OpTok, b.RightNode)
}

func (s *StringNode) String() string {
	return fmt.Sprintf("%v", s.Tok)
}

func (s *StringNode) PosStart() *Position {
	return s.Tok.PosStart
}

func (s *StringNode) PosEnd() *Position {
	return s.Tok.PosEnd
}

func (b *BinOpNode) PosStart() *Position {
	return b.LeftNode.PosStart()
}

func (b *BinOpNode) PosEnd() *Position {
	return b.RightNode.PosEnd()
}

func (u *UnaryOpNode) PosStart() *Position {
	return u.Position
}

func (u *UnaryOpNode) PosEnd() *Position {
	return u.Position
}

func (u *UnaryOpNode) String() string {
	return fmt.Sprintf("(%v, %v)", u.OpTok, u.Node)
}

func (v *VarAccessNode) PosStart() *Position {
	return v.VarNameTok.PosStart
}

func (v *VarAccessNode) PosEnd() *Position {
	return v.VarNameTok.PosEnd
}

func (v *VarAccessNode) String() string {
	return fmt.Sprintf("(%v)", v.VarNameTok)
}

func (i *IfCaseNode) PosStart() *Position {
	return i.Condition.PosStart()
}

func (i *IfCaseNode) PosEnd() *Position {
	return i.Condition.PosEnd()
}

func (i *IfCaseNode) String() string {
	return fmt.Sprintf("%v, %v", i.Expr, i.Condition)
}

func (e *ElseCaseNode) PosStart() *Position {
	return e.Expr.PosStart()
}

func (e *ElseCaseNode) PosEnd() *Position {
	return e.Expr.PosEnd()
}

func (e *ElseCaseNode) String() string {
	return fmt.Sprintf("%v, %v", e.Expr, e.Flag)
}

func (i *IndexNode) PosStart() *Position {
	return i.PositionStart
}

func (i *IndexNode) PosEnd() *Position {
	return i.PositionEnd
}

func (i *IndexNode) String() string {
	return fmt.Sprintf("(%v, %v)", i.VarAccessNode, i.IndexNode)
}

func (p *PackageMethod) PosStart() *Position {
	return p.PositionStart
}

func (p *PackageMethod) PosEnd() *Position {
	return p.PositionEnd
}

func (p *PackageMethod) String() string {
	return fmt.Sprintf("(%v, %v, %v)", p.PackageName, p.MethodName, p.CallNode)
}
