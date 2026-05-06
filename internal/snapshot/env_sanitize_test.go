package snapshot

import (
	"testing"
)

var baseSanitizeSnap = Snapshot{
	"APP_NAME":  "  hello   world  ",
	"DB_PASS":   "\x00secret\x1F",
	"API_KEY":   `"my-api-key"`,
	"SAFE_KEY":  "normal-value",
	"SKIP_ME":   "\x01should-stay",
}

func TestSanitize_StripControlChars(t *testing.T) {
	opts := DefaultSanitizeOptions()
	out := Sanitize(baseSanitizeSnap, opts)

	if got := out["DB_PASS"]; got != "secret" {
		t.Errorf("expected control chars stripped, got %q", got)
	}
	if got := out["SAFE_KEY"]; got != "normal-value" {
		t.Errorf("safe key should be unchanged, got %q", got)
	}
}

func TestSanitize_TrimQuotes(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.TrimQuotes = true
	out := Sanitize(baseSanitizeSnap, opts)

	if got := out["API_KEY"]; got != "my-api-key" {
		t.Errorf("expected quotes trimmed, got %q", got)
	}
	if got := out["SAFE_KEY"]; got != "normal-value" {
		t.Errorf("unquoted value should be unchanged, got %q", got)
	}
}

func TestSanitize_CollapseWhitespace(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.CollapseWhitespace = true
	out := Sanitize(baseSanitizeSnap, opts)

	if got := out["APP_NAME"]; got != " hello world " {
		t.Errorf("expected collapsed whitespace, got %q", got)
	}
}

func TestSanitize_ByPrefix(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.Prefixes = []string{"DB_"}
	out := Sanitize(baseSanitizeSnap, opts)

	if got := out["DB_PASS"]; got != "secret" {
		t.Errorf("expected DB_ key sanitized, got %q", got)
	}
	// SKIP_ME has control char but no DB_ prefix — should be unchanged
	if got := out["SKIP_ME"]; got != "\x01should-stay" {
		t.Errorf("expected non-prefixed key untouched, got %q", got)
	}
}

func TestSanitize_ExcludeKeys(t *testing.T) {
	opts := DefaultSanitizeOptions()
	opts.ExcludeKeys = []string{"DB_PASS"}
	out := Sanitize(baseSanitizeSnap, opts)

	if got := out["DB_PASS"]; got != "\x00secret\x1F" {
		t.Errorf("excluded key should not be sanitized, got %q", got)
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	original := Snapshot{"KEY": "\x00value"}
	opts := DefaultSanitizeOptions()
	_ = Sanitize(original, opts)

	if original["KEY"] != "\x00value" {
		t.Error("Sanitize must not mutate the input snapshot")
	}
}

func TestSanitize_DefaultOptions(t *testing.T) {
	opts := DefaultSanitizeOptions()
	if !opts.StripControlChars {
		t.Error("expected StripControlChars to be true by default")
	}
	if opts.TrimQuotes {
		t.Error("expected TrimQuotes to be false by default")
	}
	if opts.CollapseWhitespace {
		t.Error("expected CollapseWhitespace to be false by default")
	}
}
