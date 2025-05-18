package tinylib

import (
	"time"
)

//----------------------------------------------------------------------------

const (
	// Ab dieser Dauer werden Hold-Events erzeugt. Oder anders formuliert: ein
	// Druecken und Loslassen innerhalb dieser Zeit erzeugt ein Pressed Event.
	defHoldThreshold = 500 * time.Millisecond
	// In diesem Intervall werden die Hold-Events erzeugt.
	defHoldCallRate = 100 * time.Millisecond
)

// Enthaelt alle wichtigen Konfigurationseinstellungen zu einem Push-Button
type ButtonConfig struct {
	// Ab dieser Dauer werden Hold-Events erzeugt. Oder anders formuliert: ein
	// Druecken und Loslassen innerhalb dieser Zeit erzeugt ein Pressed Event.
	HoldThreshold time.Duration
	// In diesem zeitlichen Abstand wird bei konstantem Druecken des Buttons
	// der Hold-Callback aufgerufen.
	HoldCallRate time.Duration
}

// Funktionstyp der Callback-Handler fuer die Events Pressed, Push und Release.
type ButtonCallback func()

// Funktionstyp des Callback-Handlers fuer das Event Hold.
type ButtonHoldCallback func(firstCall bool)

// Mit diesem Typ koennen Push-Buttons auf vielfaeltige Weise angesteuert
// werden. Wichtig: die Entprellung des Buttons muss hardwareseitig erfolgen!
// Es gibt 4 Events, fuer welche entsprechende Callbacks hinterlegt werden
// koennen.
//
//	Push   : Druecken des Buttons (Aufruf: 1-mal)
//	Release: Loslassen des Buttons (Aufruf: 1-mal)
//	Pressed: Druecken und Loslassen innerhalb einer bestimmten Zeit
//	         (Aufruf: 1-mal)
//	Hold   : Druecken und Halten ueber einen bestimmten Zeitraum hinweg
//	         (Aufruf: alle btnPollRate Millisekunden so lange der Button
//	         gedrueckt bleibt).
type Button struct {
	holdThreshold, holdCallRate  time.Duration
	pushTime, lastHoldCall       time.Time
	isHolding                    bool
	pressedCB, pushCB, releaseCB ButtonCallback
	holdCB                       ButtonHoldCallback
}

// Erzeugt ein neues Button-Objekt und verwendet pin als Input. Der Button
// muss Active High konfiguriert werden.
func (b *Button) Configure(cfg ButtonConfig) {
	if cfg.HoldThreshold == 0 {
		cfg.HoldThreshold = defHoldThreshold
	}
	if cfg.HoldCallRate == 0 {
		cfg.HoldCallRate = defHoldCallRate
	}
	b.holdThreshold = cfg.HoldThreshold
	b.holdCallRate = cfg.HoldCallRate
}

// Setzt cb als Callback-Handler fuer das Pressed-Event.
func (b *Button) SetOnPressed(cb ButtonCallback) {
	b.pressedCB = cb
}

// Setzt cb als Callback-Handler fuer das Push-Event.
func (b *Button) SetOnPush(cb ButtonCallback) {
	b.pushCB = cb
}

// Setzt cb als Callback-Handler fuer das Release-Event.
func (b *Button) SetOnRelease(cb ButtonCallback) {
	b.releaseCB = cb
}

// Setzt cb als Callback-Handler fuer das Hold-Event. Mit dem Bool-Argument
// wird dem Handler mitgeteilt, ob der Aufruf der erste einer Serie ist.
func (b *Button) SetOnHold(cb ButtonHoldCallback) {
	b.holdCB = cb
}

// In dieser Methode steckt die ganze Logik hinter dem Button. Diese Methode
// kann entweder direkt alle btnPollRate Millisekunden aufgerufen werden
// oder durch einen Task (siehe Methode Task()).
func (b *Button) Process(buttonDown bool) {
	now := time.Now().Truncate(time.Millisecond)
	if buttonDown {
		if b.pushTime.IsZero() {
			b.pushTime = now
			if b.pushCB != nil {
				b.pushCB()
			}
		} else {
			if !b.isHolding && now.Sub(b.pushTime) >= b.holdThreshold {
				if b.holdCB != nil {
					b.holdCB(true)
				}
				b.isHolding = true
				b.lastHoldCall = now
			}
			if b.isHolding && now.Sub(b.lastHoldCall) >= b.holdCallRate {
				if b.holdCB != nil {
					b.holdCB(false)
				}
				b.lastHoldCall = now
			}
		}
	} else {
		if !b.pushTime.IsZero() {
			if b.releaseCB != nil {
				b.releaseCB()
			}
			if now.Sub(b.pushTime) < b.holdThreshold {
				if b.pressedCB != nil {
					b.pressedCB()
				}
			}
			b.pushTime = time.Time{}
			b.lastHoldCall = time.Time{}
			b.isHolding = false
		}
	}
}
