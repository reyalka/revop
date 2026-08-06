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

// Start runs the read-eval-print loop until in is exhausted. It returns an
// error when reading from in or writing to out fails.
func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		if _, err := fmt.Fprint(out, PROMPT); err != nil {
			return fmt.Errorf("write prompt: %w", err)
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("read input: %w", err)
			}
			return nil
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			if err := printParserErrors(out, p.Errors()); err != nil {
				return err
			}
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			if _, err := fmt.Fprintf(out, "%s\n", evaluated.Inspect()); err != nil {
				return fmt.Errorf("write result: %w", err)
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) error {
	if _, err := fmt.Fprintln(out, "Parser errors found"); err != nil {
		return fmt.Errorf("write parser errors: %w", err)
	}

	for _, msg := range errors {
		if _, err := fmt.Fprintf(out, "\t%s\n", msg); err != nil {
			return fmt.Errorf("write parser errors: %w", err)
		}
	}

	return nil
}
