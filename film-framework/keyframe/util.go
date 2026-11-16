package keyframe

import (
	"strconv"
	"time"
)

// seconds turns a number of seconds into a duration, for the places that mix
// the two.
func seconds(v float64) time.Duration { return time.Duration(v * float64(time.Second)) }

// secondsText writes a key's time the shortest way that keeps two nearby keys
// apart: 8, 7.98, 3.7. One decimal would print a 20 ms cut — the whole point of
// a cut — as the same time twice.
func secondsText(t float64) string {
	ms := float64(int64(t*1000 + 0.5))
	return strconv.FormatFloat(ms/1000, 'f', -1, 64)
}
