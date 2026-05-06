package snapshot

import "fmt"

// RollbackOptions controls which keys are rolled back and how.
type RollbackOptions struct {
	// Keys limits rollback to specific keys; empty means all keys.
	Keys []string
	// Prefixes limits rollback to keys with any of these prefixes.
	Prefixes []string
	// DryRun reports what would change without modifying the snapshot.
	DryRun bool
}

// RollbackResult describes the outcome of a rollback operation.
type RollbackResult struct {
	Snapshot  Snapshot
	Restored  map[string]string // key -> value restored from baseline
	Dropped   []string          // keys present in current but absent in baseline
	Unchanged []string          // keys whose values already matched baseline
}

// Rollback reverts keys in current to their values in baseline according to opts.
// Keys absent in baseline are dropped from the result unless DryRun is set.
func Rollback(current, baseline Snapshot, opts RollbackOptions) (RollbackResult, error) {
	if current.Vars == nil {
		return RollbackResult{}, fmt.Errorf("rollback: current snapshot is nil")
	}
	if baseline.Vars == nil {
		return RollbackResult{}, fmt.Errorf("rollback: baseline snapshot is nil")
	}

	result := RollbackResult{
		Restored: make(map[string]string),
	}

	out := make(map[string]string)
	for k, v := range current.Vars {
		out[k] = v
	}

	for k := range out {
		if !shouldRollback(k, opts) {
			continue
		}
		baseVal, inBase := baseline.Vars[k]
		if !inBase {
			result.Dropped = append(result.Dropped, k)
			if !opts.DryRun {
				delete(out, k)
			}
			continue
		}
		if out[k] == baseVal {
			result.Unchanged = append(result.Unchanged, k)
			continue
		}
		result.Restored[k] = baseVal
		if !opts.DryRun {
			out[k] = baseVal
		}
	}

	result.Snapshot = Snapshot{Vars: out}
	return result, nil
}

func shouldRollback(key string, opts RollbackOptions) bool {
	if len(opts.Keys) > 0 {
		for _, k := range opts.Keys {
			if k == key {
				return true
			}
		}
		return false
	}
	if len(opts.Prefixes) > 0 {
		return hasAnyPrefix(key, opts.Prefixes)
	}
	return true
}
