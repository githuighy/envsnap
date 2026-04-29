package snapshot_test

import (
	"path/filepath"
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func saveChainSnap(t *testing.T, dir, name string, s snapshot.Snapshot) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := snapshot.Save(s, p); err != nil {
		t.Fatalf("save: %v", err)
	}
	return p
}

func TestChain_RoundTripFromDisk(t *testing.T) {
	dir := t.TempDir()
	p1 := saveChainSnap(t, dir, "layer0.json", snapshot.Snapshot{"DB_HOST": "prod", "LOG": "info"})
	p2 := saveChainSnap(t, dir, "layer1.json", snapshot.Snapshot{"DB_HOST": "default", "TIMEOUT": "30"})

	l0, err := snapshot.Load(p1)
	if err != nil {
		t.Fatalf("load layer0: %v", err)
	}
	l1, err := snapshot.Load(p2)
	if err != nil {
		t.Fatalf("load layer1: %v", err)
	}

	chain := snapshot.NewChain(l0, l1)
	resolved := chain.Resolve()

	if resolved["DB_HOST"] != "prod" {
		t.Errorf("expected DB_HOST=prod, got %s", resolved["DB_HOST"])
	}
	if resolved["TIMEOUT"] != "30" {
		t.Errorf("expected TIMEOUT=30, got %s", resolved["TIMEOUT"])
	}
	if resolved["LOG"] != "info" {
		t.Errorf("expected LOG=info, got %s", resolved["LOG"])
	}
}

func TestChain_ExplainMatchesResolve(t *testing.T) {
	l0 := snapshot.Snapshot{"A": "1"}
	l1 := snapshot.Snapshot{"A": "2", "B": "3"}
	chain := snapshot.NewChain(l0, l1)

	resolved := chain.Resolve()
	origin := chain.Explain()

	for k := range resolved {
		if _, ok := origin[k]; !ok {
			t.Errorf("key %s in resolved but not in explain", k)
		}
	}
}

func TestChain_SingleLayer_ReturnsItself(t *testing.T) {
	snap := snapshot.Snapshot{"ONLY": "value"}
	chain := snapshot.NewChain(snap)
	resolved := chain.Resolve()
	if resolved["ONLY"] != "value" {
		t.Errorf("expected value, got %s", resolved["ONLY"])
	}
}
