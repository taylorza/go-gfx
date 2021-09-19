package animation

import "github.com/taylorza/go-gfx/pkg/gfx"

type rect struct {
	x, y, w, h int
}

// Animation represents an animation sequence made up of 1 or more frames
type Animation struct {
	t         *gfx.Texture
	frames    []rect
	frameTime float64
	bidi      bool
	reverse   bool
}

// Option is the signature of a configuration function for an animation
type Option func(a *Animation)

// New creates an instance of an animation with frames extracted from a texture
func New(t *gfx.Texture, opts ...Option) *Animation {
	if t == nil {
		panic("texture cannot be nil for animation")
	}
	a := &Animation{
		t:         t,
		frameTime: 0.1,
		bidi:      false,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

// FrameSlicer is an Option function that slices the texture into individual frames for the animation
func FrameSlicer(offsetX, offsetY, frameWidth, frameHeight, colCount, rowCount int) Option {
	return func(a *Animation) {
		for row := 0; row < rowCount; row++ {
			for column := 0; column < colCount; column++ {
				a.frames = append(a.frames, rect{
					x: offsetX + column*frameWidth,
					y: offsetY + row*frameHeight,
					w: frameWidth,
					h: frameHeight,
				})
			}
		}
	}
}

// Frame is an Option function that defines a single frame from the texture
func Frame(x, y, w, h int) Option {
	return func(a *Animation) {
		a.frames = append(a.frames, rect{
			x: x,
			y: y,
			w: w,
			h: h,
		})
	}
}

// Bidi is an Option function make the animation run in a cycle the switches direction when it reaches either end
func Bidi() Option {
	return func(a *Animation) {
		a.bidi = true
	}
}

// Fps is an Option function that sets the frame rate that the animation will run at
func Fps(fps int) Option {
	return func(a *Animation) {
		a.frameTime = 1 / float64(fps)
	}
}

// Reverse is an Option function that sets the animation to run in reverse
func Reverse() Option {
	return func(a *Animation) {
		a.reverse = true
	}
}
