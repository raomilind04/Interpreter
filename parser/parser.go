package parser

import (
	"fmt"
	"interpreter/ast"
	"interpreter/lexer"
	"interpreter/token"
)

type Parser struct {
    lexer *lexer.Lexer
    
    errors []string

    currToken token.Token
    peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
    parser := &Parser{
        lexer:  lexer,
        errors: []string{},
    }

    parser.nextToken()
    parser.nextToken()

    return parser
}

func (parser *Parser) Errors() []string {
    return parser.errors
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
    default:
        return nil
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
