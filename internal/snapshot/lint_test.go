package snapshot

import (
	"testing"
)

func TestLint_LowercaseKey(t *testing.T) {
	snap := map[string]string{"my_var": "value", "GOOD_VAR": "ok"}
	issues := Lint(snap, LintOptions{WarnLowercase: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != RuleLowercaseKey {
		t.Errorf("expected rule %q, got %q", RuleLowercaseKey, issues[0].Rule)
	}
	if issues[0].Key != "my_var" {
		t.Errorf("expected key 'my_var', got %q", issues[0].Key)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	snap := map[string]string{"KEY_A": "", "KEY_B": "   ", "KEY_C": "set"}
	issues := Lint(snap, LintOptions{WarnEmpty: true})
	if len(issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(issues))
	}
	for _, issue := range issues {
		if issue.Rule != RuleEmptyValue {
			t.Errorf("expected rule %q, got %q", RuleEmptyValue, issue.Rule)
		}
	}
}

func TestLint_WhitespaceKey(t *testing.T) {
	snap := map[string]string{"BAD KEY": "val", "GOOD": "val"}
	issues := Lint(snap, LintOptions{WarnWhitespace: true})
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Rule != RuleWhitespaceKey {
		t.Errorf("expected rule %q, got %q", RuleWhitespaceKey, issues[0].Rule)
	}
}

func TestLint_NoIssues(t *testing.T) {
	snap := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	issues := Lint(snap, DefaultLintOptions())
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %+v", len(issues), issues)
	}
}

func TestLint_DisabledRules(t *testing.T) {
	snap := map[string]string{"lowercase": "", "BAD KEY": "val"}
	// all rules disabled
	issues := Lint(snap, LintOptions{})
	if len(issues) != 0 {
		t.Errorf("expected no issues with all rules disabled, got %d", len(issues))
	}
}

func TestLint_MultipleRulesOnSameKey(t *testing.T) {
	snap := map[string]string{"bad key": ""}
	issues := Lint(snap, DefaultLintOptions())
	// should trigger whitespace + lowercase + empty
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues for 'bad key' with empty value, got %d", len(issues))
	}
}
