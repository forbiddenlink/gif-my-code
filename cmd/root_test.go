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

func TestOverflowWarning(t *testing.T) {
	if got := overflowWarning(0, 36); got != "" {
		t.Errorf("no overflow should not warn, got %q", got)
	}
	got := overflowWarning(7, 200)
	for _, want := range []string{"7 line", "--padding 200", "cut off"} {
		if !strings.Contains(got, want) {
			t.Errorf("warning %q missing %q", got, want)
		}
	}
}

func TestCheckThemeKnown(t *testing.T) {
	if err := checkTheme("dracula"); err != nil {
		t.Fatalf("dracula should be valid: %v", err)
	}
}

func TestCheckThemeUnknownSuggests(t *testing.T) {
	err := checkTheme("drakula")
	if err == nil {
		t.Fatal("expected error for unknown theme")
	}
	msg := err.Error()
	if !strings.Contains(msg, `unknown theme "drakula"`) {
		t.Errorf("message missing theme name: %s", msg)
	}
	if !strings.Contains(msg, "dracula") {
		t.Errorf("message should suggest dracula: %s", msg)
	}
	if !strings.Contains(msg, "--list-themes") {
		t.Errorf("message should point at --list-themes: %s", msg)
	}
}

func TestPrintThemesOnePerLine(t *testing.T) {
	var b strings.Builder
	printThemes(&b)
	lines := strings.Split(strings.TrimSpace(b.String()), "\n")
	if len(lines) < 50 {
		t.Errorf("expected at least 50 themes, got %d", len(lines))
	}
}
