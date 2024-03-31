package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
    "strconv"
)

const (
    _ int = iota
    LOWEST
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

type (
    prefixParseFn func() ast.Expression
    infixParseFn  func() ast.Expression
)

type Parser struct {
    lexer *lexer.Lexer
    
    errors []string

    currToken token.Token
    peekToken token.Token

    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
    parser := &Parser{
        lexer:  lexer,
        errors: []string{},
    }
    
    parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    parser.registerPrefix(token.IDENT, parser.parseIdentifier)
    parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
    parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
    parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)

    parser.nextToken()
    parser.nextToken()

    return parser
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    parser.prefixParseFns[tokenType] = fn
}

func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    parser.infixParseFns[tokenType] = fn 
}

func (parser *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
}

func (parser *Parser) Errors() []string {
    return parser.errors
}

func (parser *Parser) noPrefixFnError(token token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function found for %s", token)
    parser.errors = append(parser.errors, msg)
}

func (parser *Parser) peekError(tokenType token.TokenType) {
    err := fmt.Sprintf("Next token should be %s but got %s", tokenType, parser.peekToken.Type)
    parser.errors = append(parser.errors, err)
}

func (parser *Parser) nextToken() {
    parser.currToken = parser.peekToken
    parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for parser.currToken.Type != token.EOF {
        statement := parser.parseStatement()

        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        parser.nextToken()
    }
    return program
}

func (parser *Parser) parseStatement() ast.Statement {
    switch parser.currToken.Type {
    case token.LET:
        return parser.parseLetStatement()
    case token.RETURN:
        return parser.parserReturnStatement()
    default:
        return parser.parseExpressionStatement()
    }
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
    statement := &ast.LetStatement{Token: parser.currToken}

    if !parser.expectPeek(token.IDENT) {
        return nil
    }

    statement.Name = &ast.Identifier{
        Token: parser.currToken,
        Value: parser.currToken.Literal,
    }

    for !parser.currTokenIs(token.SEMICOLON) {
        parser.nextToken()
    }
    
    return statement
}

func (parser *Parser) parserReturnStatement() *ast.ReturnStatement {
    statement := &ast.ReturnStatement{Token: parser.currToken}

    parser.nextToken()

    for !parser.currTokenIs(token.SEMICOLON) {
        parser.nextToken()
    }

    return statement
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    statement := &ast.ExpressionStatement{Token: parser.currToken}

    statement.Expression = parser.parseExpression(LOWEST)

    if parser.peekTokenIs(token.SEMICOLON) {
        parser.nextToken()
    }

    return statement
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
    prefix := parser.prefixParseFns[parser.currToken.Type]
    if prefix == nil {
        parser.noPrefixFnError(parser.currToken.Type)
        return nil
    }
    leftExp := prefix()

    return leftExp; 
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
    il := &ast.IntegerLiteral{Token: parser.currToken}

    value, err := strconv.ParseInt(parser.currToken.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("unable to parse literal %q as int", parser.currToken.Literal)
        parser.errors = append(parser.errors, msg)
        return nil
    }
    il.Value = value

    return il
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token:    parser.currToken,
        Operator: parser.currToken.Literal,
    }

    parser.nextToken()

    expression.Right = parser.parseExpression(PREFIX)

    return expression
}

func (parser *Parser) currTokenIs(tokenType token.TokenType) bool {
    return parser.currToken.Type == tokenType
}

func (parser *Parser) peekTokenIs(tokenType token.TokenType) bool {
    return parser.peekToken.Type == tokenType
}

func (parser *Parser) expectPeek(tokenType token.TokenType) bool {
    if parser.peekTokenIs(tokenType) {
        parser.nextToken()
        return true
    } else {
        parser.peekError(tokenType)
        return false
    }
}

