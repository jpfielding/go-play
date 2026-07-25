package stl

import (
	"testing"
	"time"
)

func TestAt(t *testing.T) {
	sig := Signal{
		{T: time.Unix(0, 0), V: 1},
		{T: time.Unix(1, 0), V: 2},
		{T: time.Unix(2, 0), V: 3},
	}

	tests := []struct {
		name string
		t    time.Time
		want float64
		ok   bool
	}{
		{
			name: "before first sample",
			t:    time.Unix(-1, 0),
			want: 0,
			ok:   false,
		},
		{
			name: "at first sample",
			t:    time.Unix(0, 0),
			want: 1,
			ok:   true,
		},
		{
			name: "between samples",
			t:    time.Unix(1, 500000000), // 1.5 seconds
			want: 2,
			ok:   true,
		},
		{
			name: "at last sample",
			t:    time.Unix(2, 0),
			want: 3,
			ok:   true,
		},
		{
			name: "after last sample",
			t:    time.Unix(3, 0),
			want: 3,
			ok:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := sig.At(tt.t)
			if got != tt.want || ok != tt.ok {
				t.Errorf("At() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.ok)
			}
		})
	}
}
