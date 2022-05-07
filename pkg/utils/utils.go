package utils

import (
	"os"
)

func GetParseableStreamURL(streamName string) string {
	return os.Getenv("PARSEABLE_URL") + "/api/v1/stream/" + streamName
}

func GetParseableQueryURL() string {
	return os.Getenv("PARSEABLE_URL") + "/api/v1/query"
}
