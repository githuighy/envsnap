package snapshot

import "fmt"

// PromoteOptions controls how keys are promoted between environments.
type PromoteOptions struct {
	// Prefixes restricts promotion to keys matching any of these prefixes.
	Prefixes []string
	// Exclude skips keys matching any of these prefixes.
	Exclude []string
	// Overwrite allows destination keys to be overwritten.
	Overwrite bool
	// DryRun returns what would change without modifying the destination.
	DryRun bool
}

// PromoteResult records what happened to each key during promotion.
type PromoteResult struct {
	Key      string
	Action   string // "promoted", "skipped_exists", "skipped_filter"
	OldValue string
	NewValue string
}

// Promote copies keys from src into dst according to opts.
// It returns the list of results and the updated destination snapshot.
func Promote(src, dst Snapshot, opts PromoteOptions) ([]PromoteResult, Snapshot, error) {
	out := make(Snapshot, len(dst))
	for k, v := range dst {
		out[k] = v
	}

	var results []PromoteResult

	for key, val := range src {
		if len(opts.Prefixes) > 0 && !hasAnyPrefix(key, opts.Prefixes) {
			results = append(results, PromoteResult{Key: key, Action: "skipped_filter"})
			continue
		}
		if len(opts.Exclude) > 0 && hasAnyPrefix(key, opts.Exclude) {
			results = append(results, PromoteResult{Key: key, Action: "skipped_filter"})
			continue
		}
		existing, exists := out[key]
		if exists && !opts.Overwrite {
			results = append(results, PromoteResult{Key: key, Action: "skipped_exists", OldValue: existing, NewValue: val})
			continue
		}
		if !opts.DryRun {
			out[key] = val
		}
		results = append(results, PromoteResult{Key: key, Action: "promoted", OldValue: existing, NewValue: val})
	}

	if opts.DryRun {
		return results, dst, nil
	}
	return results, out, nil
}

// PromoteSummary returns a human-readable summary of promotion results.
func PromoteSummary(results []PromoteResult) string {
	promoted, skippedExists, skippedFilter := 0, 0, 0
	for _, r := range results {
		switch r.Action {
		case "promoted":
			promoted++
		case "skipped_exists":
			skippedExists++
		case "skipped_filter":
			skippedFilter++
		}
	}
	return fmt.Sprintf("promoted=%d skipped_exists=%d skipped_filter=%d", promoted, skippedExists, skippedFilter)
}
