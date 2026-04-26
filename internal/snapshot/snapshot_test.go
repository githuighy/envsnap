package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCapture_ContainsEnvVar(t *testing.T) {
	const key = "ENVSNAP_TEST_VAR"
	const val = "hello"
	t.Setenv(key, val)

	s := Capture("test")

	if s.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", s.Name)
	}
	if got, ok := s.Vars[key]; !ok || got != val {
		t.Errorf("expected var %s=%s, got %s", key, val, got)
	}
	if s.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	t.Setenv("ENVSNAP_RT", "roundtrip")

	orig := Capture("rt")

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := Save(orig, tmp); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.Name != orig.Name {
		t.Errorf("name mismatch: %q vs %q", loaded.Name, orig.Name)
	}
	if loaded.Vars["ENVSNAP_RT"] != "roundtrip" {
		t.Errorf("var mismatch for ENVSNAP_RT")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.json"))
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
