package main

type Binary int

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

type StringNode struct {
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

type IndexNode struct {
	VarAccessNode *VarAccessNode
	IndexNode     *NumberNode
	PositionStart *Position
	PositionEnd   *Position
}

// VarAssignNode represents a variable assignment node
type VarAssignNode struct {
	VarNameTok    *Token
	ValueNode     Node
	isConst       bool
	declaration   bool
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
	AdvanceCount   int
	ToReverseCount int
	Error          *Error
	Node           Node
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
	Value              *Value
	Error              *RuntimeError
	FuncReturnValue    *Value
	LoopShouldContinue bool
	LoopShouldBreak    bool
}

// SymbolTable represents a symbol table in the interpreter.
type SymbolTable struct {
	symbols   map[string]*Value
	constants map[string]*Value
	buildIn   map[string]*Value
	packages  map[string]*Package
	parent    *SymbolTable
}

type IfCaseNode struct {
	Condition Node
	Expr      Node
	Flag      bool
}

type ElseCaseNode struct {
	Expr Node
	Flag bool
}

type IfNode struct {
	Cases         []*IfCaseNode
	ElseCase      *ElseCaseNode
	PositionStart *Position
	PositionEnd   *Position
}

type WhileNode struct {
	ConditionNode Node
	BodyNode      Node
	PositionStart *Position
	PositionEnd   *Position
	Flag          bool
}

type ForNode struct {
	VarNameTok     *Token
	StartValueNode Node
	EndValueNode   Node
	StepValueNode  Node
	BodyNode       Node
	PositionStart  *Position
	PositionEnd    *Position
	Flag           bool
}

type FuncDefNode struct {
	VarNameTok    *Token
	ArgNameToks   []*Token
	BodyNode      Node
	PositionStart *Position
	PositionEnd   *Position
	Flag          bool
}

type CallNode struct {
	NodeToCall    Node
	ArgNodes      []Node
	PositionStart *Position
	PositionEnd   *Position
}

type ArrayNode struct {
	ElementNodes  []Node
	PositionStart *Position
	PositionEnd   *Position
}

type ReturnNode struct {
	NodeToReturn  Node
	PositionStart *Position
	PositionEnd   *Position
}

type BreakNode struct {
	PositionStart *Position
	PositionEnd   *Position
}

type ContinueNode struct {
	PositionStart *Position
	PositionEnd   *Position
}

type ImportNode struct {
	ImportNames                []*Token
	PackageNames               []*Token
	PositionStart, PositionEnd *Position
}

type ReferenceNode struct {
	Target        Node
	PositionStart *Position
	PositionEnd   *Position
}

type DereferenceNode struct {
	Target        Node
	PositionStart *Position
	PositionEnd   *Position
}

// Function represents a function value.
type Function struct {
	BodyNode *Node
	ArgNames []string
	Base     *BaseFunction
	Flag     bool
}

type BuildInFunction struct {
	Base    *BaseFunction
	Methods map[string]Method
}

type StdLibFunction struct {
	Base        *BaseFunction
	PackageName string
	Function    func(args []*Value) *RTResult
}

type Method struct {
	ArgsNames []string
	Fn        func(ctx *Context) *RTResult
}

type BaseFunction struct {
	Name                       string
	PositionStart, PositionEnd *Position
	Context                    *Context
}

// Number represents a numeric value.
type Number struct {
	ValueField                 interface{}
	PositionStart, PositionEnd *Position
	Context                    *Context
}

// String represents a String value.
type String struct {
	ValueField                 string
	PositionStart, PositionEnd *Position
	Context                    *Context
	Type                       string
}

type Array struct {
	Elements                   []*Value
	PositionStart, PositionEnd *Position
	Context                    *Context
}

type Null struct {
	PositionStart, PositionEnd *Position
	Context                    *Context
}

type Boolean struct {
	PositionStart, PositionEnd *Position
	Context                    *Context
	Binary                     Binary
}

type Value struct {
	Number          *Number
	Function        *Function
	BuildInFunction *BuildInFunction
	StdLibFunction  *StdLibFunction
	String          *String
	Array           *Array
	Null            *Null
	Boolean         *Boolean
	Pointer         *Pointer
	Dereference     *Dereference
	ByteArray       *ByteArray
	VariadicArray   *VariadicArray
}

type Package struct {
	Methods map[string]*Value
}

type PackageMethod struct {
	PositionStart, PositionEnd *Position
	PackageName                string
	MethodName                 string
	CallNode                   *Node
}

type Pointer struct {
	Addr          string
	PositionStart *Position
	PositionEnd   *Position
}

type Dereference struct {
	Value         Value
	PositionStart *Position
	PositionEnd   *Position
}

type ByteArray struct {
	ValueField    []byte
	PositionStart *Position
	PositionEnd   *Position
	Context       *Context
}

type VariadicArray struct {
	Array         []*Value
	PositionStart *Position
	PositionEnd   *Position
	Context       *Context
}
