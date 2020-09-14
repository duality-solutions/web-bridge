package util

import (
	"crypto/sha256"
	"encoding/binary"
	"time"
)

// UniqueId uses the input and the current nano epoch time to create a unique string id
func UniqueId(data []byte) string {
	timeByte := make([]byte, 8)
	nanoSec := time.Now().UnixNano()
	binary.LittleEndian.PutUint64(timeByte, uint64(nanoSec))
	data = append(data, timeByte...)
	sha256 := sha256.Sum256([]byte(data))
	return EncodeBase58(sha256[:])
}
