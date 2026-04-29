package snapshot

import (
	"fmt"
	"strings"
)

// LintRule represents a single lint check.
type LintRule string

const (
	RuleLowercaseKey   LintRule = "lowercase-key"
	RuleEmptyValue     LintRule = "empty-value"
	RuleWhitespaceKey  LintRule = "whitespace-key"
	RuleDuplicateKey   LintRule = "duplicate-key"
)

// LintIssue describes a single lint violation.
type LintIssue struct {
	Key     string
	Rule    LintRule
	Message string
}

// LintOptions controls which rules are applied.
type LintOptions struct {
	WarnLowercase  bool
	WarnEmpty      bool
	WarnWhitespace bool
}

// DefaultLintOptions returns sensible defaults.
func DefaultLintOptions() LintOptions {
	return LintOptions{
		WarnLowercase:  true,
		WarnEmpty:      true,
		WarnWhitespace: true,
	}
}

// Lint checks a snapshot for common environment variable issues.
func Lint(snap map[string]string, opts LintOptions) []LintIssue {
	var issues []LintIssue
	seen := make(map[string]bool)

	for k, v := range snap {
		normalized := strings.ToUpper(k)
		if seen[normalized] {
			issues = append(issues, LintIssue{
				Key:     k,
				Rule:    RuleDuplicateKey,
				Message: fmt.Sprintf("key %q appears more than once (case-insensitive)", k),
			})
		}
		seen[normalized] = true

		if opts.WarnWhitespace && (strings.TrimSpace(k) != k || strings.ContainsAny(k, " \t")) {
			issues = append(issues, LintIssue{
				Key:     k,
				Rule:    RuleWhitespaceKey,
				Message: fmt.Sprintf("key %q contains whitespace", k),
			})
		}

		if opts.WarnLowercase && k != strings.ToUpper(k) {
			issues = append(issues, LintIssue{
				Key:     k,
				Rule:    RuleLowercaseKey,
				Message: fmt.Sprintf("key %q is not uppercase", k),
			})
		}

		if opts.WarnEmpty && strings.TrimSpace(v) == "" {
			issues = append(issues, LintIssue{
				Key:     k,
				Rule:    RuleEmptyValue,
				Message: fmt.Sprintf("key %q has an empty or blank value", k),
			})
		}
	}

	return issues
}
