package object

import (
	"revop/ast"
	"revop/token"
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	hello1Hash := hello1.HashKey()
	hello2Hash := hello2.HashKey()
	diff1Hash := diff1.HashKey()
	diff2Hash := diff2.HashKey()

	if hello1Hash != hello2Hash {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1Hash != diff2Hash {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1Hash == diff1Hash {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}
	if one1.HashKey() == two.HashKey() {
		t.Errorf("integers with different values have same hash keys")
	}
	if one1.HashKey().Type != INTEGER {
		t.Errorf("hash key type wrong. got=%s, want=%s", one1.HashKey().Type, INTEGER)
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	falsy := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("booleans with same value have different hash keys")
	}
	if true1.HashKey() == falsy.HashKey() {
		t.Errorf("booleans with different values have same hash keys")
	}
	if true1.HashKey().Value != 1 || falsy.HashKey().Value != 0 {
		t.Errorf("boolean hash values wrong. true=%d, false=%d", true1.HashKey().Value, falsy.HashKey().Value)
	}
}

func TestHashKeyIsTypeScoped(t *testing.T) {
	integer := &Integer{Value: 1}
	boolean := &Boolean{Value: true}

	if integer.HashKey() == boolean.HashKey() {
		t.Errorf("integer 1 and true must not share a hash key")
	}
}

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

	key := &String{Value: "a"}
	hashMap := &HashMap{Pairs: map[HashKey]HashPair{
		key.HashKey(): {Key: key, Value: &Integer{Value: 1}},
	}}

	tests := []struct {
		obj             Object
		expectedType    Type
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
		{&HashMap{Pairs: map[HashKey]HashPair{}}, HASHMAP, "{}"},
		{hashMap, HASHMAP, "{a: 1}"},
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

func TestTypeString(t *testing.T) {
	tests := []struct {
		objectType Type
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
		{HASHMAP, "HASHMAP"},
		{Type(0), "Type(0)"},
		{Type(100), "Type(100)"},
	}

	for _, tt := range tests {
		if got := tt.objectType.String(); got != tt.expected {
			t.Errorf("Type(%d).String() wrong. got=%q, want=%q", tt.objectType, got, tt.expected)
		}
	}
}
