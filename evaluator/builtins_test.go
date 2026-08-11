package evaluator

import (
	"io"
	"os"
	"testing"

	"github.com/reyalka/revop/object"
)

func TestLenBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len([1, 2, 3])`, 3},
		{`len(#{foo: 1})`, 1},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestPushBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"push([1, 2, 3], 4)", &object.Array{Elements: []object.Object{
			&object.Integer{Value: 1},
			&object.Integer{Value: 2},
			&object.Integer{Value: 3},
			&object.Integer{Value: 4},
		}}},
		{"push(1, 2)", "first argument to `push` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case *object.Array:
				assertArrayEquals(t, evaluated, expected)
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestMapBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"map([1, 2, 3], fn(x) { x * 2 })", &object.Array{Elements: []object.Object{
			&object.Integer{Value: 2},
			&object.Integer{Value: 4},
			&object.Integer{Value: 6},
		}}},
		{"map([1, 2, 3], fn(x) { x + 1 })", &object.Array{Elements: []object.Object{
			&object.Integer{Value: 2},
			&object.Integer{Value: 3},
			&object.Integer{Value: 4},
		}}},
		{"map(1, fn(x) { x * 2 })", "first argument to `map` must be ARRAY, got INTEGER"},
		{"map([1, 2, 3], 1)", "second argument to `map` must be FUNCTION, got INTEGER"},
		{"map([1, 2, 3], fn(x, y) { x * 2 })", "function passed to `map` must have exactly 1 parameters, got 2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case *object.Array:
				assertArrayEquals(t, evaluated, expected)
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestFilterBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"filter([1, 2, 3], fn(x) { x > 1 })", &object.Array{Elements: []object.Object{
			&object.Integer{Value: 2},
			&object.Integer{Value: 3},
		}}},
		{"filter([1, 2, 3], fn(x) { x < 0 })", &object.Array{Elements: []object.Object{}}},
		{"filter(1, fn(x) { x > 1 })", "first argument to `filter` must be ARRAY, got INTEGER"},
		{"filter([1, 2, 3], 1)", "second argument to `filter` must be FUNCTION, got INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case *object.Array:
				assertArrayEquals(t, evaluated, expected)
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestReduceBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"reduce([1, 2, 3], fn(acc, x) { acc + x }, 0)", 6},
		{"reduce(1, fn(acc, x) { acc + x }, 0)", "first argument to `reduce` must be ARRAY, got INTEGER"},
		{"reduce([1, 2, 3], 1, 0)", "second argument to `reduce` must be FUNCTION, got INTEGER"},
		{"reduce([1, 2, 3], fn(acc, x, y) { acc + x }, 0)", "function passed to `reduce` must have exactly 2 parameters, got 3"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestPopBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"pop([1, 2, 3])", &object.Array{Elements: []object.Object{
			&object.Integer{Value: 1},
			&object.Integer{Value: 2},
		}}},
		{"pop([])", &object.Array{Elements: []object.Object{}}},
		{"pop(1)", "argument to `pop` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case *object.Array:
				assertArrayEquals(t, evaluated, expected)
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestSumBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"sum([1, 2, 3])", 6},
		{"sum([1, 2, 3, 4])", 10},
		{"sum(1)", "argument to `sum` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestReverseBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"reverse([1, 2, 3])", &object.Array{Elements: []object.Object{
			&object.Integer{Value: 3},
			&object.Integer{Value: 2},
			&object.Integer{Value: 1},
		}}},
		{"reverse(1)", "argument to `reverse` must be ARRAY, got INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case *object.Array:
				assertArrayEquals(t, evaluated, expected)
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				t.Fatalf("unsupported expected type %T", tt.expected)
			}
		})
	}
}

func TestEchoBuiltin(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	evaluated := testEval(`echo("Hello", 42)`)

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close stdout writer: %v", err)
	}

	output, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read stdout: %v", err)
	}

	if string(output) != "\"Hello\", 42\n" {
		t.Fatalf("wrong stdout output. got=%q", string(output))
	}

	testNullObject(t, evaluated)
}

func TestForBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		// return value is Null
		{"for([1, 2, 3], fn(x) { x * 2 })", nil},
		{"for([1, 2, 3], fn(x) { echo(x) })", nil},
		{"for(1, fn(x) { x * 2 })", "first argument to `for` must be ARRAY, got INTEGER"},
		{"for([1, 2, 3], 1)", "second argument to `for` must be FUNCTION, got INTEGER"},
		{"let mut x = 0; for([1, 2, 3], fn(y) { x = x + y }); x", 6},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				testNullObject(t, evaluated)
			}
		})
	}
}

func TestWhileBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		// return value is Null
		{"let mut x = 0; while(fn() { x < 3 }, fn() { x = x + 1 })", nil},
		{"let mut x = 0; while(fn() { x < 3 }, fn() { x = x + 1 }); x", 3},
		{"let mut x = 0; while(fn() { x < 3 }, fn() { echo(\"Hello\"); x = x + 1 })", nil},
		{"while(1, fn(x) { x < 3 })", "first argument to `while` must be FUNCTION, got INTEGER"},
		{"while(fn() { x < 3 }, 1)", "second argument to `while` must be FUNCTION, got INTEGER"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaluated, int64(expected))
			case string:
				assertErrorMessage(t, evaluated, expected)
			default:
				testNullObject(t, evaluated)
			}
		})
	}
}

func assertArrayEquals(t *testing.T, obj object.Object, expected *object.Array) {
	t.Helper()

	result, ok := obj.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", obj, obj)
	}

	if len(result.Elements) != len(expected.Elements) {
		t.Fatalf("array has wrong number of elements. got=%d, want=%d", len(result.Elements), len(expected.Elements))
	}

	for i := range expected.Elements {
		if result.Elements[i].Type() != expected.Elements[i].Type() {
			t.Fatalf("array element %d has wrong type. got=%s, want=%s", i, result.Elements[i].Type(), expected.Elements[i].Type())
		}

		if expectedInteger, ok := expected.Elements[i].(*object.Integer); ok {
			resultInteger, ok := result.Elements[i].(*object.Integer)
			if !ok || resultInteger.Value != expectedInteger.Value {
				t.Fatalf("array element %d has wrong value. got=%v, want=%v", i, result.Elements[i], expected.Elements[i])
			}
			continue
		}

		if result.Elements[i].Inspect() != expected.Elements[i].Inspect() {
			t.Fatalf("array element %d has wrong value. got=%s, want=%s", i, result.Elements[i].Inspect(), expected.Elements[i].Inspect())
		}
	}
}

func assertErrorMessage(t *testing.T, obj object.Object, expected string) {
	t.Helper()

	errObj, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", obj, obj)
	}

	if errObj.Message != expected {
		t.Fatalf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
	}
}
