package store

import "time"

var poNameTime = make(map[string]time.Time)

func GetTime(poName string) time.Time {
	return poNameTime[poName]
}

func PutPoNameTime(poName string, time time.Time) {
	poNameTime[poName] = time
}
