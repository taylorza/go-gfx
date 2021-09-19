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

// Application defines the interface that must be implemented by an application using the graphics interface.
type Application interface {
	// Load called before the application loop starts. Can be used to preload assets like textures etc.
	Load()

	// Update is called for each update frame of the application. The d argument is the time delta in seconds since the last call to update.
	Update(d float64)

	// Unload is called when the application shutsdown. Can be used to cleanup/flush any open resources.
	Unload()
}

// Font represents a raster font that can be used to render text
type Font struct {
	W, H                int
	FirstChar, LastChar byte
	data                []byte
}

// Init initialized the graphics system, creates the platform specific window and related graphics devices.
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

// Run starts the application running and executes the platform specific event loop. This function blocks.
func Run(app Application) {
	app.Load()
	go run(app)
	driver.StartEventLoop()
	app.Unload()
}

// Width returns the pixel width of graphics surface. This is an unscaled value.
func Width() float64 {
	return width
}

// Height returns the pixel height of the graphics surface. This is an unscaled value.
func Height() float64 {
	return height
}

// Fps returns the number of update frames executed in the last second
func Fps() int {
	return fps
}

// Clear clears the graphics surface using the specified color
func Clear(c Color) {
	driver.Clear(c)
}

// SetPixel draws a pixel at the specified coordinates using the passed color
func SetPixel(x, y float64, c Color) {
	driver.SetPixel(int(x+0.5), int(y+0.5), c)
}

// KeyPressed returns true if the passed key is currently pressed. KeyPressed can also be used to check the state of the mouse buttons.
func KeyPressed(key Key) bool {
	return iomgr.keyPressed(key)
}

// KeyJustPressed returns true if the key was just pressed. This will not continue to return true if the key is held down.
// KeyJustPressed can be used to check the state of the mouse buttons.
// This can be used for one shot key presses, that require the key to be released and repressed for each interaction.
func KeyJustPressed(key Key) bool {
	return iomgr.keyJustPressed(key)
}

// MouseXY returns the coordinates of the mouse
func MouseXY() (float64, float64) {
	return iomgr.mouseXY()
}

// Color represents a RGB color
type Color uint32

// Rgb creates a new color using the specified RGB color components. Each component is a value between 0 and 255.
func Rgb(r, g, b int) Color {
	return Color((uint32(0xff) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b))
}

// Rgba creates a new color using the specified RGBA color components. Each component is a value between 0 and 255.
func Rgba(r, g, b, a int) Color {
	return Color((uint32(a) << 24) | (uint32(r) << 16) | (uint32(g) << 8) | uint32(b))
}

// R returns the red component of the color
func (c Color) R() int {
	return (int(c) >> 16) & 0xff
}

// G returns the green component of the color
func (c Color) G() int {
	return (int(c) >> 8) & 0xff
}

// B returns the blue component of the color
func (c Color) B() int {
	return int(c) & 0xff
}

// A returns the alpha component of the color
func (c Color) A() int {
	return (int(c) >> 24) & 0xff
}

// RandomColor returns a random color
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

// Add combines two colors and returns the resulting color
func (c Color) Add(o Color) Color {
	return Rgb(clamp(c.R()+o.R(), 0, 255),
		clamp(c.G()+o.G(), 0, 255),
		clamp(c.B()+o.B(), 0, 255))
}

// Blend blends the color with a second color
func (c Color) Blend(o Color) Color {
	a := c.A() + 1
	ia := 256 - c.A()
	r := (a*c.R() + ia*o.R()) >> 8
	g := (a*c.G() + ia*o.G()) >> 8
	b := (a*c.B() + ia*o.B()) >> 8
	na := (a*c.A() + ia*o.A()) >> 8
	return Rgba(r, g, b, na)
}

// Predefined colors
var (
	Transparent     = Rgba(0, 0, 0, 0)
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
}

func shutdown() {
	running = false
}
