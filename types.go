package main

// UnaryOpNode represents a unary operation node.
type UnaryOpNode struct {
	OpTok *Token
	Node  Node
}

// BinOpNode represents a binary operation node.
type BinOpNode struct {
	LeftNode  Node
	OpTok     *Token
	RightNode Node
}

type NumberNode struct {
	Tok *Token
}

// Lexer represents a lexer for tokenizing the code.
type Lexer struct {
	Fn          string
	Text        string
	Pos         *Position
	CurrentChar byte
}

// Token represents a token in the code.
type Token struct {
	Type     TokenTypes
	Value    interface{}
	PosStart *Position // Add PosStart field
	PosEnd   *Position // Add PosEnd field
}

// IllegalCharError represents an error for illegal characters.
type IllegalCharError struct {
	Error
}

// Position represents a position in the code.
type Position struct {
	Idx  int
	Ln   int
	Col  int
	Fn   string
	Ftxt string
}

// TokenTypes represents the different types of tokens.
type TokenTypes string

// Node represents a generic node.
type Node interface {
	String() string
}

type Parser struct {
	Tokens  []*Token
	TokIdx  int
	Current *Token
}

type ParseResult struct {
	Error *Error
	Node  Node
}

type Error struct {
	PosStart  Position
	PosEnd    Position
	ErrorName string
	Details   string
}
