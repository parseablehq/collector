package store

import "time"

var poNameTime = make(map[string]time.Time)

func LastTimestamp(poName string) time.Time {
	return poNameTime[poName]
}

func SetLastTimestamp(poName string, time time.Time) {
	poNameTime[poName] = time
}
