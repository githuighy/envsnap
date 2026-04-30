package snapshot

import "strings"

// DefaultSensitiveKeys contains common environment variable name substrings
// that are likely to contain sensitive values.
var DefaultSensitiveKeys = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"AUTH",
	"CREDENTIAL",
	"ACCESS_KEY",
}

// RedactOptions controls how sensitive values are redacted.
type RedactOptions struct {
	// SensitiveKeys is the list of substrings to match against variable names.
	// Matching is case-insensitive.
	SensitiveKeys []string
	// Placeholder is the string used to replace sensitive values.
	// Defaults to "[REDACTED]" if empty.
	Placeholder string
}

// Redact returns a new Snapshot with sensitive variable values replaced by a
// placeholder string. The original snapshot is not modified.
func Redact(snap map[string]string, opts RedactOptions) map[string]string {
	if len(opts.SensitiveKeys) == 0 {
		opts.SensitiveKeys = DefaultSensitiveKeys
	}
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "[REDACTED]"
	}

	result := make(map[string]string, len(snap))
	for k, v := range snap {
		if isSensitive(k, opts.SensitiveKeys) {
			result[k] = placeholder
		} else {
			result[k] = v
		}
	}
	return result
}

// RedactedKeys returns the list of keys from snap that would be considered
// sensitive and redacted using the given options. This is useful for auditing
// which variables will be masked before performing a full redaction.
func RedactedKeys(snap map[string]string, opts RedactOptions) []string {
	if len(opts.SensitiveKeys) == 0 {
		opts.SensitiveKeys = DefaultSensitiveKeys
	}
	var keys []string
	for k := range snap {
		if isSensitive(k, opts.SensitiveKeys) {
			keys = append(keys, k)
		}
	}
	return keys
}

// isSensitive reports whether the key contains any of the given substrings,
// using a case-insensitive comparison.
func isSensitive(key string, sensitiveKeys []string) bool {
	upper := strings.ToUpper(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}
