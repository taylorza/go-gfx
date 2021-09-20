package main

import "math"

const (
	TWO_PI = 2 * math.Pi
)

type wave struct {
	amplitude, period, phase float64
}

func newWave(amplitude, period, phase float64) *wave {
	return &wave{
		amplitude: amplitude,
		period:    period,
		phase:     phase,
	}
}

func (w *wave) evaluate(x float64) float64 {
	return math.Sin(w.phase+TWO_PI*x/w.period) * w.amplitude
}

func (w *wave) shiftPhase(delta float64) {
	w.phase += delta / w.period
}
