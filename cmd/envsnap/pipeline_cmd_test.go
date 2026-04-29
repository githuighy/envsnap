package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writePipelineSnap(t *testing.T, dir string, name string, s snapshot.Snapshot) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := snapshot.Save(s, path); err != nil {
		t.Fatalf("save snap: %v", err)
	}
	return path
}

func TestRunPipeline_FilterStep(t *testing.T) {
	dir := t.TempDir()
	snap := snapshot.Snapshot{
		"APP_ENV":  "prod",
		"DB_HOST":  "localhost",
		"APP_PORT": "8080",
	}
	path := writePipelineSnap(t, dir, "snap.json", snap)

	// Capture stdout by redirecting; here we just check no error.
	err := runPipeline([]string{path}, "filter:APP", "text", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPipeline_RedactStep(t *testing.T) {
	dir := t.TempDir()
	snap := snapshot.Snapshot{
		"API_KEY": "super-secret",
		"APP_ENV": "staging",
	}
	path := writePipelineSnap(t, dir, "snap.json", snap)
	outPath := filepath.Join(dir, "out.json")

	err := runPipeline([]string{path}, "redact", "text", outPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := snapshot.Load(outPath)
	if err != nil {
		t.Fatalf("load output: %v", err)
	}
	if result["API_KEY"] == "super-secret" {
		t.Error("expected API_KEY to be redacted")
	}
}

func TestRunPipeline_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	snap := snapshot.Snapshot{"FOO": "bar"}
	path := writePipelineSnap(t, dir, "snap.json", snap)

	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runPipeline([]string{path}, "", "json", "")
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	if err := json.NewDecoder(r).Decode(&out); err != nil {
		t.Fatalf("decode json: %v", err)
	}
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", out["FOO"])
	}
}

func TestRunPipeline_UnknownStep_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := writePipelineSnap(t, dir, "snap.json", snapshot.Snapshot{"X": "1"})
	err := runPipeline([]string{path}, "nonexistent", "text", "")
	if err == nil {
		t.Fatal("expected error for unknown step")
	}
}

func TestRunPipeline_NoArgs_ReturnsError(t *testing.T) {
	err := runPipeline([]string{}, "redact", "text", "")
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}
