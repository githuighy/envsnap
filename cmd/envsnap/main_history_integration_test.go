package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func TestHistoryRoundTrip(t *testing.T) {
	dir := t.TempDir()
	histDir := filepath.Join(dir, ".envsnap", "history")

	// Simulate two deployments
	snap1 := map[string]string{
		"APP_ENV":  "staging",
		"DB_HOST":  "db1.internal",
		"LOG_LEVEL": "info",
	}
	snap2 := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "db2.internal",
		"LOG_LEVEL": "warn",
		"NEW_RELIC": "enabled",
	}

	store := snapshot.NewHistoryStore(histDir)
	idA, err := store.Add(snap1, "v1.0")
	if err != nil {
		t.Fatalf("Add snap1: %v", err)
	}
	idB, err := store.Add(snap2, "v1.1")
	if err != nil {
		t.Fatalf("Add snap2: %v", err)
	}

	// Verify list order
	entries, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Verify diff via cmd
	old := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	err = runHistoryDiff(histDir, idA, idB, "text")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("runHistoryDiff: %v", err)
	}

	// Verify get by id returns correct vars
	entryA, err := store.Get(idA)
	if err != nil {
		t.Fatalf("Get idA: %v", err)
	}
	if entryA.Vars["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %q", entryA.Vars["APP_ENV"])
	}
	entryB, err := store.Get(idB)
	if err != nil {
		t.Fatalf("Get idB: %v", err)
	}
	if entryB.Vars["NEW_RELIC"] != "enabled" {
		t.Errorf("expected NEW_RELIC=enabled, got %q", entryB.Vars["NEW_RELIC"])
	}
}
