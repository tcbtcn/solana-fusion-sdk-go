package time

import "time"

// Now returns current UNIX timestamp in seconds
func Now() int64 {
	return time.Now().Unix()
}
