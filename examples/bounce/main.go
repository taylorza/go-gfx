package main

import (
	"github.com/taylorza/go-gfx/pkg/gfx"
)

type myapp struct {
	x, y, r, dx, dy float64
	speed           float64
}

func (app *myapp) Load() {
	app.x = gfx.Width() / 2
	app.y = gfx.Height() / 2
	app.r = 10
	app.dx = 1
	app.dy = 1
	app.speed = 500
}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Cyan)

	gfx.FillCircle(app.x, app.y, app.r, gfx.Blue)
	app.x += app.dx * app.speed * delta
	app.y += app.dy * app.speed * delta

	if app.x-app.r <= 0 || app.x+app.r >= gfx.Width() {
		app.x -= app.dx * app.speed * delta
		app.dx = -app.dx
	}

	if app.y-app.r <= 0 || app.y+app.r >= gfx.Height() {
		app.y -= app.dy * app.speed * delta
		app.dy = -app.dy
	}
}

func (app *myapp) Unload() {

}

func main() {
	if gfx.Init("GFX Bounce", 10, 10, 320, 240, 2, 2) {
		gfx.Run(&myapp{})
	}
}
