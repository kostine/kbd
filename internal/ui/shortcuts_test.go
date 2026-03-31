package ui

import "testing"

func TestBuildShortcuts(t *testing.T) {
	tests := []struct {
		name       string
		assignees  []ShortcutSource
		types      []ShortcutSource
		statuses   []ShortcutSource
		wantLen    int
		wantChecks func(t *testing.T, got []Shortcut)
	}{
		{
			name:      "empty inputs returns empty slice",
			assignees: nil,
			types:     nil,
			statuses:  nil,
			wantLen:   0,
		},
		{
			name: "groups ordered assignees then types then statuses",
			assignees: []ShortcutSource{
				{Category: "assignee", Value: "alice", Label: "alice", Count: 5},
			},
			types: []ShortcutSource{
				{Category: "type", Value: "bug", Label: "bug", Count: 3},
			},
			statuses: []ShortcutSource{
				{Category: "status", Value: "open", Label: "open", Count: 2},
			},
			wantLen: 3,
			wantChecks: func(t *testing.T, got []Shortcut) {
				t.Helper()
				wantOrder := []struct {
					category string
					value    string
				}{
					{"assignee", "alice"},
					{"type", "bug"},
					{"status", "open"},
				}
				for i, w := range wantOrder {
					if got[i].Category != w.category {
						t.Errorf("slot %d: category = %q, want %q", i, got[i].Category, w.category)
					}
					if got[i].Value != w.value {
						t.Errorf("slot %d: value = %q, want %q", i, got[i].Value, w.value)
					}
				}
			},
		},
		{
			name: "keys are numbered 0 through N-1 sequentially",
			assignees: []ShortcutSource{
				{Category: "assignee", Value: "a", Label: "a", Count: 10},
				{Category: "assignee", Value: "b", Label: "b", Count: 8},
			},
			types: []ShortcutSource{
				{Category: "type", Value: "epic", Label: "epic", Count: 5},
			},
			statuses: nil,
			wantLen:  3,
			wantChecks: func(t *testing.T, got []Shortcut) {
				t.Helper()
				for i, s := range got {
					if s.Key != i {
						t.Errorf("slot %d: key = %d, want %d", i, s.Key, i)
					}
				}
			},
		},
		{
			name: "caps at 10 total shortcuts",
			assignees: []ShortcutSource{
				{Category: "assignee", Value: "a1", Label: "a1", Count: 50},
				{Category: "assignee", Value: "a2", Label: "a2", Count: 40},
				{Category: "assignee", Value: "a3", Label: "a3", Count: 30},
				{Category: "assignee", Value: "a4", Label: "a4", Count: 20},
			},
			types: []ShortcutSource{
				{Category: "type", Value: "t1", Label: "t1", Count: 15},
				{Category: "type", Value: "t2", Label: "t2", Count: 10},
				{Category: "type", Value: "t3", Label: "t3", Count: 5},
			},
			statuses: []ShortcutSource{
				{Category: "status", Value: "s1", Label: "s1", Count: 12},
				{Category: "status", Value: "s2", Label: "s2", Count: 8},
				{Category: "status", Value: "s3", Label: "s3", Count: 4},
				{Category: "status", Value: "s4", Label: "s4", Count: 2},
			},
			wantLen: 10,
			wantChecks: func(t *testing.T, got []Shortcut) {
				t.Helper()
				if got[9].Key != 9 {
					t.Errorf("last slot key = %d, want 9", got[9].Key)
				}
				// s4 (index 10 in the combined list) should be excluded
				for _, s := range got {
					if s.Value == "s4" {
						t.Error("s4 should have been truncated but was present")
					}
				}
			},
		},
		{
			name: "preserves sort order within each group",
			assignees: []ShortcutSource{
				{Category: "assignee", Value: "top", Label: "top", Count: 100},
				{Category: "assignee", Value: "mid", Label: "mid", Count: 50},
				{Category: "assignee", Value: "low", Label: "low", Count: 10},
			},
			types:    nil,
			statuses: nil,
			wantLen:  3,
			wantChecks: func(t *testing.T, got []Shortcut) {
				t.Helper()
				wantValues := []string{"top", "mid", "low"}
				for i, w := range wantValues {
					if got[i].Value != w {
						t.Errorf("slot %d: value = %q, want %q", i, got[i].Value, w)
					}
				}
			},
		},
		{
			name:      "single entry works",
			assignees: nil,
			types: []ShortcutSource{
				{Category: "type", Value: "bug", Label: "bug", Count: 1},
			},
			statuses: nil,
			wantLen:  1,
			wantChecks: func(t *testing.T, got []Shortcut) {
				t.Helper()
				if got[0].Key != 0 {
					t.Errorf("key = %d, want 0", got[0].Key)
				}
				if got[0].Label != "bug" {
					t.Errorf("label = %q, want %q", got[0].Label, "bug")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildShortcuts(tc.assignees, tc.types, tc.statuses, nil)
			if len(got) != tc.wantLen {
				t.Fatalf("len = %d, want %d", len(got), tc.wantLen)
			}
			if tc.wantChecks != nil {
				tc.wantChecks(t, got)
			}
		})
	}
}

func TestBuildShortcuts_WithCommands(t *testing.T) {
	assignees := []ShortcutSource{
		{Category: "assignee", Value: "alice", Label: "alice", Count: 10},
	}
	types := []ShortcutSource{
		{Category: "type", Value: "task", Label: "task", Count: 5},
	}
	commands := []ShortcutSource{
		{Category: "command", Value: "context", Label: "context"},
	}
	got := BuildShortcuts(assignees, types, nil, commands)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	// Commands should be last
	if got[2].Category != "command" {
		t.Errorf("slot 2 category = %q, want 'command'", got[2].Category)
	}
	if got[2].Value != "context" {
		t.Errorf("slot 2 value = %q, want 'context'", got[2].Value)
	}
}

func TestStatusColor(t *testing.T) {
	theme := DefaultTheme()
	tests := []struct {
		status string
		want   string // color name for comparison
	}{
		{"open", "green"},
		{"in_progress", "yellow"},
		{"blocked", "red"},
		{"closed", "gray"},
		{"unknown", "default"},
	}
	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			got := theme.StatusColor(tt.status)
			if got == 0 && tt.want != "default" {
				t.Errorf("StatusColor(%q) returned zero", tt.status)
			}
		})
	}
}

func TestKeyBindings_HandleRune(t *testing.T) {
	kb := NewKeyBindings()
	called := false
	kb.BindRune('x', "test", func() { called = true })

	hints := kb.Hints()
	if len(hints) != 1 {
		t.Fatalf("expected 1 hint, got %d", len(hints))
	}
	if hints[0].Key != "x" {
		t.Errorf("hint key = %q, want 'x'", hints[0].Key)
	}
	if hints[0].Description != "test" {
		t.Errorf("hint desc = %q, want 'test'", hints[0].Description)
	}
	// Note: can't call Handle without a real tcell.EventKey, but we tested Hints
	_ = called
}

func TestLineCount(t *testing.T) {
	tests := []struct {
		name   string
		counts []TypeStatusCounts
		want   int
	}{
		{
			name:   "empty counts returns 1",
			counts: nil,
			want:   1,
		},
		{
			name: "single type returns 1",
			counts: []TypeStatusCounts{
				{Type: "bug", Statuses: []StatusCount{{Status: "open", Count: 3}}},
			},
			want: 1,
		},
		{
			name: "two types returns 2",
			counts: []TypeStatusCounts{
				{Type: "bug", Statuses: []StatusCount{{Status: "open", Count: 3}}},
				{Type: "epic", Statuses: []StatusCount{{Status: "closed", Count: 1}}},
			},
			want: 2,
		},
		{
			name: "five types returns 5",
			counts: []TypeStatusCounts{
				{Type: "bug"},
				{Type: "epic"},
				{Type: "task"},
				{Type: "story"},
				{Type: "spike"},
			},
			want: 5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := LineCount(tc.counts)
			if got != tc.want {
				t.Errorf("LineCount() = %d, want %d", got, tc.want)
			}
		})
	}
}
