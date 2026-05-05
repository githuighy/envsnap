package snapshot

import (
	"testing"
)

func baseRenameSnap() Snapshot {
	return Snapshot{
		"APP_HOST":    "localhost",
		"APP_PORT":    "8080",
		"APP_SECRET":  "s3cr3t",
		"DB_HOST":     "db.local",
		"DB_PASSWORD": "pass",
	}
}

func TestRename_DirectMap(t *testing.T) {
	snap := baseRenameSnap()
	res, err := Rename(snap, RenameOptions{
		Map: map[string]string{"APP_HOST": "SERVICE_HOST"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snap["SERVICE_HOST"]; !ok {
		t.Error("expected SERVICE_HOST to exist after rename")
	}
	if _, ok := res.Snap["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be removed after rename")
	}
	if len(res.Renamed) != 1 {
		t.Errorf("expected 1 renamed, got %d", len(res.Renamed))
	}
}

func TestRename_StripPrefix(t *testing.T) {
	snap := baseRenameSnap()
	res, err := Rename(snap, RenameOptions{StripPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snap["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := res.Snap["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_ prefix")
	}
	if _, ok := res.Snap["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to remain unchanged")
	}
}

func TestRename_AddPrefix(t *testing.T) {
	snap := Snapshot{"HOST": "localhost", "PORT": "8080"}
	res, err := Rename(snap, RenameOptions{AddPrefix: "SVC_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snap["SVC_HOST"]; !ok {
		t.Error("expected SVC_HOST")
	}
	if _, ok := res.Snap["SVC_PORT"]; !ok {
		t.Error("expected SVC_PORT")
	}
}

func TestRename_SkipsMissingKey(t *testing.T) {
	snap := baseRenameSnap()
	res, err := Rename(snap, RenameOptions{
		Map: map[string]string{"NONEXISTENT": "NEW_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in skipped, got %v", res.Skipped)
	}
}

func TestRename_ConflictFailOnConflict(t *testing.T) {
	snap := Snapshot{"OLD_KEY": "val", "NEW_KEY": "other"}
	_, err := Rename(snap, RenameOptions{
		Map:            map[string]string{"OLD_KEY": "NEW_KEY"},
		FailOnConflict: true,
	})
	if err == nil {
		t.Error("expected error on conflict with FailOnConflict=true")
	}
}

func TestRename_StripAndAddPrefix(t *testing.T) {
	snap := Snapshot{"OLD_HOST": "localhost", "OLD_PORT": "9000"}
	res, err := Rename(snap, RenameOptions{
		StripPrefix: "OLD_",
		AddPrefix:   "NEW_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := res.Snap["NEW_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected NEW_HOST=localhost, got %q", v)
	}
	if v, ok := res.Snap["NEW_PORT"]; !ok || v != "9000" {
		t.Errorf("expected NEW_PORT=9000, got %q", v)
	}
}
