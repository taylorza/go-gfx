package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

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
	app.py = gfx.Height() - float64(app.paddleTexture.H) - 80
	app.ps = 500

	app.bx = (gfx.Width() - float64(app.ballTexture.W)) / 2
	app.by = (gfx.Height() - float64(app.ballTexture.H)) / 2
	app.dx = 1
	app.dy = 1
	app.bs = 300
}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Cyan)

	gfx.DrawTexture(app.bx, app.by, app.ballTexture)
	gfx.DrawTexture(app.px, app.py, app.paddleTexture)

	gfx.DrawRect(app.bx, app.by, float64(app.ballTexture.W), float64(app.ballTexture.H), gfx.Red)
	gfx.DrawRect(app.px, app.py, float64(app.paddleTexture.W), float64(app.paddleTexture.H), gfx.Red)

	app.bx += app.dx * app.bs * delta
	app.by += app.dy * app.bs * delta
	if app.bx <= 0 || app.bx >= gfx.Width()-float64(app.ballTexture.W) {
		app.bx -= app.dx * app.bs * delta
		app.dx = -app.dx
	}

	if gfx.KeyPressed(gfx.KeyLeft) && app.px > 0 {
		app.px -= app.ps * delta
	}

	if gfx.KeyPressed(gfx.KeyRight) && app.px < gfx.Width()-float64(app.paddleTexture.W) {
		app.px += app.ps * delta
	}

	if app.by <= 0 || app.by >= gfx.Height()-float64(app.ballTexture.H) {
		app.by -= app.dy * app.bs * delta
		app.dy = -app.dy
	}

	// simple collision check between the paddle and the ball
	rcBall := rc(app.bx, app.by, app.ballTexture.W, app.ballTexture.H).shrink(2)
	rcPaddle := rc(app.px, app.py, app.paddleTexture.W, app.paddleTexture.H).shrink(3)
	if app.dy > 0 && rcBall.intersect(rcPaddle) {
		app.by -= app.dy * app.bs * delta
		app.bx -= app.dx * app.bs * delta
		app.dy = -app.dy
		app.score++
	}

	gfx.DrawString(gfx.Font8x16, 8, 8, fmt.Sprintf("Score %v", app.score), gfx.Transparent, gfx.Black)
}

func (app *myapp) Unload() {

}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile = flag.String("memprofile", "", "write memory profile to file")
)

func main() {
	flag.Parse()
	go func() {
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
		time.Sleep(20 * time.Second)
	}()

	if gfx.Init("GFX Paddle Game", 10, 10, 400, 300, 2, 2) {
		gfx.Run(&myapp{})
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}

type rect struct {
	x1, y1, x2, y2 float64
}

func rc(x, y float64, w, h int) rect {
	return rect{x1: x, y1: y, x2: x + float64(w), y2: y + float64(h)}
}

func (r rect) shrink(d float64) rect {
	d /= 2
	return rect{r.x1 + d, r.y1 + d, r.x2 - d, r.y2 - d}
}

func (r rect) intersect(o rect) bool {
	return r.x1 <= o.x2 && r.x2 >= o.x1 &&
		r.y1 <= o.y2 && r.y2 >= o.y1
}
