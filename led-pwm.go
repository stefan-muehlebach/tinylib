//go:build !arduino_mega2560 && !pico2

package tinylib

import (
	"machine"
	// "math"
)

// Mit diesem Typ kann eine dimmbare und gamma-korrigierte LED dargestellt
// werden.
type LED struct {
	Pin machine.Pin
	Pwm *machine.TCC
	ch  uint8
	val uint8
	// gamma []uint8
}

type LEDConfig struct {
    // still empty
}

// Erstellt ein neues LED-Objekt, welches ueber Pin pin angesprochen wird und
// und fuer das Dimment wird pwm verwendet. Die Verwendung von bestimmten
// pwm-Objekten (bspwl. machine.TCC0) mit bestimmten Pins ist nicht frei
// waehlbar! Siehe Dokumentation von tinygo betr. den moeglichen Paarungen.
// Eine zukuenftige Version sollte dies verbergen.
func (l *LED) Configure(cfg LEDConfig) {
	var err error

	l.Pin.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	l.Pwm.Configure(machine.PWMConfig{})
	if l.ch, err = l.Pwm.Channel(l.Pin); err != nil {
		println(err.Error())
	}
}

// Setzt die Gamma-Korrektur der LED auf gamma.
//func (l *LED) SetGamma(gamma float64) {
//	for i := range 256 {
//		t := float64(i) / 255.0
//		l.gamma[i] = uint8(255.0 * math.Pow(t, gamma))
//	}
//}

// Retourniert den aktuellen Wert der LED.
func (l *LED) Get() uint8 {
	return l.val
}

// Setzt den Wert der LED auf val und stellt ihn auch gleich dar.
func (l *LED) Set(val uint8) {
	l.val = val
    l.show()
}

func (l *LED) Toggle() {
    if l.Get() > 128 {
        l.Set(0)
    } else {
        l.Set(255)
    }
}

type FadeDir int
const (
    In FadeDir = iota
    Out
)

func (l *LED) Fade(dir FadeDir) {
    switch dir {
    case In:
        if l.val < 255 {
            l.Set(l.val + 1)
        }
    case Out:
        if l.val > 0 {
            l.Set(l.val + 1)
        }
    }
}

func (l *LED) show() {
	l.Pwm.Set(l.ch, l.Pwm.Top()*uint32(l.val)/255)
	// l.pwm.Set(l.ch, l.pwm.Top()*uint32(l.gamma[l.val])/255)
}
