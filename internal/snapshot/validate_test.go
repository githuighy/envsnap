package snapshot

import (
	"strings"
	"testing"
)

var baseValidateSnap = Snapshot{
	"APP_ENV":    "production",
	"DB_HOST":    "db.example.com",
	"SECRET_KEY": "abc123",
}

func TestValidate_RequiredKeys_AllPresent(t *testing.T) {
	err := Validate(baseValidateSnap, ValidateOptions{
		RequiredKeys: []string{"APP_ENV", "DB_HOST"},
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_RequiredKeys_Missing(t *testing.T) {
	err := Validate(baseValidateSnap, ValidateOptions{
		RequiredKeys: []string{"APP_ENV", "MISSING_KEY"},
	})
	if err == nil {
		t.Fatal("expected error for missing required key")
	}
	if !strings.Contains(err.Error(), "MISSING_KEY") {
		t.Errorf("expected error to mention MISSING_KEY, got: %v", err)
	}
}

func TestValidate_RequiredKeys_EmptyValue(t *testing.T) {
	snap := Snapshot{"APP_ENV": "   "}
	err := Validate(snap, ValidateOptions{
		RequiredKeys: []string{"APP_ENV"},
	})
	if err == nil {
		t.Fatal("expected error for empty required key value")
	}
}

func TestValidate_ForbiddenKeys_Present(t *testing.T) {
	err := Validate(baseValidateSnap, ValidateOptions{
		ForbiddenKeys: []string{"SECRET_KEY"},
	})
	if err == nil {
		t.Fatal("expected error for forbidden key")
	}
	if !strings.Contains(err.Error(), "SECRET_KEY") {
		t.Errorf("expected error to mention SECRET_KEY, got: %v", err)
	}
}

func TestValidate_ForbiddenKeys_Absent(t *testing.T) {
	err := Validate(baseValidateSnap, ValidateOptions{
		ForbiddenKeys: []string{"NOT_HERE"},
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_MaxKeyLength(t *testing.T) {
	snap := Snapshot{"SHORT": "ok", "A_VERY_LONG_KEY_NAME_INDEED": "val"}
	err := Validate(snap, ValidateOptions{MaxKeyLength: 10})
	if err == nil {
		t.Fatal("expected error for key exceeding max length")
	}
}

func TestValidate_MultipleIssues(t *testing.T) {
	snap := Snapshot{"FORBIDDEN": "x"}
	err := Validate(snap, ValidateOptions{
		RequiredKeys:  []string{"REQUIRED_KEY"},
		ForbiddenKeys: []string{"FORBIDDEN"},
	})
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d: %v", len(ve.Issues), ve.Issues)
	}
}
