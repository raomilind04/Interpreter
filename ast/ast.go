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
    output.WriteString(ls.Name.String())
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
func (i *Identifier) String() string {
    return i.Value
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
func (rs *ReturnStatement) String() string {
    var output bytes.Buffer

    output.WriteString(rs.TokenLiteral() + " ")

    if rs.ReturnValue != nil {
        output.WriteString(rs.ReturnValue.String())
    }

    output.WriteString(";")

    return output.String()
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
func (es *ExpressionStatement) String() string {
    if es.Expression != nil {
        return es.Expression.String()
    }
    return ""
}

// Integer Literal

type IntegerLiteral struct {
    Token   token.Token
    Value   int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

// Prefix Expression

type PrefixExpression struct {
    Token     token.Token
    Operator  string
    Right     Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
    var output bytes.Buffer

    output.WriteString("(")
    output.WriteString(pe.Operator)
    output.WriteString(pe.Right.String())
    output.WriteString(")")

    return output.String()
}

// Infix Expression

type InfixExpression struct {
    Token      token.Token
    Left       Expression
    Operator   string
    Right      Expression
}

func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string {
    return oe.Token.Literal
}
func (oe *InfixExpression) String() string {
    var output bytes.Buffer

    output.WriteString("(")
    output.WriteString(oe.Left.String())
    output.WriteString(" " + oe.Operator + " ")
    output.WriteString(oe.Right.String())
    output.WriteString(")")

    return output.String()
}

// Boolean Literal

type Boolean struct {
    Token token.Token
    Value bool
}

func (b* Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}
func (b *Boolean) String() string {
    return b.Token.Literal
}
