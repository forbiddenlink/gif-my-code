package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveInputStdinPiped(t *testing.T) {
	code, lang, name, err := resolveInput(nil, strings.NewReader("print('hi')\n"), false, "python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "print('hi')\n" {
		t.Errorf("code = %q", code)
	}
	if lang != "python" {
		t.Errorf("lang = %q, want python", lang)
	}
	if name != "stdin" {
		t.Errorf("name = %q, want stdin", name)
	}
}

func TestResolveInputStdinNoLangDefaultsToText(t *testing.T) {
	_, lang, _, err := resolveInput(nil, strings.NewReader("x"), false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lang != "text" {
		t.Errorf("lang = %q, want text", lang)
	}
}

func TestResolveInputStdinEmpty(t *testing.T) {
	if _, _, _, err := resolveInput(nil, strings.NewReader(""), false, "go"); err == nil {
		t.Fatal("expected error for empty stdin")
	}
}

func TestResolveInputTerminalNoFile(t *testing.T) {
	_, _, _, err := resolveInput(nil, strings.NewReader(""), true, "")
	if err == nil {
		t.Fatal("expected usage error")
	}
	if !strings.Contains(err.Error(), "no input") {
		t.Errorf("error should explain missing input, got: %v", err)
	}
}

func TestResolveInputFileUnchanged(t *testing.T) {
	p := filepath.Join(t.TempDir(), "a.go")
	if err := os.WriteFile(p, []byte("package a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, lang, name, err := resolveInput([]string{p}, strings.NewReader("ignored"), false, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != "package a\n" || lang != "go" || name != "a.go" {
		t.Errorf("got %q %q %q", code, lang, name)
	}
}
