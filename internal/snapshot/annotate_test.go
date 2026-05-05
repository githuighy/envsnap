package snapshot

import (
	"strings"
	"testing"
)

func baseAnnotateSnap() Snapshot {
	return Snapshot{
		"APP_ENV":  "production",
		"APP_PORT": "8080",
	}
}

func TestAnnotate_AddsAnnotations(t *testing.T) {
	snap := baseAnnotateSnap()
	annotations := []Annotation{
		{Key: "author", Value: "alice"},
		{Key: "version", Value: "v1.2.3"},
	}

	out, err := Annotate(snap, annotations)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out[annotationPrefix+"author"] != "alice" {
		t.Errorf("expected annotation author=alice, got %q", out[annotationPrefix+"author"])
	}
	if out[annotationPrefix+"version"] != "v1.2.3" {
		t.Errorf("expected annotation version=v1.2.3, got %q", out[annotationPrefix+"version"])
	}
	// original keys preserved
	if out["APP_ENV"] != "production" {
		t.Errorf("original key APP_ENV should be preserved")
	}
}

func TestAnnotate_EmptyKeyReturnsError(t *testing.T) {
	snap := baseAnnotateSnap()
	_, err := Annotate(snap, []Annotation{{Key: "", Value: "oops"}})
	if err == nil {
		t.Fatal("expected error for empty annotation key")
	}
}

func TestAnnotate_InvalidKeyCharacters(t *testing.T) {
	snap := baseAnnotateSnap()
	_, err := Annotate(snap, []Annotation{{Key: "bad=key", Value: "val"}})
	if err == nil {
		t.Fatal("expected error for key containing '='")
	}
}

func TestAnnotate_DoesNotMutateOriginalSnapshot(t *testing.T) {
	snap := baseAnnotateSnap()
	originalLen := len(snap)

	_, err := Annotate(snap, []Annotation{{Key: "author", Value: "alice"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(snap) != originalLen {
		t.Errorf("Annotate mutated the original snapshot: expected %d keys, got %d", originalLen, len(snap))
	}
}

func TestGetAnnotations_ReturnsAll(t *testing.T) {
	snap := baseAnnotateSnap()
	snap[annotationPrefix+"env"] = "staging"
	snap[annotationPrefix+"build"] = "42"

	anns := GetAnnotations(snap)
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
	keys := make(map[string]string)
	for _, a := range anns {
		keys[a.Key] = a.Value
	}
	if keys["env"] != "staging" {
		t.Errorf("expected env=staging")
	}
	if keys["build"] != "42" {
		t.Errorf("expected build=42")
	}
}

func TestStripAnnotations_RemovesAnnotationKeys(t *testing.T) {
	snap := baseAnnotateSnap()
	snap[annotationPrefix+"author"] = "bob"

	stripped := StripAnnotations(snap)
	for k := range stripped {
		if strings.HasPrefix(k, annotationPrefix) {
			t.Errorf("stripped snapshot still contains annotation key %q", k)
		}
	}
	if stripped["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to remain after stripping")
	}
}

func TestTimestampAnnotation_HasCapturedAtKey(t *testing.T) {
	a := TimestampAnnotation()
	if a.Key != "captured_at" {
		t.Errorf("expected key 'captured_at', got %q", a.Key)
	}
	if a.Value == "" {
		t.Error("expected non-empty timestamp value")
	}
}
