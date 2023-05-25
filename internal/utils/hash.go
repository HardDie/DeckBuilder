package utils

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"time"
)

func HashForTime(t *time.Time) string {
	if t == nil {
		t = &time.Time{}
	}
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(t.Nanosecond()))
	hashByte := md5.Sum(buf)
	return hex.EncodeToString(hashByte[:])
}
