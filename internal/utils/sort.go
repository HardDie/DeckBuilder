package utils

import (
	"sort"
	"strings"
	"time"
)

type ISortable interface {
	GetName() string
	GetCreatedAt() time.Time
}

func Sort[T ISortable](items *[]T, field string) {
	field = strings.ToLower(field)
	sort.SliceStable(*items, func(i, j int) bool {
		switch field {
		default: // "name"
			if (*items)[i].GetName() == (*items)[j].GetName() {
				// If the names are equal, in this case order by date of creation
				return (*items)[i].GetCreatedAt().Before((*items)[j].GetCreatedAt())
			}
			return (*items)[i].GetName() < (*items)[j].GetName()
		case "name_desc":
			if (*items)[i].GetName() == (*items)[j].GetName() {
				// If the names are the same, then in this case the reverse order by date of creation
				return (*items)[i].GetCreatedAt().After((*items)[j].GetCreatedAt())
			}
			return (*items)[i].GetName() > (*items)[j].GetName()
		case "created":
			return (*items)[i].GetCreatedAt().Before((*items)[j].GetCreatedAt())
		case "created_desc":
			return (*items)[i].GetCreatedAt().After((*items)[j].GetCreatedAt())
		}
	})
}
