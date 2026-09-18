package utils

import (
	"testing"
	"time"
)

func TestNormalizeTimestamps(t *testing.T) {
	t.Run("both_zero", func(t *testing.T) {
		created, updated := NormalizeTimestamps(time.Time{}, time.Time{})
		if created.IsZero() || updated.IsZero() {
			t.Fatal("expected filled timestamps")
		}
		if !created.Equal(updated) {
			t.Fatal("updatedAt should copy createdAt")
		}
	})

	t.Run("created_only", func(t *testing.T) {
		want := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
		created, updated := NormalizeTimestamps(want, time.Time{})
		if !created.Equal(want) || !updated.Equal(want) {
			t.Fatalf("got %v %v", created, updated)
		}
	})

	t.Run("already_set", func(t *testing.T) {
		createdAt := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)
		updatedAt := time.Date(2021, 1, 2, 3, 4, 5, 0, time.UTC)
		created, updated := NormalizeTimestamps(createdAt, updatedAt)
		if !created.Equal(createdAt) || !updated.Equal(updatedAt) {
			t.Fatalf("got %v %v", created, updated)
		}
	})
}
