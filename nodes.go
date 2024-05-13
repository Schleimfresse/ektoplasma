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
