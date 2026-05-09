package snapshot

import (
	"testing"
)

func basePatchSnap() Snapshot {
	return Snapshot{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/mydb",
	}
}

func TestPatch_SetAddsAndOverwrites(t *testing.T) {
	snap := basePatchSnap()
	ops := []PatchOp{
		{Op: "set", Key: "APP_PORT", Value: "9090"},
		{Op: "set", Key: "NEW_KEY", Value: "newval"},
	}
	res, err := Patch(snap, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Snapshot["APP_PORT"] != "9090" {
		t.Errorf("expected APP_PORT=9090, got %s", res.Snapshot["APP_PORT"])
	}
	if res.Snapshot["NEW_KEY"] != "newval" {
		t.Errorf("expected NEW_KEY=newval, got %s", res.Snapshot["NEW_KEY"])
	}
	if len(res.Applied) != 2 {
		t.Errorf("expected 2 applied, got %d", len(res.Applied))
	}
}

func TestPatch_DeleteRemovesKey(t *testing.T) {
	snap := basePatchSnap()
	ops := []PatchOp{
		{Op: "delete", Key: "DB_URL"},
	}
	res, err := Patch(snap, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snapshot["DB_URL"]; ok {
		t.Error("expected DB_URL to be deleted")
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(res.Applied))
	}
}

func TestPatch_Delete_MissingKeyIsSkipped(t *testing.T) {
	snap := basePatchSnap()
	ops := []PatchOp{
		{Op: "delete", Key: "NONEXISTENT"},
	}
	res, err := Patch(snap, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "NONEXISTENT" {
		t.Errorf("expected NONEXISTENT in skipped, got %v", res.Skipped)
	}
}

func TestPatch_RenameMovesKey(t *testing.T) {
	snap := basePatchSnap()
	ops := []PatchOp{
		{Op: "rename", Key: "APP_HOST", To: "SERVICE_HOST"},
	}
	res, err := Patch(snap, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snapshot["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be removed after rename")
	}
	if res.Snapshot["SERVICE_HOST"] != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %s", res.Snapshot["SERVICE_HOST"])
	}
}

func TestPatch_UnknownOp_ReturnsError(t *testing.T) {
	snap := basePatchSnap()
	ops := []PatchOp{
		{Op: "upsert", Key: "APP_HOST", Value: "x"},
	}
	_, err := Patch(snap, ops)
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestPatch_DoesNotMutateInput(t *testing.T) {
	snap := basePatchSnap()
	ops := []PatchOp{
		{Op: "set", Key: "APP_PORT", Value: "1111"},
	}
	_, err := Patch(snap, ops)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap["APP_PORT"] != "8080" {
		t.Error("original snapshot was mutated")
	}
}
