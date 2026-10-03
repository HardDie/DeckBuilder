package utils

import (
	"reflect"
	"testing"
	"time"

	entitiesGame "github.com/HardDie/DeckBuilder/internal/entities/game"
)

var t0 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

func game(name string, minute int) *entitiesGame.Game {
	return &entitiesGame.Game{ID: name, Name: name, CreatedAt: t0.Add(time.Duration(minute) * time.Minute)}
}

// sorted runs Sort and returns the names with their minute, e.g. "Apple@2".
func sorted(items []*entitiesGame.Game, field string) []string {
	Sort(items, field)
	out := make([]string, len(items))
	for i, item := range items {
		out[i] = item.Name + "@" + item.CreatedAt.Format("4")
	}
	return out
}

func TestSort(t *testing.T) {
	items := func() []*entitiesGame.Game {
		return []*entitiesGame.Game{game("cherry", 1), game("Banana", 3), game("apple", 5), game("Apple", 2)}
	}
	tests := []struct {
		field string
		want  []string
	}{
		{"name", []string{"Apple@2", "apple@5", "Banana@3", "cherry@1"}},
		{"NAME", []string{"Apple@2", "apple@5", "Banana@3", "cherry@1"}},
		{"name_desc", []string{"cherry@1", "Banana@3", "apple@5", "Apple@2"}},
		{"created", []string{"cherry@1", "Apple@2", "Banana@3", "apple@5"}},
		{"created_desc", []string{"apple@5", "Banana@3", "Apple@2", "cherry@1"}},
		{"unknown", []string{"Apple@2", "apple@5", "Banana@3", "cherry@1"}},
		{"", []string{"Apple@2", "apple@5", "Banana@3", "cherry@1"}},
	}
	for _, tt := range tests {
		t.Run("field_"+tt.field, func(t *testing.T) {
			if got := sorted(items(), tt.field); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterByName(t *testing.T) {
	items := []*entitiesGame.Game{game("Monster", 1), game("Treasure", 2), game("Mon Ami", 3)}
	tests := []struct {
		search string
		want   []string
	}{
		{"", []string{"Monster", "Treasure", "Mon Ami"}},
		{"mon", []string{"Monster", "Mon Ami"}},
		{"MON", []string{"Monster", "Mon Ami"}},
		{"sure", []string{"Treasure"}},
		{"dragon", nil},
	}
	for _, tt := range tests {
		t.Run("search_"+tt.search, func(t *testing.T) {
			var got []string
			for _, item := range FilterByName(items, tt.search) {
				got = append(got, item.Name)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

// GetName is the real name; the helpers above do their own case folding.
func TestGetNameKeepsCase(t *testing.T) {
	if got := game("Monster", 1).GetName(); got != "Monster" {
		t.Fatalf("GetName %q, want Monster", got)
	}
}
