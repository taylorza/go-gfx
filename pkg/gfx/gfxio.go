package gfx

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
	KeySpace       = 0x20
	KeyEnd         = 0x23
	KeyHome        = 0x24
	KeyLeft        = 0x25
	KeyUp          = 0x26
	KeyRight       = 0x27
	KeyDown        = 0x28
	KeyInsert      = 0x2D
	KeyDelete      = 0x2E
	KeyA           = 0x41
	KeyB           = 0x42
	KeyC           = 0x43
	KeyD           = 0x44
	KeyE           = 0x45
	KeyF           = 0x46
	KeyG           = 0x47
	KeyH           = 0x48
	KeyI           = 0x49
	KeyJ           = 0x4A
	KeyK           = 0x4B
	KeyL           = 0x4C
	KeyM           = 0x4D
	KeyN           = 0x4E
	KeyO           = 0x4F
	KeyP           = 0x50
	KeyQ           = 0x51
	KeyR           = 0x52
	KeyS           = 0x53
	KeyT           = 0x54
	KeyU           = 0x55
	KeyV           = 0x56
	KeyW           = 0x57
	KeyX           = 0x58
	KeyY           = 0x59
	KeyZ           = 0x5A
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

type ioManager struct {
	keymap [255]Key

	keysPhysical [255]bool
	keysLogical  [255]keyState
	mouseX       float64
	mouseY       float64
}

type keyState struct {
	justPressed bool
	pressed     bool
}

func (io *ioManager) setKeyMapping(scancode byte, key Key) {
	io.keymap[scancode] = key
}

func (io *ioManager) setKeyPressed(key Key, state bool) {
	io.keysPhysical[key] = state
}

func (io *ioManager) setMappedKeyPressed(scanCode byte, state bool) {
	io.keysPhysical[io.keymap[scanCode]] = state
}

func (io *ioManager) updateMouse(x, y int) {
	io.mouseX = float64(x) / scaleX
	io.mouseY = float64(y) / scaleY
}

func (io *ioManager) keyJustPressed(key Key) bool {
	if io.keysLogical[key].justPressed {
		io.keysLogical[key].justPressed = false
		return true
	}
	return false
}

func (io *ioManager) keyPressed(key Key) bool {
	return io.keysLogical[key].pressed
}

func (io *ioManager) mouseXY() (float64, float64) {
	return io.mouseX, io.mouseY
}

func (io *ioManager) update(delta float64) {
	for i, pressed := range iomgr.keysPhysical {
		if pressed {
			if !io.keysLogical[i].pressed {
				io.keysLogical[i].justPressed = true
			}
			io.keysLogical[i].pressed = true
		} else {
			io.keysLogical[i].pressed = false
			io.keysLogical[i].justPressed = false
		}
	}
}
