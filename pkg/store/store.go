package store

import (
	"time"
)

var PoNameTime = make(map[string]time.Time)

func LastTimestamp(poName string) time.Time {
	return PoNameTime[poName]
}

func SetLastTimestamp(poName string, time time.Time) {
	PoNameTime[poName] = time
}

func IsStoreEmpty(poName string) bool {
	if len(PoNameTime) == 0 {
		return true
	}
	return false
}
