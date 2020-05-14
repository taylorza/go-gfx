package gfx

var (
	iomgr   *ioManager
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
