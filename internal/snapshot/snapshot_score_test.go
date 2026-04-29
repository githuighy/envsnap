package snapshot

import (
	"strings"
	"testing"
)

func baseScoreSnap() map[string]string {
	return map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_URL":   "postgres://localhost/db",
	}
}

func TestScore_PerfectSnap(t *testing.T) {
	snap := baseScoreSnap()
	res := Score(snap, ScoreOptions{
		RequiredKeys: []string{"APP_HOST", "APP_PORT"},
	})
	if res.Score != 100 {
		t.Errorf("expected 100, got %d (deductions: %v)", res.Score, res.Deductions)
	}
	if len(res.Deductions) != 0 {
		t.Errorf("expected no deductions, got %v", res.Deductions)
	}
}

func TestScore_MissingRequiredKey(t *testing.T) {
	snap := baseScoreSnap()
	res := Score(snap, ScoreOptions{
		RequiredKeys: []string{"APP_HOST", "MISSING_KEY"},
	})
	if res.Score != 90 {
		t.Errorf("expected 90, got %d", res.Score)
	}
	if len(res.Deductions) != 1 {
		t.Errorf("expected 1 deduction, got %v", res.Deductions)
	}
}

func TestScore_ForbiddenKeyPresent(t *testing.T) {
	snap := baseScoreSnap()
	snap["DEBUG"] = "true"
	res := Score(snap, ScoreOptions{
		ForbiddenKeys: []string{"DEBUG"},
	})
	if res.Score != 95 {
		t.Errorf("expected 95, got %d", res.Score)
	}
	found := false
	for _, d := range res.Deductions {
		if strings.Contains(d, "DEBUG") {
			found = true
		}
	}
	if !found {
		t.Error("expected deduction mentioning DEBUG")
	}
}

func TestScore_LintIssuesPenalised(t *testing.T) {
	snap := map[string]string{
		"lowercase_key": "value",
		"APP_HOST":      "",
	}
	lo := DefaultLintOptions()
	res := Score(snap, ScoreOptions{LintOptions: &lo})
	if res.Score >= 100 {
		t.Errorf("expected score below 100 due to lint issues, got %d", res.Score)
	}
}

func TestScore_CapsAt100(t *testing.T) {
	snap := baseScoreSnap()
	res := Score(snap, ScoreOptions{})
	if res.Score > 100 {
		t.Errorf("score should not exceed 100, got %d", res.Score)
	}
}

func TestScore_NeverBelowZero(t *testing.T) {
	snap := map[string]string{"a": ""}
	required := make([]string, 10)
	for i := range required {
		required[i] = "missing_"+string(rune('a'+i))
	}
	res := Score(snap, ScoreOptions{RequiredKeys: required})
	if res.Score < 0 {
		t.Errorf("score should not be negative, got %d", res.Score)
	}
}
