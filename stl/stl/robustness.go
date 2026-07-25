package stl

import (
	"math"
	"sort"
	"time"
)

// Formula represents a temporal logic formula that can be evaluated against a trace of signals.
type Formula interface {
	Robustness(trace Trace, t time.Time) float64
}

type Op int

const (
	GT Op = iota
	LT
	GE
	LE
)

// Atom is the leaf of recursive formula evaluation
type Atom struct {
	Channel ChannelID // The channel ID of the signal to evaluate
	Op      Op        // The comparison operator to use (GT, LT, GE, LE)
	Bound   float64   //
}

func (a Atom) Robustness(tr Trace, t time.Time) float64 {
	sig, ok := tr[a.Channel] // get the signal for the channel
	if !ok {
		return math.NaN()
	}
	v, ok := sig.At(t) // get the value of the signal at time t
	if !ok {
		return math.Inf(-1) // no sample at time t, so return negative infinity
	}
	switch a.Op {
	case GT, GE: // v should be greater than or equal to the bound (so v - bound should be positive)
		return v - a.Bound
	case LT, LE: // v should be less than or equal to the bound (so bound - v should be positive)
		return a.Bound - v
	}
	return math.NaN()
}

// Not computes the negation of a formula
type Not struct {
	F Formula
}

// Robustness at time t, which represents the overall robustness of the negation.
func (n Not) Robustness(tr Trace, t time.Time) float64 {
	return -n.F.Robustness(tr, t)
}

// And computes the minimum robustness of all child formulas
type And struct {
	Children []Formula
}

// Robustness at time t, which represents the overall robustness of the conjunction.
func (a And) Robustness(tr Trace, t time.Time) float64 {
	rho := math.Inf(+1)
	for _, child := range a.Children {
		rho = math.Min(rho, child.Robustness(tr, t))
	}
	return rho
}

// Or computes the maximum robustness of all child formulas
type Or struct {
	Children []Formula
}

// Robustness at time t, which represents the overall robustness of the conjunction.
func (o Or) Robustness(tr Trace, t time.Time) float64 {
	rho := math.Inf(-1)
	for _, child := range o.Children {
		rho = math.Max(rho, child.Robustness(tr, t))
	}
	return rho
}

// Implies is sugar: φ → ψ  ≡  ¬φ ∨ ψ.
func Implies(p, q Formula) Formula {
	return Or{Children: []Formula{Not{F: p}, q}}
}

type Always struct {
	A, B time.Duration
	F    Formula
}

func (g Always) Robustness(tr Trace, t time.Time) float64 {
	win := tr.sampleTimesInWindow(t.Add(g.A), t.Add(g.B))
	rho := math.Inf(+1)
	for _, x := range win {
		rho = math.Min(rho, g.F.Robustness(tr, x))
	}
	return rho
}

type Eventually struct {
	A, B time.Duration
	F    Formula
}

func (f Eventually) Robustness(tr Trace, t time.Time) float64 {
	win := tr.sampleTimesInWindow(t.Add(f.A), t.Add(f.B))
	rho := math.Inf(-1)
	for _, x := range win {
		rho = math.Max(rho, f.F.Robustness(tr, x))
	}
	return rho
}

type Until struct {
	A, B time.Duration
	F, G Formula
}

func (u Until) Robustness(tr Trace, t time.Time) float64 {
	win := tr.sampleTimesInWindow(t.Add(u.A), t.Add(u.B))
	for _, x := range win {
		rhoG := u.G.Robustness(tr, x)
		if rhoG > 0 {
			rhoF := math.Inf(+1)
			for _, y := range win {
				if y.After(x) {
					break
				}
				rhoF = math.Min(rhoF, u.F.Robustness(tr, y))
			}
			return math.Min(rhoF, rhoG)
		}
	}
	return math.NaN()
}

// sampleTimesInWindow returns the ascending union of sample times within
// [lo, hi] across all channels in tr. For tutorial clarity we sweep every
// signal; in production you'd cache a merged timeline.
func (tr Trace) sampleTimesInWindow(lo, hi time.Time) []time.Time {
	found := map[time.Time]struct{}{}
	for _, sig := range tr {
		// O(log n) search for the first sample in sig that is not before lo, then iterate until we pass hi.
		loSig := sort.Search(len(sig), func(i int) bool { return !sig[i].T.Before(lo) })
		for i := loSig; i < len(sig) && !sig[i].T.After(hi); i++ {
			found[sig[i].T] = struct{}{}
		}
	}
	window := make([]time.Time, 0, len(found))
	for t := range found {
		window = append(window, t)
	}
	sort.Slice(window, func(i, j int) bool { return window[i].Before(window[j]) })
	return window
}
