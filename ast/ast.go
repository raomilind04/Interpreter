package ast

import (
	"bytes"
	"interpreter/token"
	"strings"
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

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
    return b.Token.Literal
}
func (b *Boolean) String() string {
    return b.Token.Literal
}

// If Expression

type IfExpression struct {
    Token          token.Token
    Condition      Expression
    Consequence    *BlockStatement
    Alternative    *BlockStatement
}

func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string {
    return ie.Token.Literal
}
func (ie *IfExpression) String() string {
    var output bytes.Buffer

    output.WriteString("if")
    output.WriteString(ie.Condition.String())
    output.WriteString(" ")
    output.WriteString(ie.Consequence.String())

    if ie.Alternative != nil {
        output.WriteString("else")
        output.WriteString(ie.Alternative.String())
    }

    return output.String()
}

type BlockStatement struct {
    Token       token.Token
    Statements  []Statement
}

func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
    return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
    var output bytes.Buffer

    for _, s := range bs.Statements {
        output.WriteString(s.String())
    }
    return output.String()
}

// Functional Literal

type FunctionLiteral struct {
    Token       token.Token
    Parameters  []*Identifier
    Body        *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string {
    return fl.Token.Literal
}
func (fl *FunctionLiteral) String() string {
    var output bytes.Buffer

    params := []string{}
    for _, p := range fl.Parameters {
        params = append(params, p.String())
    }

    output.WriteString(fl.TokenLiteral())
    output.WriteString("(")
    output.WriteString(strings.Join(params, ","))
    output.WriteString(") ")
    output.WriteString(fl.Body.String())

    return output.String()
}

type CallExpression struct {
    Token      token.Token
    Function   Expression
    Arguments  []Expression
}

func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
    return ce.Token.Literal
}
func (ce *CallExpression) String() string {
    var output bytes.Buffer

    args := []string{}

    for _, a := range ce.Arguments {
        args = append(args, a.String())
    }

    output.WriteString(ce.Function.String())
    output.WriteString("(")
    output.WriteString(strings.Join(args, ", "))
    output.WriteString(")")

    return output.String()
}

type StringLiteral struct {
    Token token.Token
    Value string
}
func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) TokenLiteral() string {
    return sl.Token.Literal
}
func (sl *StringLiteral) String() string {
    return sl.Token.Literal
}

// Arrays
type ArrayLiteral struct {
    Token token.Token
    Elements []Expression
}
func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) TokenLiteral() string {
    return al.Token.Literal
}
func (al *ArrayLiteral) String() string {
    var output bytes.Buffer

    elements := []string{}
    for _, e := range al.Elements {
        elements = append(elements, e.String())
    }

    output.WriteString("[")
    output.WriteString(strings.Join(elements, ", "))
    output.WriteString("]") 

    return output.String()
}

type IndexExpression struct {
    Token token.Token
    Left  Expression
    Index Expression
}
func (ie *IndexExpression) expressionNode() {}
func (ie *IndexExpression) TokenLiteral() string {
    return ie.Token.Literal
}
func (ie *IndexExpression) String() string {
    var output bytes.Buffer

    output.WriteString("(")
    output.WriteString(ie.Left.String())
    output.WriteString("[")
    output.WriteString(ie.Index.String())
    output.WriteString("])")

    return output.String()
}

// Hash Maps
type HashLiteral struct {
    Token token.Token
    Pairs map[Expression]Expression
}
func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) TokenLiteral() string {
    return hl.Token.Literal
}
func (hl *HashLiteral) String() string {
    var output bytes.Buffer

    pairs := []string{}
    for key, value := range hl.Pairs {
        pairs = append(pairs, key.String() + ":" + value.String())
    }
    
    output.WriteString("{")
    output.WriteString(strings.Join(pairs, ", "))
    output.WriteString("}") 

    return output.String()
}
