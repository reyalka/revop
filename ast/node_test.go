package ast

import (
	"revop/token"
	"testing"
)

func ident(name string) *Identifier {
	return &Identifier{Token: token.Token{Type: token.IDENT, Literal: name}, Value: name}
}

func intLit(literal string, value int64) *IntegerLiteral {
	return &IntegerLiteral{Token: token.Token{Type: token.INT, Literal: literal}, Value: value}
}

func block(stmts ...Statement) *Block {
	return &Block{Token: token.Token{Type: token.LBRACE, Literal: "{"}, Statements: stmts}
}

func TestNodeString(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{
			"empty program",
			&Program{},
			"",
		},
		{
			"let statement without value",
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name:  ident("x"),
			},
			"let x = ;",
		},
		{
			"return statement",
			&ReturnStatement{
				Token:       token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: intLit("5", 5),
			},
			"return 5;",
		},
		{
			"return statement without value",
			&ReturnStatement{Token: token.Token{Type: token.RETURN, Literal: "return"}},
			"return ;",
		},
		{
			"expression statement",
			&ExpressionStatement{
				Token:      token.Token{Type: token.IDENT, Literal: "x"},
				Expression: ident("x"),
			},
			"x",
		},
		{
			"expression statement without expression",
			&ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "x"}},
			"",
		},
		{
			"integer literal",
			intLit("5", 5),
			"5",
		},
		{
			"prefix expression",
			&PrefixExpression{
				Token:    token.Token{Type: token.BANG, Literal: "!"},
				Operator: "!",
				Right:    ident("x"),
			},
			"(!x)",
		},
		{
			"infix expression",
			&InfixExpression{
				Token:    token.Token{Type: token.PLUS, Literal: "+"},
				Left:     intLit("1", 1),
				Operator: "+",
				Right:    intLit("2", 2),
			},
			"(1 + 2)",
		},
		{
			"boolean",
			&Boolean{Token: token.Token{Type: token.TRUE, Literal: "true"}, Value: true},
			"true",
		},
		{
			"if expression",
			&IfExpression{
				Token:       token.Token{Type: token.IF, Literal: "if"},
				Condition:   ident("x"),
				Consequence: block(&ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "y"}, Expression: ident("y")}),
			},
			"if (x) { y; }",
		},
		{
			"if else expression",
			&IfExpression{
				Token:       token.Token{Type: token.IF, Literal: "if"},
				Condition:   ident("x"),
				Consequence: block(&ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "y"}, Expression: ident("y")}),
				Alternative: block(&ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "z"}, Expression: ident("z")}),
			},
			"if (x) { y; }else { z; }",
		},
		{
			"block with non expression statement",
			block(&ReturnStatement{
				Token:       token.Token{Type: token.RETURN, Literal: "return"},
				ReturnValue: intLit("1", 1),
			}),
			"{ return 1; }",
		},
		{
			"function literal",
			&FunctionLiteral{
				Token:      token.Token{Type: token.FUNCTION, Literal: "fn"},
				Parameters: []*Identifier{ident("x"), ident("y")},
				Body:       block(&ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "x"}, Expression: ident("x")}),
			},
			"fn(x, y) { x; }",
		},
		{
			"function literal without parameters",
			&FunctionLiteral{
				Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
				Body:  block(),
			},
			"fn() {  }",
		},
		{
			"call expression",
			&CallExpression{
				Token:     token.Token{Type: token.LPAREN, Literal: "("},
				Function:  ident("add"),
				Arguments: []Expression{intLit("1", 1), ident("x")},
			},
			"<add>(1, x)",
		},
		{
			"string literal",
			&StringLiteral{Token: token.Token{Type: token.STRING, Literal: "hello"}, Value: "hello"},
			"hello",
		},
		{
			"array literal",
			&ArrayLiteral{
				Token:    token.Token{Type: token.LBRACKET, Literal: "["},
				Elements: []Expression{intLit("1", 1), ident("x")},
			},
			"[1, x]",
		},
		{
			"hash map literal",
			&HashMapLiteral{
				Token: token.Token{Type: token.LBRACE, Literal: "{"},
				Pairs: map[Expression]Expression{
					&StringLiteral{Token: token.Token{Type: token.STRING, Literal: "a"}, Value: "a"}: intLit("1", 1),
				},
			},
			"{a:1}",
		},
		{
			"empty hash map literal",
			&HashMapLiteral{Token: token.Token{Type: token.LBRACE, Literal: "{"}},
			"{}",
		},
		{
			"index expression",
			&IndexExpression{
				Token: token.Token{Type: token.LBRACKET, Literal: "["},
				Left:  ident("arr"),
				Index: intLit("0", 0),
			},
			"(arr[0])",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.String(); got != tt.expected {
				t.Errorf("String() wrong. got=%q, want=%q", got, tt.expected)
			}
		})
	}
}

func TestTokenLiteral(t *testing.T) {
	tests := []struct {
		name     string
		node     Node
		expected string
	}{
		{"let statement", &LetStatement{Token: token.Token{Type: token.LET, Literal: "let"}}, "let"},
		{"identifier", ident("x"), "x"},
		{"return statement", &ReturnStatement{Token: token.Token{Type: token.RETURN, Literal: "return"}}, "return"},
		{"expression statement", &ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "x"}}, "x"},
		{"integer literal", intLit("5", 5), "5"},
		{"prefix expression", &PrefixExpression{Token: token.Token{Type: token.BANG, Literal: "!"}}, "!"},
		{"infix expression", &InfixExpression{Token: token.Token{Type: token.PLUS, Literal: "+"}}, "+"},
		{"boolean", &Boolean{Token: token.Token{Type: token.TRUE, Literal: "true"}}, "true"},
		{"if expression", &IfExpression{Token: token.Token{Type: token.IF, Literal: "if"}}, "if"},
		{"block", block(), "{"},
		{"function literal", &FunctionLiteral{Token: token.Token{Type: token.FUNCTION, Literal: "fn"}}, "fn"},
		{"call expression", &CallExpression{Token: token.Token{Type: token.LPAREN, Literal: "("}}, "("},
		{"string literal", &StringLiteral{Token: token.Token{Type: token.STRING, Literal: "hello"}}, "hello"},
		{"array literal", &ArrayLiteral{Token: token.Token{Type: token.LBRACKET, Literal: "["}}, "["},
		{"index expression", &IndexExpression{Token: token.Token{Type: token.LBRACKET, Literal: "["}}, "["},
		{"hash map literal", &HashMapLiteral{Token: token.Token{Type: token.LBRACE, Literal: "{"}}, "{"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.TokenLiteral(); got != tt.expected {
				t.Errorf("TokenLiteral() wrong. got=%q, want=%q", got, tt.expected)
			}
		})
	}
}

func TestProgramTokenLiteral(t *testing.T) {
	empty := &Program{}
	if got := empty.TokenLiteral(); got != "" {
		t.Errorf("empty program TokenLiteral() wrong. got=%q, want=%q", got, "")
	}

	program := &Program{
		Statements: []Statement{
			&LetStatement{Token: token.Token{Type: token.LET, Literal: "let"}, Name: ident("x")},
			&ReturnStatement{Token: token.Token{Type: token.RETURN, Literal: "return"}},
		},
	}
	if got := program.TokenLiteral(); got != "let" {
		t.Errorf("program TokenLiteral() wrong. got=%q, want=%q", got, "let")
	}
}
