package keyframe

import "time"

// seconds turns a number of seconds into a duration, for the places that mix
// the two.
func seconds(v float64) time.Duration { return time.Duration(v * float64(time.Second)) }
