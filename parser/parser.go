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

var precedences = map[token.TokenType]int {
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
}

type (
    prefixParseFn func() ast.Expression
    infixParseFn  func(ast.Expression) ast.Expression
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
    parser.registerPrefix(token.TRUE, parser.parseBoolean)
    parser.registerPrefix(token.FALSE, parser.parseBoolean)
    parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)

    parser.infixParseFns = make(map[token.TokenType]infixParseFn)
    parser.registerInfix(token.PLUS, parser.parseInfixExpression)
    parser.registerInfix(token.MINUS, parser.parseInfixExpression)
    parser.registerInfix(token.SLASH, parser.parseInfixExpression)
    parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
    parser.registerInfix(token.EQ, parser.parseInfixExpression)
    parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
    parser.registerInfix(token.LT, parser.parseInfixExpression)
    parser.registerInfix(token.GT, parser.parseInfixExpression)


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

    for !parser.peekTokenIs(token.SEMICOLON) && precedence < parser.peekPrecedence() {
        inflix := parser.infixParseFns[parser.peekToken.Type]
        if inflix == nil {
            return leftExp
        }

        parser.nextToken()

        leftExp = inflix(leftExp)
    }

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

func (parser* Parser) parseBoolean() ast.Expression {
    return &ast.Boolean{
        Token: parser.currToken, 
        Value: parser.currTokenIs(token.TRUE),
    }
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

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token:      parser.currToken,
        Operator:   parser.currToken.Literal,
        Left:       left,
    }
    precedence := parser.currPrecedence()
    parser.nextToken()
    expression.Right = parser.parseExpression(precedence)
    
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

func (parser *Parser) parseGroupedExpression() ast.Expression {
    parser.nextToken()

    expression := parser.parseExpression(LOWEST)

    if !parser.expectPeek(token.RPAREN) {
        return nil
    }

    return expression
}

func (parser *Parser) peekPrecedence() int {
    if p, ok := precedences[parser.peekToken.Type]; ok {
        return p
    }
    return LOWEST
}

func (parser *Parser) currPrecedence() int {
    if p, ok := precedences[parser.currToken.Type]; ok {
        return p
    }
    return LOWEST
}

