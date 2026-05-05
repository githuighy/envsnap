package snapshot

import (
	"fmt"
	"strings"
)

// MaskOptions controls how values are masked in a snapshot.
type MaskOptions struct {
	// Keys is the explicit list of keys to mask.
	Keys []string
	// Prefixes masks any key matching one of these prefixes.
	Prefixes []string
	// ShowLength appends the character count of the original value.
	ShowLength bool
	// Placeholder overrides the default mask string.
	Placeholder string
	// VisibleChars shows this many trailing chars of the original value (0 = none).
	VisibleChars int
}

// Mask returns a new snapshot with selected values obscured.
// Keys not targeted by opts are left unchanged.
func Mask(snap Snapshot, opts MaskOptions) Snapshot {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	masked := make(Snapshot, len(snap))
	for k, v := range snap {
		if shouldMask(k, opts) {
			masked[k] = buildMasked(v, placeholder, opts)
		} else {
			masked[k] = v
		}
	}
	return masked
}

// MaskedKeys returns the sorted list of keys that would be masked given opts.
func MaskedKeys(snap Snapshot, opts MaskOptions) []string {
	var keys []string
	for k := range snap {
		if shouldMask(k, opts) {
			keys = append(keys, k)
		}
	}
	return keys
}

func shouldMask(key string, opts MaskOptions) bool {
	for _, k := range opts.Keys {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	for _, p := range opts.Prefixes {
		if strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

func buildMasked(original, placeholder string, opts MaskOptions) string {
	result := placeholder
	if opts.VisibleChars > 0 && len(original) > opts.VisibleChars {
		result = placeholder + original[len(original)-opts.VisibleChars:]
	}
	if opts.ShowLength {
		result = fmt.Sprintf("%s(%d)", result, len(original))
	}
	return result
}
