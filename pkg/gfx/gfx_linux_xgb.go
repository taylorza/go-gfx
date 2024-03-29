//go:build linux
// +build linux

package gfx

import (
	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// BigRequest extentions
func putImageRequest(Format byte, Drawable xproto.Drawable, Gc xproto.Gcontext,
	Width, Height uint16,
	DstX, DstY int16,
	LeftPad byte,
	Depth byte,
	imageByteCount int) (cmd []byte, imageOffset int) {

	size := xgb.Pad(28 + xgb.Pad(imageByteCount))
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

	//	copy(buf[b:], Data)
	//	b += len(Data)

	return buf, b
}

func putImage(c *xgb.Conn, cmd []byte) xproto.PutImageCookie {

	cookie := c.NewCookie(true, false)
	c.NewRequest(cmd, cookie)
	return xproto.PutImageCookie{Cookie: cookie}
}

/*
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
*/
