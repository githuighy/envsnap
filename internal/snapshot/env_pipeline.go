package snapshot

import "fmt"

// PipelineStep represents a named transformation applied to a snapshot.
type PipelineStep struct {
	Name    string
	ApplyFn func(Snapshot) (Snapshot, error)
}

// Pipeline is an ordered sequence of steps applied to a snapshot.
type Pipeline struct {
	steps []PipelineStep
}

// NewPipeline creates an empty pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// AddStep appends a named step to the pipeline.
func (p *Pipeline) AddStep(name string, fn func(Snapshot) (Snapshot, error)) {
	p.steps = append(p.steps, PipelineStep{Name: name, ApplyFn: fn})
}

// PipelineResult holds the snapshot after each step.
type PipelineResult struct {
	StepName string
	Snap     Snapshot
}

// Run executes all steps in order, returning per-step results.
// Execution stops and returns an error if any step fails.
func (p *Pipeline) Run(initial Snapshot) ([]PipelineResult, error) {
	results := make([]PipelineResult, 0, len(p.steps))
	current := initial
	for _, step := range p.steps {
		out, err := step.ApplyFn(current)
		if err != nil {
			return results, fmt.Errorf("pipeline step %q failed: %w", step.Name, err)
		}
		results = append(results, PipelineResult{StepName: step.Name, Snap: out})
		current = out
	}
	return results, nil
}

// Final returns the snapshot produced by the last step, or the initial
// snapshot if the pipeline has no steps.
func (p *Pipeline) Final(initial Snapshot) (Snapshot, error) {
	results, err := p.Run(initial)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return initial, nil
	}
	return results[len(results)-1].Snap, nil
}
