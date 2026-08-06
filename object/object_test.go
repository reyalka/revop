package object

import (
	"revop/ast"
	"revop/token"
	"testing"
)

func TestTypeAndInspect(t *testing.T) {
	fn := &Function{
		Parameters: []*ast.Identifier{
			{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"},
			{Token: token.Token{Type: token.IDENT, Literal: "y"}, Value: "y"},
		},
		Body: &ast.Block{
			Token: token.Token{Type: token.LBRACE, Literal: "{"},
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Expression: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
				},
			},
		},
		Env: NewEnvironment(),
	}

	tests := []struct {
		obj             Object
		expectedType    ObjectType
		expectedInspect string
	}{
		{&Integer{Value: 5}, INTEGER, "5"},
		{&Integer{Value: -5}, INTEGER, "-5"},
		{&Boolean{Value: true}, BOOLEAN, "true"},
		{&Boolean{Value: false}, BOOLEAN, "false"},
		{&Null{}, NULL, "null"},
		{&ReturnValue{Value: &Integer{Value: 10}}, RETURN_VALUE, "10"},
		{&Error{Message: "boom"}, ERROR, "ERROR: boom"},
		{&String{Value: "hello world"}, STRING, "hello world"},
		{&Builtin{}, BUILTIN, "builtin function"},
		{&Array{Elements: []Object{}}, ARRAY, "[]"},
		{&Array{Elements: []Object{&Integer{Value: 1}, &String{Value: "a"}}}, ARRAY, "[1, a]"},
		{fn, FUNCTION, "fn(x, y) {\n{ x; }\n}"},
	}

	for _, tt := range tests {
		if got := tt.obj.Type(); got != tt.expectedType {
			t.Errorf("%T.Type() wrong. got=%s, want=%s", tt.obj, got, tt.expectedType)
		}
		if got := tt.obj.Inspect(); got != tt.expectedInspect {
			t.Errorf("%T.Inspect() wrong. got=%q, want=%q", tt.obj, got, tt.expectedInspect)
		}
	}
}

func TestNewError(t *testing.T) {
	err := NewError("unknown operator: %s %s", "-", "STRING")

	if err.Message != "unknown operator: - STRING" {
		t.Errorf("err.Message wrong. got=%q", err.Message)
	}
	if err.Type() != ERROR {
		t.Errorf("err.Type() wrong. got=%s", err.Type())
	}
}

func TestBuiltinRun(t *testing.T) {
	called := 0
	b := &Builtin{
		Args: 1,
		Fn: func(args ...Object) Object {
			called++
			return args[0]
		},
	}

	result := b.Run(&String{Value: "hi"})
	if str, ok := result.(*String); !ok || str.Value != "hi" {
		t.Fatalf("Run returned wrong value. got=%+v", result)
	}
	if called != 1 {
		t.Fatalf("Fn called %d times, want 1", called)
	}

	result = b.Run()
	err, ok := result.(*Error)
	if !ok {
		t.Fatalf("Run with wrong arity did not return *Error. got=%T", result)
	}
	if err.Message != "wrong number of arguments. got=0, want=1" {
		t.Errorf("err.Message wrong. got=%q", err.Message)
	}
	if called != 1 {
		t.Errorf("Fn must not be called on arity mismatch. called=%d", called)
	}
}

func TestObjectTypeString(t *testing.T) {
	tests := []struct {
		objectType ObjectType
		expected   string
	}{
		{INTEGER, "INTEGER"},
		{BOOLEAN, "BOOLEAN"},
		{NULL, "NULL"},
		{RETURN_VALUE, "RETURN_VALUE"},
		{ERROR, "ERROR"},
		{FUNCTION, "FUNCTION"},
		{STRING, "STRING"},
		{BUILTIN, "BUILTIN"},
		{ARRAY, "ARRAY"},
		{ObjectType(0), "ObjectType(0)"},
		{ObjectType(100), "ObjectType(100)"},
	}

	for _, tt := range tests {
		if got := tt.objectType.String(); got != tt.expected {
			t.Errorf("ObjectType(%d).String() wrong. got=%q, want=%q", tt.objectType, got, tt.expected)
		}
	}
}
