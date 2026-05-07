package snapshot

// IntersectOptions controls how two snapshots are intersected.
type IntersectOptions struct {
	// PreferLeft returns the value from the left snapshot when a key exists in both.
	// When false (default), the right snapshot value is used.
	PreferLeft bool

	// Prefixes restricts the intersection to keys matching any of the given prefixes.
	Prefixes []string

	// ExcludeKeys omits specific keys from the result even if they appear in both snapshots.
	ExcludeKeys []string
}

// IntersectResult holds the resulting snapshot and metadata about the operation.
type IntersectResult struct {
	Snapshot map[string]string
	CommonKeys []string
}

// Intersect returns a new snapshot containing only the keys that appear in both
// left and right. Values are taken from right unless PreferLeft is set.
func Intersect(left, right map[string]string, opts IntersectOptions) IntersectResult {
	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	result := make(map[string]string)
	var common []string

	for k := range left {
		if excluded[k] {
			continue
		}
		if len(opts.Prefixes) > 0 && !hasAnyPrefix(k, opts.Prefixes) {
			continue
		}
		if _, ok := right[k]; !ok {
			continue
		}
		common = append(common, k)
		if opts.PreferLeft {
			result[k] = left[k]
		} else {
			result[k] = right[k]
		}
	}

	return IntersectResult{
		Snapshot:   result,
		CommonKeys: common,
	}
}
