package lexer

import (
	"interpreter/token"
)

type Lexer struct {
    input         string
    position      int
    readPosition  int
    character     byte
}

func New(input string) *Lexer {
    lexer := &Lexer{input: input}
    lexer.readChar()
    return lexer
}

func (lexer *Lexer) readChar() {
    if lexer.readPosition >= len(lexer.input) {
        lexer.character = 0
    } else {
        lexer.character = lexer.input[lexer.readPosition]
    }
    lexer.position = lexer.readPosition
    lexer.readPosition += 1
}

func (lexer *Lexer) NextToken() token.Token {
    var tok token.Token

    lexer.skipWhiteSpaces()

    switch lexer.character {
    case '=':
        tok = newToken(token.ASSIGN, lexer.character)
    case ';':
        tok = newToken(token.SEMICOLON, lexer.character)
    case ',':
        tok = newToken(token.COMMA, lexer.character)
    case '+':
        tok = newToken(token.PLUS, lexer.character)
    case '(':
        tok = newToken(token.LPAREN, lexer.character)
    case ')':
        tok = newToken(token.RPAREN, lexer.character)
    case '{':
        tok = newToken(token.LBRACE, lexer.character)
    case '}':
        tok = newToken(token.RBRACE, lexer.character)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(lexer.character) {
            tok.Literal = lexer.readIdentifier()
            tok.Type = token.LookupIdent(tok.Literal)
            return tok
        } else if isDigit(lexer.character) {
            tok.Type = token.INT
            tok.Literal = lexer.readNumber()
            return tok
        }else {
            tok = newToken(token.ILLEGAL, lexer.character)
        }
    }
    lexer.readChar()
    return tok
}

func newToken(tokenType token.TokenType, character byte) token.Token {
    return token.Token{
        Type:     tokenType,
        Literal:  string(character),
    }
}

func (lexer *Lexer) readIdentifier() string {
    position := lexer.position
    for isLetter(lexer.character) {
        lexer.readChar()
    }
    return lexer.input[position:lexer.position]
}

func isLetter(character byte) bool {
    return (('a' <= character && 'z' >= character) ||
    ('A' <= character && 'Z' >= character) || character == '_')
}

func (lexer *Lexer) skipWhiteSpaces() {
    for lexer.character == ' ' || lexer.character == '\t' || lexer.character == '\n' || lexer.character == '\r' {
        lexer.readChar()
    }
}

func (lexer *Lexer) readNumber() string {
    position := lexer.position
    for isDigit(lexer.character) {
        lexer.readChar()
    }
    return lexer.input[position:lexer.position]
}

func isDigit(character byte) bool {
    return '0' <= character && '9' >= character
}
