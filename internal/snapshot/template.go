package snapshot

import (
	"fmt"
	"regexp"
	"strings"
)

// Template represents a snapshot template with required keys and default values.
type Template struct {
	Name     string            `json:"name"`
	Defaults map[string]string `json:"defaults"`
	Required []string          `json:"required"`
}

var validTemplateName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ApplyTemplate merges template defaults into a snapshot, then validates
// that all required keys are present and non-empty.
// Existing keys in snap are NOT overwritten by defaults.
func ApplyTemplate(snap Snapshot, tmpl Template) (Snapshot, error) {
	if !validTemplateName.MatchString(tmpl.Name) {
		return nil, fmt.Errorf("invalid template name %q: must match [a-zA-Z0-9_-]+", tmpl.Name)
	}

	result := make(Snapshot, len(snap))
	for k, v := range snap {
		result[k] = v
	}

	for k, v := range tmpl.Defaults {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	var missing []string
	for _, key := range tmpl.Required {
		if val, ok := result[key]; !ok || strings.TrimSpace(val) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("template %q: missing required keys: %s", tmpl.Name, strings.Join(missing, ", "))
	}

	return result, nil
}

// StripTemplateDefaults removes keys from snap whose values exactly match
// the template defaults and are not in the required list.
func StripTemplateDefaults(snap Snapshot, tmpl Template) Snapshot {
	requiredSet := make(map[string]bool, len(tmpl.Required))
	for _, k := range tmpl.Required {
		requiredSet[k] = true
	}

	result := make(Snapshot, len(snap))
	for k, v := range snap {
		defVal, isDefault := tmpl.Defaults[k]
		if isDefault && v == defVal && !requiredSet[k] {
			continue
		}
		result[k] = v
	}
	return result
}
