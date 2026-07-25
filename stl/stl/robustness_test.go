package stl

import (
	"reflect"
	"testing"
	"time"
)

func TestAtomRobustness(t *testing.T) {
	tr := Trace{
		"ch1": Signal{
			{T: time.Unix(0, 0), V: 1},
			{T: time.Unix(1, 0), V: 2},
			{T: time.Unix(2, 0), V: 3},
		},
	}

	tests := []struct {
		name    string
		formula Atom
		t       time.Time
		want    float64
	}{
		{
			name:    "GT satisfied",
			formula: Atom{Channel: "ch1", Op: GT, Bound: 1.5},
			t:       time.Unix(1, 0),
			want:    0.5,
		},
		{
			name:    "GT violated",
			formula: Atom{Channel: "ch1", Op: GT, Bound: 2.5},
			t:       time.Unix(1, 0),
			want:    -0.5,
		},
		{
			name:    "LT satisfied",
			formula: Atom{Channel: "ch1", Op: LT, Bound: 2.5},
			t:       time.Unix(1, 0),
			want:    0.5,
		},
		{
			name:    "LT violated",
			formula: Atom{Channel: "ch1", Op: LT, Bound: 1.5},
			t:       time.Unix(1, 0),
			want:    -0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.formula.Robustness(tr, tt.t)
			if got != tt.want {
				t.Errorf("Robustness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotRobustness(t *testing.T) {
	tr := Trace{
		"ch1": Signal{
			{T: time.Unix(0, 0), V: 1},
			{T: time.Unix(1, 0), V: 2},
			{T: time.Unix(2, 0), V: 3},
		},
	}

	tests := []struct {
		name    string
		formula Not
		t       time.Time
		want    float64
	}{
		{
			name:    "Not GT",
			formula: Not{F: Atom{Channel: "ch1", Op: GT, Bound: 1.5}},
			t:       time.Unix(1, 0),
			want:    -0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.formula.Robustness(tr, tt.t)
			if got != tt.want {
				t.Errorf("Robustness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAndRobustness(t *testing.T) {
	tr := Trace{
		"ch1": Signal{
			{T: time.Unix(0, 0), V: 1},
			{T: time.Unix(1, 0), V: 2},
			{T: time.Unix(2, 0), V: 3},
		},
	}

	tests := []struct {
		name    string
		formula And
		t       time.Time
		want    float64
	}{
		{
			name: "And GT and LT",
			formula: And{Children: []Formula{
				Atom{Channel: "ch1", Op: GT, Bound: 1.5},
				Atom{Channel: "ch1", Op: LT, Bound: 2.5},
			}},
			t:    time.Unix(1, 0),
			want: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.formula.Robustness(tr, tt.t)
			if got != tt.want {
				t.Errorf("Robustness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrRobustness(t *testing.T) {
	tr := Trace{
		"ch1": Signal{
			{T: time.Unix(0, 0), V: 1},
			{T: time.Unix(1, 0), V: 2},
			{T: time.Unix(2, 0), V: 3},
		},
	}

	tests := []struct {
		name    string
		formula Or
		t       time.Time
		want    float64
	}{
		{
			name: "Or GT and LT",
			formula: Or{Children: []Formula{
				Atom{Channel: "ch1", Op: GT, Bound: 1.5},
				Atom{Channel: "ch1", Op: LT, Bound: 2.5},
			}},
			t:    time.Unix(1, 0),
			want: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.formula.Robustness(tr, tt.t)
			if got != tt.want {
				t.Errorf("Robustness() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamleTimesInWindow(t *testing.T) {
	tr := Trace{
		"ch1": Signal{
			{T: time.Unix(0, 0), V: 1},
			{T: time.Unix(1, 0), V: 2},
			{T: time.Unix(2, 0), V: 3},
		},
		"ch2": Signal{
			{T: time.Unix(1, 0), V: 4},
			{T: time.Unix(3, 0), V: 5},
		},
	}

	tests := []struct {
		name string
		lo   time.Time
		hi   time.Time
		want []time.Time
	}{
		{
			name: "full range",
			lo:   time.Unix(0, 0),
			hi:   time.Unix(3, 0),
			want: []time.Time{
				time.Unix(0, 0),
				time.Unix(1, 0),
				time.Unix(2, 0),
				time.Unix(3, 0),
			},
		},
		{
			name: "partial range",
			lo:   time.Unix(1, 0),
			hi:   time.Unix(2, 0),
			want: []time.Time{
				time.Unix(1, 0),
				time.Unix(2, 0),
			},
		},
		{
			name: "no samples",
			lo:   time.Unix(4, 0),
			hi:   time.Unix(5, 0),
			want: []time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tr.sampleTimesInWindow(tt.lo, tt.hi)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sampleTimesInWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
