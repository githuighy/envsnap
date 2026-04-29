package snapshot

import (
	"math"
	"strings"
)

// ScoreOptions controls how the snapshot health score is calculated.
type ScoreOptions struct {
	// RequiredKeys that must be present and non-empty for full score.
	RequiredKeys []string
	// ForbiddenKeys that must not be present.
	ForbiddenKeys []string
	// LintOptions used to penalise lint issues.
	LintOptions *LintOptions
}

// ScoreResult holds the computed score and a breakdown of deductions.
type ScoreResult struct {
	Score      int      // 0-100
	Deductions []string // human-readable reasons for point loss
}

// Score evaluates the health of a snapshot and returns a 0-100 score.
func Score(snap map[string]string, opts ScoreOptions) ScoreResult {
	deductions := []string{}
	total := 100.0

	// Required keys: -10 per missing or empty key, capped at -40.
	missingPenalty := 0.0
	for _, k := range opts.RequiredKeys {
		v, ok := snap[k]
		if !ok || strings.TrimSpace(v) == "" {
			missingPenalty += 10
			deductions = append(deductions, "missing required key: "+k)
		}
	}
	if missingPenalty > 40 {
		missingPenalty = 40
	}
	total -= missingPenalty

	// Forbidden keys: -5 per present key, capped at -20.
	forbidPenalty := 0.0
	for _, k := range opts.ForbiddenKeys {
		if _, ok := snap[k]; ok {
			forbidPenalty += 5
			deductions = append(deductions, "forbidden key present: "+k)
		}
	}
	if forbidPenalty > 20 {
		forbidPenalty = 20
	}
	total -= forbidPenalty

	// Lint issues: -2 per issue, capped at -20.
	lo := opts.LintOptions
	if lo == nil {
		default_ := DefaultLintOptions()
		lo = &default_
	}
	issues := Lint(snap, *lo)
	lintPenalty := math.Min(float64(len(issues))*2, 20)
	for _, iss := range issues {
		deductions = append(deductions, "lint: "+iss.Message)
	}
	total -= lintPenalty

	if total < 0 {
		total = 0
	}
	return ScoreResult{
		Score:      int(math.Round(total)),
		Deductions: deductions,
	}
}
