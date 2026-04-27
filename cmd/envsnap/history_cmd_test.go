package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeSnapFile(t *testing.T, dir string, vars map[string]string) string {
	t.Helper()
	f := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(vars, f); err != nil {
		t.Fatalf("writeSnapFile: %v", err)
	}
	return f
}

func TestRunHistoryAdd(t *testing.T) {
	dir := t.TempDir()
	snapFile := writeSnapFile(t, dir, map[string]string{"FOO": "bar"})
	histDir := filepath.Join(dir, "history")

	if err := runHistoryAdd([]string{snapFile}, "my-label", histDir); err != nil {
		t.Fatalf("runHistoryAdd: %v", err)
	}

	store := snapshot.NewHistoryStore(histDir)
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Label != "my-label" {
		t.Errorf("expected label 'my-label', got %q", entries[0].Label)
	}
}

func TestRunHistoryAdd_NoArgs(t *testing.T) {
	err := runHistoryAdd([]string{}, "", t.TempDir())
	if err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestRunHistoryList_JSON(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewHistoryStore(dir)
	_, _ = store.Add(map[string]string{"A": "1"}, "alpha")
	_, _ = store.Add(map[string]string{"B": "2"}, "beta")

	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runHistoryList(dir, "json")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("runHistoryList: %v", err)
	}

	var entries []map[string]interface{}
	if err := json.NewDecoder(r).Decode(&entries); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}
}

func TestRunHistoryDiff(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewHistoryStore(dir)
	idA, _ := store.Add(map[string]string{"X": "1", "Y": "old"}, "")
	idB, _ := store.Add(map[string]string{"X": "1", "Y": "new", "Z": "added"}, "")

	if err := runHistoryDiff(dir, idA, idB, "text"); err != nil {
		t.Fatalf("runHistoryDiff: %v", err)
	}
}
