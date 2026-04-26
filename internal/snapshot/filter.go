package snapshot

import "strings"

// FilterOptions controls which environment variables are included or excluded.
type FilterOptions struct {
	// Prefixes limits results to variables whose names start with any of these.
	Prefixes []string
	// Exclude removes variables whose names start with any of these.
	Exclude []string
}

// Filter returns a new Snapshot containing only the variables that match the
// given FilterOptions. If opts.Prefixes is empty, all variables are considered
// candidates before exclusions are applied.
func Filter(snap Snapshot, opts FilterOptions) Snapshot {
	result := make(Snapshot)

	for key, val := range snap {
		if len(opts.Prefixes) > 0 && !hasAnyPrefix(key, opts.Prefixes) {
			continue
		}
		if len(opts.Exclude) > 0 && hasAnyPrefix(key, opts.Exclude) {
			continue
		}
		result[key] = val
	}

	return result
}

// hasAnyPrefix reports whether s starts with any of the given prefixes.
func hasAnyPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}
