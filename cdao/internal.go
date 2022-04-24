package cdao

import (
	"time"
)

func formatSec(dur time.Duration) int64 {
	if dur <= 0 {
		return -1
	}
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}

func formatMs(dur time.Duration) int64 {
	if dur <= 0 {
		return -1
	}
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}
