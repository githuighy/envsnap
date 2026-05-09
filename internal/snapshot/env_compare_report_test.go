package snapshot

import (
	"encoding/json"
	"strings"
	"testing"
)

func makeDiffEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "APP_HOST", Status: "unchanged", OldValue: "localhost", NewValue: "localhost"},
		{Key: "APP_PORT", Status: "changed", OldValue: "8080", NewValue: "9090"},
		{Key: "APP_NEW", Status: "added", OldValue: "", NewValue: "yes"},
		{Key: "DB_PASS", Status: "removed", OldValue: "secret", NewValue: ""},
		{Key: "DB_HOST", Status: "unchanged", OldValue: "db", NewValue: "db"},
	}
}

func TestBuildCompareReport_Counts(t *testing.T) {
	diffs := makeDiffEntries()
	rep := BuildCompareReport(diffs, false)

	if rep.TotalKeys != 5 {
		t.Errorf("TotalKeys: got %d, want 5", rep.TotalKeys)
	}
	if rep.Added != 1 {
		t.Errorf("Added: got %d, want 1", rep.Added)
	}
	if rep.Removed != 1 {
		t.Errorf("Removed: got %d, want 1", rep.Removed)
	}
	if rep.Changed != 1 {
		t.Errorf("Changed: got %d, want 1", rep.Changed)
	}
	if rep.Unchanged != 2 {
		t.Errorf("Unchanged: got %d, want 2", rep.Unchanged)
	}
}

func TestBuildCompareReport_Similarity(t *testing.T) {
	diffs := makeDiffEntries()
	rep := BuildCompareReport(diffs, false)
	// 2 unchanged out of 5 = 40%
	if rep.Similarity < 39.9 || rep.Similarity > 40.1 {
		t.Errorf("Similarity: got %.2f, want ~40.0", rep.Similarity)
	}
}

func TestBuildCompareReport_ByPrefix(t *testing.T) {
	diffs := makeDiffEntries()
	rep := BuildCompareReport(diffs, true)

	if rep.ByPrefix == nil {
		t.Fatal("expected ByPrefix map to be populated")
	}
	if rep.ByPrefix["APP"] != 3 {
		t.Errorf("APP prefix count: got %d, want 3", rep.ByPrefix["APP"])
	}
	if rep.ByPrefix["DB"] != 2 {
		t.Errorf("DB prefix count: got %d, want 2", rep.ByPrefix["DB"])
	}
}

func TestBuildCompareReport_NoPrefixGrouping(t *testing.T) {
	diffs := makeDiffEntries()
	rep := BuildCompareReport(diffs, false)
	if rep.ByPrefix != nil {
		t.Error("expected ByPrefix to be nil when groupByPrefix=false")
	}
}

func TestRenderCompareReport_Text(t *testing.T) {
	diffs := makeDiffEntries()
	rep := BuildCompareReport(diffs, false)
	out, err := RenderCompareReport(rep, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"Total keys", "Added", "Removed", "Changed", "Similarity"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestRenderCompareReport_JSON(t *testing.T) {
	diffs := makeDiffEntries()
	rep := BuildCompareReport(diffs, false)
	out, err := RenderCompareReport(rep, "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded CompareReport
	if err := json.Unmarshal([]byte(out), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded.TotalKeys != rep.TotalKeys {
		t.Errorf("JSON TotalKeys mismatch: got %d, want %d", decoded.TotalKeys, rep.TotalKeys)
	}
}
