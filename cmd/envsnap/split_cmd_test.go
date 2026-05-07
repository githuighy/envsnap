package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeSplitSnap(t *testing.T, dir string, s snapshot.Snapshot) string {
	t.Helper()
	p := filepath.Join(dir, "split_input.json")
	if err := snapshot.Save(s, p); err != nil {
		t.Fatalf("writeSplitSnap: %v", err)
	}
	return p
}

func TestRunSplit_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeSplitSnap(t, dir, snapshot.Snapshot{
		"DB_HOST":   "localhost",
		"DB_PORT":   "5432",
		"APP_DEBUG": "true",
	})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSplit([]string{p, "--bucket", "db=DB_", "--format", "json"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]map[string]string
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if out["db"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST in db bucket")
	}
}

func TestRunSplit_NoArgs(t *testing.T) {
	if err := runSplit([]string{}); err == nil {
		t.Fatal("expected error with no args")
	}
}

func TestRunSplit_MissingFile(t *testing.T) {
	err := runSplit([]string{"/nonexistent/snap.json", "--bucket", "db=DB_"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRunSplit_InvalidBucketFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeSplitSnap(t, dir, snapshot.Snapshot{"DB_HOST": "localhost"})
	err := runSplit([]string{p, "--bucket", "noequalssign"})
	if err == nil {
		t.Fatal("expected error for invalid bucket format")
	}
}

func TestRunSplit_RemainderWritten(t *testing.T) {
	dir := t.TempDir()
	p := writeSplitSnap(t, dir, snapshot.Snapshot{
		"DB_HOST":   "localhost",
		"APP_DEBUG": "true",
	})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runSplit([]string{p, "--bucket", "db=DB_", "--remainder", "other", "--format", "json"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]map[string]string
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if out["other"]["APP_DEBUG"] != "true" {
		t.Errorf("expected APP_DEBUG in other (remainder) bucket")
	}
}
