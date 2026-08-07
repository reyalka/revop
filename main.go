package main

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/reyalka/revop/evaluator"
	"github.com/reyalka/revop/lexer"
	"github.com/reyalka/revop/object"
	"github.com/reyalka/revop/parser"
	"github.com/reyalka/revop/repl"
)

func main() {
	// 引数が指定されていない、または "repl" の場合は REPL を起動
	if len(os.Args) < 2 || os.Args[1] == "repl" {
		startRepl()
		return
	}

	// ファイルまたはディレクトリパスが指定された場合
	runFile(os.Args[1])
}

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Revop programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}

func runFile(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// ディレクトリが渡された場合は <path>/main.rv を探す
	if fi.IsDir() {
		path = filepath.Join(path, "main.rv")
	}

	code, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %s: %s\n", path, err)
		os.Exit(1)
	}

	env := object.NewEnvironment()
	l := lexer.New(string(code))
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil && evaluated.Type() == object.ERROR {
		fmt.Fprintf(os.Stderr, "Runtime error: %s\n", evaluated.Inspect())
		os.Exit(1)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = io.WriteString(out, "Parser errors found:\n")
	for _, msg := range errors {
		_, _ = io.WriteString(out, "\t"+msg+"\n")
	}
}
