package snapshot

import (
	"strings"
)

// TrimOptions controls how values are trimmed in a snapshot.
type TrimOptions struct {
	// Prefixes restricts trimming to keys with any of these prefixes.
	Prefixes []string
	// Keys restricts trimming to exactly these keys.
	Keys []string
	// TrimPrefix removes a leading string from each value.
	TrimPrefix string
	// TrimSuffix removes a trailing string from each value.
	TrimSuffix string
	// TrimSpace trims leading and trailing whitespace from values.
	TrimSpace bool
}

// Trim applies value trimming operations to a snapshot according to opts.
// It returns a new snapshot and leaves the original unchanged.
func Trim(snap map[string]string, opts TrimOptions) (map[string]string, error) {
	result := make(map[string]string, len(snap))
	for k, v := range snap {
		result[k] = v
	}

	for k, v := range result {
		if !shouldTrim(k, opts) {
			continue
		}
		if opts.TrimSpace {
			v = strings.TrimSpace(v)
		}
		if opts.TrimPrefix != "" {
			v = strings.TrimPrefix(v, opts.TrimPrefix)
		}
		if opts.TrimSuffix != "" {
			v = strings.TrimSuffix(v, opts.TrimSuffix)
		}
		result[k] = v
	}
	return result, nil
}

// shouldTrim returns true if the key is targeted by the trim options.
func shouldTrim(key string, opts TrimOptions) bool {
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
