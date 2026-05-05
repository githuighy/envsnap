package snapshot

import (
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	snap := map[string]string{"path": "/usr/bin", "home": "/root"}
	opts := NormalizeOptions{UppercaseKeys: true}
	out := Normalize(snap, opts)
	if _, ok := out["PATH"]; !ok {
		t.Error("expected PATH key after uppercase normalization")
	}
	if _, ok := out["path"]; ok {
		t.Error("original lowercase key should not be present")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	snap := map[string]string{"KEY": "  hello  ", "OTHER": "\tworld\n"}
	opts := NormalizeOptions{TrimValues: true}
	out := Normalize(snap, opts)
	if out["KEY"] != "hello" {
		t.Errorf("expected trimmed value 'hello', got %q", out["KEY"])
	}
	if out["OTHER"] != "world" {
		t.Errorf("expected trimmed value 'world', got %q", out["OTHER"])
	}
}

func TestNormalize_RemoveEmptyValues(t *testing.T) {
	snap := map[string]string{"PRESENT": "value", "EMPTY": "", "SPACES": "   "}
	opts := NormalizeOptions{TrimValues: true, RemoveEmptyValues: true}
	out := Normalize(snap, opts)
	if _, ok := out["EMPTY"]; ok {
		t.Error("expected EMPTY key to be removed")
	}
	if _, ok := out["SPACES"]; ok {
		t.Error("expected SPACES key to be removed after trim")
	}
	if out["PRESENT"] != "value" {
		t.Error("expected PRESENT key to remain")
	}
}

func TestNormalize_SanitizeKeys(t *testing.T) {
	snap := map[string]string{"my-key": "a", "1START": "b", "ok_key": "c"}
	opts := NormalizeOptions{SanitizeKeys: true}
	out := Normalize(snap, opts)
	if _, ok := out["my_key"]; !ok {
		t.Error("expected hyphen replaced with underscore")
	}
	if _, ok := out["_START"]; !ok {
		t.Error("expected leading digit replaced with underscore")
	}
	if _, ok := out["ok_key"]; !ok {
		t.Error("expected valid key to remain unchanged")
	}
}

func TestNormalize_DefaultOptions(t *testing.T) {
	snap := map[string]string{"my-var": "  trimmed  ", "ALREADY": "fine"}
	out := Normalize(snap, DefaultNormalizeOptions())
	if _, ok := out["MY_VAR"]; !ok {
		t.Error("expected MY_VAR after default normalization")
	}
	if out["MY_VAR"] != "trimmed" {
		t.Errorf("expected trimmed value, got %q", out["MY_VAR"])
	}
}

func TestNormalize_NoOptions_PassesThrough(t *testing.T) {
	snap := map[string]string{"key": "  value  "}
	out := Normalize(snap, NormalizeOptions{})
	if out["key"] != "  value  " {
		t.Error("expected value unchanged when no options set")
	}
}
