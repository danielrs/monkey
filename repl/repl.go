package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/danielrs/monkey/evaluator"
	"github.com/danielrs/monkey/lexer"
	"github.com/danielrs/monkey/object"
	"github.com/danielrs/monkey/parser"
)

const PROMPT = ">> "

// Starts is the REPL loop that goes forever.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprintf(out, PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(env, program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func Run(str string, out io.Writer) {
	env := object.NewEnvironment()

	l := lexer.New(str)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(env, program)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
