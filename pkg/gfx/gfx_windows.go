// +build windows

package gfx

import (
	"reflect"
	"syscall"
	"unsafe"

	"github.com/taylorza/w32"
)

func init() {
	driver = &windowsDriver{
		windowClassName: syscall.StringToUTF16Ptr("GO-GRAFIX-WINDOW"),
		dibPixels:       make([]Color, 0, 0),
		renderPeriod:    1.0 / 60.0,
	}
}

type keyState struct {
	justPressed bool
	pressed     bool
}

type windowsDriver struct {
	windowClassName *uint16
	hMainWnd        w32.HWND
	surfaceDC       w32.HDC
	width           int
	height          int
	backBuffer      []Color
	dibPixels       []Color
	scaleX          float64
	scaleY          float64

	keysPhysical [255]bool
	keysLogical  [255]keyState
	mouseX       float64
	mouseY       float64

	rendering     bool
	renderPeriod  float64
	renderElapsed float64

	dib    w32.HBITMAP
	olddib w32.HGDIOBJ
}

func (e *windowsDriver) Clear(c Color) {
	// Exponentially copy more pixels into the buffer
	// [1] seed value at postion 0
	// [1] copy to position 1 -> [1, 1]
	// [1,1] copy to position 2 -> [1,1,1,1]
	// [1,1,1,1] copy to position 4 [1,1,1,1,1,1,1,1] ...
	// [1,1,1,1,1,1,1,1] copy to position 8 [1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1] ...
	//*(*Color)(unsafe.Pointer(&e.backBuffer[0])) = c
	e.backBuffer[0] = c
	for i := 1; i < e.width*e.height; i *= 2 {
		copy(e.backBuffer[i:], e.backBuffer[0:i])
	}
}

func (e *windowsDriver) SetPixel(x, y int, c Color) {
	if x < 0 || x >= e.width || y < 0 || y >= e.height {
		return
	}
	e.backBuffer[(y*e.width)+x] = c
}

func (e *windowsDriver) FillRect(x, y, w, h int, c Color) {
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

	sptr := uintptr(unsafe.Pointer(&e.backBuffer[0])) + uintptr(y*e.width*4+x*4)
	for i := uintptr(0); i < uintptr(w); i++ {
		*(*Color)(unsafe.Pointer(sptr + i*4)) = c
	}
	for i := 1; i < h; i++ {
		copy(e.backBuffer[(y+i)*e.width+x:], e.backBuffer[(y+i-1)*e.width+x:(y+i-1)*e.width+x+w])
	}
}

func (e *windowsDriver) DrawTexture(x, y int, srcX, srcY, srcW, srcH int, t *Texture) {
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
	bufferRowOffset := uintptr(y*e.width + x)

	tptr := uintptr(unsafe.Pointer(&t.pixels[0]))
	sptr := uintptr(unsafe.Pointer(&e.backBuffer[0]))

	for ty := y1; ty < y2; ty++ {
		i := textureRowOffset
		j := bufferRowOffset
		for tx := x1; tx < x2; tx++ {
			c := *(*uint32)(unsafe.Pointer(tptr + i*4))
			if c != uint32(transparent) {
				*(*uint32)(unsafe.Pointer(sptr + j*4)) = c
			}
			i++
			j++
		}
		textureRowOffset += uintptr(t.W)
		bufferRowOffset += uintptr(e.width)
	}
}

func (e *windowsDriver) HLine(x1, x2, y int, c Color) {
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

	pixels := uintptr(unsafe.Pointer(&e.backBuffer[0]))
	pixels += uintptr((y * e.width * 4) + x1*4)
	for i := uintptr(0); i <= uintptr(x2-x1); i++ {
		*(*Color)(unsafe.Pointer(pixels + i*4)) = c
	}
}

func (e *windowsDriver) VLine(x, y1, y2 int, c Color) {
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

	pixels := uintptr(unsafe.Pointer(&e.backBuffer[0]))
	pixels += uintptr((y1 * e.width * 4) + x*4)
	for i := uintptr(0); i <= uintptr(y2-y1); i++ {
		*(*Color)(unsafe.Pointer(pixels + uintptr(e.width)*i*4)) = c
	}
}
func (e *windowsDriver) CreateWindow(x, y, w, h, xscale, yscale int) bool {
	e.width = w
	e.height = h
	hInst := w32.GetModuleHandle("")

	wcex := &w32.WNDCLASSEX{
		Size:       uint32(unsafe.Sizeof(w32.WNDCLASSEX{})),
		Style:      w32.CS_HREDRAW | w32.CS_VREDRAW | w32.CS_OWNDC,
		ClsExtra:   0,
		WndExtra:   0,
		Instance:   hInst,
		Icon:       w32.LoadIcon(0, w32.MakeIntResource(w32.IDI_APPLICATION)),
		Cursor:     w32.LoadCursor(0, w32.MakeIntResource(w32.IDC_ARROW)),
		MenuName:   nil,
		ClassName:  e.windowClassName,
		Background: w32.HBRUSH(w32.COLOR_WINDOW + 1),
		WndProc:    syscall.NewCallback(e.wndProc),
	}

	wc := w32.RegisterClassEx(wcex)
	if wc == 0 {
		panic("Failed to register window class")
	}

	exStyle := uint(w32.WS_EX_APPWINDOW | w32.WS_EX_WINDOWEDGE)
	style := uint(w32.WS_CAPTION | w32.WS_SYSMENU | w32.WS_BORDER)

	rc := w32.RECT{Left: 0, Top: 0, Right: int32(w * xscale), Bottom: int32(h * yscale)}
	w32.AdjustWindowRectEx(&rc, style, false, exStyle)

	e.hMainWnd = w32.CreateWindowEx(
		exStyle,
		e.windowClassName,
		syscall.StringToUTF16Ptr(""),
		style,
		x, y, int(rc.Right-rc.Left), int(rc.Bottom-rc.Top),
		w32.HWND(0),
		w32.HMENU(0),
		hInst,
		unsafe.Pointer(nil))
	if e.hMainWnd == 0 {
		panic("Failed to create window")
	}

	w32.ShowWindow(e.hMainWnd, w32.SW_SHOWDEFAULT)
	w32.UpdateWindow(e.hMainWnd)

	return true
}

func (e *windowsDriver) CreateDevice() bool {
	hdc := w32.GetDC(e.hMainWnd)
	if hdc == 0 {
		return false
	}
	defer w32.ReleaseDC(e.hMainWnd, hdc)

	e.surfaceDC = w32.CreateCompatibleDC(hdc)
	if e.surfaceDC == 0 {
		return false
	}

	pbmi := &w32.BITMAPINFO{
		BmiHeader: w32.BITMAPINFOHEADER{
			BiSize:          uint32(unsafe.Sizeof(w32.BITMAPINFOHEADER{})),
			BiWidth:         int32(e.width),
			BiHeight:        int32(-e.height),
			BiPlanes:        1,
			BiBitCount:      32,
			BiCompression:   w32.BI_RGB,
			BiSizeImage:     0,
			BiXPelsPerMeter: 0,
			BiYPelsPerMeter: 0,
			BiClrUsed:       0,
			BiClrImportant:  0,
		},
	}

	var pbits unsafe.Pointer
	e.dib = w32.CreateDIBSection(e.surfaceDC, pbmi, w32.DIB_RGB_COLORS, &pbits, 0, 0)
	if e.dib == 0 {
		return false
	}

	var bmp w32.BITMAP
	w32.GetObject(w32.HGDIOBJ(e.dib), unsafe.Sizeof(w32.BITMAP{}), unsafe.Pointer(&bmp))
	bytes := bmp.BmWidth * bmp.BmHeight

	e.backBuffer = make([]Color, bytes, bytes)

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&e.dibPixels))
	sh.Len = int(bytes)
	sh.Cap = int(bytes)
	sh.Data = uintptr(pbits)

	e.olddib = w32.SelectObject(e.surfaceDC, w32.HGDIOBJ(e.dib))

	return true
}

func (e *windowsDriver) cleanup() {
	if e.olddib != 0 && e.surfaceDC != 0 {
		w32.SelectObject(e.surfaceDC, e.olddib)
		w32.DeleteDC(e.surfaceDC)
		w32.DeleteObject(w32.HGDIOBJ(e.dib))
	}
}

func (e *windowsDriver) StartEventLoop() {
	var msg w32.MSG
	for w32.GetMessage(&msg, 0, 0, 0) {
		w32.TranslateMessage(&msg)
		w32.DispatchMessage(&msg)
	}
	e.cleanup()
}

func (e *windowsDriver) SetWindowTitle(title string) {
	w32.SetWindowText(e.hMainWnd, title)
}

func (e *windowsDriver) Update(delta float64) {
	for i, pressed := range e.keysPhysical {
		if pressed {
			if !e.keysLogical[i].pressed {
				e.keysLogical[i].justPressed = true
			}
			e.keysLogical[i].pressed = true
		} else {
			e.keysLogical[i].pressed = false
			e.keysLogical[i].justPressed = false
		}
	}
}

func (e *windowsDriver) KeyJustPressed(key Key) bool {
	if e.keysLogical[key].justPressed {
		e.keysLogical[key].justPressed = false
		return true
	}
	return false
}

func (e *windowsDriver) KeyPressed(key Key) bool {
	return e.keysLogical[key].pressed
}

func (e *windowsDriver) MouseXY() (float64, float64) {
	return e.mouseX / e.scaleX, e.mouseY / e.scaleY
}

func (e *windowsDriver) Render(delta float64) {
	e.renderElapsed += delta
	if e.rendering {
		return
	}

	if e.renderElapsed >= e.renderPeriod && !e.rendering {
		e.rendering = true
		copy(e.dibPixels, e.backBuffer)
		w32.PostMessage(e.hMainWnd, w32.WM_USER+0x100, 0, 0)
		e.renderElapsed -= e.renderPeriod
	}
}

func (e *windowsDriver) wndProc(hwnd w32.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case w32.WM_DESTROY:
		shutdown()
		w32.PostQuitMessage(0)
	case w32.WM_SIZE:
		l := int(lParam)
		e.scaleX = float64(l&0xffff) / float64(e.width)
		e.scaleY = float64((l>>16)&0xffff) / float64(e.height)
	case w32.WM_MOUSEMOVE:
		l := int(lParam)
		e.mouseX = float64(l & 0xffff)
		e.mouseY = float64((l >> 16) & 0xffff)
	case w32.WM_LBUTTONDOWN:
		e.keysPhysical[1] = true
	case w32.WM_LBUTTONUP:
		e.keysPhysical[1] = false
	case w32.WM_RBUTTONDOWN:
		e.keysPhysical[2] = true
	case w32.WM_RBUTTONUP:
		e.keysPhysical[2] = false
	case w32.WM_MBUTTONDOWN:
		e.keysPhysical[3] = true
	case w32.WM_MBUTTONUP:
		e.keysPhysical[3] = false
	case w32.WM_KEYDOWN:
		e.keysPhysical[byte(wParam)] = true
	case w32.WM_KEYUP:
		e.keysPhysical[byte(wParam)] = false

	case w32.WM_USER + 0x100:
		hdc := w32.GetDC(hwnd)
		rc := w32.GetClientRect(hwnd)
		w32.StretchBlt(hdc, 0, 0, int(rc.Right-rc.Left), int(rc.Bottom-rc.Top), e.surfaceDC, 0, 0, e.width, e.height, w32.SRCCOPY)
		w32.ReleaseDC(hwnd, hdc)
		e.rendering = false
	default:
		return w32.DefWindowProc(hwnd, msg, wParam, lParam)
	}
	return 0
}
