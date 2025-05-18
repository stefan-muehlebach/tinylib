//go:build arduino_mkrwifi1010

package conf

import (
	"machine"
)

const (
	pinRX = machine.D13
	pinTX = machine.D14

	pinSDA = machine.D11
	pinSCL = machine.D12

	pinSCK  = machine.D9
	pinMOSI = machine.D8
	pinMISO = machine.D10

	pinReset = machine.D3
	pinPwrEn = machine.D4
	pinInt   = machine.D5
	pinCS    = machine.D7
)
