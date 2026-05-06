package snapshot

import (
	"sort"
	"strings"
)

// DedupeOptions controls how duplicate detection works across snapshots.
type DedupeOptions struct {
	// ByValue removes keys whose values are identical across all provided snapshots.
	ByValue bool
	// ByPrefix collapses keys sharing the same prefix, keeping the first encountered.
	ByPrefix bool
	// Prefixes is the list of prefixes to consider when ByPrefix is true.
	Prefixes []string
}

// DedupeResult holds the output snapshot and metadata about removed keys.
type DedupeResult struct {
	Snapshot map[string]string
	RemovedKeys []string
}

// Dedupe removes duplicate or redundant keys from a snapshot according to opts.
// When ByValue is true, keys whose values already appear under another key are
// removed (the lexicographically first key is kept).
// When ByPrefix is true, for each prefix only the first key (sorted) is kept.
func Dedupe(snap map[string]string, opts DedupeOptions) DedupeResult {
	out := make(map[string]string, len(snap))
	for k, v := range snap {
		out[k] = v
	}

	var removed []string

	if opts.ByValue {
		// Build value -> first key mapping (sorted for determinism).
		keys := sortedMapKeys(out)
		seen := make(map[string]string) // value -> keeper key
		for _, k := range keys {
			v := out[k]
			if keeper, exists := seen[v]; exists && keeper != k {
				removed = append(removed, k)
				delete(out, k)
			} else if !exists {
				seen[v] = k
			}
		}
	}

	if opts.ByPrefix && len(opts.Prefixes) > 0 {
		keys := sortedMapKeys(out)
		for _, prefix := range opts.Prefixes {
			kept := false
			for _, k := range keys {
				if _, stillPresent := out[k]; !stillPresent {
					continue
				}
				if strings.HasPrefix(k, prefix) {
					if !kept {
						kept = true
					} else {
						removed = append(removed, k)
						delete(out, k)
					}
				}
			}
		}
	}

	sort.Strings(removed)
	return DedupeResult{Snapshot: out, RemovedKeys: removed}
}

func sortedMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
