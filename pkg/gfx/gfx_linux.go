// +build linux

package gfx

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

func init() {
	driver = &xcbDriver{
		renderPeriod: 1.0 / 60.0,
	}
}

type xcbDriver struct {
	conn       *xgb.Conn
	screen     *xproto.ScreenInfo
	wid        xproto.Window
	gid        xproto.Gcontext
	pid        xproto.Pixmap
	w          int
	h          int
	sx         int
	sy         int
	backbuffer []byte

	renderElapsed float64
	renderPeriod  float64
}

func (e *xcbDriver) Clear(c Color) {
	e.backbuffer[0] = byte(c.B())
	e.backbuffer[1] = byte(c.G())
	e.backbuffer[2] = byte(c.R())
	e.backbuffer[3] = 255

	for i := 4; i < e.w*4*e.h; i *= 2 {
		copy(e.backbuffer[i:], e.backbuffer[:i])
	}
}

func (e *xcbDriver) SetPixel(x, y int, c Color) {
	offset := y*e.w*4 + x*4
	e.backbuffer[offset+0] = byte(c.B())
	e.backbuffer[offset+1] = byte(c.G())
	e.backbuffer[offset+2] = byte(c.R())
	e.backbuffer[offset+3] = 255
}

func (e *xcbDriver) FillRect(x, y, w, h int, c Color) {

}

func (e *xcbDriver) DrawTexture(x, y, srcX, srcY, srcW, srcH int, t *Texture) {

}

func (e *xcbDriver) HLine(x1, x2, y int, c Color) {

}

func (e *xcbDriver) VLine(x, y1, y2 int, c Color) {

}

func (e *xcbDriver) CreateWindow(x, y, w, h, xscale, yscale int) bool {
	var err error
	e.conn, err = xgb.NewConn()
	if err != nil {
		panic(err)
	}
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
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease |
				xproto.EventMaskExposure})
	xproto.MapWindow(e.conn, e.wid)

	e.w = w * xscale
	e.h = h * yscale
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
		xproto.Drawable(e.wid), uint16(e.w), uint16(e.h))

	bufSize := int(e.w) * 4 * int(e.h)
	e.backbuffer = make([]byte, bufSize, bufSize)

	return true
}

func (e *xcbDriver) SetWindowTitle(title string) {

}

func (e *xcbDriver) Update(delta float64) {

}

func (e *xcbDriver) KeyJustPressed(key Key) bool {
	return false
}

func (e *xcbDriver) KeyPressed(key Key) bool {
	return false
}

func (e *xcbDriver) MouseXY() (float64, float64) {
	return -42, -42
}

func (e *xcbDriver) Render(delta float64) {
	e.renderElapsed += delta
	if e.renderElapsed >= e.renderPeriod {
		xproto.PutImage(e.conn, xproto.ImageFormatZPixmap,
			xproto.Drawable(e.pid), e.gid,
			uint16(e.w), uint16(e.h), 0, 0,
			0,
			e.screen.RootDepth,
			e.backbuffer)

		e.renderElapsed -= e.renderPeriod

		event := xproto.ExposeEvent{
			Count:    1,
			Sequence: 0,
			Height:   uint16(e.h),
			Width:    uint16(e.w),
			Window:   e.wid,
			X:        0,
			Y:        0,
		}

		xproto.SendEvent(e.conn, false, e.wid, xproto.EventMaskExposure, string(event.Bytes()))
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
				xproto.CopyArea(e.conn, xproto.Drawable(e.pid),
					xproto.Drawable(e.wid), xproto.Gcontext(e.gid),
					0, 0, 0, 0, uint16(e.w), uint16(e.h))
			case xproto.KeyPressEvent:
				fmt.Printf("Pressed: %0.2x\n", evt.Detail)
			case xproto.KeyReleaseEvent:
				fmt.Printf("Released: %0.2x\n", evt.Detail)
			default:
				//fmt.Printf("Event: %v\n", evt)
			}
		}
		if xerr != nil {
			fmt.Printf("Error: %v\n", xerr)
		}
	}
}
