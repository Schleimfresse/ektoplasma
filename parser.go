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

func (pr *ParseResult) TryRegister(res *ParseResult) Node {
	if res.Error != nil {
		pr.ToReverseCount = res.AdvanceCount
		return nil
	}
	return pr.Register(res)
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
	p.UpdateCurrentTok()
	return p.Current
}

func (p *Parser) Reverse(amount *int) *Token {
	if amount == nil {
		defaultAmount := 1
		amount = &defaultAmount
	}
	p.TokIdx -= 1
	p.UpdateCurrentTok()
	return p.Current
}

func (p *Parser) UpdateCurrentTok() {
	if p.TokIdx >= 0 && p.TokIdx < len(p.Tokens) {
		p.Current = p.Tokens[p.TokIdx]
	}
}

// Parse parses the tokens into an abstract syntax tree.
func (p *Parser) Parse() *ParseResult {
	res := p.Statements()
	if res.Error != nil && p.Current.Type != TT_EOF {
		//log.Println("IMP ERR:", res.Error)
		//err := NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '+', '-', '*', '/', '^', '==', '!=', '<', '>', <=', '>=', 'AND' or 'OR'")
		//n.return res.Failure(err.Error)
	}
	return res
}

// listExpr method for Interpreter
func (p *Parser) ListExpr() *ParseResult {
	res := &ParseResult{}
	elementNodes := []Node{}
	posStart := p.Current.PosStart.Copy()

	if p.Current.Type != TT_LSQUARE {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '['").Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_RSQUARE {
		res.RegisterAdvancement()
		p.Advance()
	} else {
		elementNodes = append(elementNodes, res.Register(p.Expr()))
		if res.Error != nil {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected ']', 'VAR', 'IF', 'FOR', 'WHILE', 'FUN', int, float, identifier, '+', '-', '(', '[' or 'NOT'").Error)
		}

		for p.Current.Type == TT_COMMA {
			res.RegisterAdvancement()
			p.Advance()

			elementNodes = append(elementNodes, res.Register(p.Expr()))
			if res.Error != nil {
				return res
			}
		}

		if p.Current.Type != TT_RSQUARE {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected ',' or ']'").Error)
		}

		res.RegisterAdvancement()
		p.Advance()
	}

	return res.Success(NewArrayNode(elementNodes, posStart, p.Current.PosEnd.Copy()))
}

// ifExpr is a method of Parser that handles 'IF' expressions.
func (p *Parser) ifExpr() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	allCases := res.Register(p.ifExprCases("IF"))
	if res.Error != nil {
		return res
	}
	cases, ok := allCases.(*IfNode)
	if !ok {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "unexpected type or insufficient length for allCases").Error)
	}
	caseNodes := cases.Cases
	elseCaseNode := cases.ElseCase
	parsedCases := make([]*IfCaseNode, len(caseNodes))
	for i, c := range caseNodes {
		parsedCases[i] = c
	}
	return res.Success(NewIfNode(parsedCases, elseCaseNode))
}

// ifExprB is a method of Parser that handles 'ELIF' in if expressions.
func (p *Parser) ifExprB() *ParseResult {
	return p.ifExprCases("ELIF")
}

// ifExprC is a method of Parser that handles 'ELSE' in if expressions.
func (p *Parser) ifExprC() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}

	var elseCase *ElseCaseNode

	if p.Current.Matches(TT_KEYWORD, "ELSE") {
		res.RegisterAdvancement()
		p.Advance()

		if p.Current.Type == TT_NEWLINE {
			res.RegisterAdvancement()
			p.Advance()

			statements := res.Register(p.Statements())
			if res.Error != nil {
				return res
			}
			elseCase = NewElseCaseNode(statements, true)

			if p.Current.Matches(TT_KEYWORD, "END") {
				res.RegisterAdvancement()
				p.Advance()
			} else {
				return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'END'").Error)
			}
		} else {
			expr := res.Register(p.Expr())
			if res.Error != nil {
				return res
			}
			elseCase = NewElseCaseNode(expr, false)
		}
	}

	return res.Success(elseCase)
}

// ifExprBOrC is a method of Parser that handles 'ELIF' or 'ELSE' in if expressions.
func (p *Parser) ifExprBOrC() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	var cases []*IfCaseNode
	var elseCase *ElseCaseNode

	if p.Current.Matches(TT_KEYWORD, "ELIF") {
		allCases := res.Register(p.ifExprB())
		if res.Error != nil {
			return res
		}
		ifNode, ok := allCases.(*IfNode)
		if !ok {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "unexpected type for allCases").Error)
		}
		cases, elseCase = ifNode.Cases, ifNode.ElseCase
	} else {
		elseCaseResult := res.Register(p.ifExprC())
		if res.Error != nil {
			return res
		}
		elseCase = elseCaseResult.(*ElseCaseNode)
	}

	return res.Success(NewIfNode(cases, elseCase))
}

// ifExprCases is a method of Parser that handles cases in if expressions.
func (p *Parser) ifExprCases(caseKeyword string) *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	cases := make([]*IfCaseNode, 0)
	var elseCase *ElseCaseNode

	if !p.Current.Matches(TT_KEYWORD, caseKeyword) {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, fmt.Sprintf("Expected '%s'", caseKeyword)).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	condition := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	if !p.Current.Matches(TT_KEYWORD, "THEN") {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'THEN'").Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_NEWLINE {
		res.RegisterAdvancement()
		p.Advance()

		statements := res.Register(p.Statements())
		if res.Error != nil {
			return res
		}
		cases = append(cases, NewIfCaseNode(condition, statements, true))

		if p.Current.Matches(TT_KEYWORD, "END") {
			res.RegisterAdvancement()
			p.Advance()
		} else {
			allCases := res.Register(p.ifExprBOrC())
			if res.Error != nil {
				return res
			}
			ifNode, ok := allCases.(*IfNode)
			if !ok {
				return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "unexpected type for allCases").Error)
			}
			cases = append(cases, ifNode.Cases...)
			elseCase = ifNode.ElseCase
		}
	} else {
		expr := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		cases = append(cases, NewIfCaseNode(condition, expr, false))

		allCases := res.Register(p.ifExprBOrC())
		if res.Error != nil {
			return res
		}
		ifNode, ok := allCases.(*IfNode)
		if !ok {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "unexpected type for allCases").Error)
		}
		cases = append(cases, ifNode.Cases...)
		elseCase = ifNode.ElseCase
	}

	return res.Success(NewIfNode(cases, elseCase))
}

func (p *Parser) ForExpr() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}

	if !p.Current.Matches(TT_KEYWORD, "FOR") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'FOR'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

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

	startValue := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	if !p.Current.Matches(TT_KEYWORD, "TO") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'TO'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	endValue := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	var stepValue Node
	if p.Current.Matches(TT_KEYWORD, "STEP") {
		res.RegisterAdvancement()
		p.Advance()

		stepValue = res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
	}

	if !p.Current.Matches(TT_KEYWORD, "THEN") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'THEN'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_NEWLINE {
		body := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		if res.Error != nil {
			return res
		}
		if !p.Current.Matches(TT_KEYWORD, "END") {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'END'").Error)
		}

		res.RegisterAdvancement()
		p.Advance()

		return res.Success(NewForNode(varName, startValue, endValue, stepValue, body, true))
	}

	body := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	return res.Success(NewForNode(varName, startValue, endValue, stepValue, body, false))
}

func (p *Parser) WhileExpr() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}

	if !p.Current.Matches(TT_KEYWORD, "WHILE") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'WHILE'",
		).Error)
	}

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

	if p.Current.Type == TT_NEWLINE {
		res.RegisterAdvancement()
		p.Advance()

		body := res.Register(p.Statements())
		if res.Error != nil {
			return res
		}

		if !p.Current.Matches(TT_KEYWORD, "END") {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'END'").Error)
		}

		res.RegisterAdvancement()
		p.Advance()

		return res.Success(NewWhileNode(condition, body, true))
	}

	body := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	return res.Success(NewWhileNode(condition, body, false))
}

func (p *Parser) Call() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	atom := res.Register(p.Atom())
	if res.Error != nil {
		return res
	}

	if p.Current.Type == TT_LPAREN {
		res.RegisterAdvancement()
		p.Advance()
		var ArgNodes []Node

		if p.Current.Type == TT_RPAREN {
			res.RegisterAdvancement()
			p.Advance()
		} else {
			ArgNodes = append(ArgNodes, res.Register(p.Expr()))
			if res.Error != nil {
				return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected ')', 'VAR', 'IF', 'FOR', 'WHILE', 'FUN', int, float, identifier, '+', '-', '(', '[' or 'NOT'").Error)
			}

			for p.Current.Type == TT_COMMA {
				res.RegisterAdvancement()
				p.Advance()

				ArgNodes = append(ArgNodes, res.Register(p.Expr()))
				if res.Error != nil {
					return res
				}
			}

			if p.Current.Type != TT_RPAREN {
				return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected ',' or ')'").Error)
			}

			res.RegisterAdvancement()
			p.Advance()
		}

		return res.Success(NewCallNode(atom, ArgNodes))
	}
	return res.Success(atom)
}

func (p *Parser) FuncDef() *ParseResult {
	res := &ParseResult{AdvanceCount: 0}
	var VarNameToken *Token

	if !p.Current.Matches(TT_KEYWORD, "FUNC") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'FUNC'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_IDENTIFIER {
		VarNameToken = p.Current
		res.RegisterAdvancement()
		p.Advance()
		if p.Current.Type != TT_LPAREN {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected '('",
			).Error)
		}
	} else {
		VarNameToken = nil
		if p.Current.Type != TT_LPAREN {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected identifier or '('",
			).Error)
		}

	}
	res.RegisterAdvancement()
	p.Advance()
	var ArgNameTokens []*Token

	if p.Current.Type == TT_IDENTIFIER {
		ArgNameTokens = append(ArgNameTokens, p.Current)
		res.RegisterAdvancement()
		p.Advance()

		for p.Current.Type == TT_COMMA {
			res.RegisterAdvancement()
			p.Advance()

			if p.Current.Type != TT_IDENTIFIER {
				return res.Failure(NewInvalidSyntaxError(
					p.Current.PosStart, p.Current.PosEnd,
					"Expected identifier",
				).Error)
			}

			ArgNameTokens = append(ArgNameTokens, p.Current)
			res.RegisterAdvancement()
			p.Advance()
		}

		if p.Current.Type != TT_RPAREN {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected ',' or ')'",
			).Error)
		}
	} else {
		if p.Current.Type != TT_RPAREN {
			return res.Failure(NewInvalidSyntaxError(
				p.Current.PosStart, p.Current.PosEnd,
				"Expected ',' or ')'",
			).Error)
		}
	}
	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_ARROW {
		res.RegisterAdvancement()
		p.Advance()

		body := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}

		return res.Success(NewFuncDefNode(VarNameToken, ArgNameTokens, body, false))
	}

	if p.Current.Type != TT_NEWLINE {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '=>' or NEWLINE").Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	body := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	if !p.Current.Matches(TT_KEYWORD, "END") {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'END'").Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	return res.Success(NewFuncDefNode(VarNameToken, ArgNameTokens, body, true))
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
	} else if tok.Type == TT_STRING {
		res.RegisterAdvancement()
		p.Advance()
		return res.Success(NewStringNode(tok))
	} else if tok.Type == TT_IDENTIFIER {
		res.RegisterAdvancement()
		p.Advance()
		return res.Success(NewVarAccessNode(tok))
	} else if tok.Type == TT_LSQUARE {
		listExpr := res.Register(p.ListExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(listExpr)
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
		IfExpr := res.Register(p.ifExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(IfExpr)
	} else if tok.Matches(TT_KEYWORD, "FOR") {
		IfExpr := res.Register(p.ForExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(IfExpr)
	} else if tok.Matches(TT_KEYWORD, "WHILE") {
		WhileExpr := res.Register(p.WhileExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(WhileExpr)
	} else if tok.Matches(TT_KEYWORD, "FUNC") {
		FuncDef := res.Register(p.FuncDef())
		if res.Error != nil {
			return res
		}
		return res.Success(FuncDef)
	}

	return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected int, float, identifier, '+', '-' or '(', 'IF', 'FOR', 'WHILE', 'FUNC'").Error)
}

func (p *Parser) Power() *ParseResult {
	return p.BinOp(p.Call, []TokenTypeInfo{{TT_POW, nil}}, p.Factor)
}

// ArithExpr parses an arithmetic expression.
func (p *Parser) ArithExpr() *ParseResult {
	return p.BinOp(p.Term, []TokenTypeInfo{{TT_PLUS, nil}, {TT_MINUS, nil}}, p.Term)
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

	node := res.Register(p.BinOp(p.ArithExpr, []TokenTypeInfo{{TT_EE, nil}, {TT_NE, nil}, {TT_LT, nil}, {TT_GT, nil}, {TT_LTE, nil}, {TT_GTE, nil}}, p.ArithExpr))

	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected int, float, identifier, '+', '-', '(', '[' or 'NOT'").Error)
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
	return p.BinOp(p.Factor, []TokenTypeInfo{{TT_MUL, nil}, {TT_DIV, nil}}, p.Factor)
}

func (p *Parser) Statements() *ParseResult {
	res := ParseResult{AdvanceCount: 0}

	// Check if identifier is unexpected
	if p.Current.Type == TT_IDENTIFIER && !GlobalSymbolTable.Contains(p.Current.Value.(string)) {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, fmt.Sprintf("Unresolved reference '%s'", p.Current.Value.(string))).Error)
	}

	var statements []Node
	PosStart := p.Current.PosStart.Copy()

	for p.Current.Type == TT_NEWLINE {
		res.RegisterAdvancement()
		p.Advance()
	}

	statement := res.Register(p.Expr())
	if res.Error != nil {
		return &res
	}
	statements = append(statements, statement)

	MoreStatements := true

	for {
		NewlineCount := 0
		for p.Current.Type == TT_NEWLINE {
			res.RegisterAdvancement()
			p.Advance()
			NewlineCount++
		}
		if NewlineCount == 0 {
			MoreStatements = false
		}
		if !MoreStatements {
			break
		}
		statement = res.TryRegister(p.Expr())
		if statement == nil {
			p.Reverse(&res.ToReverseCount)
			MoreStatements = false
			continue
		}
		statements = append(statements, statement)
	}

	return res.Success(NewArrayNode(statements, PosStart, p.Current.PosEnd))
}

// Expr parses an expression.
func (p *Parser) Expr() *ParseResult {
	res := ParseResult{AdvanceCount: 0}

	if p.Current.Matches(TT_KEYWORD, "VAR") {
		res.RegisterAdvancement()
		p.Advance()

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
	node := res.Register(p.BinOp(p.CompExpr, []TokenTypeInfo{{TT_KEYWORD, &and}, {TT_KEYWORD, &or}}, nil))

	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'VAR', 'IF', 'FOR', 'WHILE', 'FUN', int, float, identifier, '+', '-', '(', '[' or 'NOT'").Error)
	}

	return res.Success(node)
}

// / BinOp parses a binary operation.
func (p *Parser) BinOp(funcAParser func() *ParseResult, ops []TokenTypeInfo, funcBParser func() *ParseResult) *ParseResult {
	res := &ParseResult{AdvanceCount: 0}

	left := res.Register(funcAParser())
	if res.Error != nil {
		return res
	}

	for ContainsTypeOrValue(ops, p.Current.Type, p.Current.Value) {
		opTok := p.Current
		res.RegisterAdvancement()
		p.Advance()

		// Parse the right side of the expression
		right := res.Register(funcBParser())

		if res.Error != nil {
			return res
		}
		left = NewBinOpNode(left, opTok, right)
	}

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
