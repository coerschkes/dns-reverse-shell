package utils

import "time"

func CurrentTimeAsString() string {
	return time.Now().Format("15:04:05")
}

func CurrentTimeAsLogFormat() string {
	return "[" + CurrentTimeAsString() + "]: "
}
