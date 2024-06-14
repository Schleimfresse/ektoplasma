package main

import (
	"fmt"
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
func (p *Parser) ArrayExpr() *ParseResult {
	res := NewParseResult()
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
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected ']', 'var', 'if', 'for', 'while', 'func', int, float, identifier, '+', '-', '(', '[' or 'NOT'").Error)
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
	res := NewParseResult()
	allCases := res.Register(p.ifExprCases("if"))
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
	return p.ifExprCases("elif")
}

// ifExprC is a method of Parser that handles 'ELSE' in if expressions.
func (p *Parser) ifExprC() *ParseResult {
	res := NewParseResult()

	var elseCase *ElseCaseNode

	if p.Current.Matches(TT_KEYWORD, "else") {
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

			if p.Current.Type == TT_RBRACE {
				res.RegisterAdvancement()
				p.Advance()
			} else {
				return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '}'").Error)
			}
		} else {
			expr := res.Register(p.Statement())
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
	res := NewParseResult()
	var cases []*IfCaseNode
	var elseCase *ElseCaseNode
	if p.Current.Matches(TT_KEYWORD, "elif") {
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
	res := NewParseResult()
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

	if p.Current.Type != TT_LBRACE {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '{'").Error)
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
		if p.Current.Type == TT_RBRACE {
			res.RegisterAdvancement()
			p.Advance()
		} else {
			allCases := res.Register(p.ifExprBOrC())
			if res.Error != nil {
				return res
			}
			ifNode := allCases.(*IfNode)

			cases = append(cases, ifNode.Cases...)
			elseCase = ifNode.ElseCase
		}
	} else {
		expr := res.Register(p.Statement())
		if res.Error != nil {
			return res
		}

		cases = append(cases, NewIfCaseNode(condition, expr, false))
		if p.Current.Type != TT_NEWLINE {
			allCases := res.Register(p.ifExprBOrC())
			if res.Error != nil {
				return res
			}
			ifNode := allCases.(*IfNode)

			cases = append(cases, ifNode.Cases...)
			elseCase = ifNode.ElseCase
		}
	}

	return res.Success(NewIfNode(cases, elseCase))
}

func (p *Parser) ForExpr() *ParseResult {
	res := NewParseResult()

	if !p.Current.Matches(TT_KEYWORD, "for") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'for'",
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

	if !p.Current.Matches(TT_KEYWORD, "to") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'to'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	endValue := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	var stepValue Node
	if p.Current.Matches(TT_KEYWORD, "step") {
		res.RegisterAdvancement()
		p.Advance()

		stepValue = res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
	}

	if p.Current.Type != TT_LBRACE {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected '{'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_NEWLINE {
		body := res.Register(p.Statements())
		if res.Error != nil {
			return res
		}
		if res.Error != nil {
			return res
		}
		if p.Current.Type != TT_RBRACE {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '}'").Error)
		}

		res.RegisterAdvancement()
		p.Advance()

		return res.Success(NewForNode(varName, startValue, endValue, stepValue, body, true))
	}

	body := res.Register(p.Statement())
	if res.Error != nil {
		return res
	}

	return res.Success(NewForNode(varName, startValue, endValue, stepValue, body, false))
}

func (p *Parser) WhileExpr() *ParseResult {
	res := NewParseResult()

	if !p.Current.Matches(TT_KEYWORD, "while") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'while'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	condition := res.Register(p.Expr())
	if res.Error != nil {
		return res
	}

	if p.Current.Type != TT_LBRACE {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected '{'",
		).Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_NEWLINE {
		res.RegisterAdvancement()
		p.Advance()

		body := res.Register(p.Statements())
		if res.Error != nil {
			return res
		}

		if p.Current.Type != TT_RBRACE {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '}'").Error)
		}

		res.RegisterAdvancement()
		p.Advance()

		return res.Success(NewWhileNode(condition, body, true))
	}

	body := res.Register(p.Statement())
	if res.Error != nil {
		return res
	}

	return res.Success(NewWhileNode(condition, body, false))
}

func (p *Parser) Call() *ParseResult {
	res := NewParseResult()
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
				return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected ')', 'var', 'if', 'for', 'while', 'func', int, float, identifier, pointer, '+', '-', '(', '[' or 'not'").Error)
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
	res := NewParseResult()
	var VarNameToken *Token

	if !p.Current.Matches(TT_KEYWORD, "func") {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected 'func'",
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

	if p.Current.Type != TT_LBRACE {
		return res.Failure(NewInvalidSyntaxError(
			p.Current.PosStart, p.Current.PosEnd,
			"Expected '{'",
		).Error)
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

		return res.Success(NewFuncDefNode(VarNameToken, ArgNameTokens, body, true))
	}

	if p.Current.Type != TT_NEWLINE {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '=>' or new line").Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	body := res.Register(p.Statements())
	if res.Error != nil {
		return res
	}

	if p.Current.Type != TT_RBRACE {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected '}'").Error)
	}

	res.RegisterAdvancement()
	p.Advance()

	return res.Success(NewFuncDefNode(VarNameToken, ArgNameTokens, body, false))
}

func (p *Parser) Atom() *ParseResult {
	res := NewParseResult()
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

		// in case of an index from an array
		if p.Current.Type == TT_LSQUARE {
			res.RegisterAdvancement()
			p.Advance()

			if p.Current.Type == TT_INT {
				index := p.Current
				indexValue, err := strconv.ParseInt(p.Current.Value.(string), 10, 64)
				index.Value = int(indexValue)

				if err != nil {
					fmt.Println(err)
				}

				res.RegisterAdvancement()
				p.Advance()

				if p.Current.Type == TT_RSQUARE {
					res.RegisterAdvancement()
					p.Advance()
					return res.Success(NewIndexNode(NewVarAccessNode(tok), NewNumberNode(index)))
				}
			} else {
				return res.Failure(NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, fmt.Sprintf("Expected int, got: %v", tok.Value)).Error)
			}
		} else {
			return res.Success(NewVarAccessNode(tok))
		}
	} else if tok.Type == TT_LSQUARE {
		listExpr := res.Register(p.ArrayExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(listExpr)
	} else if tok.Type == TT_LPAREN {
		p.Advance()
		expr := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		if p.Current.Type == TT_RPAREN {
			p.Advance()
			return res.Success(expr)
		}
		err := NewInvalidSyntaxError(tok.PosStart, tok.PosEnd, "Expected ')'")
		return res.Failure(err.Error)
	} else if tok.Matches(TT_KEYWORD, "if") {
		IfExpr := res.Register(p.ifExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(IfExpr)
	} else if tok.Matches(TT_KEYWORD, "for") {
		IfExpr := res.Register(p.ForExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(IfExpr)
	} else if tok.Matches(TT_KEYWORD, "while") {
		WhileExpr := res.Register(p.WhileExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(WhileExpr)
	} else if tok.Matches(TT_KEYWORD, "func") {
		FuncDef := res.Register(p.FuncDef())
		if res.Error != nil {
			return res
		}
		return res.Success(FuncDef)
	} else if tok.Matches(TT_KEYWORD, "import") {
		ImportExpr := res.Register(p.ImportExpr())
		if res.Error != nil {
			return res
		}
		return res.Success(ImportExpr)
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

	if p.Current.Matches(TT_KEYWORD, "not") {
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
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected int, float, identifier, '+', '-', '(', '[' or 'not'").Error)
	}

	return res.Success(node)
}

// Factor parses a factor.
func (p *Parser) Factor() *ParseResult {
	res := NewParseResult()
	tok := p.Current
	if tok.Type == TT_PLUS || tok.Type == TT_MINUS {
		p.Advance()
		factor := res.Register(p.Factor())
		return res.Success(NewUnaryOpNode(tok, factor))
	}

	return p.Power()
}

func (p *Parser) Term() *ParseResult {
	return p.BinOp(p.Factor, []TokenTypeInfo{{TT_STAR, nil}, {TT_DIV, nil}}, p.Factor)
}

func (p *Parser) Statements() *ParseResult {
	res := ParseResult{AdvanceCount: 0}

	var statements []Node
	PosStart := p.Current.PosStart.Copy()

	for p.Current.Type == TT_NEWLINE {
		res.RegisterAdvancement()
		p.Advance()
	}

	statement := res.Register(p.Statement())
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

		statement = res.TryRegister(p.Statement())
		if statement == nil {
			p.Reverse(&res.ToReverseCount)
			MoreStatements = false
			continue
		}
		statements = append(statements, statement)
	}

	return res.Success(NewArrayNode(statements, PosStart, p.Current.PosEnd))
}

func (p *Parser) Statement() *ParseResult {
	res := NewParseResult()
	PosStart := p.Current.PosStart.Copy()

	if p.Current.Matches(TT_KEYWORD, "return") {
		res.RegisterAdvancement()
		p.Advance()

		expr := res.TryRegister(p.Expr())
		if expr == nil {
			p.Reverse(&res.ToReverseCount)
		}
		return res.Success(NewReturnNode(expr, PosStart, p.Current.PosEnd.Copy()))
	}
	if p.Current.Matches(TT_KEYWORD, "continue") {
		res.RegisterAdvancement()
		p.Advance()
		return res.Success(NewContinueNode(PosStart, p.Current.PosEnd.Copy()))
	}
	if p.Current.Matches(TT_KEYWORD, "break") {
		res.RegisterAdvancement()
		p.Advance()
		return res.Success(NewBreakNode(PosStart, p.Current.PosEnd.Copy()))
	}

	expr := res.Register(p.Expr())
	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'RETURN', 'CONTINUE', 'BREAK', 'VAR', 'IF', 'FOR', 'WHILE', 'FUN', int, float, identifier, '+', '-', '(', '[' or 'NOT'").Error)
	}
	return res.Success(expr)

}

// Expr parses an expression.
func (p *Parser) Expr() *ParseResult {
	res := NewParseResult()

	if p.Current.Matches(TT_KEYWORD, "var") || p.Current.Matches(TT_KEYWORD, "const") {
		isConst := false

		if p.Current.Matches(TT_KEYWORD, "const") {
			isConst = true
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
			if isConst {
				return res.Failure(NewInvalidSyntaxError(
					p.Current.PosStart, p.Current.PosEnd,
					"Missing assignment in const declaration",
				).Error)
			} else if p.Current.Type == TT_NEWLINE {
				return res.Success(NewVarAssignNode(varName, nil, false, true))
			} else {
				return res.Failure(NewInvalidSyntaxError(
					p.Current.PosStart, p.Current.PosEnd,
					"Expected assignment or new line",
				).Error)
			}
		}

		res.RegisterAdvancement()
		p.Advance()

		if p.Current.Type == TT_AND {
			res.RegisterAdvancement()
			p.Advance()
			expr := NewReference(res.Register(p.Expr()))
			if res.Error != nil {
				return res
			}

			return res.Success(NewVarAssignNode(varName, expr, isConst, true))
		} else if p.Current.Type == TT_STAR {
			res.RegisterAdvancement()
			p.Advance()
			expr := NewDereference(res.Register(p.Expr()))
			if res.Error != nil {
				return res
			}

			return res.Success(NewVarAssignNode(varName, expr, isConst, true))
		}

		expr := res.Register(p.Expr())
		if res.Error != nil {
			return res
		}
		return res.Success(NewVarAssignNode(varName, expr, isConst, true))
	} else if p.Current.Type == TT_IDENTIFIER { // in case of a variable re-assignment, so we don't need the var keyword for each assignment, only for the initial
		varName := p.Current
		res.RegisterAdvancement()
		p.Advance()

		if p.Current.Type == TT_EQ {
			res.RegisterAdvancement()
			p.Advance()

			expr := res.Register(p.Expr())
			if res.Error != nil {
				return res
			}

			return res.Success(NewVarAssignNode(varName, expr, false, false))
		} else if p.Current.Type == TT_DOT {
			res.RegisterAdvancement()
			p.Advance()
			callNode := p.Call()
			if _, ok := callNode.Node.(*CallNode); ok {
				return res.Success(NewPackageMethod(varName, callNode.Node.(*CallNode).NodeToCall.(*VarAccessNode).VarNameTok.Value.(string), callNode.Node))
			} else if _, ok := callNode.Node.(*VarAccessNode); ok {
				return res.Success(NewPackageMethod(varName, callNode.Node.(*VarAccessNode).VarNameTok.Value.(string), callNode.Node))
			}
		} else {
			p.Reverse(&res.ToReverseCount)
		}
	} else if p.Current.Type == TT_AND {
		res.RegisterAdvancement()
		p.Advance()
		ref := NewReference(res.Register(p.Expr()))
		if res.Error != nil {
			return res
		}

		return res.Success(ref)
	} else if p.Current.Type == TT_STAR {
		res.RegisterAdvancement()
		p.Advance()
		deref := NewDereference(res.Register(p.Expr()))
		if res.Error != nil {
			return res
		}

		return res.Success(deref)
	}

	and := "AND"
	or := "OR"
	node := res.Register(p.BinOp(p.CompExpr, []TokenTypeInfo{{TT_KEYWORD, &and}, {TT_KEYWORD, &or}}, p.CompExpr))
	if res.Error != nil {
		return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'var', 'if', 'for', 'while', 'func', int, float, identifier, '+', '-', '(', '[' or 'not'").Error)
	}

	return res.Success(node)
}

func (p *Parser) ImportExpr() *ParseResult {
	res := NewParseResult()
	var functionNames []*Token
	var packageNames []*Token
	posStart := p.Current.PosStart

	res.RegisterAdvancement()
	p.Advance()

	if p.Current.Type == TT_STRING {
		var packageName = p.Current
		res.RegisterAdvancement()
		p.Advance()
		if p.Current.Type == TT_NEWLINE {
			return res.Success(NewImportNode(nil, []*Token{packageName}, posStart, p.Current.PosEnd.Copy()))
		}
	} else if p.Current.Type == TT_IDENTIFIER {
		functionNames = append(functionNames, p.Current)
		res.RegisterAdvancement()
		p.Advance()

		for p.Current.Type == TT_COMMA {
			res.RegisterAdvancement()
			p.Advance()
			functionNames = append(functionNames, p.Current)
			res.RegisterAdvancement()
			p.Advance()
		}
		if !p.Current.Matches(TT_KEYWORD, "from") {
			return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, "Expected 'from'").Error)
		}
		res.RegisterAdvancement()
		p.Advance()

		packageNames = append(packageNames, p.Current)

		posEnd := p.Current.PosEnd.Copy()

		res.RegisterAdvancement()
		p.Advance()

		return res.Success(NewImportNode(functionNames, packageNames, posStart, posEnd))
	}
	return res.Failure(NewInvalidSyntaxError(p.Current.PosStart, p.Current.PosEnd, fmt.Sprintf("Can not import %s, type %s", p.Current.Value, p.Current.Type)).Error)
}

// / BinOp parses a binary operation.
func (p *Parser) BinOp(funcAParser func() *ParseResult, ops []TokenTypeInfo, funcBParser func() *ParseResult) *ParseResult {
	res := NewParseResult()

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

func NewParseResult() *ParseResult {
	return &ParseResult{AdvanceCount: 0}
}
