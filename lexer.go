package main

import "log"

// Matches checks if the token matches the given type and value
func (t Token) Matches(tType TokenTypes, value string) bool {
	return t.Type == tType && t.Value == value
}

// NewPosition creates a new Position instance.
func NewPosition(idx, ln, col int, fn, ftxt string) *Position {
	return &Position{idx, ln, col, fn, ftxt}
}

// Advance moves the position forward based on the current character.
func (p *Position) Advance(currentChar byte) *Position {
	p.Idx++
	p.Col++

	if currentChar == '\n' {
		p.Ln++
		p.Col = 0
	}

	return p
}

// Copy creates a copy of the current position.
func (p *Position) Copy() *Position {
	return &Position{p.Idx, p.Ln, p.Col, p.Fn, p.Ftxt}
}

// NewToken creates a new Token instance.
func NewToken(typ TokenTypes, value interface{}, posStart, posEnd *Position) *Token {
	return &Token{
		Type:     typ,
		Value:    value,
		PosStart: posStart.Copy(),
		PosEnd:   posEnd.Copy(),
	}
}

// NewLexer creates a new Lexer instance.
func NewLexer(fn, text string) *Lexer {
	lexer := &Lexer{
		Fn:   fn,
		Text: text,
		Pos:  NewPosition(-1, 0, -1, fn, text),
	}
	lexer.Advance()
	return lexer
}

// Advance moves the lexer forward.
func (l *Lexer) Advance() {
	l.Pos.Advance(l.CurrentChar)
	if l.Pos.Idx < len(l.Text) {
		l.CurrentChar = l.Text[l.Pos.Idx]
	} else {
		l.CurrentChar = 0 // Null character
	}
}

// MakeTokens tokenizes the input text.
func (l *Lexer) MakeTokens() ([]*Token, *IllegalCharError) {
	tokens := []*Token{}

	for l.CurrentChar != 0 {
		if l.CurrentChar == ' ' || l.CurrentChar == '\t' {
			l.Advance()
		} else if isDigit(l.CurrentChar) {
			tokens = append(tokens, l.MakeNumber())
		} else if isLetter(l.CurrentChar) {
			tokens = append(tokens, l.MakeIdentifier())
		} else if l.CurrentChar == '+' {
			tokens = append(tokens, NewToken(TT_PLUS, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == '-' {
			tokens = append(tokens, NewToken(TT_MINUS, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == '*' {
			tokens = append(tokens, NewToken(TT_MUL, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == '/' {
			tokens = append(tokens, NewToken(TT_DIV, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == '(' {
			tokens = append(tokens, NewToken(TT_LPAREN, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == ')' {
			tokens = append(tokens, NewToken(TT_RPAREN, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == '^' {
			tokens = append(tokens, NewToken(TT_POW, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else if l.CurrentChar == '=' {
			tokens = append(tokens, NewToken(TT_EQ, nil, l.Pos.Copy(), l.Pos.Copy()))
			l.Advance()
		} else {
			posStart := l.Pos.Copy()
			char := string(l.CurrentChar)
			l.Advance()
			return []*Token{}, NewIllegalCharError(posStart, l.Pos, "'"+char+"'")
		}
	}
	tokens = append(tokens, NewToken(TT_EOF, nil, l.Pos.Copy(), l.Pos.Copy()))
	return tokens, nil
}

func (l *Lexer) MakeIdentifier() *Token {
	idStr := ""
	posStart := l.Pos.Copy()

	for l.CurrentChar != 0 && (isLetter(l.CurrentChar) || isDigit(l.CurrentChar) || l.CurrentChar == '_') {
		idStr += string(l.CurrentChar)
		l.Advance()
	}

	posEnd := l.Pos.Copy().Advance(l.CurrentChar)
	posEnd.Col = posEnd.Col - 1
	posEnd.Idx = posEnd.Idx - 1
	log.Println("LEXER:", idStr)
	tokenType := TT_IDENTIFIER
	if isKeyword(idStr) {
		tokenType = TT_KEYWORD
	}

	return NewToken(tokenType, idStr, posStart, posEnd)
}

// MakeNumber tokenizes a number.
func (l *Lexer) MakeNumber() *Token {
	numStr := ""
	dotCount := 0
	posStart := l.Pos.Copy()

	for l.CurrentChar != 0 && (isDigit(l.CurrentChar) || l.CurrentChar == '.') {
		if l.CurrentChar == '.' {
			if dotCount == 1 {
				break
			}
			dotCount++
			numStr += "."
		} else {
			numStr += string(l.CurrentChar)
		}
		l.Advance()
	}

	posEnd := l.Pos.Copy()
	posEnd.Col = posEnd.Col - 1
	posEnd.Idx = posEnd.Idx - 1

	if dotCount == 0 {
		return NewToken(TT_INT, numStr, posStart, posEnd)
	}
	return NewToken(TT_FLOAT, numStr, posStart, posEnd)
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isKeyword(str string) bool {
	for _, kw := range KEYWORDS {
		if kw == str {
			return true
		}
	}
	return false
}
