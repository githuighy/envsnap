package snapshot_test

import (
	"os"
	"strings"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func TestApplyProfile_FiltersByPrefix(t *testing.T) {
	os.Setenv("PROF_KEEP", "yes")
	os.Setenv("OTHER_DROP", "no")
	t.Cleanup(func() {
		os.Unsetenv("PROF_KEEP")
		os.Unsetenv("OTHER_DROP")
	})

	p := snapshot.Profile{
		Name:     "test",
		Prefixes: []string{"PROF_"},
	}

	snap, err := snapshot.ApplyProfile(p)
	if err != nil {
		t.Fatalf("ApplyProfile: %v", err)
	}
	if _, ok := snap["PROF_KEEP"]; !ok {
		t.Error("expected PROF_KEEP in snapshot")
	}
	if _, ok := snap["OTHER_DROP"]; ok {
		t.Error("expected OTHER_DROP to be filtered out")
	}
}

func TestApplyProfile_RedactsSensitiveKeys(t *testing.T) {
	os.Setenv("SECRET_TOKEN", "supersecret")
	t.Cleanup(func() { os.Unsetenv("SECRET_TOKEN") })

	p := snapshot.Profile{
		Name:   "test",
		Redact: []string{"SECRET_TOKEN"},
	}

	snap, err := snapshot.ApplyProfile(p)
	if err != nil {
		t.Fatalf("ApplyProfile: %v", err)
	}
	val, ok := snap["SECRET_TOKEN"]
	if !ok {
		t.Fatal("expected SECRET_TOKEN in snapshot")
	}
	if !strings.Contains(val, "***") && val == "supersecret" {
		t.Errorf("expected redacted value, got %q", val)
	}
}

func TestApplyProfile_NoOptions_ReturnsAll(t *testing.T) {
	os.Setenv("PLAIN_VAR", "hello")
	t.Cleanup(func() { os.Unsetenv("PLAIN_VAR") })

	p := snapshot.Profile{Name: "empty"}
	snap, err := snapshot.ApplyProfile(p)
	if err != nil {
		t.Fatalf("ApplyProfile: %v", err)
	}
	if _, ok := snap["PLAIN_VAR"]; !ok {
		t.Error("expected PLAIN_VAR in unrestricted snapshot")
	}
}
