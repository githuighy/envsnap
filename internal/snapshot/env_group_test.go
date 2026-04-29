package snapshot

import (
	"testing"
)

func TestGroupStore_AddAndGet(t *testing.T) {
	s := NewGroupStore()
	err := s.Add("aws", []string{"AWS_", "S3_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g, ok := s.Get("aws")
	if !ok {
		t.Fatal("expected group to exist")
	}
	if g.Name != "aws" || len(g.Patterns) != 2 {
		t.Errorf("unexpected group: %+v", g)
	}
}

func TestGroupStore_Add_InvalidName(t *testing.T) {
	s := NewGroupStore()
	if err := s.Add("", []string{"FOO_"}); err == nil {
		t.Error("expected error for empty name")
	}
	if err := s.Add("bad name", []string{"FOO_"}); err == nil {
		t.Error("expected error for name with space")
	}
}

func TestGroupStore_Add_NoPatterns(t *testing.T) {
	s := NewGroupStore()
	if err := s.Add("empty", []string{}); err == nil {
		t.Error("expected error for empty patterns")
	}
}

func TestGroupStore_Delete(t *testing.T) {
	s := NewGroupStore()
	_ = s.Add("db", []string{"DB_", "DATABASE_"})
	if err := s.Delete("db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.Get("db"); ok {
		t.Error("expected group to be deleted")
	}
}

func TestGroupStore_Delete_Missing(t *testing.T) {
	s := NewGroupStore()
	if err := s.Delete("nope"); err == nil {
		t.Error("expected error for missing group")
	}
}

func TestGroupStore_List(t *testing.T) {
	s := NewGroupStore()
	_ = s.Add("zebra", []string{"Z_"})
	_ = s.Add("alpha", []string{"A_"})
	names := s.List()
	if len(names) != 2 || names[0] != "alpha" || names[1] != "zebra" {
		t.Errorf("unexpected list: %v", names)
	}
}

func TestGroupStore_ExtractGroup(t *testing.T) {
	s := NewGroupStore()
	_ = s.Add("aws", []string{"AWS_"})
	snap := Snapshot{
		"AWS_REGION":  "us-east-1",
		"AWS_KEY":     "abc",
		"DATABASE_URL": "postgres://",
	}
	result, err := s.ExtractGroup("aws", snap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
	if _, ok := result["DATABASE_URL"]; ok {
		t.Error("DATABASE_URL should not be in aws group")
	}
}

func TestGroupStore_ExtractGroup_Missing(t *testing.T) {
	s := NewGroupStore()
	_, err := s.ExtractGroup("nope", Snapshot{})
	if err == nil {
		t.Error("expected error for missing group")
	}
}
