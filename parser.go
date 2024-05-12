package main

import (
	"fmt"
	"strconv"
)

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
		err := NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '+', '-', '*' or '/'")
		return res.Failure(err.Error)
	}
	return res
}

// Factor parses a factor.
func (p *Parser) Factor() *ParseResult {
	res := &ParseResult{}
	tok := p.Current

	if tok.Type == TT_PLUS || tok.Type == TT_MINUS {
		p.Advance() // Advance token index here
		factor := res.Register(p.Factor())
		return res.Success(NewUnaryOpNode(tok, factor))
	} else if tok.Type == TT_POW {
		p.Advance()
		base := res.Register(p.Factor())
		if res.Error != nil {
			return res
		}
		exp := res.Register(p.Factor())
		if res.Error != nil {
			return res
		}
		return res.Success(NewBinOpNode(base, tok, exp))
	} else if tok.Type == TT_INT || tok.Type == TT_FLOAT {
		p.Advance() // Advance token index here
		var err error
		if tok.Type == TT_FLOAT {
			tok.Value, err = strconv.ParseFloat(tok.Value.(string), 64)
			if err != nil {
				fmt.Println(err)
			}
		}
		if tok.Type == TT_INT {
			value, err := strconv.ParseInt(tok.Value.(string), 10, 64)
			if err != nil {
				fmt.Println(err)
			}
			tok.Value = int(value)
		}
		return res.Success(NewNumberNode(tok))
	} else if tok.Type == TT_LPAREN {
		p.Advance() // Advance token index here
		expr := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		if p.Current.Type == TT_RPAREN {
			p.Advance() // Advance token index here
			return res.Success(expr)
		}
		err := NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, "Expected ')'")
		return res.Failure(err.Error)
	}
	err := NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, "Expected int or float")

	return res.Failure(err.Error)
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []TokenTypes{TT_MUL, TT_DIV, TT_POW})
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
			return res.Failure(NewInvalidSyntaxError(opTok.PosStart, opTok.PosEnd, err.Error()).Error)
		}

		// Parse the right side of the expression
		right := res.Register(funcParser())
		if res.Error != nil {
			return res
		}
		left = NewBinOpNode(left, opTok, right)
		fmt.Println("LEFT:", left)
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
