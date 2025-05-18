//go:build ignore

package tinylib

import (
	"machine"
	"math"
	"time"
)

const (
	adcResolution  = 12
	adcShift       = (16 - adcResolution)
	adcMaxValue    = (1 << adcResolution) - 1
	refValue       = (1 << (adcResolution - 1))
	defAveragePart = 50
	joySampleRate  = 10 * time.Millisecond
)

type OutValueType interface {
	int | int16 | int32 | uint16 | uint32 | float32 | float64
}

type Joystick[T OutValueType] struct {
	XPin, YPin, BtnPin     machine.Pin
	xPoti, yPoti           machine.ADC
	xRaw, yRaw             uint16
	xAvg, yAvg, avgRatio   float32
	xVal, yVal             T
	reverseX, reverseY     bool
	xMin, xMax, yMin, yMax T
	xRatio, yRatio         float32
}

type JoyConfig[T OutValueType] struct {
	ReverseX, ReverseY     bool
	XMin, XMax, YMin, YMax T
	AverageRatio         float32
}

func (j *Joystick[T]) Configure(cfg JoyConfig[T]) {

	if cfg.XMin == 0 && cfg.XMax == 0 {
		cfg.XMax = adcMaxValue
	}
	if cfg.YMin == 0 && cfg.YMax == 0 {
		cfg.YMax = adcMaxValue
	}

	j.xMin, j.xMax, j.yMin, j.yMax = cfg.XMin, cfg.XMax, cfg.YMin, cfg.YMax
	j.reverseX = cfg.ReverseX
	j.reverseY = cfg.ReverseY

	j.xRatio = float32(j.xMax-j.xMin) / adcMaxValue
	j.yRatio = float32(j.yMax-j.yMin) / adcMaxValue

	j.avgRatio = cfg.AverageRatio

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
func (j *Joystick[T]) Sample() {
	j.xRaw = j.xPoti.Get()
	if j.reverseX {
		j.xRaw = math.MaxUint16 - j.xRaw
	}
	j.xRaw >>= adcShift

	j.yRaw = j.yPoti.Get()
	if j.reverseY {
		j.yRaw = math.MaxUint16 - j.yRaw
	}
	j.yRaw >>= adcShift

	j.xAvg = j.avgRatio*j.xAvg + (1.0-j.avgRatio)*float32(j.xRaw)
	j.yAvg = j.avgRatio*j.yAvg + (1.0-j.avgRatio)*float32(j.yRaw)

	j.xVal = j.xMin + T(float32(j.xAvg)*j.xRatio)
	j.yVal = j.yMin + T(float32(j.yAvg)*j.yRatio)
}

// Liefert die unveraenderten Werte der A/D-Wandler zurueck, an welche die
// beiden Potentiometer des Joysticks angeschlossen sind.
func (j *Joystick[T]) RawValues() (uint16, uint16) {
	return j.xRaw, j.yRaw
}

// Liefert die Positionsdaten des Joysticks mit allfaelligen Anpassungen
// zurueck.
func (j *Joystick[T]) Values() (T, T) {
	return j.xVal, j.yVal
}

func (j *Joystick[T]) Sw() bool {
	return !j.BtnPin.Get()
}

// Retourniert einen Task, welcher bei einem Dispatcher hinterlegt werden
// kann.
func (j *Joystick[T]) Task() *Task {
	return NewTask(j.Sample, TaskConfig{Interval: joySampleRate})
}
