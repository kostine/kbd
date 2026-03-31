package model

import (
	"testing"

	"github.com/kostine/kbd/internal/bd"
)

func TestComputeDeltas_AllNew(t *testing.T) {
	newIssues := []bd.Issue{{ID: "a"}, {ID: "b"}}
	d := ComputeDeltas(nil, newIssues)
	for i, k := range d.Kinds {
		if k != RowAdded {
			t.Errorf("row %d: got %d, want RowAdded", i, k)
		}
	}
}

func TestComputeDeltas_Unchanged(t *testing.T) {
	issues := []bd.Issue{{ID: "a", Title: "Hello", Status: "open"}}
	d := ComputeDeltas(issues, issues)
	if d.Kinds[0] != RowUnchanged {
		t.Errorf("got %d, want RowUnchanged", d.Kinds[0])
	}
}

func TestComputeDeltas_Updated(t *testing.T) {
	old := []bd.Issue{{ID: "a", Status: "open"}}
	new := []bd.Issue{{ID: "a", Status: "closed"}}
	d := ComputeDeltas(old, new)
	if d.Kinds[0] != RowUpdated {
		t.Errorf("got %d, want RowUpdated", d.Kinds[0])
	}
}

func TestComputeDeltas_Mixed(t *testing.T) {
	old := []bd.Issue{
		{ID: "a", Status: "open"},
		{ID: "b", Title: "Old"},
	}
	new := []bd.Issue{
		{ID: "a", Status: "open"},  // unchanged
		{ID: "b", Title: "New"},    // updated
		{ID: "c", Status: "open"},  // added
	}
	d := ComputeDeltas(old, new)
	if d.Kinds[0] != RowUnchanged {
		t.Errorf("row 0: got %d, want RowUnchanged", d.Kinds[0])
	}
	if d.Kinds[1] != RowUpdated {
		t.Errorf("row 1: got %d, want RowUpdated", d.Kinds[1])
	}
	if d.Kinds[2] != RowAdded {
		t.Errorf("row 2: got %d, want RowAdded", d.Kinds[2])
	}
}

func TestComputeDeltas_IgnoresUpdateTimestamp(t *testing.T) {
	old := []bd.Issue{{ID: "a", Title: "Same"}}
	new := []bd.Issue{{ID: "a", Title: "Same"}} // only UpdatedAt differs
	new[0].UpdatedAt = old[0].UpdatedAt.Add(60)
	d := ComputeDeltas(old, new)
	if d.Kinds[0] != RowUnchanged {
		t.Errorf("got %d, want RowUnchanged (timestamp change should be ignored)", d.Kinds[0])
	}
}
