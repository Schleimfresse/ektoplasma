package main

import (
	"bufio"
	"fmt"
	"os"
)

// ecp Ektoplasma (Ektoplasma Code Program)

const (
	TT_INT        TokenTypes = "INT"
	TT_FLOAT      TokenTypes = "FLOAT"
	TT_STRING     TokenTypes = "STRING"
	TT_IDENTIFIER TokenTypes = "IDENTIFIER"
	TT_KEYWORD    TokenTypes = "KEYWORD"
	TT_PLUS       TokenTypes = "PLUS"
	TT_MINUS      TokenTypes = "MINUS"
	TT_MUL        TokenTypes = "MUL"
	TT_DIV        TokenTypes = "DIV"
	TT_EQ         TokenTypes = "EQ"
	TT_LPAREN     TokenTypes = "LPAREN"
	TT_RPAREN     TokenTypes = "RPAREN"
	TT_LSQUARE    TokenTypes = "LSQUARE"
	TT_RSQUARE    TokenTypes = "RSQUARE"
	TT_POW        TokenTypes = "POW"
	TT_EE         TokenTypes = "EE"
	TT_NE         TokenTypes = "NE"
	TT_LT         TokenTypes = "LT"
	TT_GT         TokenTypes = "GT"
	TT_LTE        TokenTypes = "LTE"
	TT_GTE        TokenTypes = "GTE"
	TT_EOF        TokenTypes = "EOF"
	TT_COMMA      TokenTypes = "COMMA"
	TT_ARROW      TokenTypes = "ARROW"
)

var KEYWORDS = []string{"VAR", "AND", "OR", "NOT", "IF", "THEN", "ELSE", "ELIF", "FOR", "TO", "STEP", "WHILE", "FUNC"}
var GlobalSymbolTable = NewSymbolTable(nil)
var lineTEMP int

func run(fileName, text string) (*Value, *RuntimeError) {
	lexer := NewLexer(fileName, text)

	tokens, err := lexer.MakeTokens()
	if err != nil {
		fmt.Println(err.AsString())
		return nil, nil
	}

	for _, token := range tokens {
		if token.PosStart != nil && token.PosEnd != nil {
			fmt.Println(token.Type, token.Value, "START:", *token.PosStart, "END:", *token.PosEnd)
		} else {
			fmt.Println(token.Type, "START:", token.PosStart, "END:", token.PosEnd)
		}
	}
	parser := NewParser(tokens)
	ast := parser.Parse()
	if ast.Error != nil {
		fmt.Println("PARSING ERR:", ast.Error.AsString())
		return nil, nil
	} else {
		fmt.Println(ast.Node)
	}

	context := NewContext("<program>", nil, nil)
	context.SymbolTable = GlobalSymbolTable
	interpreter := NewInterpreter()
	result := interpreter.visit(ast.Node, context)

	return result.Value, result.Error
}

func main() {
	// TODO Number.null, EP 11 "own" type into Number built in
	GlobalSymbolTable.Set("null", *NewNumber(0))
	GlobalSymbolTable.Set("false", *NewNumber(0))
	GlobalSymbolTable.Set("true", *NewNumber(1))
	GlobalSymbolTable.Set("print", *NewBuildInFunction("Print"))

	fileName := "file.ecp"
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Process the line
		fmt.Println("Processing line:", line)

		// Run your function for each line
		result, err := run(fileName, line)

		if err != nil {
			fmt.Println(err.AsString())
			break
		} else if result != nil {
			if result.Number != nil {
				fmt.Println(result.Number.ValueField)
			} else if result.Function != nil {
				fmt.Println(result.Function.String())
			} else if result.String != nil {
				fmt.Println(result.String.ValueField)
			} else if result.Array != nil {
				fmt.Println(result.Array.String())
			}
		}
		lineTEMP++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}

// TODO return null wenn if nichts ausgibt
