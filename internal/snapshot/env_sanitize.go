package snapshot

import (
	"regexp"
	"strings"
)

// SanitizeOptions controls how snapshot values are sanitized.
type SanitizeOptions struct {
	// StripControlChars removes non-printable control characters from values.
	StripControlChars bool
	// TrimQuotes removes surrounding single or double quotes from values.
	TrimQuotes bool
	// CollapseWhitespace replaces runs of internal whitespace with a single space.
	CollapseWhitespace bool
	// Prefixes restricts sanitization to keys with these prefixes. Empty means all keys.
	Prefixes []string
	// ExcludeKeys skips sanitization for these exact keys.
	ExcludeKeys []string
}

// DefaultSanitizeOptions returns a SanitizeOptions with sensible defaults.
func DefaultSanitizeOptions() SanitizeOptions {
	return SanitizeOptions{
		StripControlChars:  true,
		TrimQuotes:         false,
		CollapseWhitespace: false,
	}
}

var controlCharRe = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)
var whitespaceRe = regexp.MustCompile(`[ \t]+`)

// Sanitize cleans up values in snap according to opts and returns a new snapshot.
func Sanitize(snap Snapshot, opts SanitizeOptions) Snapshot {
	excluded := make(map[string]bool, len(opts.ExcludeKeys))
	for _, k := range opts.ExcludeKeys {
		excluded[k] = true
	}

	out := make(Snapshot, len(snap))
	for k, v := range snap {
		if excluded[k] {
			out[k] = v
			continue
		}
		if len(opts.Prefixes) > 0 && !hasAnyPrefix(k, opts.Prefixes) {
			out[k] = v
			continue
		}
		out[k] = sanitizeValue(v, opts)
	}
	return out
}

func sanitizeValue(v string, opts SanitizeOptions) string {
	if opts.StripControlChars {
		v = controlCharRe.ReplaceAllString(v, "")
	}
	if opts.TrimQuotes {
		if (strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`)) ||
			(strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) {
			v = v[1 : len(v)-1]
		}
	}
	if opts.CollapseWhitespace {
		v = whitespaceRe.ReplaceAllString(v, " ")
	}
	return v
}
