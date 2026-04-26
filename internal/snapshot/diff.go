package snapshot

// ChangeKind describes the type of change between two snapshots.
type ChangeKind string

const (
	Added   ChangeKind = "added"
	Removed ChangeKind = "removed"
	Changed ChangeKind = "changed"
)

// Change represents a single variable difference between two snapshots.
type Change struct {
	Key      string     `json:"key"`
	Kind     ChangeKind `json:"kind"`
	OldValue string     `json:"old_value,omitempty"`
	NewValue string     `json:"new_value,omitempty"`
}

// Diff computes the differences between a base snapshot and a target snapshot.
// Variables present only in target are Added; only in base are Removed;
// present in both with different values are Changed.
func Diff(base, target *Snapshot) []Change {
	var changes []Change

	for k, newVal := range target.Vars {
		oldVal, exists := base.Vars[k]
		if !exists {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: newVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, oldVal := range base.Vars {
		if _, exists := target.Vars[k]; !exists {
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: oldVal})
		}
	}

	return changes
}
