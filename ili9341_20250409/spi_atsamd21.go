//go:build atsamd21

package ili9341

import (
	"device/sam"
	"machine"
)

type spiDriver struct {
	bus machine.SPI
}

func NewSPI(bus machine.SPI, dc, cs, rst machine.Pin) *Device {
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
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(b))

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write8n(b byte, n int) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for i, c := 0, n; i < c; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(b))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write8sl(b []byte) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for i, c := 0, len(b); i < c; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(b[i]))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write16(data uint16) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data >> 8)))
	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data)))

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write16n(data uint16, n int) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for i := 0; i < n; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data >> 8)))
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data)))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write16sl(data []uint16) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for i, c := 0, len(data); i < c; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data[i] >> 8)))
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data[i])))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write24(data uint32) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data >> 16)))
	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data >> 8)))
	for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
	}
	pd.bus.Bus.DATA.Set(uint32(uint8(data)))

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

func (pd *spiDriver) write24n(data uint32, n int) {
	pd.bus.Bus.CTRLB.ClearBits(sam.SERCOM_SPI_CTRLB_RXEN)

	for i := 0; i < n; i++ {
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data >> 16)))
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data >> 8)))
		for !pd.bus.Bus.INTFLAG.HasBits(sam.SERCOM_SPI_INTFLAG_DRE) {
		}
		pd.bus.Bus.DATA.Set(uint32(uint8(data)))
	}

	pd.bus.Bus.CTRLB.SetBits(sam.SERCOM_SPI_CTRLB_RXEN)
	for pd.bus.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPI_SYNCBUSY_CTRLB) {
	}
}

var (
    buf [16]byte
)

func (pd *spiDriver) read8() byte {
    pd.bus.Tx(nil, buf[:1])
    pd.bus.Tx(nil, buf[:1])
    return buf[0]
}

func (pd *spiDriver) read8sl(b []byte) {
    pd.bus.Tx(nil, buf[:2])
    pd.bus.Tx(nil, b)
}
