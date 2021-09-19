package main

import (
	"github.com/taylorza/go-gfx/pkg/gfx"
)

type myapp struct {
	pts []pt
}

type pt struct {
	x, y float64
}

func (app *myapp) Load() {
	app.pts = make([]pt, 0)
}

func (app *myapp) Unload() {}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Black)
	if gfx.KeyJustPressed(gfx.KeyMouseLeft) {
		x, y := gfx.MouseXY()
		app.pts = append(app.pts, pt{x, y})
	}

	if len(app.pts) > 0 {
		pt := app.pts[0]
		gfx.DrawCircle(pt.x, pt.y, 3, gfx.Red)
		for i := 0; i < len(app.pts)-1; i++ {
			pt1 := app.pts[i]
			pt2 := app.pts[i+1]

			gfx.DrawLine(pt1.x, pt1.y, pt2.x, pt2.y, gfx.White)
			gfx.DrawCircle(pt1.x, pt1.y, 3, gfx.Red)
			gfx.DrawCircle(pt2.x, pt2.y, 3, gfx.Red)
		}
	}
}

func main() {
	if gfx.Init("Mouse Demo", 0, 0, 640, 480, 2, 2) {
		gfx.Run(&myapp{})
	}
}
