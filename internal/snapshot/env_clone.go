package snapshot

import "fmt"

// CloneOptions controls how a snapshot is cloned.
type CloneOptions struct {
	// KeyMap renames keys during clone: old name -> new name.
	KeyMap map[string]string
	// Prefixes restricts cloning to keys with any of these prefixes.
	Prefixes []string
	// Exclude omits keys with any of these prefixes.
	Exclude []string
	// OverrideValues replaces specific key values after cloning.
	OverrideValues map[string]string
}

// Clone produces a new Snapshot derived from src, applying optional
// filtering, renaming, and value overrides.
func Clone(src map[string]string, opts CloneOptions) (map[string]string, error) {
	if src == nil {
		return nil, fmt.Errorf("clone: source snapshot is nil")
	}

	out := make(map[string]string, len(src))

	for k, v := range src {
		if len(opts.Prefixes) > 0 && !hasAnyPrefix(k, opts.Prefixes) {
			continue
		}
		if len(opts.Exclude) > 0 && hasAnyPrefix(k, opts.Exclude) {
			continue
		}

		destKey := k
		if mapped, ok := opts.KeyMap[k]; ok {
			if mapped == "" {
				return nil, fmt.Errorf("clone: empty destination key for source key %q", k)
			}
			destKey = mapped
		}

		out[destKey] = v
	}

	for k, v := range opts.OverrideValues {
		out[k] = v
	}

	return out, nil
}
