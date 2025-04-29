package helpers

import (
	"strconv"
	"time"
)

func ValidateTimeStamp(timestamp string, timeSkew float64) bool {
	//check if the timestamp is unix
	timestampUnix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	receivedTime := time.Unix(timestampUnix, 0)
	if receivedTime.IsZero() {
		return false
	}
	currTime := time.Now()
	if currTime.Sub(receivedTime).Seconds() > timeSkew {
		return false
	}

	return true
}
