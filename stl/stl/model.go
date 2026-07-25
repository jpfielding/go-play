package stl

import (
	"sort"
	"time"
)

type Sample struct {
	T time.Time
	V float64
}

// ChannelID is a unique identifier for a signal channel, represented as an unsigned 64-bit integer.
type ChannelID string

type Trace map[ChannelID]Signal

type Signal []Sample

// At returns the value of s at time t using zero-order hold interpolation.
// Returns ok=false when t is before the first sample in s.
func (s Signal) At(t time.Time) (float64, bool) {
	pos := sort.Search(len(s), func(i int) bool {
		return s[i].T.After(t)
	})
	if pos == 0 { // t is before the first sample in s
		return 0, false
	}
	// a sample exists at or before t, so return the value of the sample at pos-1
	return s[pos-1].V, true
}
