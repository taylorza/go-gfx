package main

import (
	"fmt"

	"github.com/taylorza/go-gfx/pkg/gfx"
)

type myapp struct {
	paddleTexture *gfx.Texture
	ballTexture   *gfx.Texture

	px, py, ps         float64
	bx, by, dx, dy, bs float64

	score int
}

func (app *myapp) Load() {
	app.ballTexture, _ = gfx.LoadTexture("../assets/ballBlue.png")
	app.paddleTexture, _ = gfx.LoadTexture("../assets/paddleRed.png")
	app.px = (gfx.Width() - float64(app.paddleTexture.W)) / 2
	app.py = gfx.Height() - float64(app.paddleTexture.H) - 8
	app.ps = 500

	app.bx = (gfx.Width() - float64(app.ballTexture.W)) / 2
	app.by = (gfx.Height() - float64(app.ballTexture.H)) / 2
	app.dx = 1
	app.dy = 1
	app.bs = 300
}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Cyan)

	gfx.DrawTexture(app.bx, app.by, app.ballTexture, gfx.Black)
	gfx.DrawTexture(app.px, app.py, app.paddleTexture, gfx.Black)

	if gfx.KeyPressed(gfx.KeyLeft) && app.px > 0 {
		app.px -= app.ps * delta
	}

	if gfx.KeyPressed(gfx.KeyRight) && app.px < gfx.Width()-float64(app.paddleTexture.W) {
		app.px += app.ps * delta
	}

	app.bx += app.dx * app.bs * delta
	app.by += app.dy * app.bs * delta
	if app.bx <= 0 || app.bx >= gfx.Width()-float64(app.ballTexture.W) {
		app.bx -= app.dx * app.bs * delta
		app.dx = -app.dx
	}

	if app.by <= 0 || app.by >= gfx.Height()-float64(app.ballTexture.H) {
		app.by -= app.dy * app.bs * delta
		app.dy = -app.dy
	}

	// overly simple collision check between the paddle and the ball
	if app.dy > 0 &&
		app.by+float64(app.ballTexture.H) >= app.py &&
		app.bx >= app.px &&
		app.bx+float64(app.ballTexture.W) <= app.px+float64(app.paddleTexture.W) {
		app.dy = -app.dy
		app.score++
	}

	gfx.DrawString(gfx.Font8x16, 8, 8, fmt.Sprintf("Score %v", app.score), gfx.Transparent, gfx.Black)
}

func (app *myapp) Unload() {

}

func main() {
	if gfx.Init("GFX Paddle Game", 10, 10, 800, 600, 1, 1) {
		gfx.Run(&myapp{})
	}
}
