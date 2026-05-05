package snapshot

import (
	"fmt"
	"strings"
)

// RenameOptions controls how keys are renamed in a snapshot.
type RenameOptions struct {
	// Map is a direct key-to-key rename map (old -> new).
	Map map[string]string
	// StripPrefix removes a prefix from all matching keys.
	StripPrefix string
	// AddPrefix prepends a prefix to all keys (applied after StripPrefix).
	AddPrefix string
	// FailOnConflict returns an error if a rename would overwrite an existing key.
	FailOnConflict bool
}

// RenameResult holds the outcome of a rename operation.
type RenameResult struct {
	Snap     Snapshot
	Renamed  []string
	Skipped  []string
	Conflict []string
}

// Rename applies key renaming to a snapshot according to RenameOptions.
// Direct map renames are applied first, then prefix operations.
func Rename(snap Snapshot, opts RenameOptions) (RenameResult, error) {
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		out[k] = v
	}

	result := RenameResult{}

	// Apply direct key renames.
	for oldKey, newKey := range opts.Map {
		val, exists := out[oldKey]
		if !exists {
			result.Skipped = append(result.Skipped, oldKey)
			continue
		}
		if _, conflict := out[newKey]; conflict && newKey != oldKey {
			if opts.FailOnConflict {
				return RenameResult{}, fmt.Errorf("rename conflict: key %q already exists", newKey)
			}
			result.Conflict = append(result.Conflict, newKey)
			continue
		}
		delete(out, oldKey)
		out[newKey] = val
		result.Renamed = append(result.Renamed, oldKey+"="+newKey)
	}

	// Apply prefix strip + add.
	if opts.StripPrefix != "" || opts.AddPrefix != "" {
		prefixed := make(Snapshot, len(out))
		for k, v := range out {
			newKey := k
			if opts.StripPrefix != "" && strings.HasPrefix(k, opts.StripPrefix) {
				newKey = strings.TrimPrefix(k, opts.StripPrefix)
			}
			if opts.AddPrefix != "" {
				newKey = opts.AddPrefix + newKey
			}
			if newKey != k {
				if _, conflict := prefixed[newKey]; conflict {
					if opts.FailOnConflict {
						return RenameResult{}, fmt.Errorf("prefix rename conflict: key %q already exists", newKey)
					}
					result.Conflict = append(result.Conflict, newKey)
					prefixed[k] = v
					continue
				}
				result.Renamed = append(result.Renamed, k+"="+newKey)
			}
			prefixed[newKey] = v
		}
		out = prefixed
	}

	result.Snap = out
	return result, nil
}
