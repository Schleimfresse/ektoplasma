package main

import "fmt"

// ecp Ektoplasma

const (
	TT_INT    TokenTypes = "INT"
	TT_FLOAT  TokenTypes = "FLOAT"
	TT_PLUS   TokenTypes = "PLUS"
	TT_MINUS  TokenTypes = "MINUS"
	TT_MUL    TokenTypes = "MUL"
	TT_DIV    TokenTypes = "DIV"
	TT_LPAREN TokenTypes = "LPAREN"
	TT_RPAREN TokenTypes = "RPAREN"
	TT_EOF    TokenTypes = "EOF"
)

// NewNumberNode creates a new NumberNode instance.
func NewNumberNode(tok *Token) *NumberNode {
	return &NumberNode{tok}
}

// String returns the string representation of the NumberNode.
func (n *NumberNode) String() string {
	return fmt.Sprintf("%v", n.Tok)
}

// NewBinOpNode creates a new BinOpNode instance.
func NewBinOpNode(left Node, opTok *Token, right Node) *BinOpNode {
	return &BinOpNode{left, opTok, right}
}

// String returns the string representation of the BinOpNode.
func (b *BinOpNode) String() string {
	return fmt.Sprintf("(%v, %v, %v)", b.LeftNode, b.OpTok, b.RightNode)
}

// NewUnaryOpNode creates a new UnaryOpNode instance.
func NewUnaryOpNode(opTok *Token, node Node) *UnaryOpNode {
	return &UnaryOpNode{opTok, node}
}

// String returns the string representation of the UnaryOpNode.
func (u *UnaryOpNode) String() string {
	return fmt.Sprintf("(%v, %v)", u.OpTok, u.Node)
}

func main() {

	lexer := NewLexer("test", "30 + 2 * 4 - 4 + (4*4)")
	tokens, err := lexer.MakeTokens()
	if err != nil {
		panic(err)
	}
	for _, token := range tokens {
		if token.PosStart != nil && token.PosEnd != nil {
			fmt.Println(token.Type, token.Value, "START:", *token.PosStart, "END:", *token.PosEnd)
		} else {
			fmt.Println(token.Type, "START:", token.PosStart, "END:", token.PosEnd)
		}
	}
	parser := NewParser(tokens)
	ast := parser.Parse()
	if ast.Error != nil {
		fmt.Println(ast.Error.AsString())
	} else {
		fmt.Println(ast.Node)
	}
}
