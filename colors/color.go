package colors

import (
	"image/color"
)

//----------------------------------------------------------------------------

type Color struct {
	color.NRGBA
}

func NewColor(r, g, b uint8) Color {
	c := Color{}
	c.R = r
	c.G = g
	c.B = b
    c.A = 0xFF
	return c
}

func (c Color) FadeOut(t float32) Color {
	u := 1.0 - t
	return NewColor(uint8(u*float32(c.R)), uint8(u*float32(c.G)),
		uint8(u*float32(c.B)))
}

func (c Color) Interpolate(c2 Color, t float32) Color {
	u := 1.0 - t
	return NewColor(uint8(u*float32(c.R) + t*float32(c2.R)),
		uint8(u*float32(c.G) + t*float32(c2.G)),
		uint8(u*float32(c.B) + t*float32(c2.B)))
}

//----------------------------------------------------------------------------

type TinyColor struct {
	HB, LB uint8
}

func NewTinyColor(r, g, b uint8) TinyColor {
	hb := (r & 0xF8) | ((g & 0xFC) >> 5)
	lb := ((g & 0xFC) << 3) | ((b & 0xF8) >> 3)
	return TinyColor{hb, lb}
}

func NewTinyHexColor(hex uint32) TinyColor {
	r := uint8((hex >> 16) & 0xff)
	g := uint8((hex >> 8) & 0xff)
	b := uint8((hex >> 0) & 0xff)
	return NewTinyColor(r, g, b)
}

func (c TinyColor) BitsPerPixel() int {
	return 24
}

func (c TinyColor) RGBA() (r, g, b, a uint32) {
	r = uint32(c.HB & 0xF8)
	r |= r << 8
	g = uint32(((c.HB & 0x07) << 5) | ((c.LB & 0xE0) >> 3))
	g |= g << 8
	b = uint32((c.LB << 3) & 0xF8)
	b |= b << 8
	a = 0xffff
	return
}

func tinyModel(c color.Color) color.Color {
	if _, ok := c.(TinyColor); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	if a == 0xffff {
		r = (r >> 8)
		g = (g >> 8)
		b = (b >> 8)
		return NewTinyColor(uint8(r), uint8(g), uint8(b))
	}
	if a == 0x0000 {
		return NewTinyColor(0, 0, 0)
	}
	r = (r * 0xffff) / a
	r = (r >> 8)
	g = (g * 0xffff) / a
	g = (g >> 8)
	b = (b * 0xffff) / a
	b = (b >> 8)
	return NewTinyColor(uint8(r), uint8(g), uint8(b))
}

var (
	TinyModel color.Model = color.ModelFunc(tinyModel)
)
