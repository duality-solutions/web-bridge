package util

import "time"

func GetCurrentEpoch() int64 {
	return GetCurrentEpochSeconds()
}

func GetCurrentEpochSeconds() int64 {
	return time.Now().Unix()
}
