package snapshot

import (
	"fmt"
	"strings"
	"time"
)

// Annotation holds metadata attached to a snapshot.
type Annotation struct {
	Key   string
	Value string
}

// annotationPrefix is the reserved key prefix used to store annotations.
const annotationPrefix = "__envsnap_annotation_"

// Annotate attaches key-value metadata to a snapshot by embedding annotations
// as specially prefixed keys. Keys must be non-empty and must not contain '='.
func Annotate(snap Snapshot, annotations []Annotation) (Snapshot, error) {
	out := make(Snapshot, len(snap))
	for k, v := range snap {
		out[k] = v
	}

	for _, a := range annotations {
		if a.Key == "" {
			return nil, fmt.Errorf("annotation key must not be empty")
		}
		if strings.ContainsAny(a.Key, "= \t\n") {
			return nil, fmt.Errorf("annotation key %q contains invalid characters", a.Key)
		}
		out[annotationPrefix+a.Key] = a.Value
	}

	return out, nil
}

// GetAnnotations extracts all annotations embedded in a snapshot.
func GetAnnotations(snap Snapshot) []Annotation {
	var result []Annotation
	for k, v := range snap {
		if strings.HasPrefix(k, annotationPrefix) {
			result = append(result, Annotation{
				Key:   strings.TrimPrefix(k, annotationPrefix),
				Value: v,
			})
		}
	}
	return result
}

// StripAnnotations returns a copy of the snapshot with all annotation keys removed.
func StripAnnotations(snap Snapshot) Snapshot {
	out := make(Snapshot)
	for k, v := range snap {
		if !strings.HasPrefix(k, annotationPrefix) {
			out[k] = v
		}
	}
	return out
}

// TimestampAnnotation returns an Annotation recording the current UTC time.
func TimestampAnnotation() Annotation {
	return Annotation{
		Key:   "captured_at",
		Value: time.Now().UTC().Format(time.RFC3339),
	}
}
