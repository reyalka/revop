package evaluator

import (
	"testing"

	"github.com/reyalka/revop/lexer"
	"github.com/reyalka/revop/object"
	"github.com/reyalka/revop/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"2 ^ 3", 8},
		{"2 ^ 3 ^ 2", 512},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 >= 1", true},
		{"2 <= 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestElseIfEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 10 } else if (true) { 20 } else { 30 }", 10},
		{"if (false) { 10 } else if (true) { 20 } else { 30 }", 20},
		{"if (false) { 10 } else if (false) { 20 } else { 30 }", 30},
		{"if (false) { 10 } else if (false) { 20 }", nil},
		{"if (false) { 10 } else if (true) { 20 }", 20},
		{"if (1 < 2) { 10 } else if (2 < 3) { 20 } else { 30 }", 10},
		{"if (1 > 2) { 10 } else if (2 < 3) { 20 } else { 30 }", 20},
		{"if (1 > 2) { 10 } else if (2 > 3) { 20 } else { 30 }", 30},
		{"if (1 > 2) { 10 } else if (2 > 3) { 20 } else if (3 < 4) { 30 } else { 40 }", 30},
		{"if (1 > 2) { 10 } else if (2 > 3) { 20 } else if (3 > 4) { 30 } else { 40 }", 40},
		{"if (1 > 2) { 10 } else if (2 > 3) { 20 } else if (3 > 4) { 30 }", nil},
		{"let x = if (false) { 10 } else if (true) { 20 } else { 30 }; x", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestIfWithoutElseAsExpression(t *testing.T) {
	// Statement context: missing else + false path stays NULL.
	statementTests := []string{
		"if (false) { 10 }",
		"if (false) { 10 } else if (false) { 20 }",
	}
	for _, input := range statementTests {
		evaluated := testEval(input)
		testNullObject(t, evaluated)
	}

	// Expression context with branch taken: no error even without else.
	successTests := []struct {
		input    string
		expected int64
	}{
		{"let x = if (true) { 10 }; x", 10},
		{"let x = if (false) { 10 } else if (true) { 20 }; x", 20},
	}
	for _, tt := range successTests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}

	// Expression context without else + false path: Error (issue #5).
	errorTests := []string{
		"let x = if (false) { 10 }; x",
		"let x = if (false) { 10 } else if (false) { 20 }; x",
		"let x = if (1 > 2) { 10 } else if (2 > 3) { 20 } else if (3 > 4) { 30 }; x",
	}
	for _, input := range errorTests {
		evaluated := testEval(input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for %q. got=%T(%+v)", input, evaluated, evaluated)
			continue
		}
		expectedMessage := "no else branch for if expression"
		if errObj.Message != expectedMessage {
			t.Errorf("wrong error message for %q. got=%q, want=%q", input, errObj.Message, expectedMessage)
		}
	}
}

func TestIfNonBooleanCondition(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"if (1) { 10 }", "condition for `if` must return BOOLEAN, got INTEGER"},
		{"if (1) { 10 } else { 20 }", "condition for `if` must return BOOLEAN, got INTEGER"},
		{"let x = if (1) { 10 } else { 20 }; x", "condition for `if` must return BOOLEAN, got INTEGER"},
		{"if (false) { 10 } else if (1) { 20 } else { 30 }", "condition for `if` must return BOOLEAN, got INTEGER"},
		{"if (false) { 10 } else if (1) { 20 }", "condition for `if` must return BOOLEAN, got INTEGER"},
		{`if ("a") { 10 }`, "condition for `if` must return BOOLEAN, got STRING"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for %q. got=%T(%+v)", tt.input, evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message for %q. got=%q, want=%q", tt.input, errObj.Message, tt.expectedMessage)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) { if (10 > 1) { return 10 ;} return 1; }
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBlockExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"{ 10 }", 10},
		{"{ 10; }", 10},
		{"{ 1; 2; 3; }", 3},
		{"{return 4}", 4},
		{"{return 4;}", 4},
		{"{1;2;return 4;5;}", 4},
		{
			`
			{
				if (10 > 1) {
					if (10 > 1) {
						return 10;
					}
					return 1;
				}
				return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`"Hello" + 5`,
			"type mismatch: STRING + INTEGER",
		},
		{
			`let x = 5; x + "Hello"`,
			"type mismatch: INTEGER + STRING",
		},
		{
			`#{name: "Revop"}[fn(x) { x }];`,
			"cannot use as hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let mut a = 5; a;", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestAssign(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"let mut a = 5; a = 10; a;", 10},
		{"let mut a = 5; let mut b = a; b = 10; b;", 10},
		{"a = 5;", "cannot assign to undefined variable: a"},
		{"let mut a = 5; let b = a; b = 10;", "cannot assign to immutable variable: b"},
		{`let mut a = 5; a = "Hello"; a;`, "type mismatch: cannot assign STRING to INTEGER"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "{ (x + 2); }"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosedFunctionEnvironments(t *testing.T) {
	input := `
	let mut x = 0;
	let increment = fn() {
		x = x + 1;
	};
	increment();
	increment();
	increment();
	x;
	`

	testIntegerObject(t, testEval(input), 3)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`"Hello" == "Hello"`, true},
		{`"Hello" != "Hello"`, false},
		{`"Hello" == "World"`, false},
		{`"Hello" != "World"`, true},
		{`"Hello" + " " + "World!" == "Hello World!"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		boolean, ok := tt.expected.(bool)
		if ok {
			testBooleanObject(t, evaluated, boolean)
		} else {
			t.Errorf("expected value is not boolean. got=%T (%+v)", tt.expected, tt.expected)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3, 3 == 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 4 {
		t.Fatalf("array has wrong number of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
	testBooleanObject(t, result.Elements[3], true)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			"index out of range",
		},
		{
			"[1, 2, 3][-1]",
			"index out of range",
		},
		{
			"[1, 2, 3][-4]",
			"index out of range",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestHashMapLiterals(t *testing.T) {
	input := `let two = "two";
	#{
		one: 10 - 9,
		[two]: 1 + 1,
		["thr" + "ee"]: 6 / 2,
		[4]: 4,
		[true]: 5,
		[false]: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.HashMap)
	if !ok {
		t.Fatalf("Eval didn't return HashMap. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		(&object.Boolean{Value: true}).HashKey():   5,
		(&object.Boolean{Value: false}).HashKey():  6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("HashMap has wrong number of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashMapIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`#{foo: 5}["foo"]`,
			5,
		},
		{
			`#{foo: 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; #{foo: 5}[key]`,
			5,
		},
		{
			`#{}["foo"]`,
			nil,
		},
		{
			`#{[5]: 5}[5]`,
			5,
		},
		{
			`#{[true]: 5}[true]`,
			5,
		},
		{
			`#{[false]: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		})
	}
}

func TestShortCircuitEvaluation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true || (1 / 0 == 1)", true},
		{"false && (1 / 0 == 1)", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}
