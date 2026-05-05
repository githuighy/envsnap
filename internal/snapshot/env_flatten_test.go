package snapshot

import (
	"testing"
)

func makeFlatSnap(vars map[string]string) Snapshot {
	return Snapshot{Vars: vars}
}

func TestFlatten_BasicLabels(t *testing.T) {
	a := makeFlatSnap(map[string]string{"HOST": "localhost", "PORT": "5432"})
	b := makeFlatSnap(map[string]string{"HOST": "remotehost", "USER": "admin"})

	opts := DefaultFlattenOptions()
	result, err := Flatten([]Snapshot{a, b}, []string{"DB", "CACHE"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := map[string]string{
		"DB_HOST":    "localhost",
		"DB_PORT":    "5432",
		"CACHE_HOST": "remotehost",
		"CACHE_USER": "admin",
	}
	for k, want := range expected {
		got, ok := result.Vars[k]
		if !ok {
			t.Errorf("missing key %q", k)
			continue
		}
		if got != want {
			t.Errorf("key %q: got %q, want %q", k, got, want)
		}
	}
}

func TestFlatten_PrefixApplied(t *testing.T) {
	a := makeFlatSnap(map[string]string{"KEY": "val"})
	opts := DefaultFlattenOptions()
	opts.Prefix = "SNAP"

	result, err := Flatten([]Snapshot{a}, []string{"SVC"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result.Vars["SNAP_SVC_KEY"]; !ok {
		t.Errorf("expected key SNAP_SVC_KEY, got vars: %v", result.Vars)
	}
}

func TestFlatten_SkipEmpty(t *testing.T) {
	a := makeFlatSnap(map[string]string{"PRESENT": "yes", "EMPTY": ""})
	opts := DefaultFlattenOptions()
	opts.SkipEmpty = true

	result, err := Flatten([]Snapshot{a}, []string{"SVC"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result.Vars["SVC_EMPTY"]; ok {
		t.Error("expected SVC_EMPTY to be skipped")
	}
	if _, ok := result.Vars["SVC_PRESENT"]; !ok {
		t.Error("expected SVC_PRESENT to be present")
	}
}

func TestFlatten_LabelsMismatch_ReturnsError(t *testing.T) {
	a := makeFlatSnap(map[string]string{"K": "v"})
	b := makeFlatSnap(map[string]string{"K": "v"})

	_, err := Flatten([]Snapshot{a, b}, []string{"ONLY_ONE"}, DefaultFlattenOptions())
	if err == nil {
		t.Error("expected error for mismatched labels, got nil")
	}
}

func TestFlatten_EmptyLabel_ReturnsError(t *testing.T) {
	a := makeFlatSnap(map[string]string{"K": "v"})

	_, err := Flatten([]Snapshot{a}, []string{""}, DefaultFlattenOptions())
	if err == nil {
		t.Error("expected error for empty label, got nil")
	}
}

func TestFlatten_NoLabels_UsesIndex(t *testing.T) {
	a := makeFlatSnap(map[string]string{"KEY": "val"})
	opts := DefaultFlattenOptions()

	result, err := Flatten([]Snapshot{a}, nil, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result.Vars["0_KEY"]; !ok {
		t.Errorf("expected key 0_KEY, got vars: %v", result.Vars)
	}
}
