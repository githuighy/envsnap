package snapshot

// Status constants for diff entries.
const (
	StatusAdded   = "added"
	StatusRemoved = "removed"
	StatusChanged = "changed"
)

// DiffEntry represents a single changed environment variable.
type DiffEntry struct {
	Key      string
	Status   string
	OldValue string
	NewValue string
}

// Diff compares two snapshots and returns a list of differences.
func Diff(before, after Snapshot) []DiffEntry {
	var entries []DiffEntry

	// Check for removed or changed keys.
	for k, oldVal := range before.Vars {
		newVal, exists := after.Vars[k]
		if !exists {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusRemoved,
				OldValue: oldVal,
				NewValue: "",
			})
		} else if oldVal != newVal {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusChanged,
				OldValue: oldVal,
				NewValue: newVal,
			})
		}
	}

	// Check for added keys.
	for k, newVal := range after.Vars {
		if _, exists := before.Vars[k]; !exists {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusAdded,
				OldValue: "",
				NewValue: newVal,
			})
		}
	}

	return entries
}
