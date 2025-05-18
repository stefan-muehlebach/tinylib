package tinylib

import (
	"image"
	"image/color"
    "tinylib/colors"
)

//----------------------------------------------------------------------------

const (
	bytesPerPixel = 2
)

type TinyImage struct {
	Pix    []byte
	Stride int
	Rect   image.Rectangle
}

func NewTinyImage(r image.Rectangle) *TinyImage {
	i := &TinyImage{
	    Pix: make([]byte, r.Dx() * r.Dy() * bytesPerPixel),
        Stride: r.Dx() * bytesPerPixel,
	    Rect: r,
    }
	return i
}

func (i *TinyImage) ColorModel() color.Model {
	return colors.TinyModel
}

func (i *TinyImage) Bounds() image.Rectangle {
	return i.Rect
}

func (i *TinyImage) At(x, y int) color.Color {
	return i.TinyColorAt(x, y)
}

func (i *TinyImage) TinyColorAt(x, y int) colors.TinyColor {
    idx := i.pixOffset(x, y)
    s := i.Pix[idx : idx+bytesPerPixel : idx+bytesPerPixel]
    return colors.TinyColor{s[0], s[1]}
}

func (i *TinyImage) Set(x, y int, c color.Color) {
    idx := i.pixOffset(x, y)
    c1 := colors.TinyModel.Convert(c).(colors.TinyColor)
    s := i.Pix[idx : idx+bytesPerPixel : idx+bytesPerPixel]
    s[0] = c1.HB
    s[1] = c1.LB
}

func (i *TinyImage) SetTinyColor(x, y int, c colors.TinyColor) {
    idx := i.pixOffset(x, y)
    s := i.Pix[idx : idx+bytesPerPixel : idx+bytesPerPixel]
    s[0] = c.HB
    s[1] = c.LB
}

func (i *TinyImage) pixOffset(x, y int) int {
    return (y - i.Rect.Min.Y) * i.Stride + (x - i.Rect.Min.X) * bytesPerPixel
}
