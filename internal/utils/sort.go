package utils

import (
	"slices"
	"strings"
	"time"
)

type ISortable interface {
	GetName() string
	GetCreatedAt() time.Time
}

// FilterByName keeps items whose name contains search, ignoring case.
// An empty search returns items unchanged.
func FilterByName[T ISortable](items []T, search string) []T {
	if search == "" {
		return items
	}
	search = strings.ToLower(search)
	var filtered []T
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.GetName()), search) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// Sort orders items in place by field: "name" (the default), "name_desc", "created", or "created_desc".
// Names compare ignoring case. Equal names fall back to the creation time:
// oldest first for "name", newest first for "name_desc".
func Sort[T ISortable](items []T, field string) {
	byName := func(a, b T) int {
		return strings.Compare(strings.ToLower(a.GetName()), strings.ToLower(b.GetName()))
	}
	byCreated := func(a, b T) int {
		return a.GetCreatedAt().Compare(b.GetCreatedAt())
	}

	var cmp func(a, b T) int
	switch strings.ToLower(field) {
	case "name_desc":
		cmp = func(a, b T) int {
			if c := byName(b, a); c != 0 {
				return c
			}
			return byCreated(b, a)
		}
	case "created":
		cmp = byCreated
	case "created_desc":
		cmp = func(a, b T) int { return byCreated(b, a) }
	default:
		cmp = func(a, b T) int {
			if c := byName(a, b); c != 0 {
				return c
			}
			return byCreated(a, b)
		}
	}
	slices.SortStableFunc(items, cmp)
}
