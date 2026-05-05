package snapshot

import (
	"strings"
	"unicode"
)

// NormalizeOptions controls how snapshot keys and values are normalized.
type NormalizeOptions struct {
	// UppercaseKeys converts all keys to uppercase.
	UppercaseKeys bool
	// TrimValues strips leading and trailing whitespace from values.
	TrimValues bool
	// RemoveEmptyValues drops keys whose values are empty after trimming.
	RemoveEmptyValues bool
	// SanitizeKeys replaces invalid characters in keys with underscores.
	SanitizeKeys bool
}

// DefaultNormalizeOptions returns the recommended normalization defaults.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		UppercaseKeys:     true,
		TrimValues:        true,
		RemoveEmptyValues: false,
		SanitizeKeys:      true,
	}
}

// Normalize applies the given options to a snapshot, returning a new
// snapshot with normalized keys and values. The original is not modified.
func Normalize(snap map[string]string, opts NormalizeOptions) map[string]string {
	out := make(map[string]string, len(snap))
	for k, v := range snap {
		if opts.SanitizeKeys {
			k = sanitizeKey(k)
		}
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		if opts.RemoveEmptyValues && v == "" {
			continue
		}
		out[k] = v
	}
	return out
}

// sanitizeKey replaces any character that is not a letter, digit, or
// underscore with an underscore, and strips leading digits.
func sanitizeKey(k string) string {
	var b strings.Builder
	for i, r := range k {
		switch {
		case unicode.IsLetter(r) || r == '_':
			b.WriteRune(r)
		case unicode.IsDigit(r):
			if i == 0 {
				b.WriteRune('_')
			} else {
				b.WriteRune(r)
			}
		default:
			b.WriteRune('_')
		}
	}
	return b.String()
}
