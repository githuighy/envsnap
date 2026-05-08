package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeReorderSnap(t *testing.T, dir string, vars map[string]string) string {
	t.Helper()
	path := filepath.Join(dir, "reorder_input.json")
	if err := snapshot.Save(snapshot.Snapshot{Vars: vars}, path); err != nil {
		t.Fatalf("save snap: %v", err)
	}
	return path
}

func TestRunReorder_Alpha(t *testing.T) {
	dir := t.TempDir()
	path := writeReorderSnap(t, dir, map[string]string{
		"ZEBRA": "1", "ALPHA": "2", "MIDDLE": "3",
	})
	if err := runReorder([]string{path, "--alpha"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunReorder_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	path := writeReorderSnap(t, dir, map[string]string{
		"B_KEY": "b", "A_KEY": "a",
	})
	out := filepath.Join(dir, "out.json")
	if err := runReorder([]string{path, "--alpha", "--out", out}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(out)
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("output not valid JSON: %v", err)
	}
	if m["A_KEY"] != "a" || m["B_KEY"] != "b" {
		t.Errorf("unexpected values: %v", m)
	}
}

func TestRunReorder_Priority(t *testing.T) {
	dir := t.TempDir()
	path := writeReorderSnap(t, dir, map[string]string{
		"DB_HOST": "db", "APP_PORT": "80", "MISC": "x",
	})
	if err := runReorder([]string{path, "--alpha", "--priority", "DB_,APP_"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunReorder_NoArgs(t *testing.T) {
	if err := runReorder([]string{}); err == nil {
		t.Error("expected error for no args")
	}
}

func TestRunReorder_MissingFile(t *testing.T) {
	if err := runReorder([]string{"/nonexistent/snap.json"}); err == nil {
		t.Error("expected error for missing file")
	}
}
