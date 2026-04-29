package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

const groupStoreDir = ".envsnap/groups"

func loadGroupStore() (*snapshot.GroupStore, error) {
	s := snapshot.NewGroupStore()
	if err := s.LoadGroups(groupStoreDir); err != nil {
		return nil, err
	}
	return s, nil
}

func runGroupAdd(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap group add <name> <PREFIX1> [PREFIX2...]")
	}
	name := args[0]
	patterns := args[1:]
	s, err := loadGroupStore()
	if err != nil {
		return err
	}
	if err := s.Add(name, patterns); err != nil {
		return err
	}
	return s.SaveGroups(groupStoreDir)
}

func runGroupDelete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap group delete <name>")
	}
	s, err := loadGroupStore()
	if err != nil {
		return err
	}
	if err := s.Delete(args[0]); err != nil {
		return err
	}
	return s.SaveGroups(groupStoreDir)
}

func runGroupList(args []string) error {
	s, err := loadGroupStore()
	if err != nil {
		return err
	}
	names := s.List()
	if len(names) == 0 {
		fmt.Println("no groups defined")
		return nil
	}
	for _, name := range names {
		g, _ := s.Get(name)
		fmt.Printf("%-20s %s\n", name, strings.Join(g.Patterns, ", "))
	}
	return nil
}

func runGroupExtract(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap group extract <name> <snapshot-file>")
	}
	groupName, snapFile := args[0], args[1]
	snap, err := snapshot.Load(snapFile)
	if err != nil {
		return fmt.Errorf("load snapshot: %w", err)
	}
	s, err := loadGroupStore()
	if err != nil {
		return err
	}
	result, err := s.ExtractGroup(groupName, snap)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}
