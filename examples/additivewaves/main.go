// Additive Waves - go-gfx implementation of demonstration presented on The Coding Train
// youtube channel
// 3.7: Additive Waves - The Nature of Code - https://www.youtube.com/watch?v=okfZRl4Xw-c
package main

import (
	"math/rand"

	"github.com/taylorza/go-gfx/pkg/gfx"
)

type myapp struct {
	waves []*wave
}

func rnd(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (app *myapp) Load() {
	// Switch to a fixed frame rate, default is 60 frames per second
	gfx.EnableFixedFrameRate(true)

	app.waves = make([]*wave, 5)
	for i := 0; i < 5; i++ {
		app.waves[i] = newWave(rnd(20, 80), rnd(100, gfx.Width()), rnd(0, TWO_PI))
	}
}

func (app *myapp) Unload() {}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Black)
	for x := 0.0; x < gfx.Width(); x += 10 {
		y := gfx.Height() / 2
		for _, w := range app.waves {
			y += w.evaluate(x)
		}
		gfx.FillCircle(x, y, 5, gfx.White)
	}

	for _, w := range app.waves {
		w.shiftPhase(10)
	}
}

func main() {
	if gfx.Init("Additive Waves", 0, 0, 600, 400, 1, 1) {
		gfx.Run(&myapp{})
	}
}
