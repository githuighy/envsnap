package snapshot

import (
	"fmt"
	"sort"
)

// LockResult holds the outcome of a lock check for a single key.
type LockResult struct {
	Key      string
	Expected string
	Actual   string
	Locked   bool
}

// LockOptions configures which keys are locked and their expected values.
type LockOptions struct {
	// LockedKeys maps key names to their expected (locked) values.
	LockedKeys map[string]string
}

// Lock checks that the specified keys in snap match their expected values.
// It returns a slice of LockResult, one per locked key, and an error if
// any locked key is missing or has an unexpected value.
func Lock(snap Snapshot, opts LockOptions) ([]LockResult, error) {
	if len(opts.LockedKeys) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(opts.LockedKeys))
	for k := range opts.LockedKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var results []LockResult
	var violations []string

	for _, key := range keys {
		expected := opts.LockedKeys[key]
		actual, exists := snap[key]
		if !exists {
			results = append(results, LockResult{
				Key:      key,
				Expected: expected,
				Actual:   "",
				Locked:   false,
			})
			violations = append(violations, fmt.Sprintf("%s: missing (expected %q)", key, expected))
			continue
		}
		matched := actual == expected
		results = append(results, LockResult{
			Key:      key,
			Expected: expected,
			Actual:   actual,
			Locked:   matched,
		})
		if !matched {
			violations = append(violations, fmt.Sprintf("%s: got %q, want %q", key, actual, expected))
		}
	}

	if len(violations) > 0 {
		return results, fmt.Errorf("lock violations: %v", violations)
	}
	return results, nil
}
