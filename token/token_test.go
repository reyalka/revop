package token

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		tokenType Type
		ch        byte
		expected  Token
	}{
		{ASSIGN, '=', Token{Type: ASSIGN, Literal: "="}},
		{PLUS, '+', Token{Type: PLUS, Literal: "+"}},
		{LPAREN, '(', Token{Type: LPAREN, Literal: "("}},
		{ILLEGAL, 0, Token{Type: ILLEGAL, Literal: "\x00"}},
	}

	for _, tt := range tests {
		tok := New(tt.tokenType, tt.ch)
		if tok != tt.expected {
			t.Errorf("New(%q, %q) wrong. got=%+v, want=%+v", tt.tokenType, tt.ch, tok, tt.expected)
		}
	}
}

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		ident    string
		expected Type
	}{
		{"fn", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"foobar", IDENT},
		{"Let", IDENT},
		{"", IDENT},
	}

	for _, tt := range tests {
		if got := LookupIdent(tt.ident); got != tt.expected {
			t.Errorf("LookupIdent(%q) wrong. got=%q, want=%q", tt.ident, got, tt.expected)
		}
	}
}
