package tinylib

import (
	"machine"
	"time"
)

//----------------------------------------------------------------------------

const (
	// Damit wird das Intervall definiert, in welchem die Methode Tick()
	// aufgerufen werden muss.
	defButtonPollRate = 10 * time.Millisecond
)

type ButtonSoloConfig struct {
    ButtonConfig
	// Intervall, in welchem der Zustand des Buttons abgefragt werden soll.
	PollRate time.Duration
	// Ist der Anschluss 'active high' oder 'active low' angesteuert
	// (Default: 'active high')
	ActiveLow bool
}

type ButtonSolo struct {
	Button
	Pin        machine.Pin
	pollRate   time.Duration
	activeHigh bool
}

func (b *ButtonSolo) Configure(cfg ButtonSoloConfig) {
	var mode machine.PinMode

	if cfg.PollRate == 0 {
		cfg.PollRate = defButtonPollRate
	}
	b.pollRate = cfg.PollRate

	if cfg.ActiveLow {
		mode = machine.PinInputPullup
	} else {
		mode = machine.PinInput
	}
	b.Pin.Configure(machine.PinConfig{
		Mode: mode,
	})
	b.activeHigh = !cfg.ActiveLow

    b.Button.Configure(cfg.ButtonConfig)
}

// Retourniert einen Task, welcher bei einem Dispatcher hinterlegt werden
// kann und alle btnPollRate Millisekunden aufgerufen werden muss (sollte).
func (b *ButtonSolo) Task() *Task {
	return NewTask(b.Tick, TaskConfig{Interval: b.pollRate})
}

// Mit dieser Methode kann man einen ButtonSolo periodisch als Task durch den
// Dispatcher aufrufen lassen.
func (b *ButtonSolo) Tick() {
    if (b.activeHigh && b.Pin.Get()) || (!b.activeHigh && !b.Pin.Get()) {
        b.Process(true)
    } else {
        b.Process(false)
    }
}

//----------------------------------------------------------------------------
