package snapshot

import (
	"testing"
)

var baseScopeSnap = map[string]string{
	"DB_HOST": "localhost",
	"DB_PORT": "5432",
	"APP_ENV":  "development",
}

func TestScope_AddsScopePrefix(t *testing.T) {
	res, err := Scope(baseScopeSnap, ScopeOptions{Scope: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snapshot["PROD_DB_HOST"]; !ok {
		t.Error("expected PROD_DB_HOST in scoped snapshot")
	}
	if _, ok := res.Snapshot["PROD_APP_ENV"]; !ok {
		t.Error("expected PROD_APP_ENV in scoped snapshot")
	}
	if len(res.Renamed) != 3 {
		t.Errorf("expected 3 renamed keys, got %d", len(res.Renamed))
	}
}

func TestScope_CustomSeparator(t *testing.T) {
	res, err := Scope(baseScopeSnap, ScopeOptions{Scope: "staging", PrefixSeparator: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snapshot["STAGING.DB_HOST"]; !ok {
		t.Error("expected STAGING.DB_HOST")
	}
}

func TestScope_EmptyScope_ReturnsError(t *testing.T) {
	_, err := Scope(baseScopeSnap, ScopeOptions{Scope: ""})
	if err == nil {
		t.Fatal("expected error for empty scope")
	}
}

func TestScope_WhitespaceScopeName_ReturnsError(t *testing.T) {
	_, err := Scope(baseScopeSnap, ScopeOptions{Scope: "prod env"})
	if err == nil {
		t.Fatal("expected error for whitespace in scope name")
	}
}

func TestScope_StripExisting(t *testing.T) {
	snap := map[string]string{
		"DEV_DB_HOST": "localhost",
		"DEV_APP_ENV": "development",
	}
	res, err := Scope(snap, ScopeOptions{Scope: "prod", StripExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Snapshot["PROD_DB_HOST"]; !ok {
		t.Error("expected PROD_DB_HOST after strip and re-scope")
	}
	if _, ok := res.Snapshot["PROD_DEV_DB_HOST"]; ok {
		t.Error("did not expect PROD_DEV_DB_HOST — old scope should have been stripped")
	}
}

func TestUnscope_RemovesPrefix(t *testing.T) {
	scoped := map[string]string{
		"PROD_DB_HOST": "localhost",
		"PROD_DB_PORT": "5432",
		"UNRELATED":    "value",
	}
	out := Unscope(scoped, "prod", "_")
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", out["DB_HOST"])
	}
	if _, ok := out["UNRELATED"]; !ok {
		t.Error("expected UNRELATED key to be preserved")
	}
}

func TestUnscope_DefaultSeparator(t *testing.T) {
	scoped := map[string]string{"PROD_KEY": "val"}
	out := Unscope(scoped, "prod", "")
	if out["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", out["KEY"])
	}
}
