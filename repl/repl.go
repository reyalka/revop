package repl

import (
	"bufio"
	"fmt"
	"io"
	"revop/evaluator"
	"revop/lexer"
	"revop/object"
	"revop/parser"
)

const PROMPT = ">> "

// maxLineSize bounds a single input line so that overlong input is reported
// instead of silently ending the session.
const maxLineSize = 1024 * 1024

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 0, bufio.MaxScanTokenSize), maxLineSize)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			if err := scanner.Err(); err != nil {
				_, _ = fmt.Fprintf(out, "input error: %s\n", err)
			}
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

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {

			_, _ = io.WriteString(out, evaluated.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "Parser errors found")

	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
