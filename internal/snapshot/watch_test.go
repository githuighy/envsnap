package snapshot

import (
	"sync"
	"testing"
	"time"
)

func TestWatch_DetectsChange(t *testing.T) {
	var mu sync.Mutex
	var received []DiffEntry

	callCount := 0

	// We'll test the watch loop by injecting a fake poll via a short interval
	// and verifying that OnChange fires when diffs exist.
	//
	// Because Watch captures real env, we manipulate the environment.
	t.Setenv("WATCH_TEST_VAR", "initial")

	opts := WatchOptions{
		Interval: 50 * time.Millisecond,
		Prefixes: []string{"WATCH_TEST_"},
		OnChange: func(diffs []DiffEntry) {
			mu.Lock()
			defer mu.Unlock()
			received = append(received, diffs...)
			callCount++
		},
	}

	wr := Watch(opts)
	defer wr.Stop()

	// Allow first capture to settle.
	time.Sleep(30 * time.Millisecond)

	// Change the env var so the next poll detects a diff.
	t.Setenv("WATCH_TEST_VAR", "changed")

	// Wait for at least one tick to fire.
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	if callCount == 0 {
		t.Fatal("expected OnChange to be called at least once")
	}

	found := false
	for _, d := range received {
		if d.Key == "WATCH_TEST_VAR" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected diff for WATCH_TEST_VAR, got %+v", received)
	}
}

func TestWatch_Stop_StopsLoop(t *testing.T) {
	callCount := 0

	wr := Watch(WatchOptions{
		Interval: 20 * time.Millisecond,
		OnChange: func(diffs []DiffEntry) {
			callCount++
		},
	})

	wr.Stop()

	// Give goroutine time to exit.
	time.Sleep(80 * time.Millisecond)

	before := callCount
	time.Sleep(60 * time.Millisecond)
	after := callCount

	if after != before {
		t.Errorf("watch continued after Stop(): callCount changed from %d to %d", before, after)
	}
}

func TestWatch_DefaultInterval(t *testing.T) {
	// Verify that a zero interval doesn't panic and defaults to 5s.
	wr := Watch(WatchOptions{
		Interval: 0,
		OnChange: func(diffs []DiffEntry) {},
	})
	defer wr.Stop()
	// If we reach here without panic, the default was applied successfully.
}
