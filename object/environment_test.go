package object

import "testing"

func TestEnvironmentSetAndGet(t *testing.T) {
	env := NewEnvironment()

	returned := env.Set("a", &Integer{Value: 1})
	if integer, ok := returned.(*Integer); !ok || integer.Value != 1 {
		t.Fatalf("Set did not return the stored value. got=%+v", returned)
	}

	obj, ok := env.Get("a")
	if !ok {
		t.Fatalf("Get(\"a\") returned ok=false")
	}
	if integer, ok := obj.(*Integer); !ok || integer.Value != 1 {
		t.Fatalf("Get(\"a\") wrong value. got=%+v", obj)
	}

	if _, ok := env.Get("missing"); ok {
		t.Fatalf("Get(\"missing\") returned ok=true")
	}
}

func TestEnvironmentOverwrite(t *testing.T) {
	env := NewEnvironment()
	env.Set("a", &Integer{Value: 1})
	env.Set("a", &Integer{Value: 2})

	obj, ok := env.Get("a")
	if !ok {
		t.Fatalf("Get(\"a\") returned ok=false")
	}
	if integer, ok := obj.(*Integer); !ok || integer.Value != 2 {
		t.Fatalf("Get(\"a\") wrong value. got=%+v", obj)
	}
}

func TestEnclosedEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("outerOnly", &Integer{Value: 1})
	outer.Set("shadowed", &Integer{Value: 2})

	inner := NewEnclosedEnvironment(outer)
	inner.Set("shadowed", &Integer{Value: 3})
	inner.Set("innerOnly", &Integer{Value: 4})

	tests := []struct {
		name     string
		expected int64
	}{
		{"outerOnly", 1},
		{"shadowed", 3},
		{"innerOnly", 4},
	}

	for _, tt := range tests {
		obj, ok := inner.Get(tt.name)
		if !ok {
			t.Fatalf("inner.Get(%q) returned ok=false", tt.name)
		}
		integer, ok := obj.(*Integer)
		if !ok {
			t.Fatalf("inner.Get(%q) not *Integer. got=%T", tt.name, obj)
		}
		if integer.Value != tt.expected {
			t.Errorf("inner.Get(%q) wrong. got=%d, want=%d", tt.name, integer.Value, tt.expected)
		}
	}

	if _, ok := outer.Get("innerOnly"); ok {
		t.Errorf("outer.Get(\"innerOnly\") returned ok=true; inner bindings must not leak outward")
	}

	obj, _ := outer.Get("shadowed")
	if integer, ok := obj.(*Integer); !ok || integer.Value != 2 {
		t.Errorf("outer binding was mutated by inner Set. got=%+v", obj)
	}
}
