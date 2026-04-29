package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeChainSnap(t *testing.T, dir string, name string, snap snapshot.Snapshot) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := snapshot.Save(snap, path); err != nil {
		t.Fatalf("save %s: %v", name, err)
	}
	return path
}

func TestRunChain_MergesLayers(t *testing.T) {
	dir := t.TempDir()
	high := writeChainSnap(t, dir, "high.json", snapshot.Snapshot{"A": "high", "B": "high"})
	low := writeChainSnap(t, dir, "low.json", snapshot.Snapshot{"A": "low", "C": "low"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runChain([]string{high, low}, "env")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 512)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if !contains(out, "A=high") {
		t.Errorf("expected A=high in output, got: %s", out)
	}
	if !contains(out, "C=low") {
		t.Errorf("expected C=low in output, got: %s", out)
	}
}

func TestRunChain_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeChainSnap(t, dir, "snap.json", snapshot.Snapshot{"FOO": "bar"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := runChain([]string{p}, "json")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	json.NewDecoder(r).Decode(&result)
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", result["FOO"])
	}
}

func TestRunChain_NoArgs(t *testing.T) {
	if err := runChain([]string{}, "text"); err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunChain_MissingFile(t *testing.T) {
	if err := runChain([]string{"/nonexistent/snap.json"}, "text"); err == nil {
		t.Error("expected error for missing file")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
