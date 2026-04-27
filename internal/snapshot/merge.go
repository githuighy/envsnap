package snapshot

// MergeOptions controls how two snapshots are merged.
type MergeOptions struct {
	// Prefer controls which snapshot wins on key conflicts.
	// "base" keeps the base value, "override" keeps the override value.
	// Defaults to "override".
	Prefer string
}

// Snapshot is referenced from snapshot.go; we use the same type here.
// Merge combines base and override snapshots into a new snapshot.
// Keys present in both snapshots are resolved according to opts.Prefer.
// Keys unique to either snapshot are always included.
func Merge(base, override map[string]string, opts MergeOptions) map[string]string {
	if opts.Prefer == "" {
		opts.Prefer = "override"
	}

	result := make(map[string]string, len(base)+len(override))

	// Copy all base keys first.
	for k, v := range base {
		result[k] = v
	}

	// Apply override keys according to preference.
	for k, v := range override {
		if _, exists := result[k]; !exists {
			// Key only in override — always include.
			result[k] = v
			continue
		}
		if opts.Prefer == "override" {
			result[k] = v
		}
		// If prefer == "base", keep the existing base value.
	}

	return result
}
