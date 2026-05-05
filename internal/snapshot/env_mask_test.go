package snapshot

import (
	"strings"
	"testing"
)

var baseMaskSnap = Snapshot{
	"DATABASE_URL":  "postgres://user:secret@host/db",
	"API_KEY":       "abcdef123456",
	"APP_NAME":      "envsnap",
	"SECRET_TOKEN":  "topsecret",
	"LOG_LEVEL":     "info",
}

func TestMask_ExplicitKeys(t *testing.T) {
	opts := MaskOptions{Keys: []string{"API_KEY", "DATABASE_URL"}}
	out := Mask(baseMaskSnap, opts)

	if out["API_KEY"] != "***" {
		t.Errorf("expected API_KEY masked, got %q", out["API_KEY"])
	}
	if out["DATABASE_URL"] != "***" {
		t.Errorf("expected DATABASE_URL masked, got %q", out["DATABASE_URL"])
	}
	if out["APP_NAME"] != "envsnap" {
		t.Errorf("expected APP_NAME unchanged, got %q", out["APP_NAME"])
	}
}

func TestMask_ByPrefix(t *testing.T) {
	opts := MaskOptions{Prefixes: []string{"SECRET_"}}
	out := Mask(baseMaskSnap, opts)

	if out["SECRET_TOKEN"] != "***" {
		t.Errorf("expected SECRET_TOKEN masked, got %q", out["SECRET_TOKEN"])
	}
	if out["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL unchanged, got %q", out["LOG_LEVEL"])
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	opts := MaskOptions{
		Keys:        []string{"API_KEY"},
		Placeholder: "[REDACTED]",
	}
	out := Mask(baseMaskSnap, opts)

	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected custom placeholder, got %q", out["API_KEY"])
	}
}

func TestMask_ShowLength(t *testing.T) {
	opts := MaskOptions{
		Keys:       []string{"API_KEY"},
		ShowLength: true,
	}
	out := Mask(baseMaskSnap, opts)

	// "abcdef123456" is 12 chars
	if !strings.Contains(out["API_KEY"], "(12)") {
		t.Errorf("expected length suffix in %q", out["API_KEY"])
	}
}

func TestMask_VisibleChars(t *testing.T) {
	opts := MaskOptions{
		Keys:         []string{"API_KEY"},
		VisibleChars: 4,
	}
	out := Mask(baseMaskSnap, opts)

	// original "abcdef123456" — last 4 chars "3456"
	if !strings.HasSuffix(out["API_KEY"], "3456") {
		t.Errorf("expected trailing chars visible, got %q", out["API_KEY"])
	}
}

func TestMask_UnchangedKeys(t *testing.T) {
	opts := MaskOptions{Keys: []string{"API_KEY"}}
	out := Mask(baseMaskSnap, opts)

	if out["LOG_LEVEL"] != baseMaskSnap["LOG_LEVEL"] {
		t.Errorf("unmasked key should be unchanged")
	}
}

func TestMaskedKeys_ReturnsList(t *testing.T) {
	opts := MaskOptions{Prefixes: []string{"SECRET_", "API_"}}
	keys := MaskedKeys(baseMaskSnap, opts)

	found := map[string]bool{}
	for _, k := range keys {
		found[k] = true
	}
	if !found["API_KEY"] || !found["SECRET_TOKEN"] {
		t.Errorf("expected API_KEY and SECRET_TOKEN in masked keys, got %v", keys)
	}
}
