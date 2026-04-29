package snapshot

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaFieldType represents the expected type of an environment variable value.
type SchemaFieldType string

const (
	FieldTypeString SchemaFieldType = "string"
	FieldTypeInt    SchemaFieldType = "int"
	FieldTypeBool   SchemaFieldType = "bool"
	FieldTypeURL    SchemaFieldType = "url"
)

// SchemaField defines constraints for a single environment variable.
type SchemaField struct {
	// Type is the expected value type. Defaults to string if empty.
	Type SchemaFieldType `json:"type,omitempty"`
	// Required indicates the key must be present and non-empty.
	Required bool `json:"required,omitempty"`
	// Pattern is an optional regex the value must match.
	Pattern string `json:"pattern,omitempty"`
	// AllowedValues restricts the value to one of the listed options.
	AllowedValues []string `json:"allowed_values,omitempty"`
}

// Schema maps environment variable names to their field definitions.
type Schema map[string]SchemaField

// SchemaViolation describes a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// ValidateSchema checks a snapshot against a schema definition and returns
// any violations found. An empty violation slice means the snapshot is valid.
func ValidateSchema(snap Snapshot, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for key, field := range schema {
		val, exists := snap[key]

		// Check required presence.
		if field.Required && (!exists || strings.TrimSpace(val) == "") {
			violations = append(violations, SchemaViolation{
				Key:     key,
				Message: "required key is missing or empty",
			})
			continue
		}

		// Skip further checks if the key is absent and not required.
		if !exists {
			continue
		}

		// Check type constraints.
		if err := checkFieldType(key, val, field.Type); err != nil {
			violations = append(violations, *err)
		}

		// Check regex pattern.
		if field.Pattern != "" {
			matched, err := regexp.MatchString(field.Pattern, val)
			if err != nil {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: fmt.Sprintf("invalid pattern %q: %v", field.Pattern, err),
				})
			} else if !matched {
				violations = append(violations, SchemaViolation{
					Key:     key,
					Message: fmt.Sprintf("value %q does not match pattern %q", val, field.Pattern),
				})
			}
		}

		// Check allowed values.
		if len(field.AllowedValues) > 0 && !containsString(field.AllowedValues, val) {
			violations = append(violations, SchemaViolation{
				Key:     key,
				Message: fmt.Sprintf("value %q is not one of allowed values %v", val, field.AllowedValues),
			})
		}
	}

	return violations
}

// checkFieldType validates that a string value conforms to the expected type.
func checkFieldType(key, val string, ft SchemaFieldType) *SchemaViolation {
	switch ft {
	case FieldTypeInt:
		if !regexp.MustCompile(`^-?[0-9]+$`).MatchString(val) {
			return &SchemaViolation{Key: key, Message: fmt.Sprintf("value %q is not a valid integer", val)}
		}
	case FieldTypeBool:
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return &SchemaViolation{Key: key, Message: fmt.Sprintf("value %q is not a valid boolean (true/false/1/0)", val)}
		}
	case FieldTypeURL:
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			return &SchemaViolation{Key: key, Message: fmt.Sprintf("value %q is not a valid URL (must start with http:// or https://)", val)}
		}
	}
	return nil
}

// containsString reports whether slice contains target.
func containsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
