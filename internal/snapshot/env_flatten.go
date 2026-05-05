package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// FlattenOptions controls how nested key structures are flattened.
type FlattenOptions struct {
	// Separator is the delimiter used to join nested key segments (default: "_").
	Separator string
	// Uppercase converts all resulting keys to uppercase.
	Uppercase bool
	// Prefix is prepended to every flattened key.
	Prefix string
	// SkipEmpty drops entries whose value is an empty string.
	SkipEmpty bool
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: "_",
		Uppercase: true,
	}
}

// Flatten merges multiple snapshots into a single snapshot, rewriting keys
// so that each key is prefixed with its source snapshot index (or a provided
// label) and the configured separator. Duplicate keys from later snapshots
// overwrite earlier ones.
//
// Example: label "APP", separator "_", key "DB_HOST" → "APP_DB_HOST"
func Flatten(snaps []Snapshot, labels []string, opts FlattenOptions) (Snapshot, error) {
	if opts.Separator == "" {
		opts.Separator = "_"
	}

	if len(labels) > 0 && len(labels) != len(snaps) {
		return Snapshot{}, fmt.Errorf("flatten: labels length %d does not match snapshots length %d", len(labels), len(snaps))
	}

	out := make(map[string]string)

	for i, snap := range snaps {
		label := fmt.Sprintf("%d", i)
		if len(labels) > 0 {
			label = labels[i]
		}
		if label == "" {
			return Snapshot{}, fmt.Errorf("flatten: label at index %d must not be empty", i)
		}

		for k, v := range snap.Vars {
			if opts.SkipEmpty && v == "" {
				continue
			}

			parts := []string{label, k}
			joined := strings.Join(parts, opts.Separator)

			if opts.Prefix != "" {
				joined = opts.Prefix + opts.Separator + joined
			}
			if opts.Uppercase {
				joined = strings.ToUpper(joined)
			}

			out[joined] = v
		}
	}

	// Deterministic key ordering for the returned snapshot.
	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = out[k]
	}

	return Snapshot{Vars: ordered}, nil
}
