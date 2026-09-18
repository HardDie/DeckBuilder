package utils

import "time"

// NormalizeTimestamps fills catalog createdAt/updatedAt for older on-disk JSON.
// Preview builds could omit them, store JSON null, or leave a zero time.
// New writes always store both: missing createdAt becomes now, missing
// updatedAt copies createdAt. Disk is left unchanged; this only adjusts
// values returned to services.
func NormalizeTimestamps(createdAt, updatedAt time.Time) (time.Time, time.Time) {
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	return createdAt, updatedAt
}
