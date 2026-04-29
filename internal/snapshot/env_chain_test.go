package snapshot

import (
	"testing"
)

func TestChain_Resolve_PriorityOrder(t *testing.T) {
	high := Snapshot{"A": "high", "B": "high"}
	low := Snapshot{"A": "low", "C": "low"}
	chain := NewChain(high, low)
	resolved := chain.Resolve()

	if resolved["A"] != "high" {
		t.Errorf("expected A=high, got %s", resolved["A"])
	}
	if resolved["B"] != "high" {
		t.Errorf("expected B=high, got %s", resolved["B"])
	}
	if resolved["C"] != "low" {
		t.Errorf("expected C=low, got %s", resolved["C"])
	}
}

func TestChain_ResolveKey_FindsLayer(t *testing.T) {
	layer0 := Snapshot{"X": "from-layer0"}
	layer1 := Snapshot{"Y": "from-layer1"}
	chain := NewChain(layer0, layer1)

	v, idx := chain.ResolveKey("X")
	if v != "from-layer0" || idx != 0 {
		t.Errorf("expected from-layer0/0, got %s/%d", v, idx)
	}

	v, idx = chain.ResolveKey("Y")
	if v != "from-layer1" || idx != 1 {
		t.Errorf("expected from-layer1/1, got %s/%d", v, idx)
	}
}

func TestChain_ResolveKey_Missing(t *testing.T) {
	chain := NewChain(Snapshot{"A": "1"})
	v, idx := chain.ResolveKey("NOPE")
	if v != "" || idx != -1 {
		t.Errorf("expected empty/-1, got %q/%d", v, idx)
	}
}

func TestChain_Explain(t *testing.T) {
	layer0 := Snapshot{"A": "1"}
	layer1 := Snapshot{"A": "2", "B": "3"}
	chain := NewChain(layer0, layer1)
	origin := chain.Explain()

	if origin["A"] != 0 {
		t.Errorf("expected A from layer 0, got %d", origin["A"])
	}
	if origin["B"] != 1 {
		t.Errorf("expected B from layer 1, got %d", origin["B"])
	}
}

func TestChain_Validate_NilLayer(t *testing.T) {
	chain := NewChain(Snapshot{"A": "1"}, nil)
	if err := chain.Validate(); err == nil {
		t.Error("expected error for nil layer")
	}
}

func TestChain_Validate_OK(t *testing.T) {
	chain := NewChain(Snapshot{"A": "1"}, Snapshot{"B": "2"})
	if err := chain.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
