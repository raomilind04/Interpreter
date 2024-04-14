package repl

import (
	"bufio"
	"fmt"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"io"
)

const PROMPT = ">> "

func Start(input io.Reader, output io.Writer) {
    scanner := bufio.NewScanner(input)
    env := object.NewEnvironment()

    for {
        fmt.Printf(PROMPT)
        scanned := scanner.Scan()
        if !scanned {
            return
        }

        line := scanner.Text()
        lex := lexer.New(line)
        parser := parser.New(lex)

        program := parser.ParseProgram()
        if len(parser.Errors()) != 0 {
            printParserErrors(output, parser.Errors())
            continue
        }

        evaluated := evaluator.Eval(program, env)
        if evaluated != nil {
            io.WriteString(output, evaluated.Inspect())
            io.WriteString(output, "\n")
        }
        // io.WriteString(output, program.String())
        // io.WriteString(output, "\n")
        
        // for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
        //     fmt.Printf("%+v\n", tok)
        // }
    }
}

func printParserErrors(output io.Writer, errors []string) {
    for _, msg := range errors {
        io.WriteString(output, "\t" + msg + "\n") 
    }
}
