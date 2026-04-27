package snapshot

import (
	"fmt"
	"strings"
)

// ValidationError holds a list of issues found during snapshot validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("snapshot validation failed: %s", strings.Join(e.Issues, "; "))
}

// ValidateOptions controls which checks are applied during validation.
type ValidateOptions struct {
	// RequiredKeys lists keys that must be present and non-empty.
	RequiredKeys []string
	// ForbiddenKeys lists keys that must not be present.
	ForbiddenKeys []string
	// MaxKeyLength rejects keys longer than this value (0 = no limit).
	MaxKeyLength int
}

// Validate checks a Snapshot against the provided options and returns a
// *ValidationError if any issues are found, or nil on success.
func Validate(snap Snapshot, opts ValidateOptions) error {
	var issues []string

	for _, key := range opts.RequiredKeys {
		val, ok := snap[key]
		if !ok || strings.TrimSpace(val) == "" {
			issues = append(issues, fmt.Sprintf("required key %q is missing or empty", key))
		}
	}

	for _, key := range opts.ForbiddenKeys {
		if _, ok := snap[key]; ok {
			issues = append(issues, fmt.Sprintf("forbidden key %q is present", key))
		}
	}

	if opts.MaxKeyLength > 0 {
		for key := range snap {
			if len(key) > opts.MaxKeyLength {
				issues = append(issues, fmt.Sprintf("key %q exceeds max length of %d", key, opts.MaxKeyLength))
			}
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
