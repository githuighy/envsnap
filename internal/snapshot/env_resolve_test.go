package snapshot

import (
	"os"
	"testing"
)

func TestResolve_ExpandsInternalRef(t *testing.T) {
	snap := Snapshot{
		"BASE": "/usr/local",
		"BIN":  "${BASE}/bin",
	}
	out, err := Resolve(snap, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["BIN"]; got != "/usr/local/bin" {
		t.Errorf("BIN = %q, want /usr/local/bin", got)
	}
}

func TestResolve_FallsBackToOsEnv(t *testing.T) {
	os.Setenv("_ENVSNAP_TEST_HOME", "/home/tester")
	defer os.Unsetenv("_ENVSNAP_TEST_HOME")

	snap := Snapshot{
		"CONFIG": "${_ENVSNAP_TEST_HOME}/.config",
	}
	out, err := Resolve(snap, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["CONFIG"]; got != "/home/tester/.config" {
		t.Errorf("CONFIG = %q, want /home/tester/.config", got)
	}
}

func TestResolve_AllowMissing_LeavesUnresolved(t *testing.T) {
	snap := Snapshot{
		"VAL": "${UNDEFINED_XYZ}/suffix",
	}
	opts := DefaultResolveOptions()
	opts.AllowMissing = true

	out, err := Resolve(snap, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// os.Expand replaces ${UNDEFINED_XYZ} with "", so value becomes "/suffix".
	if got := out["VAL"]; got != "/suffix" {
		t.Errorf("VAL = %q, want /suffix", got)
	}
}

func TestResolve_NoRefs_Unchanged(t *testing.T) {
	snap := Snapshot{
		"PLAIN": "hello",
		"NUM":   "42",
	}
	out, err := Resolve(snap, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range snap {
		if out[k] != v {
			t.Errorf("%s = %q, want %q", k, out[k], v)
		}
	}
}

func TestResolve_ChainedRefs(t *testing.T) {
	snap := Snapshot{
		"A": "alpha",
		"B": "${A}-beta",
		"C": "${B}-gamma",
	}
	out, err := Resolve(snap, DefaultResolveOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["C"]; got != "alpha-beta-gamma" {
		t.Errorf("C = %q, want alpha-beta-gamma", got)
	}
}

func TestResolve_MaxDepthExceeded(t *testing.T) {
	// Construct a self-referencing-like chain longer than MaxDepth.
	snap := Snapshot{
		"A": "${B}",
		"B": "${C}",
		"C": "${D}",
		"D": "${E}",
		"E": "${F}",
		"F": "${G}",
		"G": "${H}",
		"H": "${I}",
		"I": "${J}",
		"J": "${K}",
		"K": "${A}", // cycle
	}
	opts := DefaultResolveOptions()
	opts.MaxDepth = 5

	_, err := Resolve(snap, opts)
	if err == nil {
		t.Fatal("expected error for max depth exceeded, got nil")
	}
}
