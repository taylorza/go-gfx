package animation

import "github.com/taylorza/go-gfx/pkg/gfx"

type rect struct {
	x, y, w, h int
}

type Animation struct {
	t         *gfx.Texture
	frames    []rect
	frameTime float64
	bidi      bool
	reverse   bool
}

type Option func(a *Animation)

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

func Bidi() Option {
	return func(a *Animation) {
		a.bidi = true
	}
}

func Fps(fps int) Option {
	return func(a *Animation) {
		a.frameTime = 1 / float64(fps)
	}
}

func Reverse() Option {
	return func(a *Animation) {
		a.reverse = true
	}
}
