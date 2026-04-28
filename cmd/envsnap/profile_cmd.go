package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/user/envsnap/internal/snapshot"
)

func profileStore() *snapshot.ProfileStore {
	dir := os.Getenv("ENVSNAP_PROFILE_DIR")
	if dir == "" {
		dir = ".envsnap/profiles"
	}
	return snapshot.NewProfileStore(dir)
}

func runProfileSave(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap profile save <name> [--prefix P1,P2] [--exclude E1] [--redact R1]")
	}
	name := args[0]
	p := snapshot.Profile{Name: name}

	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--prefix":
			i++
			if i < len(args) {
				p.Prefixes = strings.Split(args[i], ",")
			}
		case "--exclude":
			i++
			if i < len(args) {
				p.Exclude = strings.Split(args[i], ",")
			}
		case "--redact":
			i++
			if i < len(args) {
				p.Redact = strings.Split(args[i], ",")
			}
		}
	}

	if err := profileStore().Save(p); err != nil {
		return fmt.Errorf("save profile: %w", err)
	}
	fmt.Printf("profile %q saved\n", name)
	return nil
}

func runProfileList() error {
	names, err := profileStore().List()
	if err != nil {
		return fmt.Errorf("list profiles: %w", err)
	}
	if len(names) == 0 {
		fmt.Println("no profiles found")
		return nil
	}
	for _, n := range names {
		fmt.Println(n)
	}
	return nil
}

func runProfileShow(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap profile show <name>")
	}
	p, err := profileStore().Load(args[0])
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(p, "", "  ")
	fmt.Println(string(data))
	return nil
}

func runProfileDelete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap profile delete <name>")
	}
	if err := profileStore().Delete(args[0]); err != nil {
		return err
	}
	fmt.Printf("profile %q deleted\n", args[0])
	return nil
}
