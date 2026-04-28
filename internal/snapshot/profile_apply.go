package snapshot

// ApplyProfile captures a snapshot and applies the profile's filter and redact
// settings in one step.
func ApplyProfile(p Profile) (Snapshot, error) {
	snap := Capture()

	fopts := FilterOptions{}
	if len(p.Prefixes) > 0 {
		fopts.Prefixes = p.Prefixes
	}
	if len(p.Exclude) > 0 {
		fopts.Exclude = p.Exclude
	}

	if len(fopts.Prefixes) > 0 || len(fopts.Exclude) > 0 {
		snap = Filter(snap, fopts)
	}

	if len(p.Redact) > 0 {
		snap = Redact(snap, RedactOptions{
			SensitiveKeys: p.Redact,
		})
	}

	return snap, nil
}
