package snapshot

import (
	"testing"
	"time"
)

func makeSnap(vars map[string]string) *Snapshot {
	return &Snapshot{Name: "test", Timestamp: time.Now(), Vars: vars}
}

func TestDiff_Added(t *testing.T) {
	base := makeSnap(map[string]string{"A": "1"})
	target := makeSnap(map[string]string{"A": "1", "B": "2"})

	changes := Diff(base, target)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != Added || changes[0].Key != "B" || changes[0].NewValue != "2" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiff_Removed(t *testing.T) {
	base := makeSnap(map[string]string{"A": "1", "B": "2"})
	target := makeSnap(map[string]string{"A": "1"})

	changes := Diff(base, target)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != Removed || changes[0].Key != "B" || changes[0].OldValue != "2" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestDiff_Changed(t *testing.T) {
	base := makeSnap(map[string]string{"A": "old"})
	target := makeSnap(map[string]string{"A": "new"})

	changes := Diff(base, target)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	c := changes[0]
	if c.Kind != Changed || c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected change: %+v", c)
	}
}

func TestDiff_NoChanges(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	base := makeSnap(vars)
	target := makeSnap(map[string]string{"A": "1", "B": "2"})

	changes := Diff(base, target)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}

func TestDiff_MultipleChanges(t *testing.T) {
	base := makeSnap(map[string]string{"A": "1", "B": "old", "C": "3"})
	target := makeSnap(map[string]string{"A": "1", "B": "new", "D": "4"})

	changes := Diff(base, target)
	// Expect: B changed, C removed, D added
	if len(changes) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(changes))
	}

	kinds := map[ChangeKind]int{}
	for _, c := range changes {
		kinds[c.Kind]++
	}
	if kinds[Added] != 1 || kinds[Removed] != 1 || kinds[Changed] != 1 {
		t.Errorf("unexpected change kinds: %v", kinds)
	}
}
