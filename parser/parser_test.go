package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
    "fmt"
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

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program does not have enough statements, got = %d\n", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
       t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement, got = %T\n", program.Statements[0])
    }

    ident, ok := stmt.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp not *ast.Identifier, got=%T\n", stmt.Expression)
    }

    if ident.Value != "foobar" {
        t.Fatalf("ident.Value not %s, got=%s\n", "foobar", ident.Value)
    }

    if ident.TokenLiteral() != "foobar" {
        t.Errorf("ident.TokenLiteral not %s, got = %s\n", "foobar", ident.TokenLiteral())
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParserErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program does not have enough statements, got=%d\n", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T\n", program.Statements[0])
    }

    literal, ok := stmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("exp not *ast.IntegerLiteral. got=%T\n", stmt.Expression)
    }

    if literal.Value != 5 {
        t.Errorf("literal.Value not %d, got=%d\n", 5, literal.Value)
    }

    if literal.TokenLiteral() != "5" {
        t.Errorf("literal.TokenLiteral not %s, got = %s\n", "5", literal.TokenLiteral())
    } 
}

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct{
        input string
        operator string
        value interface{}
    }{
        {"!5;", "!", 5},
        {"-15;", "-", 15},
        {"!true", "!", true},
        {"!false", "!", false},
    }

    for _, tt := range prefixTests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements, got %d\n", 1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T\n", program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.PrefixExpression. got = %T\n", stmt.Expression)
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s', got = %s", tt.operator, exp.Operator)
        }

    }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integ, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il not *ast.IntegerLiteral, got=%T\n", il)
        return false
    }

    if integ.Value != value {
        t.Errorf("integ.Value not %d, got=%s", value, integ.TokenLiteral())
        return false
    }

    return true
}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct{
        input string
        leftValue interface{} 
        operator string
        rightValue interface{}
    }{
        {"5 + 5", 5, "+",5},
        {"5 - 5", 5, "-",5},
        {"5 * 5", 5, "*",5},
        {"5 / 5", 5, "/",5},
        {"5 > 5", 5, ">",5},
        {"5 > 5", 5, ">",5},
        {"5 < 5", 5, "<",5},
        {"5 == 5", 5, "==",5},
        {"5 != 5", 5, "!=",5},
        {"true == true", true, "==", true},
        {"true != false", true, "!=", false},
        {"false == false", false, "==", false},
    }


    for _, tt := range infixTests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements, got=%d", 1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("exp is not ast.InfixExpression, got=%T", stmt.Expression)
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not %s, got=%s", tt.operator, exp.Operator)
        }

        if !testLiteralExpression(t, exp.Left, tt.leftValue) {
            return
        }
        if !testLiteralExpression(t, exp.Right, tt.rightValue) {
            return
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct{
        input string
        expected string
    }{
        {
            "-a*b", 
            "((-a) * b)",
        },
        {
            "!-a",
            "(!(-a))",
        },
        {
            "a + b + c", 
            "((a + b) + c)",
        },
        {
            "a + b - c", 
            "((a + b) - c)"},
        {
            "a * b * c",
            "((a * b) * c)",
        },
        {
            "a * b / c",
            "((a * b) / c)",
        },
        {
            "a + b / c",
            "(a + (b / c))",
        },
        {
            "a + b * c + d / e - f",
            "(((a + (b * c)) + (d / e)) - f)",
        },
        {
            "3 + 4; -5 * 5",
            "(3 + 4)((-5) * 5)",
        },
        {
            "5 > 4 == 3 < 4",
            "((5 > 4) == (3 < 4))",
        },
        {
            "5 < 4 != 3 > 4",
            "((5 < 4) != (3 > 4))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "true",
            "true",
        },
        {
            "false",
            "false",
        },
        {
            "3 > 5 == false", 
            "((3 > 5) == false)",
        },
        {
            "3 < 5 == true",
            "((3 < 5) == true)",
        },
    }

    for _, tt := range tests{
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)

        actual := program.String()

        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
    ident, ok := exp.(*ast.Identifier)
    if !ok {
        t.Errorf("exp not *ast.Identifierm got=%T", exp)
        return false
    }

    if ident.Value != value {
        t.Errorf("ident.Value not %s, got=%s", value, ident.Value)
        return false
    }

    if ident.TokenLiteral() != value {
        t.Errorf("ident.TokenLiteral not %s, got=%s", value, ident.TokenLiteral())
        return false
    }

    return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
    switch v := expected.(type) {
    case int:
        return testIntegerLiteral(t, exp, int64(v))
    case int64:
        return testIntegerLiteral(t, exp, v)
    case string:
        return testIdentifier(t, exp, v)
    case bool:
        return testBooleanLiteral(t, exp, v)
    }

    t.Errorf("type of exp not handled. got=%T", exp)
    return false
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
    bo, ok := exp.(*ast.Boolean)
    if !ok {
        t.Errorf("exp not *ast.Boolean, got=%T", exp)
        return false
    }

    if bo.Value != value {
        t.Errorf("bo.Value is not %t, got=%t", value, bo.Value)
        return false
    }

    if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
        t.Errorf("bo.TokenLiteral not %t, got=%s", value, bo.TokenLiteral())
        return false
    }

    return true
}
