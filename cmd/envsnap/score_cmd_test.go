package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeScoreSnap(t *testing.T, dir string, data map[string]string) string {
	t.Helper()
	path := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(data, path); err != nil {
		t.Fatalf("save snap: %v", err)
	}
	return path
}

func TestRunScore_TextFormat(t *testing.T) {
	dir := t.TempDir()
	path := writeScoreSnap(t, dir, map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runScore([]string{path, "--required", "APP_HOST,APP_PORT"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := make([]byte, 512)
	n, _ := r.Read(buf)
	out := string(buf[:n])
	if len(out) == 0 {
		t.Error("expected output, got empty string")
	}
}

func TestRunScore_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	path := writeScoreSnap(t, dir, map[string]string{
		"APP_HOST": "localhost",
	})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runScore([]string{path, "--format", "json"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result struct {
		Score      int      `json:"score"`
		Deductions []string `json:"deductions"`
	}
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if result.Score < 0 || result.Score > 100 {
		t.Errorf("score out of range: %d", result.Score)
	}
}

func TestRunScore_NoArgs(t *testing.T) {
	err := runScore([]string{})
	if err == nil {
		t.Error("expected error for missing args")
	}
}

func TestRunScore_MissingFile(t *testing.T) {
	err := runScore([]string{"/nonexistent/snap.json"})
	if err == nil {
		t.Error("expected error for missing file")
	}
}
