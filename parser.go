package main

import (
	"fmt"
	"log"
	"strconv"
)

// Register registers the result of a parsing operation.
func (pr *ParseResult) Register(res *ParseResult) Node {
	pr.AdvanceCount += res.AdvanceCount
	if res.Error != nil {
		pr.Error = res.Error
	}
	return res.Node
}

// Success returns a successful parsing result.
func (pr *ParseResult) Success(node Node) *ParseResult {
	pr.Node = node
	return pr
}

// Failure returns a failed parsing result.
func (pr *ParseResult) Failure(err Error) *ParseResult {
	if pr.Error == nil || pr.AdvanceCount == 0 {
		pr.Error = &err
	}
	return pr
}

// Failure returns a failed parsing result.
func (pr *ParseResult) RegisterAdvancement() int {
	pr.AdvanceCount += 1
	return pr.AdvanceCount
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
func (p *Parser) Advance() *Token {
	p.TokIdx++
	if p.TokIdx < len(p.Tokens) {
		p.Current = p.Tokens[p.TokIdx]
		return nil
	}
	return p.Current
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

func (p *Parser) Atom() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	tok := p.Current

	if tok.Type == TT_INT || tok.Type == TT_FLOAT {
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
	} else if tok.Type == TT_IDENTIFIER {
		res.RegisterAdvancement()
		p.Advance()
		return res.Success(NewVarAccessNode(tok))
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
	log.Println(p.Current.Type, p.Current.Value, p.Current.PosStart, p.Current.PosEnd)
	return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected int, float, identifier, '+', '-' or '('").Error)
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Atom, []TokenTypes{TT_POW})
}

// Factor parses a factor.
func (p *Parser) Factor() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	tok := p.Current
	log.Println(tok.Type, p.Tokens)
	if tok.Type == TT_PLUS || tok.Type == TT_MINUS {
		p.Advance() // Advance token index here
		factor := res.Register(p.Factor())
		return res.Success(NewUnaryOpNode(tok, factor))
	} /*else if tok.Type == TT_POW {
		/*p.Advance()
		base := res.Register(p.Factor())
		if res.Error != nil {
			fmt.Println("1")
			return res.Failure(NewInvalidSyntaxError(base.PosStart(), base.PosEnd(), "Expected int or float").Error)
		}
		exp := res.Register(p.Factor())
		log.Println(exp)
		if res.Error != nil {
			fmt.Println("2")
			return res.Failure(NewInvalidSyntaxError(exp.PosStart(), exp.PosEnd(), "Expected int or float").Error)
		}
		return res.Success(NewBinOpNode(base, tok, exp))
		return p.BinOp(p.Factor, []TokenTypes{TT_POW})
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
	}*/
	//err := NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, "Expected int or float")

	//return res.Failure(err.Error)
	return p.Power()
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []TokenTypes{TT_MUL, TT_DIV, TT_POW})
}

// Expr parses an expression.
func (p *Parser) Expr() *ParseResult {
	res := ParseResult{}
	log.Println(p.Current.Matches(TT_KEYWORD, "VAR"), p.Current)
	if p.Current.Matches(TT_KEYWORD, "VAR") {
		res.RegisterAdvancement()
		p.Advance()

		log.Println("IMP ", p.Current, p.Tokens)
		if p.Current.Type != TT_IDENTIFIER {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected identifier",
			).Error)
		}

		varName := p.Current
		res.RegisterAdvancement()
		p.Advance()
		log.Println("ADVANCE EP4:", p.Current)
		if p.Current.Type != TT_EQ {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected '='",
			).Error)
		}

		res.RegisterAdvancement()
		p.Advance()

		expr := res.Register(p.Expr())
		if res.Error != nil {
			return &res
		}
		return res.Success(NewVarAssignNode(varName, expr))
	}

	node := res.Register(p.BinOp(p.Term, []TokenTypes{TT_PLUS, TT_MINUS}))

	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'VAR', int, float, identifier, '+', '-' or '('").Error)
	}

	return res.Success(node)
}

// BinOp parses a binary operation.
func (p *Parser) BinOp(funcParser func() *ParseResult, ops []TokenTypes) *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	left := res.Register(funcParser())
	if res.Error != nil {
		return res
	}

	for ContainsType(ops, p.Current.Type) {
		opTok := p.Current
		res.RegisterAdvancement()
		p.Advance()

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
