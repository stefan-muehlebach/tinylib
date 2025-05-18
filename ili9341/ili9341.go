package ili9341

import (
	"image"
	"image/draw"
	"machine"
	"time"
)

const (
	bytesPerPixel       = 2
	pixfmt        uint8 = 0x05
)

type Config struct {
	Width            int16
	Height           int16
	Rotation         Rotation
	DisplayInversion bool
}

type Device struct {
	width    int16
	height   int16
	rotation Rotation
	driver   driver

	x0, x1 int16 // cached address window; prevents useless/expensive
	y0, y1 int16 // syscalls to PASET and CASET

	// ili *ILIImage

	cs  machine.Pin
	dc  machine.Pin
	rst machine.Pin
	rd  machine.Pin
}

// type Image = pixel.Image[pixel.RGB565BE]

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
	PIXFMT, 1, pixfmt,
	FRMCTR1, 2, 0x00, 0x18,
	DFUNCTR, 3, 0x08, 0x82, 0x27, // Display Function Control
	GAMMA_3G, 1, 0x02, // 3Gamma Function Disable
	GAMMASET, 1, 0x01, // Gamma curve selected
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

	w, h := d.Size()
	d.SetWindow(0, 0, w, h)
	d.SendCommand(RAMWR, nil)
	d.startWrite()
	d.driver.write16n(0x0000, int(w)*int(h))
	d.endWrite()

	// Because of (known) memory limitations, we cannot draw the whole in once
	// but can use this very limited space in order to cleverly decide the
	// smalles
	//
	// d.ili = NewILIImageBySize(50 * 320)
}

func (d *Device) Bounds() image.Rectangle {
	w, h := d.Size()
	return image.Rect(0, 0, int(w), int(h))
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
func (d *Device) Rotation() Rotation {
	return d.rotation
}

// SetRotation changes the rotation of the device (clock-wise).
func (d *Device) SetRotation(rotation Rotation) {
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
}

func (d *Device) SetScrollMargins(top, bottom int16) {
	if top+bottom <= d.height {
		middle := d.height - (top + bottom)
		cmdBuf[0] = uint8(top >> 8)
		cmdBuf[1] = uint8(top)
		cmdBuf[2] = uint8(middle >> 8)
		cmdBuf[3] = uint8(middle)
		cmdBuf[4] = uint8(bottom >> 8)
		cmdBuf[5] = uint8(bottom)
		d.SendCommand(VSCRDEF, cmdBuf[:6])
	}
}

// SetScroll sets the vertical scroll address of the display.
func (d *Device) ScrollTo(line int16) {
	cmdBuf[0] = uint8(line >> 8)
	cmdBuf[1] = uint8(line)
	d.SendCommand(VSCRSADD, cmdBuf[:2])
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

// Damit wird das Bild hinter img auf dem Bildschirm dargestellt. Mit rect
// kann bewirkt werden, dass nur ein Teil des Bildes aktualisiert wird und
// der Bildaufbau damit massiv beschleunigt werden kann.
// Rect muss innerhalb von img.Bounds() liegen; wird das Rechteck (0,0)-(0,0)
// (und somit ein leeres Rechteck) uebergeben, dann wird das gesamte Bild
// img auf dem Display dargestellt.
// func (d *Device) WriteImage(img image.Image, rect image.Rectangle) {
// 	var len int

// 	if rect.Empty() {
// 		rect = img.Bounds()
// 	}
// 	if !rect.In(img.Bounds()) {
// 		println("Is NOT inside; returning")
// 		return
// 	}
// 	convMin := rect.Min
// 	done := false
// 	for !done {
// 		convHeight := d.ili.Redim(convMin, rect.Dx(), 0)
// 		if !d.ili.Rect.In(rect) {
// 			d.ili.Rect = d.ili.Rect.Intersect(rect)
// 			done = true
// 		}

// 		d.ili.Convert(img.(*image.RGBA))
// 		d.SetWindow(int16(d.ili.Rect.Min.X), int16(d.ili.Rect.Min.Y),
// 			int16(d.ili.Rect.Dx()), int16(d.ili.Rect.Dy()))
// 		d.SendCommand(RAMWR, nil)
// 		d.startWrite()
// 		len = d.ili.Rect.Dy() * d.ili.Stride
// 		d.driver.write8sl(d.ili.Pix[:len:len])
// 		d.endWrite()
// 		convMin = convMin.Add(image.Point{0, convHeight})
// 	}
// }

func (d *Device) WriteImage(img image.Image, rect image.Rectangle) {
	buf := []byte{0x00, 0x00}

	if rect.Empty() {
		rect = img.Bounds()
	}
	if !rect.In(img.Bounds()) {
		println("Is NOT inside; returning")
		return
	}
	src := img.(*image.RGBA)
	d.SetWindow(int16(rect.Min.X), int16(rect.Min.Y),
		int16(rect.Dx()), int16(rect.Dy()))
	d.SendCommand(RAMWR, nil)
	d.startWrite()
	baseIdx := src.PixOffset(rect.Min.X, rect.Min.Y)
	for range rect.Dy() {
		idx := baseIdx
		for range rect.Dx() {
			s := src.Pix[idx : idx+3 : idx+3]
			r := s[0] & 0xF8
			g := s[1] & 0xFC
			b := s[2] & 0xF8
			buf[0] = (r) | (g >> 5)
			buf[1] = (g << 3) | (b >> 3)
			d.driver.write8sl(buf)
			idx += 4
		}
		baseIdx += src.Stride
	}
	d.endWrite()
}

func (d *Device) ReadImage(rect image.Rectangle, dst draw.Image) {
	dstImg := dst.(*image.RGBA)
	dstImg.Rect = rect
	dstImg.Stride = rect.Dx() * 3
	dataLen := rect.Dx() * rect.Dy() * 3
	if len(dstImg.Pix) < dataLen {
		println("ReadImage(): destination image is too small")
		return
	}
	d.SetWindow(int16(rect.Min.X), int16(rect.Min.Y), int16(rect.Dx()), int16(rect.Dy()))
	d.SendCommand(RAMRD, nil)
	d.startWrite()
	for i := range dataLen {
		dstImg.Pix[i] = d.driver.read8()
	}
	// d.driver.read8sl(dstImg.Pix[:dataLen])
	d.endWrite()
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
	read8() byte
	read8sl(b []byte)
}

func delay(m int) {
	t := time.Now().UnixNano() + int64(time.Duration(m*1000)*time.Microsecond)
	for time.Now().UnixNano() < t {
	}
}

func abs[T ~int | ~int16 | ~float64](v T) T {
	if v < 0 {
		return -v
	} else {
		return v
	}
}
