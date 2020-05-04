package gfx

import (
	"math/rand"
	"runtime"
	"time"
)

var (
	width  float64
	height float64
	scaleX float64
	scaleY float64

	fps int
)

type Application interface {
	Load()
	Update(d float64)
	Unload()
}

type Font struct {
	W, H                int
	FirstChar, LastChar byte
	data                []byte
}

func Init(title string, x, y, w, h, xscale, yscale int) bool {
	runtime.LockOSThread()

	width = float64(w)
	height = float64(h)
	scaleX = float64(xscale)
	scaleY = float64(yscale)

	iomgr = &ioManager{}
	driver.Init()

	if !driver.CreateWindow(x, y, w, h, xscale, yscale) {
		return false
	}

	if !driver.CreateDevice() {
		return false
	}

	driver.SetWindowTitle(title)

	return true
}

func Run(app Application) {
	go run(app)
	driver.StartEventLoop()
}

func Width() float64 {
	return width
}

func Height() float64 {
	return height
}

func Fps() int {
	return fps
}

func Clear(c Color) {
	driver.Clear(c)
}

func SetPixel(x, y float64, c Color) {
	driver.SetPixel(int(x), int(y), c)
}

func KeyPressed(key Key) bool {
	return iomgr.keyPressed(key)
}

func KeyJustPressed(key Key) bool {
	return iomgr.keyJustPressed(key)
}

func MouseXY() (float64, float64) {
	return iomgr.mouseXY()
}

type Color uint32

func Rgb(r, g, b int) Color {
	return Color((r << 16) | (g << 8) | b)
}

func (c Color) A() int {
	return (int(c) >> 24) & 0xff
}

func (c Color) R() int {
	return (int(c) >> 16) & 0xff
}

func (c Color) G() int {
	return (int(c) >> 8) & 0xff
}

func (c Color) B() int {
	return int(c) & 0xff
}

func RandomColor() Color {
	return Rgb(rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (c Color) Add(o Color) Color {
	return Rgb(clamp(c.R()+o.R(), 0, 255),
		clamp(c.G()+o.G(), 0, 255),
		clamp(c.B()+o.B(), 0, 255))
}

var (
	Transparent     = Color(255 << 24)
	Black           = Rgb(0, 0, 0)
	BrightBlue      = Rgb(0, 0, 255)
	BrightGreen     = Rgb(0, 255, 0)
	BrightCyan      = Rgb(0, 255, 255)
	BrightRed       = Rgb(255, 0, 0)
	BrightMagenta   = Rgb(255, 0, 255)
	BrightYellow    = Rgb(255, 255, 0)
	White           = Rgb(255, 255, 255)
	Blue            = Rgb(0, 0, 192)
	Green           = Rgb(0, 192, 0)
	Cyan            = Rgb(0, 192, 192)
	Red             = Rgb(192, 0, 0)
	Magenta         = Rgb(192, 0, 192)
	Yellow          = Rgb(192, 192, 0)
	Grey            = Rgb(192, 192, 192)
	DarkBlue        = Rgb(0, 0, 128)
	DarkGreen       = Rgb(0, 128, 0)
	DarkCyan        = Rgb(0, 128, 128)
	DarkRed         = Rgb(128, 0, 0)
	DarkMagenta     = Rgb(128, 0, 128)
	DarkYellow      = Rgb(128, 128, 0)
	DarkGrey        = Rgb(128, 128, 128)
	VeryDarkBlue    = Rgb(0, 0, 64)
	VeryDarkGreen   = Rgb(0, 64, 0)
	VeryDarkCyan    = Rgb(0, 64, 64)
	VeryDarkRed     = Rgb(64, 0, 0)
	VeryDarkMagenta = Rgb(64, 0, 64)
	VeryDarkYellow  = Rgb(64, 64, 0)
	VeryDarkGrey    = Rgb(64, 64, 64)

	ZXBlack         = Black              // 0
	ZXBlue          = Rgb(0, 0, 205)     // 1
	ZXRed           = Rgb(205, 0, 0)     // 2
	ZXMagenta       = Rgb(205, 0, 205)   // 3
	ZXGreen         = Rgb(0, 205, 0)     // 4
	ZXCyan          = Rgb(0, 205, 205)   // 5
	ZXYellow        = Rgb(205, 205, 0)   // 6
	ZXWhite         = Rgb(205, 205, 205) // 7
	ZXBrightBlack   = Black              // 0 + Bright
	ZXBrightBlue    = BrightBlue         // 1 + Bright
	ZXBrightRed     = BrightRed          // 2 + Bright
	ZXBrightMagenta = BrightMagenta      // 3 + Bright
	ZXBrightGreen   = BrightGreen        // 4 + Bright
	ZXBrightCyan    = BrightCyan         // 5 + Bright
	ZXBrightYellow  = BrightYellow       // 6 + Bright
	ZXBrightWhite   = White              // 7 + Bright

	C64Black      = Black              // 0
	C64White      = White              // 1
	C64Red        = Rgb(146, 74, 64)   // 2
	C64Cyan       = Rgb(132, 197, 204) // 3
	C64Purple     = Rgb(147, 81, 182)  // 4
	C64Green      = Rgb(114, 177, 75)  // 5
	C64Blue       = Rgb(72, 58, 170)   // 6
	C64Yellow     = Rgb(213, 223, 124) // 7
	C64Orange     = Rgb(153, 105, 45)  // 8
	C64Brown      = Rgb(103, 82, 0)    // 9
	C64LightRed   = Rgb(193, 129, 120) // 10
	C64DarkGrey   = Rgb(96, 96, 96)    // 11
	C64Grey       = Rgb(138, 138, 138) // 12
	C64LightGreen = Rgb(179, 236, 145) // 13
	C64LightBlue  = Rgb(134, 122, 222) // 14
	C64LightGrey  = Rgb(179, 179, 179) // 15
)

func run(app Application) {
	app.Load()
	running = true
	lastFrame := time.Now()
	frameTimer := 0.0
	frameCount := 0

	for running {
		startUpdate := time.Now()
		elapsedTime := startUpdate.Sub(lastFrame).Seconds()
		lastFrame = startUpdate

		driver.Update(elapsedTime)
		iomgr.update(elapsedTime)
		app.Update(elapsedTime)
		driver.Render(elapsedTime)

		frameTimer += elapsedTime
		frameCount++
		if frameTimer >= 1 {
			fps = frameCount
			frameTimer--
			frameCount = 0
		}
	}
	app.Unload()
	done <- true
}

func shutdown() {
	running = false
	<-done
}
