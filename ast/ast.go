package ast

import (
    "interpreter/token"
)

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (prog *Program) TokenLiteral() string {
    if len(prog.Statements) > 0 {
        return prog.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

type LetStatement struct {
    Token   token.Token
    Name    *Identifier
    Value   Expression 
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}

type Identifier struct {
    Token   token.Token
    Value   string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}


