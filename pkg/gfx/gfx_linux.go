// +build linux

package gfx

import (
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/bigreq"
	"github.com/BurntSushi/xgb/xproto"
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
	backBufPtr   uintptr
	renderBuffer []byte
	scaleBuffer  []byte

	renderElapsed float64
	renderPeriod  float64
	rendering     int32
}

func (e *xcbDriver) Init() error {
	iomgr.keymap[0x6f] = KeyUp
	iomgr.keymap[0x71] = KeyLeft
	iomgr.keymap[0x72] = KeyRight
	iomgr.keymap[0x74] = KeyDown
	return nil
}

func (e *xcbDriver) Clear(c Color) {
	e.backBuffer[0] = byte(c.B())
	e.backBuffer[1] = byte(c.G())
	e.backBuffer[2] = byte(c.R())
	e.backBuffer[3] = 255

	for i := 4; i < e.width*4*e.height; i *= 2 {
		copy(e.backBuffer[i:], e.backBuffer[:i])
	}
}

func (e *xcbDriver) idx(x, y int) int {
	return (y * e.width * 4) + (x * 4)
}

func (e *xcbDriver) idxPtr(x, y int) uintptr {
	return uintptr(e.backBufPtr + uintptr((y*e.width*4)+(x*4)))
}

func (e *xcbDriver) SetPixel(x, y int, c Color) {
	if x < 0 || x > e.width || y < 0 || y > e.height {
		return
	}
	e.setPixelFast(e.idxPtr(x, y), c)
}

func (e *xcbDriver) setPixelFast(offset uintptr, c Color) {
	*(*byte)(unsafe.Pointer(offset + 0)) = byte(c.B())
	*(*byte)(unsafe.Pointer(offset + 1)) = byte(c.G())
	*(*byte)(unsafe.Pointer(offset + 2)) = byte(c.R())
	*(*byte)(unsafe.Pointer(offset + 3)) = 255
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

	offset := e.idxPtr(x, y)
	for i := 0; i < w; i++ {
		e.setPixelFast(offset+uintptr(i*4), c)
	}

	for i := 1; i < h; i++ {
		copy(e.backBuffer[e.idx(x, y+i):], e.backBuffer[e.idx(x, y+i-1):e.idx(x+w, y+i-1)])
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

	transparent := t.pixels[0]
	textureRowOffset := uintptr((srcY+y1)*t.W + (srcX + x1))
	bufferRowOffset := e.idxPtr(x, y)

	tptr := uintptr(unsafe.Pointer(&t.pixels[0]))
	for ty := y1; ty < y2; ty++ {
		i := textureRowOffset
		j := bufferRowOffset
		for tx := x1; tx < x2; tx++ {
			c := Color(*(*uint32)(unsafe.Pointer(tptr + i*4)))
			if c != transparent {
				e.setPixelFast(j, c)
			}
			i++
			j += 4
		}
		textureRowOffset += uintptr(t.W)
		bufferRowOffset += uintptr(e.width * 4)
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

	offset := e.idxPtr(x1, y)
	for i := uintptr(0); i <= uintptr(x2-x1); i++ {
		e.setPixelFast(offset, c)
		offset += 4
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

	offset := e.idxPtr(x, y1)
	for i := uintptr(0); i <= uintptr(y2-y1); i++ {
		e.setPixelFast(offset, c)
		offset += uintptr(e.width * 4)
	}
}

func (e *xcbDriver) CreateWindow(x, y, w, h, xscale, yscale int) bool {
	var err error
	e.conn, err = xgb.NewConn()
	if err != nil {
		panic(err)
	}

	bigreq.Init(e.conn)
	bigreq.Enable(e.conn)

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
				xproto.EventMaskExposure})
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
	e.backBuffer = make([]byte, bufSize, bufSize)
	e.backBufPtr = uintptr(unsafe.Pointer(&e.backBuffer[0]))
	e.renderBuffer = make([]byte, bufSize, bufSize)

	bufSize = int((e.width * e.sx * 4) * (e.height * e.sy))
	e.scaleBuffer = make([]byte, bufSize, bufSize)

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
			copy(e.renderBuffer, e.backBuffer)
			e.renderElapsed -= e.renderPeriod

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
		if ev == nil && xerr != nil {
			fmt.Println(xerr)
			return
		}
		if ev != nil {
			switch evt := ev.(type) {
			case xproto.ExposeEvent:
				e.scaleImage()
				xproto.CopyArea(e.conn, xproto.Drawable(e.pid),
					xproto.Drawable(e.wid), xproto.Gcontext(e.gid),
					0, 0, 0, 0, uint16(e.width*e.sx), uint16(e.height*e.sy))
				atomic.StoreInt32(&e.rendering, 0)
			case xproto.MotionNotifyEvent:
				iomgr.mouseX = float64(evt.EventX)
				iomgr.mouseY = float64(evt.EventY)
			case xproto.ButtonPressEvent:
				fmt.Printf("MB Pressed: %0.2x\n", evt.Detail)
				switch evt.Detail {
				case 1:
					iomgr.keysPhysical[KeyMouseLeft] = true
				case 2:
					iomgr.keysPhysical[KeyMouseMiddle] = true
				case 3:
					iomgr.keysPhysical[KeyMouseRight] = true
				}
			case xproto.ButtonReleaseEvent:
				switch evt.Detail {
				case 1:
					iomgr.keysPhysical[KeyMouseLeft] = false
				case 2:
					iomgr.keysPhysical[KeyMouseMiddle] = false
				case 3:
					iomgr.keysPhysical[KeyMouseRight] = false
				}
			case xproto.KeyPressEvent:
				fmt.Printf("Pressed: %0.2x\n", evt.Detail)
				iomgr.keysPhysical[iomgr.keymap[int(evt.Detail)]] = true
			case xproto.KeyReleaseEvent:
				iomgr.keysPhysical[iomgr.keymap[int(evt.Detail)]] = false
			default:
				//fmt.Printf("Event: %v\n", evt)
			}
		}
		if xerr != nil {
			fmt.Printf("Error: %v\n", xerr)
		}
	}
}

func (e *xcbDriver) scaleImage() {
	dst := 0
	src := 0
	for y := 0; y < e.height*e.sy; y += e.sy {
		for x := 0; x < e.width*e.sx; x++ {
			xi := x / e.sx
			e.scaleBuffer[dst+(x*4)+0] = e.renderBuffer[src+(xi*4)+0]
			e.scaleBuffer[dst+(x*4)+1] = e.renderBuffer[src+(xi*4)+1]
			e.scaleBuffer[dst+(x*4)+2] = e.renderBuffer[src+(xi*4)+2]
			e.scaleBuffer[dst+(x*4)+3] = e.renderBuffer[src+(xi*4)+3]
		}
		row := dst
		dst += e.width * e.sx * 4
		copy(e.scaleBuffer[dst:], e.scaleBuffer[row:row+(e.width*4*e.sx)*(e.sy-1)])

		dst += (e.width * e.sx * 4) * (e.sy - 1)
		src += e.width * 4
	}

	putImage(e.conn, xproto.ImageFormatZPixmap,
		xproto.Drawable(e.pid), e.gid,
		uint16(e.width*e.sx), uint16(e.height*e.sy), 0, 0,
		0,
		e.screen.RootDepth,
		e.scaleBuffer)
}
