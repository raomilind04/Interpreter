package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
    input := `
    let x = 5;
    let y = 10;
    let foobar = 8080;
    `
    lexer := lexer.New(input)
    parser := New(lexer)

    program := parser.ParseProgram() 
    checkParserErrors(t, parser)

    if program == nil {
        t.Fatalf("Parse returned nil")
    }

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements has %d Statements, wanted 3", len(program.Statements))
    }

    tests := []struct {
        expectedIdentifier string
    }{
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i, tt := range tests {
        statement := program.Statements[i]
        if !testLetStatement(t, statement, tt.expectedIdentifier) {
            return 
        }
    }
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("statement.TokenLiteral is not let, got: %q", s.TokenLiteral())
        return false
    }

    letStatement, ok := s.(*ast.LetStatement)
    if !ok {
        t.Errorf("s not *ast.LetStatement, got: %T", s)
        return false
    }

    if letStatement.Name.Value != name {
        t.Errorf("letStatement.Name.Value should be %s but got %s", name, letStatement.Name.Value)
        return false
    }

    if letStatement.Name.TokenLiteral() != name {
        t.Errorf("s.Name should be %s but got %s", name, letStatement.Name)
        return false
    }
    return true
}

func checkParserErrors(t *testing.T, parser *Parser) {
    errors := parser.Errors()

    if len(errors) == 0 {
        return 
    }

    t.Errorf("parser has %d errors", len(errors))
    for _, err := range errors {
        t.Errorf("parser error: %q", err)
    }
    t.FailNow()
}

func TestReturnStatements(t *testing.T) {
    input := `
    return 5;
    return 10;
    return 998877;
    `
    lexer := lexer.New(input)
    parser := New(lexer)

    program := parser.ParseProgram()
    checkParserErrors(t, parser)

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements should have 3 statement but had %d", len(program.Statements))
    }

    for _, statement := range program.Statements {
        returnStatement, ok := statement.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("statement is not *ast.ReturnStatement, got: %T", statement)
            continue
        }
        if returnStatement.TokenLiteral() != "return" {
            t.Errorf("returnStatement.TokenLiteral should be 'return but was %q", returnStatement.TokenLiteral())
        }
    }
}
