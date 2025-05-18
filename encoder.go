//go:build !arduino_mega2560

//go:generate pioasm -o go encoder.pio encoder_pio.go

package tinylib

import (
	"machine"
	"time"
)

type Direction int

const (
	CCW Direction = iota
	CW
	Unspecified
)

func (d Direction) String() string {
	switch d {
	case CCW:
		return "CCW"
	case CW:
		return "CW"
	default:
		return "Unspecified"
	}
}

const (
	// defEncoderChangeDelay = 40 * time.Millisecond
	// defEncoderRepeatDelay = 30 * time.Millisecond
	defEncoderPollRate = 30 * time.Millisecond
)

type RotationCallback func(dir Direction, steps int)

type EncoderConfig struct {
	PollRate time.Duration
}

// Mit diesem Typ kann ein inkrementeller Rotations-Encoder einfach ausgelesen
// werden.
type Encoder struct {
	PinA, PinB            machine.Pin
	pollRate              time.Duration
	rotateCB              RotationCallback
	position, oldPosition int
	state                 byte
}

func (e *Encoder) Configure(conf EncoderConfig) {
	if conf.PollRate == 0 {
		conf.PollRate = defEncoderPollRate
	}
	e.pollRate = conf.PollRate
	e.PinA.Configure(machine.PinConfig{Mode: machine.PinInput})
	e.PinB.Configure(machine.PinConfig{Mode: machine.PinInput})
	e.PinA.SetInterrupt(machine.PinToggle, e.newIsr)
	e.PinB.SetInterrupt(machine.PinToggle, e.newIsr)
}

func (e *Encoder) SetOnRotate(cb RotationCallback) {
	e.rotateCB = cb
}

func (e *Encoder) Task() *Task {
	return NewTask(e.Tick, TaskConfig{Interval: e.pollRate})
}

func (e *Encoder) Tick() {
	pos := e.position
	if pos == e.oldPosition {
		return
	}
	diff := (pos - e.oldPosition) / 4
	if diff == 0 {
		return
	}
	if e.rotateCB != nil {
		if diff > 0 {
			e.rotateCB(CCW, diff)
		} else {
			e.rotateCB(CW, -diff)
		}
	}
	e.oldPosition = pos
}

func (e *Encoder) newIsr(pin machine.Pin) {
	s := e.state & 0x03
	if e.PinA.Get() {
		s |= 0x04
	}
	if e.PinB.Get() {
		s |= 0x08
	}
	switch s {
	case 1, 7, 8, 14:
		e.position++
	case 2, 4, 11, 13:
		e.position--
	case 3, 12:
		e.position += 2
	case 6, 9:
		e.position -= 2
	}
	e.state = (s >> 2)
	// println("pos: ", e.position)
}
