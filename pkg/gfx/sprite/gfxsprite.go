package sprite

import (
	"github.com/taylorza/go-gfx/pkg/gfx"
	"github.com/taylorza/go-gfx/pkg/gfx/animation"
)

// Sprite represents a sprite
type Sprite struct {
	X, Y            float64
	ox, oy          float64
	t               *gfx.Texture
	animations      map[string]*animation.Animator
	currentAnimator *animation.Animator
}

// Option is the signature of a configuration function for a Sprite
type Option func(s *Sprite)

// New creates a new sprite. If there are no animations the entire texture is used to render the sprite
func New(t *gfx.Texture, opts ...Option) *Sprite {
	s := &Sprite{
		t: t,
	}

	for _, opt := range opts {
		opt(s)
	}

	if t == nil && s.animations == nil {
		panic("sprite must have a texture or animations")
	}

	return s
}

// Animation is an Option function that adds a named animation to the sprite
func Animation(name string, a *animation.Animation) Option {
	return func(s *Sprite) {
		if s.animations == nil {
			s.animations = make(map[string]*animation.Animator)
		}
		p := animation.NewAnimator(name, a)
		s.animations[name] = p
		if s.currentAnimator == nil {
			s.currentAnimator = p
		}
	}
}

// Origin is an Option function that sets the sprites origin
func Origin(x, y int) Option {
	return func(s *Sprite) {
		s.ox = float64(x)
		s.oy = float64(y)
	}
}

// Update will updates the sprite, including the currently active animation
func (s *Sprite) Update(delta float64) {
	if len(s.animations) == 0 {
		gfx.DrawTexture(s.X-s.ox, s.Y-s.oy, s.t)
		return
	} else if s.currentAnimator != nil {
		x, y, w, h := s.currentAnimator.CurrentFrame()
		gfx.DrawTextureRect(s.X-s.ox, s.Y-s.oy, x, y, w, h, s.t)
		s.currentAnimator.Update(delta)
	}
}

// PlayAnimation selects the animation to play for the sprite
func (s *Sprite) PlayAnimation(name string, restart bool) {
	if s.animations != nil {
		if p, ok := s.animations[name]; ok {
			if p != s.currentAnimator {
				s.currentAnimator.Stop()
				s.currentAnimator = p
				if restart {
					p.Restart()
				}
			}
			p.Play()
		}
	}
}

// StopAnimation stops the sprites animation
func (s *Sprite) StopAnimation() {
	if s.currentAnimator != nil {
		s.currentAnimator.Stop()
	}
}

// ResetAnimation resets the animation to the first frame
func (s *Sprite) ResetAnimation() {
	if s.currentAnimator != nil {
		s.currentAnimator.Restart()
	}
}
