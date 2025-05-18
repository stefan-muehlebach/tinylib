//go:build pico2

package conf

import (
	"machine"
)

const (
	PinRX = machine.GP1
	PinTX = machine.GP0

	PinSDA = machine.GP8
	PinSCL = machine.GP9

	PinSCK  = machine.GP18
	PinMOSI = machine.GP19
	PinMISO = machine.GP16

	PinReset = machine.GP13
	PinPwrEn = machine.GP14
	PinInt   = machine.GP15
	PinCS    = machine.GP17

	PinBtnGrp = machine.ADC0

	PinTFTSCK    = machine.GP10
	PinTFTMOSI   = machine.GP11
	PinTFTMISO   = machine.GP12
	PinTFTCS     = machine.GP5
	PinTFTDatCmd = machine.GP3

	PinRotA = machine.GP21
	PinRotB = machine.GP22

	PinIRRecv = machine.GP20
)
