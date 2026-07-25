package main

import (
	"fmt"
	"stl/stl"
	"time"
)

// simple STL implementation for personal edification. This is not intended to be a production-ready
// implementation, but rather a learning exercise to understand the concepts behind STL and how to
// implement them in Go.

func main() {
	// G_[0,1s] ( speed > 60  →  F_[0,500ms] brake > 0.5 )
	// One-second moving window so the sweep can react to local behavior.
	spec := stl.Always{A: 0, B: 1 * time.Second, F: stl.Implies(
		stl.Atom{Channel: "speed", Op: stl.GT, Bound: 60},
		stl.Eventually{A: 0, B: 500 * time.Millisecond,
			F: stl.Atom{Channel: "brake", Op: stl.GT, Bound: 0.5}},
	)}

	t0 := time.Now()
	tr := demoTrace(t0)

	// Robustness signal sampled every 1s across the trace.
	times := make([]time.Time, 0, 6)
	for d := time.Duration(0); d <= 5*time.Second; d += time.Second {
		times = append(times, t0.Add(d))
	}
	for _, s := range EvaluateAlong(spec, tr, times) {
		fmt.Printf("t=%2ds  rho=%+.3f\n", s.T.Sub(t0)/time.Second, s.V)
	}
}

// EvaluateAlong returns the robustness signal of spec over the given times.
func EvaluateAlong(spec stl.Formula, tr stl.Trace, times []time.Time) stl.Signal {
	out := make(stl.Signal, len(times))
	for i, t := range times {
		out[i] = stl.Sample{T: t, V: spec.Robustness(tr, t)} // rho is itself a signal, so we can sample it at t.
	}
	return out
}

// demoTrace builds a 2-channel trace with three phases:
//
//	[0,1s)   speeding, brake never crosses 0.5  → violation
//	[2s,4s)  speeding, brake responds <500ms     → satisfied
//	[5s,6s]  speeding, brake stays under 0.5     → violation
//
// Sweeping the Always_[0,1s] spec across t=0..5s produces a
// rho signal that starts negative, recovers, then dips negative again.
func demoTrace(t0 time.Time) stl.Trace {
	ms := func(d int) time.Time {
		return t0.Add(time.Duration(d) * time.Millisecond)
	}
	return stl.Trace{
		"speed": stl.Signal{
			// phase 1 — speeding, no brake
			{T: ms(0), V: 70},
			{T: ms(300), V: 70},
			{T: ms(600), V: 70},
			// phase 2 — compliant
			{T: ms(1000), V: 70}, // late brake response arrives
			{T: ms(2000), V: 50}, // back under the limit
			{T: ms(3000), V: 60}, // speeding up
			{T: ms(3200), V: 70},
			{T: ms(4000), V: 80}, // speeding again
			// phase 3 — speeding, brake stays under
			{T: ms(5000), V: 70},
			{T: ms(5300), V: 70},
			{T: ms(5500), V: 70},
			{T: ms(6000), V: 40},
		},
		"brake": stl.Signal{
			// phase 1
			{T: ms(0), V: 0.0},
			{T: ms(300), V: 0.0},
			{T: ms(600), V: 0.0},
			// phase 2
			{T: ms(1000), V: 0.9}, // arrives just outside [0,500ms]
			{T: ms(2000), V: 0.0},
			{T: ms(3000), V: 0.0},
			{T: ms(3200), V: 0.9}, // responds within 500ms of t=3s
			{T: ms(4000), V: 0.9},
			// phase 3 — brake reaches for 0.5 but never crosses
			{T: ms(5000), V: 0.30},
			{T: ms(5300), V: 0.40},
			{T: ms(5500), V: 0.45},
			{T: ms(6000), V: 0.0},
		},
	}
}
