package snapshot

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a single env value.
type TransformFunc func(key, value string) (string, error)

// TransformOptions controls which keys are transformed and how.
type TransformOptions struct {
	// Keys is an explicit list of keys to transform. If empty, all keys are transformed.
	Keys []string
	// Prefix restricts transformation to keys with the given prefix.
	Prefix string
	// SkipErrors causes transform errors to be skipped (original value kept) instead of failing.
	SkipErrors bool
}

// Transform applies fn to selected keys in snap, returning a new snapshot.
func Transform(snap map[string]string, fn TransformFunc, opts TransformOptions) (map[string]string, error) {
	if fn == nil {
		return nil, fmt.Errorf("transform: fn must not be nil")
	}

	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = struct{}{}
	}

	result := make(map[string]string, len(snap))
	for k, v := range snap {
		if !shouldTransform(k, keySet, opts.Prefix) {
			result[k] = v
			continue
		}
		transformed, err := fn(k, v)
		if err != nil {
			if opts.SkipErrors {
				result[k] = v
				continue
			}
			return nil, fmt.Errorf("transform: key %q: %w", k, err)
		}
		result[k] = transformed
	}
	return result, nil
}

func shouldTransform(key string, keySet map[string]struct{}, prefix string) bool {
	if len(keySet) > 0 {
		_, ok := keySet[key]
		return ok
	}
	if prefix != "" {
		return strings.HasPrefix(key, prefix)
	}
	return true
}
