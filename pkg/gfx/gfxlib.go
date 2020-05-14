package gfx

import (
	"image"
	_ "image/gif"  // imported to register the gif image decoder used to load textures from image files of this type
	_ "image/jpeg" // imported to register the jpeg image decoder used to load textures from image files of this type
	_ "image/png"  // imported to register the png image decoder used to load textures from image files of this type
	"math"
	"os"
)

// Texture in memory representation of a texture
type Texture struct {
	W, H   int
	pixels []Color
}

// LoadTexture loads an image from a file and creates a texture from it
func LoadTexture(filename string) (*Texture, error) {
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	b := m.Bounds()
	pixels := make([]Color, b.Dx()*b.Dy(), b.Dx()*b.Dy())
	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			c := Rgba(int((float64(r)/65535.)*255), int((float64(g)/65535.)*255), int((float64(b)/65535.)*255), int((float64(a)/65535.)*255))
			pixels[i] = c
			i++
		}
	}
	return &Texture{
		W:      b.Dx(),
		H:      b.Dy(),
		pixels: pixels,
	}, nil
}

// DrawTexture draws a texture to the screen at the specified location.
func DrawTexture(x, y float64, t *Texture) {
	if t == nil {
		panic("texture cannot be nil")
	}
	driver.DrawTexture(int(x), int(y), 0, 0, t.W, t.H, t)
}

// DrawTextureRect extracts a sub-rectangle from a texture and draws it to the screen.
func DrawTextureRect(x, y float64, srcX, srcY, srcW, srcH int, t *Texture) {
	if t == nil {
		panic("texture cannot be nil")
	}
	driver.DrawTexture(int(x), int(y), srcX, srcY, srcW, srcH, t)
}

func fmin4(a, b, c, d float64) float64 {
	return math.Min(math.Min(math.Min(a, b), c), d)
}

func fmax4(a, b, c, d float64) float64 {
	return math.Max(math.Max(math.Max(a, b), c), d)
}

// DrawTextureRotate draws a rotated texture to the screen
func DrawTextureRotate(x, y float64, srcX, srcY, srcW, srcH, cx, cy int, sx, sy, angle float64, t *Texture) {
	if t == nil {
		panic("texture cannot be nil")
	}

	ix := int(x)
	iy := int(y)

	if x > width || y > height || ix+srcW < 0 || iy+srcH < 0 || srcX >= t.W || srcY >= t.H {
		return
	}

	if srcX < 0 || srcX+srcW > t.W || srcY < 0 || srcY+srcH > t.H {
		return
	}

	cosine := math.Cos(angle)
	sine := math.Sin(angle)

	// Calculate the coordinates of the 4 corners of the rectangle around the sprite
	// the top, left corner is translated so that the cx and cy are at the centre of the rectangle
	w := float64(srcW) * sx
	h := float64(srcH) * sy
	rx1 := float64(-cx) * sx
	ry1 := float64(-cy) * sy
	rx2 := rx1 + w
	ry2 := ry1
	rx3 := rx1 + w
	ry3 := ry1 + h
	rx4 := rx1
	ry4 := ry1 + h

	// Rotate the rectangle
	x1 := cosine*rx1 - sine*ry1
	y1 := sine*rx1 + cosine*ry1
	x2 := cosine*rx2 - sine*ry2
	y2 := sine*rx2 + cosine*ry2
	x3 := cosine*rx3 - sine*ry3
	y3 := sine*rx3 + cosine*ry3
	x4 := cosine*rx4 - sine*ry4
	y4 := sine*rx4 + cosine*ry4

	// Find the new axis aligned rectangle that encompasses the rotated rectangle
	minx := int(math.Floor(fmin4(x1, x2, x3, x4)) * 1.1)
	miny := int(math.Floor(fmin4(y1, y2, y3, y4)) * 1.1)
	maxx := int(math.Ceil(fmax4(x1, x2, x3, x4)) * 1.1)
	maxy := int(math.Ceil(fmax4(y1, y2, y3, y4)) * 1.1)
	//DrawRect(float64(minx+ix), float64(miny+iy), float64(maxx-minx), float64(maxy-miny), Red)

	// Scan the axis aligned rectangle left to right, top to bottom and map the
	// point into the source texture and plot it if it is a valid source pixel
	transparent := t.pixels[0]
	sx = 1 / sx
	sy = 1 / sy
	for dstY := miny; dstY < maxy; dstY++ {
		for dstX := minx; dstX < maxx; dstX++ {
			tx := srcX + int(0.5+(cosine*float64(dstX)+sine*float64(dstY))*sx) + cx
			ty := srcY - int(0.5+(sine*float64(dstX)-cosine*float64(dstY))*sy) + cy

			if tx >= srcX && tx < srcX+srcW && ty >= srcY && ty < srcY+srcH {
				if t.pixels[tx+ty*t.W] != transparent {
					driver.SetPixel(ix+dstX, iy+dstY, t.pixels[tx+ty*t.W])
				}
			}
		}
	}
}

func drawHLine(x1, x2, y int, c Color) {
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	driver.HLine(x1, x2, y, c)
}

func drawVLine(x, y1, y2 int, c Color) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	driver.VLine(x, y1, y2, c)
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawLine draws a line
func DrawLine(x1, y1, x2, y2 float64, c Color) {
	ix1 := int(x1)
	iy1 := int(y1)
	ix2 := int(x2)
	iy2 := int(y2)
	dy := iabs(iy2 - iy1)
	if dy == 0 {
		drawHLine(ix1, ix2, iy1, c)
		return
	}

	dx := iabs(ix2 - ix1)
	if dx == 0 {
		drawVLine(ix1, iy1, iy2, c)
		return
	}

	sx := 1
	if x2 < x1 {
		sx = -1
	}

	sy := 1
	if y2 < y1 {
		sy = -1
	}
	dy = -dy
	err := dx + dy
	for ix1 != ix2 && iy1 != iy2 {
		driver.SetPixel(ix1, iy1, c)
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			ix1 += sx
		}
		if e2 <= dx {
			err += dx
			iy1 += sy
		}
	}
}

// DrawRect draws a rectangle
func DrawRect(x, y, w, h float64, c Color) {
	ix := int(x)
	iy := int(y)
	iw := int(w)
	ih := int(h)

	drawHLine(ix, ix+iw, iy, c)
	drawHLine(ix, ix+iw, iy+ih, c)
	drawVLine(ix, iy, iy+ih, c)
	drawVLine(ix+iw, iy, iy+ih, c)
}

// FillRect draws a filled rectangle
func FillRect(x, y, w, h float64, c Color) {
	ix := int(x)
	iy := int(y)
	iw := int(w)
	ih := int(h)
	driver.FillRect(ix, iy, iw, ih, c)
	// for i := 0; i <= ih; i++ {
	// 	drawHLine(ix, ix+iw, iy+i, c)
	// }
}

// DrawCircle draws a circle
func DrawCircle(x, y, r float64, c Color) {
	if r <= 0 {
		return
	}

	ix := int(x)
	iy := int(y)
	ir := int(r)

	ex := 0
	ey := ir
	err := 3 - 2*ir
	for ey >= ex {
		driver.SetPixel(ix+ex, iy-ey, c)
		driver.SetPixel(ix+ey, iy-ex, c)
		driver.SetPixel(ix+ey, iy+ex, c)
		driver.SetPixel(ix+ex, iy+ey, c)
		driver.SetPixel(ix-ex, iy+ey, c)
		driver.SetPixel(ix-ey, iy+ex, c)
		driver.SetPixel(ix-ey, iy-ex, c)
		driver.SetPixel(ix-ex, iy-ey, c)
		if err > 0 {
			err += 4*(ex-ey) + 10
			ey--
		} else {
			err += 4*ex + 6
		}
		ex++
	}
}

// FillCircle draws a filled circle
func FillCircle(x, y, r float64, c Color) {
	if r <= 0 {
		return
	}

	ix := int(x)
	iy := int(y)
	ir := int(r)

	ex := 0
	ey := ir
	err := 3 - 2*ir
	for ey >= ex {
		drawHLine(ix-ex, ix+ex, iy-ey, c)
		drawHLine(ix-ey, ix+ey, iy-ex, c)
		drawHLine(ix-ey, ix+ey, iy+ex, c)
		drawHLine(ix-ex, ix+ex, iy+ey, c)
		if err >= 0 {
			err += 4*(ex-ey) + 10
			ey--
		} else {
			err += 4*ex + 6
		}
		ex++
	}
}

// DrawEllipse draws an ellipse
func DrawEllipse(x, y, rx, ry float64, c Color) {
	a2 := int(rx * rx)
	b2 := int(ry * ry)
	fa2 := int(4 * a2)
	fb2 := int(4 * b2)

	ix := int(x)
	iy := int(y)
	irx := int(rx)
	iry := int(ry)

	ex := 0
	ey := iry
	sigma := 2*b2 + a2*(1-2*iry)
	for b2*ex <= a2*ey {
		driver.SetPixel(ix+ex, iy+ey, c)
		driver.SetPixel(ix-ex, iy+ey, c)
		driver.SetPixel(ix+ex, iy-ey, c)
		driver.SetPixel(ix-ex, iy-ey, c)
		if sigma >= 0 {
			sigma += fa2 * (1 - ey)
			ey--
		}
		sigma += b2 * ((4 * ex) + 6)
		ex++
	}

	ex = irx
	ey = 0
	sigma = 2*a2 + b2*(1-2*irx)
	for a2*ey <= b2*ex {
		driver.SetPixel(ix+ex, iy+ey, c)
		driver.SetPixel(ix-ex, iy+ey, c)
		driver.SetPixel(ix+ex, iy-ey, c)
		driver.SetPixel(ix-ex, iy-ey, c)
		if sigma >= 0 {
			sigma += fb2 * (1 - ex)
			ex--
		}
		sigma += a2 * ((4 * ey) + 6)
		ey++
	}
}

// FillEllipse draws a filled ellipse
func FillEllipse(x, y, rx, ry float64, c Color) {
	a2 := int(rx * rx)
	b2 := int(ry * ry)
	fa2 := int(4 * a2)
	fb2 := int(4 * b2)

	ix := int(x)
	iy := int(y)
	irx := int(rx)
	iry := int(ry)

	ex := 0
	ey := iry
	sigma := 2*b2 + a2*(1-2*iry)
	for b2*ex <= a2*ey {
		driver.HLine(ix-ex, ix+ex, iy+ey, c)
		driver.HLine(ix-ex, ix+ex, iy-ey, c)
		if sigma >= 0 {
			sigma += fa2 * (1 - ey)
			ey--
		}
		sigma += b2 * ((4 * ex) + 6)
		ex++
	}

	ex = irx
	ey = 0
	sigma = 2*a2 + b2*(1-2*irx)
	for a2*ey <= b2*ex {
		driver.HLine(ix-ex, ix+ex, iy+ey, c)
		driver.HLine(ix-ex, ix+ex, iy-ey, c)
		if sigma >= 0 {
			sigma += fb2 * (1 - ex)
			ex--
		}
		sigma += a2 * ((4 * ey) + 6)
		ey++
	}
}

func drawChar(font *Font, x, y int, ch byte, bk, fg Color) {
	if ch < font.FirstChar {
		ch = font.FirstChar
	}
	if ch > font.LastChar {
		ch = font.LastChar
	}

	firstByte := (int(ch) - int(font.FirstChar)) * font.H
	for i := 0; i < font.H; i++ {
		b := font.data[firstByte+i]
		for j := 0; j < font.W; j++ {
			if b&0x80 != 0 {
				driver.SetPixel(x+j, y+i, fg)
			} else if bk != Transparent {
				driver.SetPixel(x+j, y+i, bk)
			}
			b <<= 1
		}
	}
}

// DrawChar renders a character using the specified font. A Transparent color can be used for the background.
func DrawChar(font *Font, x, y float64, ch byte, bk, fg Color) {
	ix := int(x)
	iy := int(y)

	if ch < font.FirstChar {
		ch = font.FirstChar
	}
	if ch > font.LastChar {
		ch = font.LastChar
	}

	firstByte := (int(ch) - int(font.FirstChar)) * font.H
	for i := 0; i < font.H; i++ {
		b := font.data[firstByte+i]
		for j := 0; j < font.W; j++ {
			if b&0x80 != 0 {
				driver.SetPixel(ix+j, iy+i, fg)
			} else if bk != Transparent {
				driver.SetPixel(ix+j, iy+i, bk)
			}
			b <<= 1
		}
	}
}

// DrawString renders a string using the specified font. The background color can be Transparent
func DrawString(font *Font, x, y float64, str string, bk, fg Color) {
	ix := int(x)
	iy := int(y)
	fw := float64(font.W)
	sw := width
	if x > width || y > height || ix+len(str)*font.W < 0 || iy+font.H < 0 {
		return
	}
	for _, ch := range str {
		DrawChar(font, x, y, byte(ch), bk, fg)
		x += fw
		if x > sw {
			break
		}
	}
}
