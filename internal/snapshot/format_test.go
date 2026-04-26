package snapshot

import (
	"strings"
	"testing"
)

func makeDiffs() []DiffEntry {
	return []DiffEntry{
		{Key: "APP_ENV", Status: StatusChanged, OldValue: "staging", NewValue: "production"},
		{Key: "NEW_KEY", Status: StatusAdded, OldValue: "", NewValue: "hello"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "bye", NewValue: ""},
	}
}

func TestRenderDiff_Text(t *testing.T) {
	var sb strings.Builder
	err := RenderDiff(&sb, makeDiffs(), FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, want := range []string{"APP_ENV", "changed", "staging", "production", "NEW_KEY", "OLD_KEY"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderDiff_JSON(t *testing.T) {
	var sb strings.Builder
	err := RenderDiff(&sb, makeDiffs(), FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Errorf("expected JSON array, got: %s", out)
	}
	for _, want := range []string{"\"status\"", "\"key\"", "APP_ENV", "changed"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected JSON to contain %q, got:\n%s", want, out)
		}
	}
}

func TestRenderDiff_Env(t *testing.T) {
	var sb strings.Builder
	err := RenderDiff(&sb, makeDiffs(), FormatEnv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "+ NEW_KEY=hello") {
		t.Errorf("expected '+ NEW_KEY=hello' in env output, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD_KEY=bye") {
		t.Errorf("expected '- OLD_KEY=bye' in env output, got:\n%s", out)
	}
	if !strings.Contains(out, "~ APP_ENV=production") {
		t.Errorf("expected '~ APP_ENV=production' in env output, got:\n%s", out)
	}
}

func TestRenderDiff_Text_NoDiffs(t *testing.T) {
	var sb strings.Builder
	err := RenderDiff(&sb, []DiffEntry{}, FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", sb.String())
	}
}
