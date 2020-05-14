package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/taylorza/go-gfx/pkg/gfx"
)

const (
	newBalls = 10
)

type ball struct {
	x, y, r, dx, dy float64
	speed           float64
	color           gfx.Color
}

func newBall() *ball {
	return &ball{
		x:     gfx.Width() / 2,
		y:     gfx.Height() / 2,
		r:     10,
		speed: 300 + (rand.Float64() * 100),
		dx:    math.Sin(rand.Float64() * math.Pi * 2),
		dy:    math.Cos(rand.Float64() * math.Pi * 2),
		color: gfx.RandomColor(),
	}
}

func (b *ball) update(delta float64) {
	gfx.FillCircle(b.x, b.y, b.r, b.color)
	b.x += b.dx * b.speed * delta
	b.y += b.dy * b.speed * delta

	if b.x-b.r <= 0 || b.x+b.r >= gfx.Width() {
		b.x -= b.dx * b.speed * delta
		b.dx = -b.dx
	}

	if b.y-b.r <= 0 || b.y+b.r >= gfx.Height() {
		b.y -= b.dy * b.speed * delta
		b.dy = -b.dy
	}
}

type myapp struct {
	balls []*ball
}

func (app *myapp) Load() {
	for i := 0; i < newBalls; i++ {
		app.balls = append(app.balls, newBall())
	}
}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Cyan)

	for _, b := range app.balls {
		b.update(delta)
	}

	if gfx.KeyJustPressed(gfx.KeySpace) {
		for i := 0; i < newBalls; i++ {
			app.balls = append(app.balls, newBall())
		}
	}

	gfx.DrawString(gfx.Font8x8, 8, 8, "Press SPACE to add balls", gfx.Black, gfx.Grey)
	gfx.DrawString(gfx.Font8x8, 8, 16, fmt.Sprintf("Balls: %v", len(app.balls)), gfx.Black, gfx.Grey)
}

func (app *myapp) Unload() {

}

func main() {
	if gfx.Init("GFX Big Bounce", 10, 10, 320, 240, 2, 2) {
		gfx.Run(&myapp{})
	}
}
