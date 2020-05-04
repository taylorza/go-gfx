package gfx

var (
	driver  platformDriver
	running bool
	done    = make(chan bool)
)

type platformDriver interface {
	Init() error
	CreateWindow(x, y, w, h, xscale, yscale int) bool
	CreateDevice() bool
	StartEventLoop()
	Render(delta float64)
	SetWindowTitle(title string)

	Update(delta float64)

	Clear(c Color)
	SetPixel(x, y int, c Color)
	FillRect(x, y, w, h int, c Color)
	DrawTexture(x, y int, srcX, srcY, srcW, srcH int, t *Texture)
	HLine(x1, x2, y int, c Color)
	VLine(x, y1, y2 int, c Color)
}

type keyState struct {
	justPressed bool
	pressed     bool
}

type ioManager struct {
	keymap [255]Key

	keysPhysical [255]bool
	keysLogical  [255]keyState
	mouseX       float64
	mouseY       float64
}

func (io *ioManager) update(delta float64) {
	for i, pressed := range iomgr.keysPhysical {
		if pressed {
			if !io.keysLogical[i].pressed {
				io.keysLogical[i].justPressed = true
			}
			io.keysLogical[i].pressed = true
		} else {
			io.keysLogical[i].pressed = false
			io.keysLogical[i].justPressed = false
		}
	}
}

func (io *ioManager) KeyJustPressed(key Key) bool {
	if io.keysLogical[key].justPressed {
		io.keysLogical[key].justPressed = false
		return true
	}
	return false
}

func (io *ioManager) KeyPressed(key Key) bool {
	return io.keysLogical[key].pressed
}

func (io *ioManager) MouseXY() (float64, float64) {
	return io.mouseX / scaleX, io.mouseY / scaleY
}
