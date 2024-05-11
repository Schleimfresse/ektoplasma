package main

import "fmt"

// Register registers the result of a parsing operation.
func (pr *ParseResult) Register(res *ParseResult) Node {
	if res != nil {
		if res.Error != nil {
			pr.Error = res.Error
		}
		return res.Node
	}
	return nil
}

// Success returns a successful parsing result.
func (pr *ParseResult) Success(node Node) *ParseResult {
	pr.Node = node
	return pr
}

// Failure returns a failed parsing result.
func (pr *ParseResult) Failure(err Error) *ParseResult {
	pr.Error = &err
	return pr
}

// NewParser creates a new Parser instance.
func NewParser(tokens []*Token) *Parser {
	parser := &Parser{
		Tokens: tokens,
		TokIdx: -1,
	}
	parser.Advance()
	return parser
}

// Advance moves the parser to the next token.
func (p *Parser) Advance() error {
	p.TokIdx++
	if p.TokIdx < len(p.Tokens) {
		p.Current = p.Tokens[p.TokIdx]
		return nil
	}
	return fmt.Errorf("reached end of tokens")
}

// Parse parses the tokens into an abstract syntax tree.
func (p *Parser) Parse() *ParseResult {
	res := p.Expr()
	if res.Error == nil && p.Current.Type != TT_EOF {
		err := &InvalidSyntaxError{
			Error{PosStart: *p.Current.PosStart,
				PosEnd:    *p.Current.PosEnd,
				Details:   "Expected '+', '-', '*' or '/'",
				ErrorName: "InvalidSyntaxError"}}
		return res.Failure(err.Error)
	}
	fmt.Println("LOG: ", p.Current.Type, res.Error)
	return res
}

// Factor parses a factor.
func (p *Parser) Factor() *ParseResult {
	res := &ParseResult{}
	tok := p.Current

	if tok.Type == TT_PLUS || tok.Type == TT_MINUS {
		factor := res.Register(p.Factor())
		return res.Success(NewUnaryOpNode(tok, factor))
	} else if tok.Type == TT_INT || tok.Type == TT_FLOAT {
		if err := p.Advance(); err != nil {
			return res.Failure(NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, err.Error()).Error)
		}
		return res.Success(NewNumberNode(tok))
	} else if tok.Type == TT_LPAREN {
		if err := p.Advance(); err != nil {
			return res.Failure(NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, err.Error()).Error)
		}
		expr := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		if p.Current.Type == TT_RPAREN {
			if err := p.Advance(); err != nil {
				return res.Failure(NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, err.Error()).Error)
			}
			return res.Success(expr)
		}
		err := NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, "Expected ')'")
		return res.Failure(err.Error)
	}
	err := &InvalidSyntaxError{Error{
		PosStart: *tok.PosStart,
		PosEnd:   *tok.PosEnd,
		Details:  "Expected int or float"},
	}
	return res.Failure(err.Error)
}

// Term parses a term.
func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []TokenTypes{TT_MUL, TT_DIV})
}

// Expr parses an expression.
func (p *Parser) Expr() *ParseResult {
	return p.BinOp(p.Term, []TokenTypes{TT_PLUS, TT_MINUS})
}

// BinOp parses a binary operation.
func (p *Parser) BinOp(funcParser func() *ParseResult, ops []TokenTypes) *ParseResult {
	res := &ParseResult{}
	left := res.Register(funcParser())
	if res.Error != nil {
		return res
	}

	for ContainsType(ops, p.Current.Type) {
		opTok := p.Current
		if err := p.Advance(); err != nil {
			return res.Failure(NewInvalidSyntaxError(opTok.PosStart, opTok.PosEnd, "Error advancing parser: "+err.Error()).Error)
		}

		// Parse the right side of the expression
		right := res.Register(funcParser())
		if res.Error != nil {
			return res
		}
		left = NewBinOpNode(left, opTok, right)
	}

	return res.Success(left)
}

// ContainsType checks if the token type exists in the list.
func ContainsType(types []TokenTypes, typ TokenTypes) bool {
	for _, t := range types {
		if t == typ {
			return true
		}
	}
	return false
}
