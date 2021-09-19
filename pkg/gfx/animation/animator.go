package animation

// Animator drives an animation
type Animator struct {
	name    string
	a       *Animation
	frame   int
	elapsed float64
	dir     int
	playing bool
}

// NewAnimator creates a named animator for an animation
func NewAnimator(name string, a *Animation) *Animator {
	animator := &Animator{
		name: name,
		a:    a,
		dir:  1,
	}
	if a.reverse {
		animator.frame = len(a.frames) - 1
		animator.dir = -1
	}

	return animator
}

// Update must be called every render frame to update the animation state based on the setup of the attached animation
func (p *Animator) Update(delta float64) {
	if !p.playing || len(p.a.frames) < 2 {
		return
	}

	p.elapsed += delta
	if p.elapsed >= p.a.frameTime {
		p.elapsed -= p.a.frameTime
		p.frame += p.dir
		if (p.dir == 1 && p.frame == len(p.a.frames)) || (p.dir == -1 && p.frame == -1) {
			if p.a.bidi {
				p.dir = -p.dir
				p.frame += 2 * p.dir
			} else {
				if p.a.reverse {
					p.frame = len(p.a.frames) - 1
				} else {
					p.frame = 0
				}
			}
		}
	}
}

// Play starts the animation playing
func (p *Animator) Play() {
	p.playing = true
}

// Stop the animation
func (p *Animator) Stop() {
	p.playing = false
}

// Restart the animation from the start frame and direction
func (p *Animator) Restart() {
	p.elapsed = 0
	if p.a.reverse {
		p.frame = len(p.a.frames) - 1
		p.dir = -1
	} else {
		p.frame = 0
		p.dir = 1
	}
}

// CurrentFrame returns the coordinates of the image to show for the current animation state
func (p *Animator) CurrentFrame() (x, y, w, h int) {
	return p.a.frames[p.frame].x, p.a.frames[p.frame].y, p.a.frames[p.frame].w, p.a.frames[p.frame].h
}
