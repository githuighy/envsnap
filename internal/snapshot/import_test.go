package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestImport_EnvFile(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "snap.env", "APP_ENV=staging\nPORT=9090\n# comment\n")

	snap, err := Import(p, ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Vars["APP_ENV"] != "staging" {
		t.Errorf("expected APP_ENV=staging, got %s", snap.Vars["APP_ENV"])
	}
	if snap.Vars["PORT"] != "9090" {
		t.Errorf("expected PORT=9090, got %s", snap.Vars["PORT"])
	}
}

func TestImport_JSONFile(t *testing.T) {
	dir := t.TempDir()
	m := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	b, _ := json.Marshal(m)
	p := writeFile(t, dir, "snap.json", string(b))

	snap, err := Import(p, ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Vars["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %s", snap.Vars["DB_HOST"])
	}
}

func TestImport_SkipsComments(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "snap.env", "# this is a comment\nKEY=value\n")

	snap, err := Import(p, ImportOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := snap.Vars["# this is a comment"]; ok {
		t.Error("comment line should not be parsed as a key")
	}
	if snap.Vars["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", snap.Vars["KEY"])
	}
}

func TestImport_MissingFile(t *testing.T) {
	_, err := Import("/nonexistent/path/snap.env", ImportOptions{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestImport_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, "snap.env", "KEY=val")
	_, err := Import(p, ImportOptions{Format: "yaml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
