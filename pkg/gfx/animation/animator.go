package animation

type Animator struct {
	name    string
	a       *Animation
	frame   int
	elapsed float64
	dir     int
	playing bool
}

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

func (p *Animator) Play() {
	p.playing = true
}

func (p *Animator) Stop() {
	p.playing = false
}

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

func (p *Animator) CurrentFrame() (x, y, w, h int) {
	return p.a.frames[p.frame].x, p.a.frames[p.frame].y, p.a.frames[p.frame].w, p.a.frames[p.frame].h
}
