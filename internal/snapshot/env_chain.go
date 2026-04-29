package snapshot

import "fmt"

// Chain represents an ordered list of snapshot sources to resolve variables from.
// Earlier entries take precedence over later ones (highest priority first).
type Chain struct {
	Layers []Snapshot
}

// NewChain creates a Chain from the provided snapshots in priority order.
func NewChain(layers ...Snapshot) *Chain {
	return &Chain{Layers: layers}
}

// Resolve returns a single merged Snapshot where earlier layers win.
func (c *Chain) Resolve() Snapshot {
	result := Snapshot{}
	// Iterate in reverse so higher-priority layers overwrite lower ones.
	for i := len(c.Layers) - 1; i >= 0; i-- {
		for k, v := range c.Layers[i] {
			result[k] = v
		}
	}
	return result
}

// ResolveKey returns the value and the layer index that provided it.
// Returns "", -1 if the key is not found in any layer.
func (c *Chain) ResolveKey(key string) (string, int) {
	for i, layer := range c.Layers {
		if v, ok := layer[key]; ok {
			return v, i
		}
	}
	return "", -1
}

// Explain returns a map of key -> layer index showing which layer each key
// originates from in the resolved snapshot.
func (c *Chain) Explain() map[string]int {
	origin := map[string]int{}
	for i := len(c.Layers) - 1; i >= 0; i-- {
		for k := range c.Layers[i] {
			origin[k] = i
		}
	}
	return origin
}

// Validate checks that all layers are non-nil maps and returns an error if not.
func (c *Chain) Validate() error {
	for i, l := range c.Layers {
		if l == nil {
			return fmt.Errorf("chain layer %d is nil", i)
		}
	}
	return nil
}
