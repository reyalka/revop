package main

import "revop/lexer"

// for lexer test
func main() {
	input := `let i = 5;`;
	l := lexer.New(input)
	// show all tokens
	for tok := l.NextToken(); tok.Type != "EOF"; tok = l.NextToken() {
		println(tok.Type, tok.Literal)
	}
}

