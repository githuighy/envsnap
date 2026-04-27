package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

func runHistoryAdd(args []string, label, historyDir string) error {
	if len(args) < 1 {
		return fmt.Errorf("history add: snapshot file required")
	}
	snap, err := snapshot.Load(args[0])
	if err != nil {
		return fmt.Errorf("history add: load snapshot: %w", err)
	}
	store := snapshot.NewHistoryStore(historyDir)
	id, err := store.Add(snap, label)
	if err != nil {
		return fmt.Errorf("history add: %w", err)
	}
	fmt.Printf("Saved snapshot as history entry: %s\n", id)
	return nil
}

func runHistoryList(historyDir, format string) error {
	store := snapshot.NewHistoryStore(historyDir)
	entries, err := store.List()
	if err != nil {
		return fmt.Errorf("history list: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("No history entries found.")
		return nil
	}
	if format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(entries)
	}
	for _, e := range entries {
		label := e.Label
		if label == "" {
			label = "(no label)"
		}
		fmt.Printf("%s  %s  %s  (%d vars)\n",
			e.ID, e.Timestamp.Format("2006-01-02 15:04:05"), label, len(e.Vars))
	}
	return nil
}

func runHistoryDiff(historyDir, idA, idB, format string) error {
	store := snapshot.NewHistoryStore(historyDir)
	a, err := store.Get(idA)
	if err != nil {
		return fmt.Errorf("history diff: get %s: %w", idA, err)
	}
	b, err := store.Get(idB)
	if err != nil {
		return fmt.Errorf("history diff: get %s: %w", idB, err)
	}
	diffs := snapshot.Diff(a.Vars, b.Vars)
	out, err := snapshot.RenderDiff(diffs, format)
	if err != nil {
		return fmt.Errorf("history diff: render: %w", err)
	}
	fmt.Print(out)
	return nil
}
