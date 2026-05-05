package snapshot

import (
	"fmt"
	"os"
	"strings"
)

// ResolveOptions controls how variable references are expanded.
type ResolveOptions struct {
	// AllowMissing, when true, leaves unresolvable references unchanged
	// instead of returning an error.
	AllowMissing bool
	// MaxDepth limits recursive expansion to prevent infinite loops.
	MaxDepth int
}

// DefaultResolveOptions returns sensible defaults for resolution.
func DefaultResolveOptions() ResolveOptions {
	return ResolveOptions{
		AllowMissing: false,
		MaxDepth:     10,
	}
}

// Resolve expands ${VAR} and $VAR references within snapshot values using
// the snapshot itself as the source, falling back to os.Getenv for keys
// not present in the snapshot.
func Resolve(snap Snapshot, opts ResolveOptions) (Snapshot, error) {
	if opts.MaxDepth <= 0 {
		opts.MaxDepth = DefaultResolveOptions().MaxDepth
	}

	out := make(Snapshot, len(snap))
	for k, v := range snap {
		resolved, err := expandValue(v, snap, opts, 0)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

// expandValue recursively expands ${VAR} / $VAR references in a single value.
func expandValue(value string, snap Snapshot, opts ResolveOptions, depth int) (string, error) {
	if depth > opts.MaxDepth {
		return "", fmt.Errorf("max expansion depth %d exceeded", opts.MaxDepth)
	}

	expanded := os.Expand(value, func(key string) string {
		if v, ok := snap[key]; ok {
			return v
		}
		if env := os.Getenv(key); env != "" {
			return env
		}
		return ""
	})

	// If the result still contains references and changed, recurse.
	if expanded != value && strings.Contains(expanded, "$") {
		return expandValue(expanded, snap, opts, depth+1)
	}

	// If nothing changed and there are still unresolved refs, handle missing.
	if !opts.AllowMissing && strings.Contains(expanded, "${") {
		return "", fmt.Errorf("unresolved reference in value %q", value)
	}

	return expanded, nil
}
