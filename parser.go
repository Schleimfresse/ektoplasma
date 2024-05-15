package main

import (
	"fmt"
	"log"
	"reflect"
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
	if res.Error != nil && p.Current.Type != TT_EOF {
		log.Println("IMP ERR:", res.Error)
		err := NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '+', '-', '*', '/', '^', '==', '!=', '<', '>', <=', '>=', 'AND' or 'OR'")
		return res.Failure(err.Error)
	}
	return res
}

func (p *Parser) IfExpr() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	var cases []*IfCaseNode
	var elseCase *ParseResult

	if !p.Current.Matches(TT_KEYWORD, "IF") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'IF'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	condition := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	log.Println("condition from IdExpr: ", condition)

	if !p.Current.Matches(TT_KEYWORD, "THEN") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'THEN'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	expr := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	cases = append(cases, NewIfCaseNode(condition, expr))

	for p.Current.Matches(TT_KEYWORD, "ELIF") {
		res.RegisterAdvancement()
		p.Advance()

		condition := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}

		if !p.Current.Matches(TT_KEYWORD, "THEN") {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected 'THEN'",
			).Error)
		}

		res.RegisterAdvancement()
		p.Advance()

		expr := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		cases = append(cases, NewIfCaseNode(condition, expr))
	}

	if p.Current.Matches(TT_KEYWORD, "ELSE") {
		res.RegisterAdvancement()
		p.Advance()
		log.Println("ELSE TT:", p.Current)
		elseCase = res.Success(res.Register(p.Expr()))
		log.Println("elseCase", elseCase.Node, reflect.TypeOf(elseCase.Node))
		if res.Error != nil {
			return res
		}
	}

	return res.Success(NewIfNode(cases, elseCase))
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
	} else if tok.Matches(TT_KEYWORD, "IF") {
		IfExpr := res.Register(p.IfExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(IfExpr)
	}
	log.Println(p.Current.Type, p.Current.Value, p.Current.PosStart, p.Current.PosEnd)
	return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected int, float, identifier, '+', '-' or '('").Error)
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Atom, []TokenTypeInfo{{TT_POW, nil}})
}

// ArithExpr parses an arithmetic expression.
func (p *Parser) ArithExpr() *ParseResult {
	return p.BinOp(p.Term, []TokenTypeInfo{{TT_PLUS, nil}, {TT_MINUS, nil}})
}

// CompExpr parses a comparison expression.
func (p *Parser) CompExpr() *ParseResult {
	res := ParseResult{AdvanceCount: 0}

	if p.Current.Matches(TT_KEYWORD, "NOT") {
		opTok := p.Current
		res.RegisterAdvancement()
		p.Advance()

		node := res.Register(p.CompExpr())
		if res.Error != nil {
			return &res
		}
		return res.Success(NewUnaryOpNode(opTok, node))
	}

	node := res.Register(p.BinOp(p.ArithExpr, []TokenTypeInfo{{TT_EE, nil}, {TT_NE, nil}, {TT_LT, nil}, {TT_GT, nil}, {TT_LTE, nil}, {TT_GTE, nil}}))
	log.Println("COMP:", node)
	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected int, float, identifier, '+', '-', '(' or 'NOT'").Error)
	}

	return res.Success(node)
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
	}

	return p.Power()
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []TokenTypeInfo{{TT_MUL, nil}, {TT_DIV, nil}})
}

// Expr parses an expression.
func (p *Parser) Expr() *ParseResult {
	res := ParseResult{AdvanceCount: 0}

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

	and := "AND"
	or := "OR"
	node := res.Register(p.BinOp(p.CompExpr, []TokenTypeInfo{{TT_KEYWORD, &and}, {TT_KEYWORD, &or}}))
	log.Println("NODE in expr:", node)
	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'VAR', int, float, identifier, '+', '-', '(' or 'NOT'").Error)
	}

	return res.Success(node)
}

// / BinOp parses a binary operation.
func (p *Parser) BinOp(funcParser func() *ParseResult, ops []TokenTypeInfo) *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	left := res.Register(funcParser())
	if res.Error != nil {
		return res
	}

	for ContainsTypeOrValue(ops, p.Current.Type, p.Current.Value) {
		//log.Println("ACCEPTED", ContainsTypeOrValue(ops, p.Current.Type, p.Current.Value), p.Current.Type, p.Current.Value)
		opTok := p.Current
		res.RegisterAdvancement()
		p.Advance()

		// Parse the right side of the expression
		right := res.Register(funcParser())

		if res.Error != nil {
			return res
		}
		left = NewBinOpNode(left, opTok, right)

		//log.Println("LEFT TOTAL:", left.(*BinOpNode).LeftNode, left.(*BinOpNode).OpTok, left.(*BinOpNode).RightNode)
	}
	//log.Println("ACCEPTED 1", ContainsTypeOrValue(ops, p.Current.Type, p.Current.Value), p.Current.Type, p.Current.Value)

	return res.Success(left)
}

// ContainsTypeOrValue checks if the token type and value exist in the list of type-value pairs.
func ContainsTypeOrValue(types []TokenTypeInfo, typ TokenTypes, val interface{}) bool {
	for _, t := range types {
		if val != nil {
			if t.Value != nil {
				if *t.Value == val.(string) && t.Type == typ {
					return true
				}
			}
		} else if t.Type == typ {
			return true
		}
	}
	return false
}
