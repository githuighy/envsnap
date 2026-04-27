package snapshot

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var exportSnap = Snapshot{
	Vars: map[string]string{
		"APP_ENV": "production",
		"PORT":    "8080",
	},
}

func TestExport_EnvFormat(t *testing.T) {
	r, w, _ := os.Pipe()
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	if err := Export(exportSnap, ExportOptions{Format: ExportFormatEnv}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Close()

	var buf strings.Builder
	b := make([]byte, 256)
	for {
		n, err := r.Read(b)
		buf.Write(b[:n])
		if err != nil {
			break
		}
	}
	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got: %s", out)
	}
}

func TestExport_JSONFormat_ToFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "snap.json")

	if err := Export(exportSnap, ExportOptions{Format: ExportFormatJSON, OutFile: out}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", m["PORT"])
	}
}

func TestExport_ShellFormat_ToFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "env.sh")

	if err := Export(exportSnap, ExportOptions{Format: ExportFormatShell, OutFile: out}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}

	if !strings.Contains(string(data), "export APP_ENV=") {
		t.Errorf("expected shell export statement, got: %s", string(data))
	}
}

func TestExport_InvalidFormat(t *testing.T) {
	err := Export(exportSnap, ExportOptions{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
