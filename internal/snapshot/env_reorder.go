package snapshot

import "sort"

// ReorderOptions controls how keys are sorted in the output snapshot.
type ReorderOptions struct {
	// Alphabetical sorts all keys A→Z.
	Alphabetical bool
	// Descending reverses the sort order (requires Alphabetical).
	Descending bool
	// PrefixPriority lists prefixes whose keys should appear first, in order.
	PrefixPriority []string
}

// Reorder returns a new snapshot whose Vars are sorted according to opts.
// The original snapshot is never mutated.
func Reorder(s Snapshot, opts ReorderOptions) (Snapshot, error) {
	out := Snapshot{Vars: make(map[string]string, len(s.Vars))}
	for k, v := range s.Vars {
		out.Vars[k] = v
	}

	if !opts.Alphabetical && len(opts.PrefixPriority) == 0 {
		return out, nil
	}

	keys := make([]string, 0, len(out.Vars))
	for k := range out.Vars {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		pi := prefixRank(keys[i], opts.PrefixPriority)
		pj := prefixRank(keys[j], opts.PrefixPriority)
		if pi != pj {
			return pi < pj
		}
		if opts.Alphabetical {
			if opts.Descending {
				return keys[i] > keys[j]
			}
			return keys[i] < keys[j]
		}
		return false
	})

	ordered := make(map[string]string, len(keys))
	for _, k := range keys {
		ordered[k] = out.Vars[k]
	}
	out.Vars = ordered
	return out, nil
}

// prefixRank returns the priority index for a key (lower = higher priority).
// Keys that match no prefix get rank len(prefixes).
func prefixRank(key string, prefixes []string) int {
	for i, p := range prefixes {
		if len(key) >= len(p) && key[:len(p)] == p {
			return i
		}
	}
	return len(prefixes)
}
