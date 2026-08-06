package repl

import (
	"bytes"
	"strings"
	"testing"
)

func TestStartEvaluatesInput(t *testing.T) {
	in := strings.NewReader("let a = 5;\na + 5;\n\"hello\" + \" world\";\n[1, 2][1];\n")
	out := &bytes.Buffer{}

	Start(in, out)

	expected := []string{"10", "hello world", "2"}
	got := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(got) != len(expected) {
		t.Fatalf("wrong number of output lines. got=%d (%q), want=%d", len(got), out.String(), len(expected))
	}
	for i, want := range expected {
		if got[i] != want {
			t.Errorf("output line %d wrong. got=%q, want=%q", i, got[i], want)
		}
	}
}

func TestStartKeepsEnvironmentBetweenLines(t *testing.T) {
	in := strings.NewReader("let counter = fn(x) { x + 1; };\ncounter(1);\n")
	out := &bytes.Buffer{}

	Start(in, out)

	if !strings.Contains(out.String(), "2") {
		t.Errorf("bindings not preserved across lines. got=%q", out.String())
	}
}

func TestStartPrintsParserErrors(t *testing.T) {
	in := strings.NewReader("let = 5;\n")
	out := &bytes.Buffer{}

	Start(in, out)

	if !strings.Contains(out.String(), "Parser errors found") {
		t.Errorf("parser errors not reported. got=%q", out.String())
	}
}

func TestStartReturnsOnEmptyInput(t *testing.T) {
	out := &bytes.Buffer{}

	Start(strings.NewReader(""), out)

	if out.String() != "" {
		t.Errorf("expected no output for empty input. got=%q", out.String())
	}
}

func TestPrintParserErrors(t *testing.T) {
	out := &bytes.Buffer{}

	printParserErrors(out, []string{"first", "second"})

	want := "Parser errors found\tfirst\n\tsecond\n"
	if out.String() != want {
		t.Errorf("printParserErrors wrong. got=%q, want=%q", out.String(), want)
	}
}
