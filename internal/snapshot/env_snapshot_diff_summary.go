package snapshot

import "fmt"

// DiffSummary holds aggregate statistics for a diff result.
type DiffSummary struct {
	Added   int
	Removed int
	Changed int
	Total   int
}

// String returns a human-readable one-line summary.
func (s DiffSummary) String() string {
	return fmt.Sprintf("+%d added, -%d removed, ~%d changed (%d total)",
		s.Added, s.Removed, s.Changed, s.Total)
}

// HasChanges reports whether any differences exist.
func (s DiffSummary) HasChanges() bool {
	return s.Added > 0 || s.Removed > 0 || s.Changed > 0
}

// SummariseDiff computes aggregate statistics from a slice of DiffEntry values.
func SummariseDiff(diffs []DiffEntry) DiffSummary {
	var s DiffSummary
	for _, d := range diffs {
		switch d.Status {
		case "added":
			s.Added++
		case "removed":
			s.Removed++
		case "changed":
			s.Changed++
		}
	}
	s.Total = s.Added + s.Removed + s.Changed
	return s
}

// DiffSummaryByPrefix returns a map from key-prefix to DiffSummary,
// grouping entries by the first component of each key (split on '_').
func DiffSummaryByPrefix(diffs []DiffEntry) map[string]DiffSummary {
	result := make(map[string]DiffSummary)
	for _, d := range diffs {
		prefix := keyPrefix(d.Key)
		s := result[prefix]
		switch d.Status {
		case "added":
			s.Added++
		case "removed":
			s.Removed++
		case "changed":
			s.Changed++
		}
		s.Total = s.Added + s.Removed + s.Changed
		result[prefix] = s
	}
	return result
}

func keyPrefix(key string) string {
	for i, c := range key {
		if c == '_' {
			return key[:i]
		}
	}
	return key
}
