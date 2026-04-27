package snapshot

import (
	"time"
)

// WatchOptions configures the Watch behavior.
type WatchOptions struct {
	Interval  time.Duration
	Prefixes  []string
	Exclude   []string
	OnChange  func(diffs []DiffEntry)
	OnError   func(err error)
}

// WatchResult holds state for a running watch session.
type WatchResult struct {
	stop chan struct{}
}

// Stop halts the watch loop.
func (w *WatchResult) Stop() {
	close(w.stop)
}

// Watch polls the environment at the given interval and calls OnChange
// whenever the environment differs from the previously captured snapshot.
func Watch(opts WatchOptions) *WatchResult {
	if opts.Interval <= 0 {
		opts.Interval = 5 * time.Second
	}

	wr := &WatchResult{stop: make(chan struct{})}

	go func() {
		prev := captureFiltered(opts.Prefixes, opts.Exclude)

		ticker := time.NewTicker(opts.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-wr.stop:
				return
			case <-ticker.C:
				curr := captureFiltered(opts.Prefixes, opts.Exclude)
				diffs := Diff(prev, curr)
				if len(diffs) > 0 && opts.OnChange != nil {
					opts.OnChange(diffs)
				}
				prev = curr
			}
		}
	}()

	return wr
}

// captureFiltered captures the environment and applies optional filtering.
func captureFiltered(prefixes, exclude []string) Snapshot {
	snap := Capture()
	if len(prefixes) == 0 && len(exclude) == 0 {
		return snap
	}
	return Filter(snap, FilterOptions{
		Prefixes: prefixes,
		Exclude:  exclude,
	})
}
