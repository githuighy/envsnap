package snapshot

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHistory_AddAndList(t *testing.T) {
	dir := t.TempDir()
	store := NewHistoryStore(dir)

	snap1 := map[string]string{"FOO": "bar", "BAZ": "qux"}
	id1, err := store.Add(snap1, "first")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if id1 == "" {
		t.Fatal("expected non-empty id")
	}

	snap2 := map[string]string{"FOO": "updated"}
	_, err = store.Add(snap2, "second")
	if err != nil {
		t.Fatalf("Add second: %v", err)
	}

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Label != "first" {
		t.Errorf("expected first label 'first', got %q", entries[0].Label)
	}
}

func TestHistory_Get(t *testing.T) {
	dir := t.TempDir()
	store := NewHistoryStore(dir)

	snap := map[string]string{"KEY": "value"}
	id, err := store.Add(snap, "test-label")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	entry, err := store.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if entry.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", entry.Label)
	}
	if entry.Vars["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", entry.Vars["KEY"])
	}
}

func TestHistory_Get_Missing(t *testing.T) {
	dir := t.TempDir()
	store := NewHistoryStore(dir)

	_, err := store.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing entry")
	}
}

func TestHistory_List_Empty(t *testing.T) {
	dir := t.TempDir()
	store := NewHistoryStore(dir)

	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestHistory_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "nested", "history")
	store := NewHistoryStore(dir)

	_, err := store.Add(map[string]string{"A": "1"}, "")
	if err != nil {
		t.Fatalf("Add: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}
