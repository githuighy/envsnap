package snapshot_test

import (
	"testing"

	"github.com/user/envsnap/internal/snapshot"
)

func savePipelineSnap(t *testing.T, dir, name string, s snapshot.Snapshot) string {
	t.Helper()
	path := dir + "/" + name
	if err := snapshot.Save(s, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	return path
}

func TestPipeline_FilterThenRedact_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	original := snapshot.Snapshot{
		"APP_SECRET": "mysecret",
		"APP_ENV":    "prod",
		"DB_HOST":    "db.internal",
	}
	path := savePipelineSnap(t, dir, "snap.json", original)

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	p := snapshot.NewPipeline()
	p.AddStep("filter-APP", func(s snapshot.Snapshot) (snapshot.Snapshot, error) {
		return snapshot.Filter(s, snapshot.FilterOptions{Prefixes: []string{"APP"}}), nil
	})
	p.AddStep("redact", func(s snapshot.Snapshot) (snapshot.Snapshot, error) {
		return snapshot.Redact(s, snapshot.RedactOptions{}), nil
	})

	results, err := p.Run(loaded)
	if err != nil {
		t.Fatalf("pipeline run: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	afterFilter := results[0].Snap
	if _, ok := afterFilter["DB_HOST"]; ok {
		t.Error("DB_HOST should be absent after filter step")
	}
	if _, ok := afterFilter["APP_ENV"]; !ok {
		t.Error("APP_ENV should be present after filter step")
	}

	afterRedact := results[1].Snap
	if afterRedact["APP_SECRET"] == "mysecret" {
		t.Error("APP_SECRET should be redacted in final step")
	}
}

func TestPipeline_LintStep_FailsOnBadSnap(t *testing.T) {
	bad := snapshot.Snapshot{
		"lowercase_key": "value",
	}
	p := snapshot.NewPipeline()
	p.AddStep("lint", func(s snapshot.Snapshot) (snapshot.Snapshot, error) {
		issues := snapshot.Lint(s, snapshot.DefaultLintOptions())
		if len(issues) > 0 {
			return nil, &lintPipelineError{count: len(issues)}
		}
		return s, nil
	})
	_, err := p.Final(bad)
	if err == nil {
		t.Fatal("expected lint step to fail on lowercase key")
	}
}

type lintPipelineError struct{ count int }

func (e *lintPipelineError) Error() string {
	return "lint issues found"
}
