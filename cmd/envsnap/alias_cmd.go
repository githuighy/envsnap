package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/envsnap/internal/snapshot"
)

var aliasStore = snapshot.NewAliasStore(".envsnap/aliases")

// runAliasSet sets an alias: envsnap alias set <name> <path>
func runAliasSet(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envsnap alias set <name> <snapshot-path>")
	}
	name, path := args[0], args[1]
	if err := aliasStore.Set(name, path); err != nil {
		return fmt.Errorf("alias set: %w", err)
	}
	fmt.Printf("alias %q -> %s\n", name, path)
	return nil
}

// runAliasResolve prints the path for a given alias.
func runAliasResolve(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap alias resolve <name>")
	}
	path, err := aliasStore.Resolve(args[0])
	if err != nil {
		return err
	}
	fmt.Println(path)
	return nil
}

// runAliasDelete removes an alias by name.
func runAliasDelete(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envsnap alias delete <name>")
	}
	if err := aliasStore.Delete(args[0]); err != nil {
		return fmt.Errorf("alias delete: %w", err)
	}
	fmt.Printf("alias %q deleted\n", args[0])
	return nil
}

// runAliasList lists all aliases. Accepts --format=json flag.
func runAliasList(args []string) error {
	format := "text"
	for _, a := range args {
		if a == "--format=json" {
			format = "json"
		}
	}

	aliases, err := aliasStore.List()
	if err != nil {
		return fmt.Errorf("alias list: %w", err)
	}

	if format == "json" {
		return json.NewEncoder(os.Stdout).Encode(aliases)
	}

	if len(aliases) == 0 {
		fmt.Println("no aliases defined")
		return nil
	}
	for _, a := range aliases {
		fmt.Printf("%-20s %s\n", a.Name, a.Path)
	}
	return nil
}
