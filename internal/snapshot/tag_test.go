package snapshot_test

import (
	"testing"

	"github.com/yourorg/envsnap/internal/snapshot"
)

var baseTagSnap = snapshot.Snapshot{
	"APP_ENV":  "production",
	"LOG_LEVEL": "info",
	"PORT":     "8080",
}

func TestTagSnapshot_AddsTags(t *testing.T) {
	tags := []snapshot.Tag{
		{Name: "env", Value: "prod"},
		{Name: "version", Value: "1.2.3"},
	}

	result, err := snapshot.TagSnapshot(baseTagSnap, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["__tag.env"] != "prod" {
		t.Errorf("expected __tag.env=prod, got %q", result["__tag.env"])
	}
	if result["__tag.version"] != "1.2.3" {
		t.Errorf("expected __tag.version=1.2.3, got %q", result["__tag.version"])
	}
	if result["APP_ENV"] != "production" {
		t.Error("original keys should be preserved")
	}
}

func TestTagSnapshot_InvalidName(t *testing.T) {
	tags := []snapshot.Tag{
		{Name: "bad name!", Value: "x"},
	}
	_, err := snapshot.TagSnapshot(baseTagSnap, tags)
	if err == nil {
		t.Error("expected error for invalid tag name, got nil")
	}
}

func TestListTags_ReturnsSorted(t *testing.T) {
	snap := snapshot.Snapshot{
		"APP_ENV":       "production",
		"__tag.version": "2.0.0",
		"__tag.env":     "staging",
		"__tag.owner":   "team-a",
	}

	tags := snapshot.ListTags(snap)
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
	if tags[0].Name != "env" || tags[1].Name != "owner" || tags[2].Name != "version" {
		t.Errorf("tags not sorted: %+v", tags)
	}
}

func TestListTags_Empty(t *testing.T) {
	tags := snapshot.ListTags(baseTagSnap)
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %d", len(tags))
	}
}

func TestStripTags_RemovesTagKeys(t *testing.T) {
	snap := snapshot.Snapshot{
		"APP_ENV":       "production",
		"__tag.env":     "prod",
		"__tag.version": "1.0.0",
	}

	stripped := snapshot.StripTags(snap)
	if _, ok := stripped["__tag.env"]; ok {
		t.Error("__tag.env should have been stripped")
	}
	if stripped["APP_ENV"] != "production" {
		t.Error("APP_ENV should be preserved after strip")
	}
	if len(stripped) != 1 {
		t.Errorf("expected 1 key after strip, got %d", len(stripped))
	}
}
