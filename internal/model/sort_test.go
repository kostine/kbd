package model

import (
	"testing"
	"time"

	"github.com/kostine/kbd/internal/bd"
)

func TestSortEpics_StatusThenPriorityThenNewest(t *testing.T) {
	issues := []bd.Issue{
		{ID: "e1", Status: "closed", Priority: 1, CreatedAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		{ID: "e2", Status: "open", Priority: 2, CreatedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		{ID: "e3", Status: "open", Priority: 1, CreatedAt: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
		{ID: "e4", Status: "in_progress", Priority: 1, CreatedAt: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)},
	}
	SortDefault(issues, "epic")
	// in_progress first, then open P1 before open P2, then closed
	want := []string{"e4", "e3", "e2", "e1"}
	for i, id := range want {
		if issues[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, issues[i].ID, id)
		}
	}
}

func TestSortEpics_PriorityWithinSameStatus(t *testing.T) {
	issues := []bd.Issue{
		{ID: "e1", Status: "open", Priority: 3, CreatedAt: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
		{ID: "e2", Status: "open", Priority: 1, CreatedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		{ID: "e3", Status: "open", Priority: 0, CreatedAt: time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC)},
		{ID: "e4", Status: "open", Priority: 1, CreatedAt: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)},
	}
	SortDefault(issues, "epic")
	// P0 first, then P1 newest first, then P3
	want := []string{"e3", "e4", "e2", "e1"}
	for i, id := range want {
		if issues[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, issues[i].ID, id)
		}
	}
}

func TestSortTasks_StatusThenPriorityThenParentThenTitle(t *testing.T) {
	issues := []bd.Issue{
		{ID: "t1", Status: "open", Priority: 2, Parent: "epic-b", Title: "Zebra task"},
		{ID: "t2", Status: "open", Priority: 1, Parent: "epic-a", Title: "Beta task"},
		{ID: "t3", Status: "open", Priority: 1, Parent: "epic-a", Title: "Alpha task"},
		{ID: "t4", Status: "open", Priority: 2, Parent: "epic-b", Title: "Apple task"},
		{ID: "t5", Status: "closed", Priority: 1, Parent: "epic-a", Title: "Done task"},
	}
	SortDefault(issues, "task")
	// open first, P1 before P2, then parent, then title; closed last
	want := []string{"t3", "t2", "t4", "t1", "t5"}
	for i, id := range want {
		if issues[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, issues[i].ID, id)
		}
	}
}

func TestSortGeneric_StatusThenPriorityThenNewest(t *testing.T) {
	issues := []bd.Issue{
		{ID: "i1", Type: "bug", Status: "open", Priority: 2, CreatedAt: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
		{ID: "i2", Type: "bug", Status: "open", Priority: 1, CreatedAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)},
		{ID: "i3", Type: "bug", Status: "closed", Priority: 1, CreatedAt: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)},
	}
	SortDefault(issues, "")
	// open P1 first, then open P2, then closed
	want := []string{"i2", "i1", "i3"}
	for i, id := range want {
		if issues[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, issues[i].ID, id)
		}
	}
}

func TestSortDefault_EmptySlice(t *testing.T) {
	SortDefault(nil, "epic")
	SortDefault([]bd.Issue{}, "task")
	// no panic = pass
}

func TestSortGeneric_GroupsByType(t *testing.T) {
	issues := []bd.Issue{
		{ID: "t1", Type: "task", Status: "open", Priority: 1, Parent: "epic-b", Title: "Zebra"},
		{ID: "e1", Type: "epic", Status: "open", Priority: 1, CreatedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		{ID: "t2", Type: "task", Status: "open", Priority: 1, Parent: "epic-a", Title: "Alpha"},
		{ID: "b1", Type: "bug", Status: "open", Priority: 1, CreatedAt: time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)},
		{ID: "e2", Type: "epic", Status: "closed", Priority: 1, CreatedAt: time.Date(2026, 3, 20, 0, 0, 0, 0, time.UTC)},
	}
	SortDefault(issues, "")
	// epics first: e1 open before e2 closed
	// then bugs: b1
	// then tasks: t2 (epic-a) before t1 (epic-b)
	want := []string{"e1", "e2", "b1", "t2", "t1"}
	for i, id := range want {
		if issues[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, issues[i].ID, id)
		}
	}
}

func TestSortDefault_UnknownStatusRank(t *testing.T) {
	issues := []bd.Issue{
		{ID: "i1", Status: "custom_status", CreatedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		{ID: "i2", Status: "open", CreatedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
		{ID: "i3", Status: "closed", CreatedAt: time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)},
	}
	SortDefault(issues, "")
	// open (1) < custom (3) < closed (7)
	want := []string{"i2", "i1", "i3"}
	for i, id := range want {
		if issues[i].ID != id {
			t.Errorf("position %d: got %s, want %s", i, issues[i].ID, id)
		}
	}
}
