//go:build linux
// +build linux

package gfx

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/bigreq"
	"github.com/jezek/xgb/xproto"
)

func init() {
	driver = &xcbDriver{
		renderPeriod: 1.0 / 60.0,
	}
}

type xcbDriver struct {
	conn         *xgb.Conn
	screen       *xproto.ScreenInfo
	wid          xproto.Window
	gid          xproto.Gcontext
	pid          xproto.Pixmap
	width        int
	height       int
	sx           int
	sy           int
	backBuffer   []byte
	backBufPtr   unsafe.Pointer
	renderBuffer []byte
	scaleBuffer  []byte
	putImageCmd  []byte

	renderElapsed float64
	renderPeriod  float64
	rendering     int32
}

func (e *xcbDriver) Init() error {
	var err error
	e.conn, err = xgb.NewConn()
	if err != nil {
		panic(err)
	}

	bigreq.Init(e.conn)
	bigreq.Enable(e.conn)

	iomgr.setKeyMapping(0x09, KeyEsc)
	iomgr.setKeyMapping(0x16, KeyBack)
	iomgr.setKeyMapping(0x17, KeyTab)
	iomgr.setKeyMapping(0x18, KeyQ)
	iomgr.setKeyMapping(0x19, KeyW)
	iomgr.setKeyMapping(0x1a, KeyE)
	iomgr.setKeyMapping(0x1b, KeyR)
	iomgr.setKeyMapping(0x1c, KeyT)
	iomgr.setKeyMapping(0x1d, KeyY)
	iomgr.setKeyMapping(0x1e, KeyU)
	iomgr.setKeyMapping(0x1f, KeyI)
	iomgr.setKeyMapping(0x20, KeyO)
	iomgr.setKeyMapping(0x21, KeyP)
	iomgr.setKeyMapping(0x24, KeyReturn)
	iomgr.setKeyMapping(0x26, KeyA)
	iomgr.setKeyMapping(0x27, KeyS)
	iomgr.setKeyMapping(0x28, KeyD)
	iomgr.setKeyMapping(0x29, KeyF)
	iomgr.setKeyMapping(0x2a, KeyG)
	iomgr.setKeyMapping(0x2b, KeyH)
	iomgr.setKeyMapping(0x2c, KeyJ)
	iomgr.setKeyMapping(0x2d, KeyK)
	iomgr.setKeyMapping(0x2e, KeyL)

	iomgr.setKeyMapping(0x34, KeyZ)
	iomgr.setKeyMapping(0x35, KeyX)
	iomgr.setKeyMapping(0x36, KeyC)
	iomgr.setKeyMapping(0x37, KeyV)
	iomgr.setKeyMapping(0x38, KeyB)
	iomgr.setKeyMapping(0x39, KeyN)
	iomgr.setKeyMapping(0x3a, KeyM)

	iomgr.setKeyMapping(0x41, KeySpace)
	iomgr.setKeyMapping(0x6e, KeyHome)
	iomgr.setKeyMapping(0x6f, KeyUp)
	iomgr.setKeyMapping(0x71, KeyLeft)
	iomgr.setKeyMapping(0x72, KeyRight)
	iomgr.setKeyMapping(0x73, KeyEnd)
	iomgr.setKeyMapping(0x74, KeyDown)

	iomgr.setKeyMapping(0x76, KeyInsert)
	iomgr.setKeyMapping(0x77, KeyDelete)

	iomgr.setKeyMapping(0x7f, KeyPause)

	iomgr.setKeyMapping(0x5a, KeyNumPad0)
	iomgr.setKeyMapping(0x57, KeyNumPad1)
	iomgr.setKeyMapping(0x58, KeyNumPad2)
	iomgr.setKeyMapping(0x59, KeyNumPad3)
	iomgr.setKeyMapping(0x53, KeyNumPad4)
	iomgr.setKeyMapping(0x54, KeyNumPad5)
	iomgr.setKeyMapping(0x55, KeyNumPad6)
	iomgr.setKeyMapping(0x4f, KeyNumPad7)
	iomgr.setKeyMapping(0x50, KeyNumPad8)
	iomgr.setKeyMapping(0x51, KeyNumPad9)

	iomgr.setKeyMapping(0x3f, KeyMultiply)
	iomgr.setKeyMapping(0x52, KeySubtract)
	iomgr.setKeyMapping(0x56, KeyAdd)
	iomgr.setKeyMapping(0x5b, KeyDecimal)
	iomgr.setKeyMapping(0x6a, KeyDivide)

	iomgr.setKeyMapping(0x43, KeyF1)
	iomgr.setKeyMapping(0x44, KeyF2)
	iomgr.setKeyMapping(0x45, KeyF3)
	iomgr.setKeyMapping(0x46, KeyF4)
	iomgr.setKeyMapping(0x47, KeyF5)
	iomgr.setKeyMapping(0x48, KeyF6)
	iomgr.setKeyMapping(0x49, KeyF7)
	iomgr.setKeyMapping(0x4a, KeyF8)
	iomgr.setKeyMapping(0x4b, KeyF9)
	iomgr.setKeyMapping(0x4c, KeyF10)
	iomgr.setKeyMapping(0x5f, KeyF11)
	iomgr.setKeyMapping(0x60, KeyF12)

	return nil
}

func (e *xcbDriver) idx(x, y int) int {
	return (y*e.width + x) * 4
}

func (e *xcbDriver) idxPtr(x, y int) unsafe.Pointer {
	return unsafe.Add(e.backBufPtr, (y*e.width*4)+(x*4))
}

func (e *xcbDriver) Clear(c Color) {
	*(*Color)(unsafe.Pointer(e.backBufPtr)) = c

	for i := 4; i < e.width*4*e.height; i *= 2 {
		copy(e.backBuffer[i:], e.backBuffer[:i])
	}
}

func (e *xcbDriver) SetPixel(x, y int, c Color) {
	if x < 0 || x > e.width || y < 0 || y > e.height {
		return
	}
	ptr := e.idxPtr(x, y)
	if c.A() != 255 {
		c = c.Blend(*(*Color)(unsafe.Pointer(ptr)))
	}
	*(*Color)(unsafe.Pointer(ptr)) = c
}

func (e *xcbDriver) FillRect(x, y, w, h int, c Color) {
	if x > e.width || y > e.height || x+w < 0 || y+h < 0 {
		return
	}

	if x < 0 {
		w += x
		x = 0
	}
	if y < 0 {
		h += y
		y = 0
	}
	if x+w > e.width {
		w -= (x + w) - e.width
	}
	if y+h > e.height {
		h -= (y + h) - e.height
	}

	ptr := e.idxPtr(x, y)
	if c.A() != 255 {
		for y1 := 0; y1 < h; y1++ {
			for x1 := 0; x1 < w; x1++ {
				p := (*Color)(unsafe.Add(ptr, x1*4))
				*p = c.Blend(*p)
			}
			ptr = unsafe.Add(ptr, e.width*4)
		}
	} else {
		for i := 0; i < w; i++ {
			*(*Color)(unsafe.Add(ptr, i*4)) = c
		}
		for i := 1; i < h; i++ {
			copy(e.backBuffer[e.idx(x, y+i):], e.backBuffer[e.idx(x, y+i-1):e.idx(x+w, y+i-1)])
		}
	}
}

func (e *xcbDriver) DrawTexture(x, y, srcX, srcY, srcW, srcH int, t *Texture) {
	if x > e.width || y > e.height || x+srcW < 0 || y+srcH < 0 || srcX >= t.W || srcY >= t.H {
		return
	}

	if x < 0 {
		srcX += -x
		srcW += x
		x = 0
	}

	if y < 0 {
		srcY += -y
		srcH += y
		y = 0
	}

	if x+srcW > e.width {
		srcW -= (x + srcW) - e.width
	}

	if y+srcH > e.height {
		srcH -= (y + srcH) - e.height
	}

	if srcX < 0 {
		srcW += srcX
		srcX = 0
	}
	if srcY < 0 {
		srcH += srcY
		srcY = 0
	}
	if srcX+srcW > t.W {
		srcW = t.W - srcX
	}
	if srcY+srcH > t.H {
		srcH = t.H - srcY
	}

	x1 := 0
	y1 := 0
	x2 := srcW
	y2 := srcH

	if x < 0 {
		x1 = -x
	}
	if y < 0 {
		y1 = -y
	}
	if x+x2 > e.width {
		x2 -= (x + x2) - e.width
	}
	if y+y2 > e.height {
		y2 -= (y + y2) - e.height
	}

	textureRowOffset := ((srcY+y1)*t.W + (srcX + x1)) * 4
	bufferRowOffset := e.idx(x, y)

	tptr := unsafe.Pointer(&t.pixels[0])
	for ty := y1; ty < y2; ty++ {
		i := textureRowOffset
		j := bufferRowOffset
		for tx := x1; tx < x2; tx++ {
			tc := Color(*(*uint32)(unsafe.Add(tptr, i)))
			pbc := (*Color)(unsafe.Add(e.backBufPtr, j))
			if tc.A() != 255 {
				tc = tc.Blend(*pbc)
			}
			*pbc = tc
			i += 4
			j += 4
		}
		textureRowOffset += t.W * 4
		bufferRowOffset += e.width * 4
	}
}

func (e *xcbDriver) HLine(x1, x2, y int, c Color) {
	if y < 0 || y >= e.height {
		return
	}
	if x1 < 0 {
		x1 = 0
	}
	if x2 >= e.width {
		x2 = e.width - 1
	}
	if x1 > x2 {
		return
	}

	ptr := e.idxPtr(x1, y)
	for i := 0; i <= x2-x1; i++ {
		if c.A() != 255 {
			c = c.Blend(*(*Color)(ptr))
		}
		*(*Color)(ptr) = c
		ptr = unsafe.Add(ptr, 4)
	}
}

func (e *xcbDriver) VLine(x, y1, y2 int, c Color) {
	if x < 0 || x >= e.width {
		return
	}
	if y1 < 0 {
		y1 = 0
	}
	if y2 >= e.height {
		y2 = e.height - 1
	}
	if y1 > y2 {
		return
	}

	ptr := e.idxPtr(x, y1)
	for i := 0; i <= y2-y1; i++ {
		if c.A() != 255 {
			c = c.Blend(*(*Color)(unsafe.Pointer(ptr)))
		}
		*(*Color)(ptr) = c
		ptr = unsafe.Add(ptr, e.width*4)
	}
}

func (e *xcbDriver) CreateWindow(x, y, w, h, xscale, yscale int) bool {
	var err error

	e.wid, err = xproto.NewWindowId(e.conn)
	if err != nil {
		panic(err)
	}
	e.screen = xproto.Setup(e.conn).DefaultScreen(e.conn)
	xproto.CreateWindow(e.conn, e.screen.RootDepth, e.wid, e.screen.Root,
		int16(x), int16(y), uint16(w*xscale), uint16(h*yscale), 0,
		xproto.WindowClassInputOutput, e.screen.RootVisual,
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{0xff000000,
			xproto.EventMaskStructureNotify |
				xproto.EventMaskButtonPress |
				xproto.EventMaskButtonRelease |
				xproto.EventMaskPointerMotion |
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease |
				xproto.EventMaskExposure |
				xproto.EventMaskStructureNotify})
	xproto.MapWindow(e.conn, e.wid)

	e.width = w
	e.height = h
	e.sx = xscale
	e.sy = yscale
	return true
}

func (e *xcbDriver) CreateDevice() bool {
	var err error
	e.gid, err = xproto.NewGcontextId(e.conn)
	if err != nil {
		panic(err)
	}
	xproto.CreateGC(e.conn, e.gid,
		xproto.Drawable(e.wid),
		xproto.GcForeground,
		[]uint32{e.screen.BlackPixel})

	e.pid, err = xproto.NewPixmapId(e.conn)
	if err != nil {
		panic(err)
	}
	xproto.CreatePixmap(e.conn, e.screen.RootDepth, e.pid,
		xproto.Drawable(e.wid), uint16(e.width*e.sx), uint16(e.height*e.sy))

	bufSize := int(e.width * 4 * e.height)
	e.backBuffer = make([]byte, bufSize)
	e.backBufPtr = unsafe.Pointer(&e.backBuffer[0])
	e.renderBuffer = make([]byte, bufSize)

	bufSize = int((e.width * e.sx * 4) * (e.height * e.sy))
	var imageOffset int
	e.putImageCmd, imageOffset = putImageRequest(
		xproto.ImageFormatZPixmap,
		xproto.Drawable(e.pid),
		e.gid,
		uint16(e.width*e.sx), uint16(e.height*e.sy), 0, 0, 0,
		e.screen.RootDepth, bufSize)
	e.scaleBuffer = e.putImageCmd[imageOffset:]

	return true
}

func (e *xcbDriver) SetWindowTitle(title string) {
	xproto.ChangeProperty(e.conn, xproto.PropModeReplace, e.wid, xproto.AtomWmName, xproto.AtomString, 8, uint32(len(title)), []byte(title))
}

func (e *xcbDriver) Update(delta float64) {

}

func (e *xcbDriver) Render(delta float64) {
	e.renderElapsed += delta
	if e.renderElapsed >= e.renderPeriod {
		if atomic.CompareAndSwapInt32(&e.rendering, 0, 1) {
			for e.renderElapsed >= e.renderPeriod {
				e.renderElapsed -= e.renderPeriod
			}
			copy(e.renderBuffer, e.backBuffer)
			event := xproto.ExposeEvent{
				Count:    1,
				Sequence: 0,
				Width:    uint16(e.width * e.sx),
				Height:   uint16(e.height * e.sy),
				Window:   e.wid,
				X:        0,
				Y:        0,
			}
			xproto.SendEvent(e.conn, false, e.wid, xproto.EventMaskExposure, string(event.Bytes()))
		}
	}
}

func (e *xcbDriver) StartEventLoop() {
	for {
		ev, xerr := e.conn.WaitForEvent()
		if ev == nil && xerr == nil {
			break
		}
		if xerr != nil {
			fmt.Printf("Error: %v\n", xerr)
		}
		if ev != nil {
			switch evt := ev.(type) {
			case xproto.NoExposureEvent:
				continue
			case xproto.ExposeEvent:
				e.scaleImage()
				xproto.CopyArea(e.conn, xproto.Drawable(e.pid),
					xproto.Drawable(e.wid), xproto.Gcontext(e.gid),
					0, 0, 0, 0, uint16(e.width*e.sx), uint16(e.height*e.sy))
				atomic.StoreInt32(&e.rendering, 0)
			case xproto.MotionNotifyEvent:
				iomgr.updateMouse(int(evt.EventX), int(evt.EventY))
			case xproto.ButtonPressEvent:
				switch evt.Detail {
				case 1:
					iomgr.setKeyPressed(KeyMouseLeft, true)
				case 2:
					iomgr.setKeyPressed(KeyMouseMiddle, true)
				case 3:
					iomgr.setKeyPressed(KeyMouseRight, true)
				}
			case xproto.ButtonReleaseEvent:
				switch evt.Detail {
				case 1:
					iomgr.setKeyPressed(KeyMouseLeft, false)
				case 2:
					iomgr.setKeyPressed(KeyMouseMiddle, false)
				case 3:
					iomgr.setKeyPressed(KeyMouseRight, false)
				}
			case xproto.KeyPressEvent:
				//fmt.Printf("Pressed: %0.2x\n", evt.Detail)
				iomgr.setMappedKeyPressed(byte(evt.Detail), true)
			case xproto.KeyReleaseEvent:
				iomgr.setMappedKeyPressed(byte(evt.Detail), false)
			default:
				fmt.Printf("Event: %v\n", evt)
			}
		}
	}
	shutdown()
}

func (e *xcbDriver) scaleImage() {
	if e.sx == 1 && e.sy == 1 {
		copy(e.scaleBuffer, e.renderBuffer)
	} else {
		dst := unsafe.Pointer(&e.scaleBuffer[0])
		src := unsafe.Pointer(&e.renderBuffer[0])
		cdst := 0
		for y := 0; y < e.height; y++ {
			for x := 0; x < e.width; x++ {
				srcPixel := *(*uint32)(unsafe.Add(src, x*4))
				for xi := 0; xi < e.sx; xi++ {
					*(*uint32)(unsafe.Add(dst, (x*e.sx+xi)*4)) = srcPixel
				}
			}
			row := cdst
			for yi := 1; yi <= e.sy; yi++ {
				cdst += e.width * e.sx * 4
				copy(e.scaleBuffer[cdst:], e.scaleBuffer[row:row+(e.width*e.sx*4)])
				row += e.width * e.sx * 4
			}
			dst = unsafe.Add(dst, e.width*e.sx*e.sy*4)
			src = unsafe.Add(src, e.width*4)
		}
	}
	putImage(e.conn, e.putImageCmd)
}
