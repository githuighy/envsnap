package snapshot

import (
	"fmt"
	"sort"
)

// CompareResult holds a summary of differences between two snapshots.
type CompareResult struct {
	Added    int
	Removed  int
	Changed  int
	Unchanged int
	Score    float64 // similarity score 0.0 - 1.0
}

// Summary returns a human-readable one-line summary of the comparison.
func (c CompareResult) Summary() string {
	return fmt.Sprintf(
		"added=%d removed=%d changed=%d unchanged=%d similarity=%.1f%%",
		c.Added, c.Removed, c.Changed, c.Unchanged, c.Score*100,
	)
}

// Compare produces a CompareResult between two snapshots.
func Compare(base, other map[string]string) CompareResult {
	diffs := Diff(base, other)

	result := CompareResult{}
	for _, d := range diffs {
		switch d.Status {
		case "added":
			result.Added++
		case "removed":
			result.Removed++
		case "changed":
			result.Changed++
		}
	}

	// Count keys present in both that are unchanged
	allKeys := unionKeys(base, other)
	total := len(allKeys)
	result.Unchanged = total - result.Added - result.Removed - result.Changed

	if total > 0 {
		result.Score = float64(result.Unchanged) / float64(total)
	} else {
		result.Score = 1.0
	}

	return result
}

// unionKeys returns a sorted slice of all unique keys across both maps.
func unionKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
