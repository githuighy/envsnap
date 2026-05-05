package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func writeTransformSnap(t *testing.T, dir string, data map[string]string) string {
	t.Helper()
	p := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(p, data); err != nil {
		t.Fatalf("writeTransformSnap: %v", err)
	}
	return p
}

func TestRunTransform_Uppercase(t *testing.T) {
	dir := t.TempDir()
	p := writeTransformSnap(t, dir, map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTransform([]string{p, "uppercase"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	buf := new(strings.Builder)
	var b [4096]byte
	for {
		n, e := r.Read(b[:])
		buf.Write(b[:n])
		if e != nil {
			break
		}
	}
	if !strings.Contains(buf.String(), "LOCALHOST") {
		t.Errorf("expected LOCALHOST in output, got: %s", buf.String())
	}
}

func TestRunTransform_JSONFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeTransformSnap(t, dir, map[string]string{"KEY": "  hello  "})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := runTransform([]string{p, "trim", "--keys=KEY", "--format=json"})
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]string
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		t.Fatalf("json decode: %v", err)
	}
	if result["KEY"] != "hello" {
		t.Errorf("KEY: want %q, got %q", "hello", result["KEY"])
	}
}

func TestRunTransform_UnknownOp(t *testing.T) {
	dir := t.TempDir()
	p := writeTransformSnap(t, dir, map[string]string{"X": "y"})
	err := runTransform([]string{p, "reverse"})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestRunTransform_NoArgs(t *testing.T) {
	err := runTransform([]string{})
	if err == nil {
		t.Fatal("expected error for missing args")
	}
}

func TestRunTransform_MissingFile(t *testing.T) {
	err := runTransform([]string{"/nonexistent/snap.json", "uppercase"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
