package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/danielrs/monkey/lexer"
	"github.com/danielrs/monkey/token"
)

const PROMPT = ">> "

// Starts is the REPL loop that goes forever.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, PROMPT)

		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
