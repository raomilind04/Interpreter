package ast

import (
	"bytes"
	"interpreter/token"
)

type Node interface {
    TokenLiteral() string
    String() string
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

func (prog *Program) String() string {
    var output bytes.Buffer

    for _, s := range prog.Statements {
        output.WriteString(s.String())
    }

    return output.String()
}


// LET statements

type LetStatement struct {
    Token   token.Token
    Name    *Identifier
    Value   Expression 
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}
func (ls *LetStatement) String() string {
    var output bytes.Buffer

    output.WriteString(ls.TokenLiteral() + " ")
    output.WriteString(ls.Name.String)
    output.WriteString(" = ")

    if ls.Value != nil {
        output.WriteString(ls.Value.String())
    }

    output.WriteString(";")

    return output.String()
}

type Identifier struct {
    Token   token.Token
    Value   string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
    return i.Token.Literal
}


// Return Statement

type ReturnStatement struct {
    Token       token.Token
    ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
    return rs.Token.Literal
}


// Expression Statement

type ExpressionStatement struct {
    Token       token.Token
    Expression  Expression
}

func (es *ExpressionStatement) statementNode()  {}
func (es *ExpressionStatement) TokenLiteral() string {
    return es.Token.Literal
}
