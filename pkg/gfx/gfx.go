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

type Key byte

const (
	KeyUndefined   = 0x00
	KeyMouseLeft   = 0x01
	KeyMouseRight  = 0x02
	KeyMouseMiddle = 0x03
	KeyBack        = 0x08
	KeyTab         = 0x09
	KeyClear       = 0x0C
	KeyReturn      = 0x0D
	KeyPause       = 0x13
	KeyEsc         = 0x1B
	KeyEnd         = 0x23
	KeyHome        = 0x24
	KeyLeft        = 0x25
	KeyUp          = 0x26
	KeyRight       = 0x27
	KeyDown        = 0x28
	KeyInsert      = 0x2D
	KeyDelete      = 0x2E
	KeyNumPad0     = 0x60
	KeyNumPad1     = 0x61
	KeyNumPad2     = 0x62
	KeyNumPad3     = 0x63
	KeyNumPad4     = 0x64
	KeyNumPad5     = 0x65
	KeyNumPad6     = 0x66
	KeyNumPad7     = 0x67
	KeyNumPad8     = 0x68
	KeyNumPad9     = 0x69
	KeyMultiply    = 0x6A
	KeyAdd         = 0x6B
	KeySeparator   = 0x6C
	KeySubtract    = 0x6D
	KeyDecimal     = 0x6E
	KeyDivide      = 0x6F
	KeyF1          = 0x70
	KeyF2          = 0x71
	KeyF3          = 0x72
	KeyF4          = 0x73
	KeyF5          = 0x74
	KeyF6          = 0x75
	KeyF7          = 0x76
	KeyF8          = 0x77
	KeyF9          = 0x78
	KeyF10         = 0x79
	KeyF11         = 0x7A
	KeyF12         = 0x7B
	KeyF13         = 0x7C
	KeyF14         = 0x7D
	KeyF15         = 0x7E
	KeyF16         = 0x7F
	KeyF17         = 0x80
	KeyF18         = 0x81
	KeyF19         = 0x82
	KeyF20         = 0x83
	KeyF21         = 0x84
	KeyF22         = 0x85
	KeyF23         = 0x86
	KeyF24         = 0x87
	KeySCROLL      = 0x91
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
