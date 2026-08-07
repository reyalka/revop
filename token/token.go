package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL   Type = "ILLEGAL"
	EOF       Type = "EOF"
	IDENT     Type = "IDENT"
	INT       Type = "INT"
	STRING    Type = "STRING"
	ASSIGN    Type = "="
	PLUS      Type = "+"
	MINUS     Type = "-"
	BANG      Type = "!"
	ASTERISK  Type = "*"
	SLASH     Type = "/"
	COLON     Type = ":"
	HASH      Type = "#"
	LT        Type = "<"
	GT        Type = ">"
	LE        Type = "<="
	GE        Type = ">="
	EQ        Type = "=="
	NOT_EQ    Type = "!="
	COMMA     Type = ","
	SEMICOLON Type = ";"
	LPAREN    Type = "("
	RPAREN    Type = ")"
	LBRACE    Type = "{"
	RBRACE    Type = "}"
	LBRACKET  Type = "["
	RBRACKET  Type = "]"
	PIPE      Type = "|>"
	FUNCTION  Type = "FUNCTION"
	LET       Type = "LET"
	TRUE      Type = "TRUE"
	FALSE     Type = "FALSE"
	IF        Type = "IF"
	ELSE      Type = "ELSE"
	RETURN    Type = "RETURN"
	MUT       Type = "MUT"
)

func New(tokenType Type, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

var keywords = map[string]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"mut":    MUT,
}

// 識別子が予約語であればkeywordのTokenTypeを、そうでなければ単にIDENTを返す
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
