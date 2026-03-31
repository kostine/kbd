package model

import (
	"sort"
	"strings"

	"github.com/kostine/kbd/internal/bd"
)

// statusRank returns a sort rank for statuses (lower = higher in list).
var statusRank = map[string]int{
	"in_progress":       0,
	"open":              1,
	"blocked":           2,
	"hooked":            3,
	"awaiting_response": 4,
	"on-hold":           4,
	"on_hold":           4,
	"deferred":          5,
	"pinned":            6,
	"closed":            7,
}

func rankStatus(s string) int {
	if r, ok := statusRank[strings.ToLower(s)]; ok {
		return r
	}
	return 3 // unknown statuses sort in the middle
}

// SortDefault applies type-aware default sorting to an issue slice.
// The typeFilter indicates which view we're in (e.g. "epic", "task", or "").
func SortDefault(issues []bd.Issue, typeFilter string) {
	switch strings.ToLower(typeFilter) {
	case "epic":
		sortEpics(issues)
	case "task":
		sortTasks(issues)
	default:
		sortGeneric(issues)
	}
}

// sortEpics: status, then priority (highest first), then newest.
func sortEpics(issues []bd.Issue) {
	sort.SliceStable(issues, func(i, j int) bool {
		ri, rj := rankStatus(issues[i].Status), rankStatus(issues[j].Status)
		if ri != rj {
			return ri < rj
		}
		if issues[i].Priority != issues[j].Priority {
			return issues[i].Priority < issues[j].Priority
		}
		return issues[i].CreatedAt.After(issues[j].CreatedAt)
	})
}

// sortTasks: status, then priority, then parent group, then title.
func sortTasks(issues []bd.Issue) {
	sort.SliceStable(issues, func(i, j int) bool {
		ri, rj := rankStatus(issues[i].Status), rankStatus(issues[j].Status)
		if ri != rj {
			return ri < rj
		}
		if issues[i].Priority != issues[j].Priority {
			return issues[i].Priority < issues[j].Priority
		}
		pi, pj := issues[i].Parent, issues[j].Parent
		if pi != pj {
			return pi < pj
		}
		return strings.ToLower(issues[i].Title) < strings.ToLower(issues[j].Title)
	})
}

// typeRank gives a sort rank per issue type so they group together.
var typeRank = map[string]int{
	"epic":     0,
	"feature":  1,
	"bug":      2,
	"task":     3,
	"chore":    4,
	"decision": 5,
}

func rankType(t string) int {
	if r, ok := typeRank[strings.ToLower(t)]; ok {
		return r
	}
	return 6 // custom types sort at the end
}

// sortGeneric: group by type, then status, then priority, then type-specific.
func sortGeneric(issues []bd.Issue) {
	sort.SliceStable(issues, func(i, j int) bool {
		ti, tj := rankType(issues[i].Type), rankType(issues[j].Type)
		if ti != tj {
			return ti < tj
		}
		ri, rj := rankStatus(issues[i].Status), rankStatus(issues[j].Status)
		if ri != rj {
			return ri < rj
		}
		if issues[i].Priority != issues[j].Priority {
			return issues[i].Priority < issues[j].Priority
		}
		// Within same type/status/priority, apply type-specific tiebreaker
		typ := strings.ToLower(issues[i].Type)
		switch typ {
		case "task":
			pi, pj := issues[i].Parent, issues[j].Parent
			if pi != pj {
				return pi < pj
			}
			return strings.ToLower(issues[i].Title) < strings.ToLower(issues[j].Title)
		default:
			return issues[i].CreatedAt.After(issues[j].CreatedAt)
		}
	})
}
