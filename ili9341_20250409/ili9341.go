package ili9341

import (
	"errors"
	"image"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pixel"
)

type Config struct {
	Width            int16
	Height           int16
	Rotation         drivers.Rotation
	DisplayInversion bool
}

type Device struct {
	width    int16
	height   int16
	rotation drivers.Rotation
	driver   driver

	fillColor, lineColor color.RGBA

	x0, x1 int16 // cached address window; prevents useless/expensive
	y0, y1 int16 // syscalls to PASET and CASET

	dc  machine.Pin
	cs  machine.Pin
	rst machine.Pin
	rd  machine.Pin
}

type Image = pixel.Image[pixel.RGB565BE]

var cmdBuf [6]byte

var initCmd = []byte{
	0xEF, 3, 0x03, 0x80, 0x02,
	PWCTRLB, 3, 0x00, 0xC1, 0x30,
	PWOSEQCTR, 4, 0x64, 0x03, 0x12, 0x81,
	DRVTICTRLA, 3, 0x85, 0x00, 0x78,
	// Disable this, since it's an EXTC command and these values are the
	// default anyway
	PWCTRLA, 5, 0x39, 0x2C, 0x00, 0x34, 0x02,
	PMPRTCTR, 1, 0x20,
	DRVTICTRLB, 2, 0x00, 0x00,
	PWCTR1, 1, 0x23, // Power control VRH[5:0]
	PWCTR2, 1, 0x10, // Power control SAP[2:0];BT[3:0]
	VMCTR1, 2, 0x3e, 0x28, // VCM control
	VMCTR2, 1, 0x86, // VCM control2
	MADCTL, 1, 0x48, // Memory Access Control
	VSCRSADD, 1, 0x00, // Vertical scroll zero
	PIXFMT, 1, 0x55,
	FRMCTR1, 2, 0x00, 0x18,
	DFUNCTR, 3, 0x08, 0x82, 0x27, // Display Function Control
	GAMMA_3G, 1, 0x00, // 3Gamma Function Disable
	GAMMASET, 1, 0x01, // Gamma curve selected
	// Disabled, EXTC commands and special in any way
	GMCTRP1, 15, 0x0F, 0x31, 0x2B, 0x0C, 0x0E, 0x08,
	0x4E, 0xF1, 0x37, 0x07, 0x10, 0x03, 0x0E, 0x09, 0x00,
	GMCTRN1, 15, 0x00, 0x0E, 0x14, 0x03, 0x11, 0x07, // Set Gamma
	0x31, 0xC1, 0x48, 0x08, 0x0F, 0x0C, 0x31, 0x36, 0x0F,
}

// Configure prepares display for use
func (d *Device) Configure(config Config) {

	if config.Width == 0 {
		config.Width = TFTWIDTH
	}
	if config.Height == 0 {
		config.Height = TFTHEIGHT
	}
	d.width = config.Width
	d.height = config.Height
	d.rotation = config.Rotation

	// try to pick an initial cache miss for one of the points
	d.x0, d.x1 = -(d.width + 1), d.x0
	d.y0, d.y1 = -(d.height + 1), d.y0

	output := machine.PinConfig{machine.PinOutput}

	// configure chip select if there is one
	if d.cs != machine.NoPin {
		d.cs.Configure(output)
		d.cs.High() // deselect
	}

	d.dc.Configure(output)
	d.dc.High() // data mode

	// driver-specific configuration
	d.driver.configure(&config)

	if d.rd != machine.NoPin {
		d.rd.Configure(output)
		d.rd.High()
	}

	// reset the display
	if d.rst != machine.NoPin {
		// configure hardware reset if there is one
		d.rst.Configure(output)
		d.rst.High()
		delay(100)
		d.rst.Low()
		delay(100)
		d.rst.High()
		delay(200)
	} else {
		// if no hardware reset, send software reset
		d.SendCommand(SWRESET, nil)
		delay(150)
	}

	if config.DisplayInversion {
		initCmd = append(initCmd, INVON, 0x80)
	}

	initCmd = append(initCmd,
		SLPOUT, 0x80, // Exit Sleep
		DISPON, 0x80, // Display on
		0x00, // End of list
	)
	for i, c := 0, len(initCmd); i < c; {
		cmd := initCmd[i]
		if cmd == 0x00 {
			break
		}
		x := initCmd[i+1]
		numArgs := int(x & 0x7F)
		d.SendCommand(cmd, initCmd[i+2:i+2+numArgs])
		if x&0x80 > 0 {
			delay(150)
		}
		i += numArgs + 2
	}
	d.SetRotation(d.rotation)
}

func (d *Device) ColorModel() color.Model {
	return color.RGBAModel
}

func (d *Device) Bounds() image.Rectangle {
	w, h := d.Size()
	return image.Rect(0, 0, int(w), int(h))
}

func (d *Device) At(x, y int) color.Color {
	buf := []byte{0, 0, 0}

	d.SetWindow(int16(x), int16(y), 1, 1)
	d.SendCommand(RAMRD, nil)
	d.startWrite()
	d.driver.read8sl(buf)
	d.endWrite()
	return color.RGBA{buf[0], buf[1], buf[2], 0x0f}
}

func (d *Device) Set(x, y int, c color.Color) {
	c565 := rgbaTo565(c)
	// c666 := RGBATo666(c)
	d.SetWindow(int16(x), int16(y), 1, 1)
	d.SendCommand(RAMWR, nil)
	d.startWrite()
	d.driver.write16(c565)
	// d.driver.write24(c666)
	d.endWrite()
}

func (d *Device) FillColor() color.RGBA {
	return d.fillColor
}

func (d *Device) SetFillColor(c color.RGBA) {
	d.fillColor = c
}

func (d *Device) LineColor() color.RGBA {
	return d.lineColor
}

func (d *Device) SetLineColor(c color.RGBA) {
	d.lineColor = c
}

// Size returns the current size of the display.
func (d *Device) Size() (x, y int16) {
	switch d.rotation {
	case Rotation90, Rotation270, Rotation90Mirror, Rotation270Mirror:
		return d.height, d.width
	default: // Rotation0, Rotation180, etc
		return d.width, d.height
	}
}

// Rotation returns the current rotation of the device.
func (d *Device) Rotation() drivers.Rotation {
	return d.rotation
}

// SetRotation changes the rotation of the device (clock-wise).
func (d *Device) SetRotation(rotation drivers.Rotation) error {
	madctl := uint8(0)
	switch rotation % 8 {
	case Rotation0:
		madctl = MADCTL_MX | MADCTL_BGR
	case Rotation90:
		madctl = MADCTL_MV | MADCTL_BGR
	case Rotation180:
		madctl = MADCTL_MY | MADCTL_BGR | MADCTL_ML
	case Rotation270:
		madctl = MADCTL_MX | MADCTL_MY | MADCTL_MV | MADCTL_BGR | MADCTL_ML
	case Rotation0Mirror:
		madctl = MADCTL_BGR
	case Rotation90Mirror:
		madctl = MADCTL_MY | MADCTL_MV | MADCTL_BGR | MADCTL_ML
	case Rotation180Mirror:
		madctl = MADCTL_MX | MADCTL_MY | MADCTL_BGR | MADCTL_ML
	case Rotation270Mirror:
		madctl = MADCTL_MX | MADCTL_MY | MADCTL_MV | MADCTL_BGR | MADCTL_ML
	}
	cmdBuf[0] = madctl
	d.SendCommand(MADCTL, cmdBuf[:1])
	d.rotation = rotation
	return nil
}

// SetScroll sets the vertical scroll address of the display.
func (d *Device) SetScroll(line int16) {
	cmdBuf[0] = uint8(line >> 8)
	cmdBuf[1] = uint8(line)
	d.SendCommand(VSCRSADD, cmdBuf[:2])
}

// SetPixel modifies the internal buffer.
func (d *Device) SetPixel(x, y int16, c color.RGBA) {
	d.Set(int(x), int(y), c)
}

// FillScreen fills the screen with a given color
func (d *Device) FillScreen(c color.RGBA) {
	w, h := d.Size()
	d.fillRect(0, 0, w, h, c)
}

// DrawRectangle draws a rectangle at given coordinates with a color
func (d *Device) DrawRectangle(x, y, w, h int16) {
	d.drawFastHLine(x, y, w, d.lineColor)
	d.drawFastHLine(x, y+h-1, w, d.lineColor)
	d.drawFastVLine(x, y, h, d.lineColor)
	d.drawFastVLine(x+w-1, y, h, d.lineColor)
}

func (d *Device) DrawRoundedRectangle(x, y, w, h, r int16) {
	maxRadius := max(w, h) / 2
	if r > maxRadius {
		r = maxRadius
	}
	d.drawFastHLine(x+r, y, w-2*r, d.lineColor)
	d.drawFastHLine(x+r, y+h-1, w-2*r, d.lineColor)
	d.drawFastVLine(x, y+r, h-2*r, d.lineColor)
	d.drawFastVLine(x+w-1, y+r, h-2*r, d.lineColor)

	d.drawCircleHelper(x+r, y+r, r, upperLeft, d.lineColor)
	d.drawCircleHelper(x+w-r-1, y+r, r, upperRight, d.lineColor)
	d.drawCircleHelper(x+w-r-1, y+h-r-1, r, lowerRight, d.lineColor)
	d.drawCircleHelper(x+r, y+h-r-1, r, lowerLeft, d.lineColor)
}

// FillRectangle fills a rectangle at given coordinates with a color
func (d *Device) FillRectangle(x, y, width, height int16) {
	w, h := d.Size()
	if x < 0 || y < 0 || width <= 0 || height <= 0 ||
		x >= w || (x+width) > w || y >= h || (y+height) > h {
		return
	}
	d.fillRect(x, y, width, height, d.fillColor)
}

func (d *Device) FillRoundedRectangle(x, y, w, h, r int16) {
	maxRadius := max(w, h) / 2
	if r > maxRadius {
		r = maxRadius
	}
	d.fillRect(x+r, y, w-2*r, h, d.fillColor)
	d.fillCircleHelper(x+w-r-1, y+r, r, h-2*r-1, upperLeft, d.fillColor)
	d.fillCircleHelper(x+r, y+r, r, h-2*r-1, upperRight, d.fillColor)
}

func (d *Device) DrawTriangle(x0, y0, x1, y1, x2, y2 int16) {
	d.drawLine(x0, y0, x1, y1, d.lineColor)
	d.drawLine(x1, y1, x2, y2, d.lineColor)
	d.drawLine(x2, y2, x0, y0, d.lineColor)
}

func (d *Device) FillTriangle(x0, y0, x1, y1, x2, y2 int16) {
	var a, b, y, last int16

	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y0 > y1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	d.startWrite()
	if y0 == y2 {
		a, b = x0, x0
		if x1 < a {
			a = x1
		} else if x1 > b {
			b = x1
		}
		if x2 < a {
			a = x2
		} else if x2 > b {
			b = x2
		}
		d.drawFastHLine(a, y0, b-a+1, d.fillColor)
		d.endWrite()
		return
	}

	dx01 := x1 - x0
	dy01 := y1 - y0
	dx02 := x2 - x0
	dy02 := y2 - y0
	dx12 := x2 - x1
	dy12 := y2 - y1

	sa, sb := int32(0), int32(0)

	if y1 == y2 {
		last = y1
	} else {
		last = y1 - 1
	}

	for y = y0; y <= last; y++ {
		a = x0 + int16(sa)/dy01
		b = x0 + int16(sb)/dy02
		sa += int32(dx01)
		sb += int32(dx02)
		if a > b {
			a, b = b, a
		}
		d.drawFastHLine(a, y, b-a+1, d.fillColor)
	}

	sa = int32(dx12 * (y - y1))
	sb = int32(dx02 * (y - y0))
	for ; y <= y2; y++ {
		a = x1 + int16(sa)/dy12
		b = x0 + int16(sb)/dy02
		sa += int32(dx12)
		sb += int32(dx02)
		if a > b {
			a, b = b, a
		}
		d.drawFastHLine(a, y, b-a+1, d.fillColor)
	}
	d.endWrite()
}

func (d *Device) DrawLine(x0, y0, x1, y1 int16) {
	if x0 == x1 {
		if y0 > y1 {
			y0, y1 = y1, y0
		}
		d.drawFastVLine(x0, y0, y1-y0, d.lineColor)
	} else if y0 == y1 {
		if x0 > x1 {
			x0, x1 = x1, x0
		}
		d.drawFastHLine(x0, y0, x1-x0, d.lineColor)
	} else {
		d.drawLine(x0, y0, x1, y1, d.lineColor)
	}
}

func (d *Device) DrawCircle(x0, y0, r int16) {
	f := 1 - r
	ddF_x := int16(1)
	ddF_y := -2 * r
	x := int16(0)
	y := r

	d.SetPixel(x0, y0+r, d.lineColor)
	d.SetPixel(x0, y0-r, d.lineColor)
	d.SetPixel(x0+r, y0, d.lineColor)
	d.SetPixel(x0-r, y0, d.lineColor)

	for x < y {
		if f >= 0 {
			y--
			ddF_y += 2
			f += ddF_y
		}
		x++
		ddF_x += 2
		f += ddF_x

		d.SetPixel(x0+x, y0+y, d.lineColor)
		d.SetPixel(x0-x, y0+y, d.lineColor)
		d.SetPixel(x0+x, y0-y, d.lineColor)
		d.SetPixel(x0-x, y0-y, d.lineColor)
		d.SetPixel(x0+y, y0+x, d.lineColor)
		d.SetPixel(x0-y, y0+x, d.lineColor)
		d.SetPixel(x0+y, y0-x, d.lineColor)
		d.SetPixel(x0-y, y0-x, d.lineColor)
	}
}

func (d *Device) FillCircle(x0, y0, r int16) {
	d.drawFastVLine(x0, y0-r, 2*r, d.fillColor)
	d.fillCircleHelper(x0, y0, r, 0, upperLeft|upperRight, d.fillColor)
}

// DrawBitmap copies the bitmap to the internal buffer on the screen at the
// given coordinates. It returns once the image data has been sent completely.
func (d *Device) DrawBitmap(x, y int16, bitmap Image) error {
	width, height := bitmap.Size()
	return d.DrawRGBBitmap8(x, y, bitmap.RawBuffer(), int16(width), int16(height))
}

func (d *Device) DrawRGBBitmap8(x, y int16, data []uint8, w, h int16) error {
	k, i := d.Size()
	if x < 0 || y < 0 || w <= 0 || h <= 0 ||
		x >= k || (x+w) > k || y >= i || (y+h) > i {
		return errors.New("rectangle coordinates outside display area")
	}
	d.SetWindow(x, y, w, h)
	d.SendCommand(RAMWR, nil)
	d.startWrite()
	d.driver.write8sl(data)
	d.endWrite()
	return nil
}

// DrawRGBBitmap copies an RGB bitmap to the internal buffer at given coordinates
//
// Deprecated: use DrawBitmap instead.
func (d *Device) DrawRGBBitmap(x, y int16, data []uint16, w, h int16) error {
	k, i := d.Size()
	if x < 0 || y < 0 || w <= 0 || h <= 0 ||
		x >= k || (x+w) > k || y >= i || (y+h) > i {
		return errors.New("rectangle coordinates outside display area")
	}
	d.SetWindow(x, y, w, h)
	d.SendCommand(RAMWR, nil)
	d.startWrite()
	d.driver.write16sl(data)
	d.endWrite()
	return nil
}

// Methods for reading registers or sending commands to the ILI9341.
func (d *Device) ReadRegister(reg byte) byte {
	d.startWrite()
	d.dc.Low()
	d.driver.write8(reg)
	d.dc.High()
	b := d.driver.read8()
	d.endWrite()
	return b
}

func (d *Device) SendCommand(cmd byte, data []byte) {
	d.startWrite()
	d.dc.Low()
	d.driver.write8(cmd)
	d.dc.High()
	if data != nil {
		d.driver.write8sl(data)
	}
	d.endWrite()
}

//----------------------------------------------------------------------------
//
// Private methods

func (d *Device) fillRect(x, y, width, height int16, c color.RGBA) {
	c565 := rgbaTo565(c)
	d.SetWindow(x, y, width, height)
	d.SendCommand(RAMWR, nil)
	d.startWrite()
	d.driver.write16n(c565, int(width)*int(height))
	d.endWrite()
}

type cornerName uint8

const (
	upperLeft  cornerName = 1 << iota
	upperRight            = 1 << iota
	lowerRight            = 1 << iota
	lowerLeft             = 1 << iota
)

func (d *Device) drawCircleHelper(x0, y0, r int16, corner cornerName, c color.RGBA) {
	f := 1 - r
	ddF_x := int16(1)
	ddF_y := -2 * r
	x := int16(0)
	y := r

	for x < y {
		if f >= 0 {
			y--
			ddF_y += 2
			f += ddF_y
		}
		x++
		ddF_x += 2
		f += ddF_x
		if corner&upperLeft != 0 {
			d.SetPixel(x0-y, y0-x, c)
			d.SetPixel(x0-x, y0-y, c)
		}
		if corner&upperRight != 0 {
			d.SetPixel(x0+x, y0-y, c)
			d.SetPixel(x0+y, y0-x, c)
		}
		if corner&lowerRight != 0 {
			d.SetPixel(x0+x, y0+y, c)
			d.SetPixel(x0+y, y0+x, c)
		}
		if corner&lowerLeft != 0 {
			d.SetPixel(x0-y, y0+x, c)
			d.SetPixel(x0-x, y0+y, c)
		}
	}
}

func (d *Device) fillCircleHelper(x0, y0, r, delta int16, corners cornerName, c color.RGBA) {
	f := 1 - r
	ddF_x := int16(1)
	ddF_y := -2 * r
	x := int16(0)
	y := r
	px := x
	py := y

	delta++

	for x < y {
		if f >= 0 {
			y--
			ddF_y += 2
			f += ddF_y
		}
		x++
		ddF_x += 2
		f += ddF_x
		if x < (y + 1) {
			if corners&upperLeft != 0 {
				d.drawFastVLine(x0+x, y0-y, 2*y+delta, c)
			}
			if corners&upperRight != 0 {
				d.drawFastVLine(x0-x, y0-y, 2*y+delta, c)
			}
		}
		if y != py {
			if corners&upperLeft != 0 {
				d.drawFastVLine(x0+py, y0-px, 2*px+delta, c)
			}
			if corners&upperRight != 0 {
				d.drawFastVLine(x0-py, y0-px, 2*px+delta, c)
			}
			py = y
		}
		px = x
	}
}

func (d *Device) drawLine(x0, y0, x1, y1 int16, c color.RGBA) {
	steep := abs(y1-y0) > abs(x1-x0)
	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}
	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	dx := x1 - x0
	dy := abs(y1 - y0)
	err := dx / 2
	ystep := int16(0)

	if y0 < y1 {
		ystep = +1
	} else {
		ystep = -1
	}

	for ; x0 <= x1; x0++ {
		if steep {
			d.SetPixel(y0, x0, c)
		} else {
			d.SetPixel(x0, y0, c)
		}
		err -= dy
		if err < 0 {
			y0 += ystep
			err += dx
		}
	}
}

// DrawFastVLine draws a vertical line faster than using SetPixel
func (d *Device) drawFastVLine(x, y, h int16, c color.RGBA) {
	d.fillRect(x, y, 1, h, c)
}

// DrawFastHLine draws a horizontal line faster than using SetPixel
func (d *Device) drawFastHLine(x, y, w int16, c color.RGBA) {
	d.fillRect(x, y, w, 1, c)
}

// setWindow prepares the screen to be modified at a given rectangle
func (d *Device) SetWindow(x, y, w, h int16) {
	x1 := x + w - 1
	if x != d.x0 || x1 != d.x1 {
		cmdBuf[0] = uint8(x >> 8)
		cmdBuf[1] = uint8(x)
		cmdBuf[2] = uint8(x1 >> 8)
		cmdBuf[3] = uint8(x1)
		d.SendCommand(CASET, cmdBuf[:4])
		d.x0, d.x1 = x, x1
	}
	y1 := y + h - 1
	if y != d.y0 || y1 != d.y1 {
		cmdBuf[0] = uint8(y >> 8)
		cmdBuf[1] = uint8(y)
		cmdBuf[2] = uint8(y1 >> 8)
		cmdBuf[3] = uint8(y1)
		d.SendCommand(PASET, cmdBuf[:4])
		d.y0, d.y1 = y, y1
	}
}

//go:inline
func (d *Device) startWrite() {
	if d.cs != machine.NoPin {
		d.cs.Low()
	}
}

//go:inline
func (d *Device) endWrite() {
	if d.cs != machine.NoPin {
		d.cs.High()
	}
}

type driver interface {
	configure(config *Config)
	write8(b byte)
	write8n(b byte, n int)
	write8sl(b []byte)
	write16(data uint16)
	write16n(data uint16, n int)
	write16sl(data []uint16)
	// write24(data uint32)
	// write24n(data uint32, n int)
	read8() byte
	read8sl(b []byte)
}

func delay(m int) {
	t := time.Now().UnixNano() + int64(time.Duration(m*1000)*time.Microsecond)
	for time.Now().UnixNano() < t {
	}
}

// RGBATo565 converts a color.RGBA to uint16 used in the display
func rgbaTo565(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) |
		((g & 0xFC00) >> 5) |
		((b & 0xF800) >> 11))
}

func rgbaTo666(c color.Color) uint32 {
	r, g, b, _ := c.RGBA()
	return uint32(((r & 0xFF00) << 8) +
		(g & 0xFF00) +
		((b & 0xFF00) >> 8))
}

func abs[T ~int | ~int16 | ~float64](v T) T {
	if v < 0 {
		return -v
	} else {
		return v
	}
}
