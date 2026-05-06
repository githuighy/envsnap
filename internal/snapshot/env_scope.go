package snapshot

import (
	"fmt"
	"sort"
	"strings"
)

// ScopeOptions controls how scoping is applied to a snapshot.
type ScopeOptions struct {
	// Scope is the name of the scope (e.g. "prod", "staging").
	Scope string
	// PrefixSeparator is placed between scope and key (default "_").
	PrefixSeparator string
	// StripExisting removes any existing scope prefix before applying the new one.
	StripExisting bool
}

// ScopeResult holds the output of a Scope operation.
type ScopeResult struct {
	Snapshot map[string]string
	Renamed  map[string]string // oldKey -> newKey
}

// Scope applies a named scope prefix to all keys in the snapshot.
func Scope(snap map[string]string, opts ScopeOptions) (ScopeResult, error) {
	if opts.Scope == "" {
		return ScopeResult{}, fmt.Errorf("scope name must not be empty")
	}
	if strings.ContainsAny(opts.Scope, " \t\n") {
		return ScopeResult{}, fmt.Errorf("scope name must not contain whitespace")
	}
	sep := opts.PrefixSeparator
	if sep == "" {
		sep = "_"
	}
	prefix := strings.ToUpper(opts.Scope) + sep

	out := make(map[string]string, len(snap))
	renamed := make(map[string]string, len(snap))

	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		newKey := k
		if opts.StripExisting {
			if idx := strings.Index(k, sep); idx != -1 {
				newKey = k[idx+len(sep):]
			}
		}
		if !strings.HasPrefix(newKey, prefix) {
			newKey = prefix + newKey
		}
		out[newKey] = snap[k]
		if newKey != k {
			renamed[k] = newKey
		}
	}
	return ScopeResult{Snapshot: out, Renamed: renamed}, nil
}

// Unscope removes a scope prefix from all matching keys in the snapshot.
func Unscope(snap map[string]string, scope, separator string) map[string]string {
	if separator == "" {
		separator = "_"
	}
	prefix := strings.ToUpper(scope) + separator
	out := make(map[string]string, len(snap))
	for k, v := range snap {
		if strings.HasPrefix(k, prefix) {
			out[k[len(prefix):]] = v
		} else {
			out[k] = v
		}
	}
	return out
}
