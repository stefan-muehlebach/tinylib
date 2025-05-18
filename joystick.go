package tinylib

import (
	"machine"
	"math"
	"time"
)

const (
	adcResolution = 12
	adcShift      = (16 - adcResolution)
	adcMaxValue   = (1 << adcResolution) - 1
	refValue      = (1 << (adcResolution - 1))
	joySampleRate = 30 * time.Millisecond
)

type Joystick struct {
	XPin, YPin, BtnPin             machine.Pin
	xPoti, yPoti                   machine.ADC
	xPotiVal, yPotiVal             uint16
	xPotiDiff, yPotiDiff           int16
	xVal, yVal, xDiffVal, yDiffVal float32
	avgValRatio, avgDiffRatio      float32
	valFact, diffValFact           float32
	reverseX, reverseY             bool
}

type JoyConfig struct {
	ReverseX, ReverseY bool
	AverageRatio       float32
}

func (j *Joystick) Configure(cfg JoyConfig) {

	j.reverseX = cfg.ReverseX
	j.reverseY = cfg.ReverseY
	j.avgValRatio = cfg.AverageRatio

	j.avgDiffRatio = 0.9
	j.valFact = 1.0 / adcMaxValue
	j.diffValFact = 10.0 * j.valFact

	j.xPoti = machine.ADC{
		Pin: j.XPin,
	}
	j.xPoti.Configure(machine.ADCConfig{
		Resolution: adcResolution,
	})
	j.yPoti = machine.ADC{
		Pin: j.YPin,
	}
	j.yPoti.Configure(machine.ADCConfig{
		Resolution: adcResolution,
	})

	if j.BtnPin != machine.NoPin {
		j.BtnPin.Configure(machine.PinConfig{
			Mode: machine.PinInputPullup,
		})
	}
}

// Wenn der Joystick mehr soll als nur unmittelbare Messwerte zurueckliefern,
// dann muss die Datenerfassung auf reglmaessiger Basis erfolgen, damit ab-
// haengige Groessen laufend aktualisiert werden können. Mit Sample() wird
// die Datenerfassung und Aufbereitung durchgefuehrt.
func (j *Joystick) Sample() {
	var xPotiVal, yPotiVal uint16

	if j.reverseX {
		xPotiVal = math.MaxUint16 - j.xPoti.Get()
	} else {
		xPotiVal = j.xPoti.Get()
	}
	xPotiVal >>= adcShift

	if j.reverseY {
		yPotiVal = math.MaxUint16 - j.yPoti.Get()
	} else {
		yPotiVal = j.yPoti.Get()
	}
	yPotiVal >>= adcShift

	j.xPotiDiff = int16(xPotiVal) - int16(j.xPotiVal)
	j.yPotiDiff = int16(yPotiVal) - int16(j.yPotiVal)

	j.xPotiVal, j.yPotiVal = xPotiVal, yPotiVal

	j.xVal = j.avgValRatio*j.xVal + (1.0-j.avgValRatio)*float32(xPotiVal)*j.valFact
	j.yVal = j.avgValRatio*j.yVal + (1.0-j.avgValRatio)*float32(yPotiVal)*j.valFact

	j.xDiffVal = j.avgDiffRatio*j.xDiffVal + (1.0-j.avgDiffRatio)*float32(j.xPotiDiff)*j.diffValFact
	j.yDiffVal = j.avgDiffRatio*j.yDiffVal + (1.0-j.avgDiffRatio)*float32(j.yPotiDiff)*j.diffValFact
}

// Liefert die unveraenderten Werte der A/D-Wandler zurueck, an welche die
// beiden Potentiometer des Joysticks angeschlossen sind.
func (j *Joystick) RawValues() (uint16, uint16) {
	return j.xPotiVal, j.yPotiVal
}

// Liefert die Positionsdaten des Joysticks mit allfaelligen Anpassungen
// zurueck.
func (j *Joystick) Values() (x, y float32) {
	return j.xVal, j.yVal
}

func (j *Joystick) DiffValues() (x, y float32) {
	return j.xDiffVal, j.yDiffVal
}

func (j *Joystick) Sw() bool {
	return !j.BtnPin.Get()
}

// Retourniert einen Task, welcher bei einem Dispatcher hinterlegt werden
// kann.
func (j *Joystick) Task() *Task {
	return NewTask(j.Sample, TaskConfig{Interval: joySampleRate})
}
