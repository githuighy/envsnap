package snapshot

import (
	"testing"
)

var auditBase = Snapshot{
	"APP_ENV":    "production",
	"DB_HOST":    "db.prod.example.com",
	"LOG_LEVEL":  "warn",
	"SECRET_KEY": "abc123",
}

func TestAudit_Added(t *testing.T) {
	after := Snapshot{
		"APP_ENV":   "production",
		"DB_HOST":   "db.prod.example.com",
		"LOG_LEVEL": "warn",
		"SECRET_KEY": "abc123",
		"NEW_FEATURE_FLAG": "true",
	}
	log := Audit(auditBase, after, "deploy-42")
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Event != AuditAdded {
		t.Errorf("expected event %q, got %q", AuditAdded, e.Event)
	}
	if e.Key != "NEW_FEATURE_FLAG" {
		t.Errorf("unexpected key %q", e.Key)
	}
	if e.NewValue != "true" {
		t.Errorf("unexpected new value %q", e.NewValue)
	}
	if e.Source != "deploy-42" {
		t.Errorf("unexpected source %q", e.Source)
	}
}

func TestAudit_Removed(t *testing.T) {
	after := Snapshot{
		"APP_ENV":  "production",
		"DB_HOST":  "db.prod.example.com",
		"LOG_LEVEL": "warn",
	}
	log := Audit(auditBase, after, "")
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Event != AuditRemoved {
		t.Errorf("expected event %q, got %q", AuditRemoved, e.Event)
	}
	if e.Key != "SECRET_KEY" {
		t.Errorf("unexpected key %q", e.Key)
	}
	if e.OldValue != "abc123" {
		t.Errorf("unexpected old value %q", e.OldValue)
	}
}

func TestAudit_Changed(t *testing.T) {
	after := Snapshot{
		"APP_ENV":    "staging",
		"DB_HOST":    "db.prod.example.com",
		"LOG_LEVEL":  "warn",
		"SECRET_KEY": "abc123",
	}
	log := Audit(auditBase, after, "rollback")
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	e := log.Entries[0]
	if e.Event != AuditChanged {
		t.Errorf("expected event %q, got %q", AuditChanged, e.Event)
	}
	if e.OldValue != "production" || e.NewValue != "staging" {
		t.Errorf("unexpected values old=%q new=%q", e.OldValue, e.NewValue)
	}
}

func TestAudit_NoChanges(t *testing.T) {
	log := Audit(auditBase, auditBase, "noop")
	if len(log.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(log.Entries))
	}
}

func TestAuditLog_Summary(t *testing.T) {
	after := Snapshot{
		"APP_ENV":         "staging",
		"DB_HOST":         "db.prod.example.com",
		"NEW_VAR":         "hello",
	}
	log := Audit(auditBase, after, "")
	summary := log.Summary()
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	// Should mention counts
	expected := "audit: 1 added, 2 removed, 1 changed"
	if summary != expected {
		t.Errorf("expected %q, got %q", expected, summary)
	}
}
