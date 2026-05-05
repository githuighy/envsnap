package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/envsnap/internal/snapshot"
)

func writeCloneSnap(t *testing.T, dir string, name string, data map[string]string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := snapshot.Save(p, data); err != nil {
		t.Fatalf("writeCloneSnap: %v", err)
	}
	return p
}

var cloneBase = map[string]string{
	"APP_HOST":  "localhost",
	"APP_PORT":  "8080",
	"DB_PASS":   "secret",
	"LOG_LEVEL": "debug",
}

func TestRunClone_CopiesAll(t *testing.T) {
	dir := t.TempDir()
	src := writeCloneSnap(t, dir, "src.json", cloneBase)
	dst := filepath.Join(dir, "dst.json")

	if err := runClone([]string{src, dst}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, err := snapshot.Load(dst)
	if err != nil {
		t.Fatalf("load dst: %v", err)
	}
	if len(out) != len(cloneBase) {
		t.Errorf("expected %d keys, got %d", len(cloneBase), len(out))
	}
}

func TestRunClone_PrefixFilter(t *testing.T) {
	dir := t.TempDir()
	src := writeCloneSnap(t, dir, "src.json", cloneBase)
	dst := filepath.Join(dir, "dst.json")

	if err := runClone([]string{src, dst, "--prefix", "APP_"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, _ := snapshot.Load(dst)
	if len(out) != 2 {
		t.Errorf("expected 2 APP_ keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should be excluded")
	}
}

func TestRunClone_Rename(t *testing.T) {
	dir := t.TempDir()
	src := writeCloneSnap(t, dir, "src.json", cloneBase)
	dst := filepath.Join(dir, "dst.json")

	if err := runClone([]string{src, dst, "--rename", "APP_HOST=SERVICE_HOST"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out, _ := snapshot.Load(dst)
	if out["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %q", out["SERVICE_HOST"])
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Error("APP_HOST should be absent after rename")
	}
}

func TestRunClone_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	src := writeCloneSnap(t, dir, "src.json", cloneBase)
	dst := filepath.Join(dir, "dst.json")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runClone([]string{src, dst, "--format", "json"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result["dst"] != dst {
		t.Errorf("expected dst=%q in JSON, got %v", dst, result["dst"])
	}
}

func TestRunClone_NoArgs_ReturnsError(t *testing.T) {
	if err := runClone([]string{}); err == nil {
		t.Error("expected error for missing arguments")
	}
}

func TestRunClone_MissingSource_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	if err := runClone([]string{filepath.Join(dir, "missing.json"), filepath.Join(dir, "dst.json")}); err == nil {
		t.Error("expected error for missing source file")
	}
}
