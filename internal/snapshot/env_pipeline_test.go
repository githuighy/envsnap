package snapshot

import (
	"errors"
	"testing"
)

func basePipelineSnap() Snapshot {
	return Snapshot{
		"APP_ENV":    "production",
		"DB_PASSWORD": "secret",
		"DEBUG":      "true",
	}
}

func TestPipeline_EmptyPipeline_ReturnsInitial(t *testing.T) {
	p := NewPipeline()
	snap := basePipelineSnap()
	out, err := p.Final(snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", out["APP_ENV"])
	}
}

func TestPipeline_StepsAppliedInOrder(t *testing.T) {
	p := NewPipeline()
	p.AddStep("add-key", func(s Snapshot) (Snapshot, error) {
		out := copySnap(s)
		out["STEP1"] = "yes"
		return out, nil
	})
	p.AddStep("add-key2", func(s Snapshot) (Snapshot, error) {
		out := copySnap(s)
		out["STEP2"] = s["STEP1"] + "-done"
		return out, nil
	})
	out, err := p.Final(basePipelineSnap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["STEP2"] != "yes-done" {
		t.Errorf("expected STEP2=yes-done, got %q", out["STEP2"])
	}
}

func TestPipeline_StepError_Propagates(t *testing.T) {
	p := NewPipeline()
	p.AddStep("ok", func(s Snapshot) (Snapshot, error) { return s, nil })
	p.AddStep("fail", func(s Snapshot) (Snapshot, error) {
		return nil, errors.New("boom")
	})
	p.AddStep("never", func(s Snapshot) (Snapshot, error) { return s, nil })
	_, err := p.Final(basePipelineSnap())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestPipeline_RunReturnsPerStepResults(t *testing.T) {
	p := NewPipeline()
	p.AddStep("filter", func(s Snapshot) (Snapshot, error) {
		return Filter(s, FilterOptions{Prefixes: []string{"APP"}}), nil
	})
	results, err := p.Run(basePipelineSnap())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].StepName != "filter" {
		t.Errorf("expected step name 'filter', got %q", results[0].StepName)
	}
	if _, ok := results[0].Snap["DB_PASSWORD"]; ok {
		t.Error("DB_PASSWORD should have been filtered out")
	}
}

// copySnap makes a shallow copy of a snapshot.
func copySnap(s Snapshot) Snapshot {
	out := make(Snapshot, len(s))
	for k, v := range s {
		out[k] = v
	}
	return out
}
