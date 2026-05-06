package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func saveScopeSnap(t *testing.T, dir string, data map[string]string) string {
	t.Helper()
	p := filepath.Join(dir, "snap.json")
	if err := snapshot.Save(data, p); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p
}

func TestScope_RoundTripFromDisk(t *testing.T) {
	dir := t.TempDir()
	orig := map[string]string{
		"DB_HOST": "db.internal",
		"DB_PORT": "5432",
		"SECRET":  "s3cr3t",
	}
	p := saveScopeSnap(t, dir, orig)

	loaded, err := snapshot.Load(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	res, err := snapshot.Scope(loaded, snapshot.ScopeOptions{Scope: "prod"})
	if err != nil {
		t.Fatalf("scope: %v", err)
	}

	out := filepath.Join(dir, "scoped.json")
	if err := snapshot.Save(res.Snapshot, out); err != nil {
		t.Fatalf("save scoped: %v", err)
	}

	scoped, err := snapshot.Load(out)
	if err != nil {
		t.Fatalf("load scoped: %v", err)
	}
	for _, k := range []string{"PROD_DB_HOST", "PROD_DB_PORT", "PROD_SECRET"} {
		if _, ok := scoped[k]; !ok {
			t.Errorf("expected key %s in scoped snapshot", k)
		}
	}
}

func TestScope_UnscopeRestoresOriginal(t *testing.T) {
	orig := map[string]string{"API_KEY": "abc", "LOG_LEVEL": "info"}

	res, err := snapshot.Scope(orig, snapshot.ScopeOptions{Scope: "dev"})
	if err != nil {
		t.Fatalf("scope: %v", err)
	}
	restored := snapshot.Unscope(res.Snapshot, "dev", "_")

	for k, v := range orig {
		if restored[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, restored[k])
		}
	}
}

func TestScope_StripAndReplace_RoundTrip(t *testing.T) {
	snap := map[string]string{
		"DEV_DB_HOST": "localhost",
		"DEV_APP_ENV": "development",
	}
	res, err := snapshot.Scope(snap, snapshot.ScopeOptions{
		Scope:         "prod",
		StripExisting: true,
	})
	if err != nil {
		t.Fatalf("scope: %v", err)
	}
	if res.Snapshot["PROD_DB_HOST"] != "localhost" {
		t.Errorf("expected PROD_DB_HOST=localhost")
	}
	if _, ok := res.Snapshot["PROD_DEV_DB_HOST"]; ok {
		t.Error("should not have double-prefixed key")
	}
}
