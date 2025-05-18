//go:build !atsamd51 && !atsame5x && !atsamd21

package ili9341

import (
	"machine"

	"tinygo.org/x/drivers"
)

var buf [128]byte

type spiDriver struct {
	bus drivers.SPI
}

func NewSPI(bus drivers.SPI, dc, cs, rst machine.Pin) *Device {
	return &Device{
		dc:  dc,
		cs:  cs,
		rst: rst,
		rd:  machine.NoPin,
		driver: &spiDriver{
			bus: bus,
		},
	}
}

func (pd *spiDriver) configure(config *Config) {
}

func (pd *spiDriver) write8(b byte) {
	pd.bus.Transfer(b)
}

func (pd *spiDriver) write8n(b byte, n int) {
	for range n {
		pd.bus.Transfer(b)
	}
}

func (pd *spiDriver) write8sl(b []byte) {
    for _, ch := range b {
        pd.bus.Transfer(ch)
    }
}

func (pd *spiDriver) write16(data uint16) {
	buf[0] = uint8(data >> 8)
	buf[1] = uint8(data)
	pd.bus.Tx(buf[:2], nil)
}

func (pd *spiDriver) write16n(data uint16, n int) {
    buf[0] = uint8(data >> 8)
    buf[1] = uint8(data)
    for range n {
        pd.bus.Tx(buf[:2], nil)
    }
}

// func (pd *spiDriver) write16n(data uint16, n int) {
// 	for i := 0; i < len(buf); i += 2 {
// 		buf[i] = uint8(data >> 8)
// 		buf[i+1] = uint8(data)
// 	}

// 	for i := 0; i < (n >> 7); i++ {
// 		pd.bus.Tx(buf[:], nil)
// 	}

// 	pd.bus.Tx(buf[:n%128], nil)
// }

func (pd *spiDriver) write16sl(data []uint16) {
	for i, c := 0, len(data); i < c; i++ {
		buf[0] = uint8(data[i] >> 8)
		buf[1] = uint8(data[i])
		pd.bus.Tx(buf[:2], nil)
	}
}

func (pd *spiDriver) read8() byte {
    pd.bus.Tx(nil, buf[:1])
    return buf[0]
}

func (pd *spiDriver) read8sl(b []byte) {
    pd.bus.Tx(nil, b)
}
