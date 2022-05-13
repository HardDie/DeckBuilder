package utils

import (
	"sort"
	"strings"
	"time"
)

type Sortable interface {
	GetName() string
	GetCreatedAt() time.Time
}

func Sort[T Sortable](items *[]T, field string) {
	field = strings.ToLower(field)
	sort.SliceStable(*items, func(i, j int) bool {
		switch field {
		default: // "name"
			return (*items)[i].GetName() < (*items)[j].GetName()
		case "name_desc":
			return (*items)[i].GetName() > (*items)[j].GetName()
		case "created":
			return (*items)[i].GetCreatedAt().Before((*items)[j].GetCreatedAt())
		case "created_desc":
			return (*items)[i].GetCreatedAt().After((*items)[j].GetCreatedAt())
		}
	})
}
