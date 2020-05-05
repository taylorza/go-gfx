package main

import (
	"github.com/taylorza/go-gfx/pkg/gfx"
	"github.com/taylorza/go-gfx/pkg/gfx/animation"
	"github.com/taylorza/go-gfx/pkg/gfx/sprite"
)

type myapp struct {
	player *sprite.Sprite
}

func (app *myapp) Load() {
	t, _ := gfx.LoadTexture("../assets/character_zombie_sheet.png")

	// Create aninations using the FrameSlicer to pull images from a spritesheet

	// idle anumation is 2 frames, starting at 288, 128 each frame is
	// 96x128 pixels in size and will be animated at 2 frames per second
	idleAnimation := animation.New(t, animation.FrameSlicer(288, 128, 96, 128, 2, 1), animation.Fps(2))

	// walk animation is 8 frames, starting at 0, 512 each frame is
	// 96x128 pixels in size and will be animated at 10 frames per second
	walkAnimation := animation.New(t, animation.FrameSlicer(0, 512, 96, 128, 8, 1), animation.Fps(10))

	// Define a player sprite with animations attached to the sprite
	// Multiple sprites can reuse the same animations independently
	app.player = sprite.New(t,
		sprite.Origin(48, 64),
		sprite.Animation("idle", idleAnimation),
		sprite.Animation("walk", walkAnimation))

	app.player.X = gfx.Width() / 2
	app.player.Y = gfx.Height() / 2

	// Start the player idle animation immediately
	app.player.PlayAnimation("idle", true)
}

func (app *myapp) Update(delta float64) {
	gfx.Clear(gfx.Black)

	if gfx.KeyPressed(gfx.KeySpace) {
		app.player.PlayAnimation("walk", false)
	} else {
		app.player.PlayAnimation("idle", false)
	}

	app.player.Update(delta)

	gfx.DrawString(gfx.Font8x8Bold, 4, 8, "Press SPACE to start the walk animation", gfx.Transparent, gfx.Yellow)
}

func (app *myapp) Unload() {

}

func main() {
	if gfx.Init("GFX Animation", 10, 10, 320, 240, 2, 2) {
		gfx.Run(&myapp{})
	}
}
