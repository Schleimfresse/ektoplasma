package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
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
	TT_NEWLINE    TokenTypes = "NEWLINE"
	TT_ARROW      TokenTypes = "ARROW"
	Zero          Binary     = 0
	One           Binary     = 1
)

var KEYWORDS = []string{"VAR", "AND", "OR", "NOT", "IF", "THEN", "ELSE", "ELIF", "FOR", "TO", "STEP", "WHILE", "FUNC", "END", "IMPORT"}
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
		fmt.Println(ast.Error.AsString())
		return nil, nil
	}
	// TODO fix pos:
	//IDENTIFIER input START: {0 0 0 file.ecp input()} END: {5 0 5 file.ecp input()}
	//LPAREN <nil> START: {5 0 5 file.ecp input()} END: {5 0 5 file.ecp input()}
	//RPAREN <nil> START: {6 0 6 file.ecp input()} END: {6 0 6 file.ecp input()}
	context := NewContext("<program>", nil, nil)
	context.SymbolTable = GlobalSymbolTable
	interpreter := NewInterpreter()
	result := interpreter.visit(ast.Node, context)
	return result.Value, result.Error
}

func main() {
	GlobalSymbolTable.Set("null", NewNull())
	GlobalSymbolTable.Set("false", NewBoolean(0))
	GlobalSymbolTable.Set("true", NewBoolean(1))
	GlobalSymbolTable.Set("print", NewBuildInFunction("Print"))
	GlobalSymbolTable.Set("println", NewBuildInFunction("PrintLn"))
	GlobalSymbolTable.Set("input", NewBuildInFunction("Input"))
	GlobalSymbolTable.Set("isString", NewBuildInFunction("isString"))
	GlobalSymbolTable.Set("isNumber", NewBuildInFunction("isNumber"))
	GlobalSymbolTable.Set("isFunction", NewBuildInFunction("isFunction"))
	GlobalSymbolTable.Set("isArray", NewBuildInFunction("isArray"))
	GlobalSymbolTable.Set("append", NewBuildInFunction("append"))
	GlobalSymbolTable.Set("len", NewBuildInFunction("len"))

	if len(os.Args) >= 2 {
		filePath, _ := filepath.Abs(os.Args[1])
		fileName := path.Base(filePath)
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Invalid path, cannot open specified file.")
			return
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		ScanLine(fileName, scanner)
	} else {
		for {
			buf := make([]byte, 1024)
			n, err := syscall.Read(syscall.Stdin, buf)
			if err != nil {
				log.Fatal(err)
			}

			newReader := bufio.NewReader(bytes.NewReader(buf[:n]))
			scanner := bufio.NewScanner(newReader)
			ScanLine("<stdin>", scanner)
		}
	}
}

func ScanLine(fileName string, scanner *bufio.Scanner) {
	for scanner.Scan() {
		line := scanner.Text()

		line = strings.ReplaceAll(line, "", "")
		if line == "" {
			continue
		}

		// Process the line
		fmt.Println("Processing line:", line)

		// Run your function for each line
		result, err := run(fileName, line)

		// log.Println("RESULT:", result)
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
				if len(result.Array.Elements) == 1 && result.Array.Elements[0] != nil {
					fmt.Println(result.Array.Elements[0].Value())
				} else {
					fmt.Println(result.Array.String())
				}
			} else if result.Boolean != nil {
				fmt.Println(result.Boolean.String())
			} else {
				fmt.Println(result.Null.String())
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
