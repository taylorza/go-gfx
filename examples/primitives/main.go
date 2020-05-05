package main

import (
	"github.com/taylorza/go-gfx/pkg/gfx"
)

type myapp struct {
	screen int
}

func (app *myapp) Load() {
}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Cyan)

	gfx.DrawString(gfx.Font6x8Ati, 109, 4, "GO GFX Primitives", gfx.Transparent, gfx.Black)
	gfx.DrawString(gfx.Font6x8Ati, 73, 16, "Press SPACE to toggle display", gfx.Transparent, gfx.Black)

	if gfx.KeyJustPressed(gfx.KeySpace) {
		app.screen++
		if app.screen > 1 {
			app.screen = 0
		}
	}
	if app.screen == 0 {
		gfx.DrawCircle(gfx.Width()/2, gfx.Height()/2, gfx.Height()/3, gfx.Red)
		gfx.DrawRect(gfx.Width()/2-50, gfx.Height()/2-50, 100, 100, gfx.Yellow)
		gfx.DrawLine(0, 0, gfx.Width(), gfx.Height(), gfx.Black)
		gfx.DrawLine(gfx.Width(), 0, 0, gfx.Height(), gfx.Blue)
	} else if app.screen == 1 {
		gfx.FillCircle(gfx.Width()/2, gfx.Height()/2, gfx.Height()/3, gfx.Red)
		gfx.FillRect(gfx.Width()/2-50, gfx.Height()/2-50, 100, 100, gfx.Yellow)
		gfx.DrawLine(0, 0, gfx.Width(), gfx.Height(), gfx.Black)
		gfx.DrawLine(gfx.Width(), 0, 0, gfx.Height(), gfx.Blue)
	}
}

func (app *myapp) Unload() {

}

func main() {
	if gfx.Init("GFX Primitives", 10, 10, 320, 240, 2, 2) {
		gfx.Run(&myapp{})
	}
}
