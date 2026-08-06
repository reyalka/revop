package main

import (
	"fmt"
	"os"
	"os/user"
	"revop/repl"
)

func main() {
	name := "user"
	if u, err := user.Current(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not determine current user: %v\n", err)
	} else {
		name = u.Username
	}

	fmt.Printf("Hello %s! This is the Revop programming language!\n", name)
	fmt.Printf("Feel free to type in commands\n")

	if err := repl.Start(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "revop: %v\n", err)
		os.Exit(1)
	}
}
