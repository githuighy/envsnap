package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeScopeSnap(t *testing.T, dir string, data map[string]string) string {
	t.Helper()
	p := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(data, p); err != nil {
		t.Fatalf("save snap: %v", err)
	}
	return p
}

func TestRunScope_AddsPrefix(t *testing.T) {
	dir := t.TempDir()
	p := writeScopeSnap(t, dir, map[string]string{"DB_HOST": "localhost", "APP_ENV": "dev"})

	err := runScope([]string{p, "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunScope_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeScopeSnap(t, dir, map[string]string{"DB_HOST": "localhost"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runScope([]string{p, "staging", "--format", "json"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if _, ok := out["STAGING_DB_HOST"]; !ok {
		t.Error("expected STAGING_DB_HOST in JSON output")
	}
}

func TestRunScope_WritesFile(t *testing.T) {
	dir := t.TempDir()
	p := writeScopeSnap(t, dir, map[string]string{"KEY": "val"})
	out := filepath.Join(dir, "scoped.json")

	err := runScope([]string{p, "test", "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := snapshot.Load(out)
	if err != nil {
		t.Fatalf("load scoped: %v", err)
	}
	if loaded["TEST_KEY"] != "val" {
		t.Errorf("expected TEST_KEY=val, got %q", loaded["TEST_KEY"])
	}
}

func TestRunScope_Unscope(t *testing.T) {
	dir := t.TempDir()
	p := writeScopeSnap(t, dir, map[string]string{"PROD_DB_HOST": "localhost"})
	out := filepath.Join(dir, "unscoped.json")

	err := runScope([]string{p, "prod", "--unscope", "--out", out})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := snapshot.Load(out)
	if err != nil {
		t.Fatalf("load unscoped: %v", err)
	}
	if loaded["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", loaded["DB_HOST"])
	}
}

func TestRunScope_NoArgs_ReturnsError(t *testing.T) {
	err := runScope([]string{})
	if err == nil || !strings.Contains(err.Error(), "usage") {
		t.Errorf("expected usage error, got %v", err)
	}
}
