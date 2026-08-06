package repl

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type failingWriter struct {
	err error
}

func (w *failingWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

type failingReader struct {
	err error
}

func (r *failingReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

func TestStartPropagatesWriteError(t *testing.T) {
	want := errors.New("write failed")

	err := Start(strings.NewReader("1 + 1\n"), &failingWriter{err: want})
	if !errors.Is(err, want) {
		t.Fatalf("expected write error to be propagated, got=%v", err)
	}
}

func TestStartPropagatesReadError(t *testing.T) {
	want := errors.New("read failed")

	err := Start(&failingReader{err: want}, io.Discard)
	if !errors.Is(err, want) {
		t.Fatalf("expected read error to be propagated, got=%v", err)
	}
}

func TestStartReturnsNilOnEOF(t *testing.T) {
	var out strings.Builder

	if err := Start(strings.NewReader("1 + 1\n"), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "2") {
		t.Fatalf("expected evaluated result in output, got=%q", out.String())
	}
}

func TestStartReportsParserErrors(t *testing.T) {
	var out strings.Builder

	if err := Start(strings.NewReader("let x = ;\n"), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out.String(), "Parser errors found") {
		t.Fatalf("expected parser errors in output, got=%q", out.String())
	}
}
