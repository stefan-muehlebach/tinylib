package tinylib

import (
	"machine"
	"time"
)

const (
	defValueEpsilon    = 30
	defNumCalibSamples = 200
	defButtonsPollRate = 10 * time.Millisecond
	defADCResolution   = 12
)

var (
	defADCRightShift uint32
	defADCMaxValue   uint16
)

// Enthaelt alle wichtigen Konfigurationseinstellungen zu einer Reihe von
// Buttons, welche über einen einzigen Analog-Pin gefuehrt werden. Die Wahl
// der für jedem Button eigenen Widerstandswerte ist so zu wählen, dass die
// Zuordnung von analogem Messet zu geschlossenem Button eindeutig ist.
type ButtonGroupConfig struct {
	// Intervall, in welchem der Zustand des Buttons abgefragt werden soll.
	PollRate   time.Duration
	Resolution uint32
}

// Mit diesem Typ wird ein Button mit dem Intervall eines analogen Signals
// verbunden, das fuer diesen Button gemessen werden kann.
type AnalogButtonReadout struct {
	MeanValue, LowerBound, UpperBound uint16
	Button                            *Button
}

// Mit diesem Typ koennen Push-Buttons auf vielfaeltige Weise angesteuert
// werden. Wichtig: die Entprellung des Buttons muss hardwareseitig erfolgen!
// Es gibt 4 Events, fuer welche entsprechende Callbacks hinterlegt werden
// koennen.
//
//	Push   : Druecken des Buttons (Aufruf: 1-mal)
//	Hold   : Druecken und Halten ueber einen bestimmten Zeitraum hinweg
//	         (Aufruf: alle btnPollRate Millisekunden so lange der Button
//	         gedrueckt bleibt).
//	Release: Loslassen des Buttons (Aufruf: 1-mal)
//	Pressed: Druecken und Loslassen innerhalb einer bestimmten Zeit
//	         (Aufruf: 1-mal)
type ButtonGroup struct {
	Pin        machine.Pin
	adc        machine.ADC
	pollRate   time.Duration
	lastId     int
	buttonList []*AnalogButtonReadout
	tickFunc   func()

    // Das sind die Variablen, welche waehrend der Kalibrierung der Buttons
    // verwendet werden. Stellt sich die Frage, ob und wie man die in einen
    // weiteren, privaten Record auslagern koennte.
	buttonToCalibrate    int
	sumValues, numValues int
	minValue, maxValue   uint16
	collectingData       bool
}

func (b *ButtonGroup) Configure(cfg ButtonGroupConfig) {
	if cfg.PollRate == 0 {
		cfg.PollRate = defButtonsPollRate
	}
	if cfg.Resolution == 0 {
		cfg.Resolution = defADCResolution
	}
	defADCRightShift = 16 - cfg.Resolution
	defADCMaxValue = (1 << cfg.Resolution) - 1
	b.Pin.Configure(machine.PinConfig{
		Mode: machine.PinInput,
	})
	b.adc = machine.ADC{
		Pin: b.Pin,
	}
	b.adc.Configure(machine.ADCConfig{
		Resolution: cfg.Resolution,
	})
	b.pollRate = cfg.PollRate
	b.lastId = -1
	b.buttonList = make([]*AnalogButtonReadout, 0)
	b.tickFunc = b.workTick
}

// Damit wird ein Button der Gruppen hinzugefuegt. Mit id wird die Nummer
// des Buttons innerhalb der Gruppe angegeben, wobei die Buttons lueckenlos
// von 0 bis numButtons-1 durchnumeriert werden. val ist der Messwert des
// A/D-Wandlers, der beim Druecken dieses Buttons erwartet wird und btn
// schliesslich ein Pointer auf den Button.
func (b *ButtonGroup) AddButton(btn *Button, val uint16) {
	br := &AnalogButtonReadout{
		MeanValue:  val,
		LowerBound: val - defValueEpsilon,
		UpperBound: val + defValueEpsilon,
		Button:     btn,
	}
	b.buttonList = append(b.buttonList, br)
}

func (b *ButtonGroup) StartCalibration() {
	println("Starting calibration of analog buttons")
	b.calibrate(0)
}

// Startet die Kalibrierung der Buttons, wobei mit dem Button gestartet wird,
// dessen id als Parameter uebergben werden kann. Normalerweise startet man
// mit id=0 und durchlaeuft alle Buttons der Reihe nach.
func (b *ButtonGroup) calibrate(id int) {
	if id >= len(b.buttonList) {
		b.tickFunc = b.workTick
		return
	}
	println(">> press button", id, "and hold it!")
	b.buttonToCalibrate = id
	b.sumValues = 0
	b.minValue = 0xFFFF
	b.maxValue = 0x0000
	b.numValues = 0
	b.collectingData = true
	b.tickFunc = b.calibTick
}

// Retourniert einen Task, welcher bei einem Dispatcher hinterlegt werden
// kann und alle btnPollRate Millisekunden aufgerufen werden muss (sollte).
func (b *ButtonGroup) Task() *Task {
	return NewTask(b.Tick, TaskConfig{Interval: b.pollRate})
}

func (b *ButtonGroup) Tick() {
	b.tickFunc()
}

func (b *ButtonGroup) calibTick() {
	val := b.adc.Get() >> defADCRightShift
	if !b.collectingData {
		if val == defADCMaxValue {
			println("  collected data for button", b.buttonToCalibrate)
			avg := uint16(b.sumValues / b.numValues)
			println("    avg:", avg)
			println("    min:", b.minValue)
			println("    max:", b.maxValue)
			println("  data has been updated")
			b.buttonList[b.buttonToCalibrate].MeanValue = avg
			b.buttonList[b.buttonToCalibrate].LowerBound = avg - defValueEpsilon
			b.buttonList[b.buttonToCalibrate].UpperBound = avg + defValueEpsilon
			b.calibrate(b.buttonToCalibrate + 1)
		}
		return
	}
	if val == defADCMaxValue {
		return
	}
	b.sumValues += int(val)
	if val > b.maxValue {
		b.maxValue = val
	}
	if val < b.minValue {
		b.minValue = val
	}
	b.numValues += 1
	if b.numValues == defNumCalibSamples {
		println(">> got enough data, release button", b.buttonToCalibrate)
		b.collectingData = false
	}
}

// In dieser Methode steckt die ganze Logik hinter dem Button. Diese Methode
// kann entweder direkt alle btnPollRate Millisekunden aufgerufen werden
// oder durch einen Task (siehe Methode Task()).
func (b *ButtonGroup) workTick() {
	val := b.adc.Get() >> defADCRightShift
	for _, buttonInfo := range b.buttonList {
		if buttonInfo.LowerBound <= val && val < buttonInfo.UpperBound {
			buttonInfo.Button.Process(true)
		} else {
			buttonInfo.Button.Process(false)
		}
	}
}
