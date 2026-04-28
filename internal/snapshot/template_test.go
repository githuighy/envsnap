package snapshot

import (
	"strings"
	"testing"
)

var baseTmpl = Template{
	Name: "web-service",
	Defaults: map[string]string{
		"LOG_LEVEL": "info",
		"PORT":      "8080",
	},
	Required: []string{"DATABASE_URL", "APP_SECRET"},
}

func TestApplyTemplate_FillsDefaults(t *testing.T) {
	snap := Snapshot{"DATABASE_URL": "postgres://localhost/db", "APP_SECRET": "s3cr3t"}
	result, err := ApplyTemplate(snap, baseTmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %q", result["LOG_LEVEL"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", result["PORT"])
	}
}

func TestApplyTemplate_ExistingKeyNotOverwritten(t *testing.T) {
	snap := Snapshot{"DATABASE_URL": "postgres://localhost/db", "APP_SECRET": "s3cr3t", "PORT": "9090"}
	result, err := ApplyTemplate(snap, baseTmpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["PORT"] != "9090" {
		t.Errorf("expected PORT=9090 (not overwritten), got %q", result["PORT"])
	}
}

func TestApplyTemplate_MissingRequired(t *testing.T) {
	snap := Snapshot{"APP_SECRET": "s3cr3t"}
	_, err := ApplyTemplate(snap, baseTmpl)
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("error should mention DATABASE_URL, got: %v", err)
	}
}

func TestApplyTemplate_InvalidName(t *testing.T) {
	tmpl := Template{Name: "bad name!", Defaults: nil, Required: nil}
	_, err := ApplyTemplate(Snapshot{}, tmpl)
	if err == nil {
		t.Fatal("expected error for invalid template name")
	}
}

func TestStripTemplateDefaults_RemovesUnchangedDefaults(t *testing.T) {
	snap := Snapshot{
		"LOG_LEVEL":    "info",
		"PORT":         "9090",
		"DATABASE_URL": "postgres://localhost/db",
		"APP_SECRET":   "s3cr3t",
	}
	result := StripTemplateDefaults(snap, baseTmpl)
	if _, ok := result["LOG_LEVEL"]; ok {
		t.Error("LOG_LEVEL with default value should have been stripped")
	}
	if result["PORT"] != "9090" {
		t.Error("PORT with non-default value should be retained")
	}
}
