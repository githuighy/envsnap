package snapshot

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Tag represents a named label attached to a snapshot.
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

var validTagName = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

// TagSnapshot attaches key=value tags to a snapshot's metadata.
// Tags are stored under the "__tag.*" namespace in the snapshot.
func TagSnapshot(snap Snapshot, tags []Tag) (Snapshot, error) {
	result := make(Snapshot, len(snap))
	for k, v := range snap {
		result[k] = v
	}

	for _, tag := range tags {
		if !validTagName.MatchString(tag.Name) {
			return nil, fmt.Errorf("invalid tag name %q: must match [a-zA-Z0-9_\\-.]+", tag.Name)
		}
		key := "__tag." + tag.Name
		result[key] = tag.Value
	}

	return result, nil
}

// ListTags returns all tags embedded in a snapshot.
func ListTags(snap Snapshot) []Tag {
	var tags []Tag
	for k, v := range snap {
		if strings.HasPrefix(k, "__tag.") {
			name := strings.TrimPrefix(k, "__tag.")
			tags = append(tags, Tag{Name: name, Value: v})
		}
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})
	return tags
}

// StripTags returns a copy of the snapshot with all tag metadata removed.
func StripTags(snap Snapshot) Snapshot {
	result := make(Snapshot)
	for k, v := range snap {
		if !strings.HasPrefix(k, "__tag.") {
			result[k] = v
		}
	}
	return result
}
