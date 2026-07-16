package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "="
	PLUS      = "+"
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	FUNCTION  = "FUNCTION"
	LET       = "LET"
)

func New(tokenType TokenType, ch byte) Token {
	return Token{
		Type: tokenType,
		Literal: string(ch),
	}
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// 識別子が予約語であればkeywordのTokenTypeを、そうでなければ単にIDENTを返す
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}