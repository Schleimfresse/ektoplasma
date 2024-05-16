package main

// UnaryOpNode represents a unary operation node.
type UnaryOpNode struct {
	OpTok    *Token
	Node     Node
	Position *Position
}

// BinOpNode represents a binary operation node.
type BinOpNode struct {
	LeftNode      Node
	OpTok         *Token
	RightNode     Node
	PositionStart *Position
	PositionEnd   *Position
}

type NumberNode struct {
	Tok           *Token
	Value         interface{}
	PositionStart *Position
	PositionEnd   *Position
}

// VarAccessNode represents a variable access node
type VarAccessNode struct {
	VarNameTok    *Token
	PositionStart *Position
	PositionEnd   *Position
}

// VarAssignNode represents a variable assignment node
type VarAssignNode struct {
	VarNameTok    *Token
	ValueNode     Node
	PositionStart *Position
	PositionEnd   *Position
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

type TokenTypeInfo struct {
	Type  TokenTypes
	Value *string
}

// Node represents a generic node.
type Node interface {
	String() string
	PosStart() *Position
	PosEnd() *Position
}

type Parser struct {
	Tokens  []*Token
	TokIdx  int
	Current *Token
}

type ParseResult struct {
	AdvanceCount int
	Error        *Error
	Node         Node
}

type Error struct {
	PosStart  *Position
	PosEnd    *Position
	ErrorName string
	Details   string
}

type Interpreter struct{}

type Context struct {
	DisplayName    string
	Parent         *Context
	ParentEntryPos *Position
	SymbolTable    *SymbolTable
}

type Number struct {
	Value    interface{}
	PosStart *Position
	PosEnd   *Position
	Context  *Context
}

type RuntimeError struct {
	*Error
	Context *Context
}

// InvalidSyntaxError represents an error for invalid syntax.
type InvalidSyntaxError struct {
	Error
}

// ExpectedCharError represents an error for an expected character.
type ExpectedCharError struct {
	Error
}

// RTResult represents the result of a runtime operation.
type RTResult struct {
	Value interface{}
	Error *RuntimeError
}

// SymbolTable represents a symbol table in the interpreter.
type SymbolTable struct {
	symbols map[string]interface{}
	parent  *SymbolTable
}

type IfCaseNode struct {
	Condition Node
	Expr      Node
}

type IfNode struct {
	Cases         []*IfCaseNode
	ElseCase      *NumberNode
	PositionStart *Position
	PositionEnd   *Position
}

type WhileNode struct {
	ConditionNode Node
	BodyNode      Node
	PositionStart *Position
	PositionEnd   *Position
}

type ForNode struct {
	VarNameTok     *Token
	StartValueNode Node
	EndValueNode   Node
	StepValueNode  Node
	BodyNode       Node
	PositionStart  *Position
	PositionEnd    *Position
}
