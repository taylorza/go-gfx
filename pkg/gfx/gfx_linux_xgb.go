// +build linux

package gfx

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// BigRequest extentions

func putImageRequest(c *xgb.Conn, Format byte, Drawable xproto.Drawable, Gc xproto.Gcontext,
	Width, Height uint16,
	DstX, DstY int16,
	LeftPad byte,
	Depth byte,
	Data []byte) []byte {

	size := xgb.Pad(28 + xgb.Pad(len(Data)))
	buf := make([]byte, size)
	b := 0

	buf[b] = 72 // PutImage opcode
	b++

	buf[b] = Format
	b++

	xgb.Put16(buf[b:], 0)
	b += 2

	xgb.Put32(buf[b:], uint32(size/4))
	b += 4

	xgb.Put32(buf[b:], uint32(Drawable))
	b += 4

	xgb.Put32(buf[b:], uint32(Gc))
	b += 4

	xgb.Put16(buf[b:], Width)
	b += 2

	xgb.Put16(buf[b:], Height)
	b += 2

	xgb.Put16(buf[b:], uint16(DstX))
	b += 2

	xgb.Put16(buf[b:], uint16(DstY))
	b += 2

	buf[b] = LeftPad
	b++

	buf[b] = Depth
	b++

	b += 2

	copy(buf[b:], Data)
	b += len(Data)

	return buf
}

func putImage(c *xgb.Conn, Format byte, Drawable xproto.Drawable, Gc xproto.Gcontext,
	Width, Height uint16,
	DstX, DstY int16,
	LeftPad byte,
	Depth byte,
	Data []byte) xproto.PutImageCookie {

	cookie := c.NewCookie(true, false)
	c.NewRequest(putImageRequest(c, Format, Drawable, Gc, Width, Height, DstX, DstY, LeftPad, Depth, Data), cookie)
	return xproto.PutImageCookie{Cookie: cookie}
}
