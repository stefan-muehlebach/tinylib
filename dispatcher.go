//go:build inline

package tinylib

import (
	"time"
)

// Dieses File enthaelt alle Typen und Konstanten, um ein einfaches
// Dispatching unter TinyGo auf einem Single-Core Microcontroller zu
// realisieren. Im Wesentlichen gibt es zwei Typen:
//
// Dispatcher: ist fuer die zeitlich korrekte Ausfuehrung der ihm
//
//	zugewiesenen Tasks verantwortlich. Es darf nur eine (1)
//	Instanz dieses Typs geben!
//
// Task:       fuer jede periodisch ausfuehrbare Funktion oder Methode wird
//
//	eine Instanz dieses Typs erstellt, welche alle relevanten
//	Daten (gewuenschtes Intervall, naechster Ausfuehrungszeitpunkt,
//	etc.) und statistische Daten enthaelt.
//
// Eine typische Anwendung sieht folgendermassen aus (das Beispiel zeigt, wie
// eine LED im Sekundentakt ein- und ausgeschaltet werden kann):
//
//	type BlinkyLED struct {
//	    pin machine.Pin
//	    isOn bool
//	}
//
//	func NewBlinkyLED(pin machine.Pin) *BlinkyLED {
//	    led := &BlinkyLED{}
//	    led.pin = pin
//	    led.pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
//	    return led
//	}
//
//	func (led *BlinkyLED) Toggle() {
//	    led.isOn = !led.isOn
//	    if led.isOn {
//	        led.pin.High()
//	    } else {
//	        led.pin.Low()
//	    }
//	}
//
// led := NewBlinkyLED(machine.LED)
// ledTask := tinylib.NewTask(led.Toggle,
//
//	tinylib.TaskConfig{Interval: time.Second})
//
// tinylib.Disp.AddTask(ledTask)
//
//	for {
//	    tinylib.Disp.Tick()
//	}
var (
	// Dies ist der Default-Dispatcher, welcher in diesem Package erzeugt wird
	// und genutzt werden MUSS, d.h. es gibt keine Notwendigkeit, im Programm
	// einen eigenen Dispatcher zu erzeugen.
	// TO THINK ABOUT: wirklich die beste Loesung?
	Disp *Dispatcher
)

func init() {
	Disp = NewDispatcher()
}

// Ist eine Alternative zu time.Now(), welche die aktuelle Zeit auf
// Millisekunden genau liefert.
func Now() time.Time {
	return time.Now().Truncate(time.Millisecond)
}

func NowMS() int64 {
	return Now().UnixMilli()
}

//----------------------------------------------------------------------------

// Jede Funktion/Methode dieses Typs kann in einem Task hinterlegt und
// periodisch durch den Dispatcher aufgerufen werden.
type TaskFunc func()

type TaskConfig struct {
	ExecTime time.Time
	Interval time.Duration
}

type Task struct {
	Func                  TaskFunc
	execTime              time.Time
	interval              time.Duration
	next                  *Task
	isHalting             bool
	lastTerm, term, delay time.Duration
	numCalls              uint32
}

func NewTask(fnc TaskFunc, cfg TaskConfig) *Task {
	t := &Task{}
	t.Func = fnc
	t.execTime = cfg.ExecTime
	t.interval = cfg.Interval
	return t
}

func (t *Task) Configure(cfg TaskConfig) {
	t.execTime = cfg.ExecTime
	t.interval = cfg.Interval
}

func (t *Task) Start(now time.Time) {
	if t.isHalting {
		return
	}
	t.numCalls += 1
	t.delay += now.Sub(t.execTime)
	t0 := time.Now()
	t.Run()
	t.lastTerm = time.Since(t0)
	t.term += t.lastTerm
	t.execTime = now.Add(t.interval)
}

func (t *Task) Run() {
	t.Func()
}

func (t *Task) Halt() {
	t.isHalting = true
}

func (t *Task) Interval() time.Duration {
	return t.interval
}

func (t *Task) SetInterval(i time.Duration) {
	t.interval = i
}

func (t *Task) NumCalls() uint32 {
	return t.numCalls
}

func (t *Task) Term() time.Duration {
	return t.term
}

func (t *Task) AvgTerm() time.Duration {
	return t.term / time.Duration(t.numCalls)
}

func (t *Task) Delay() time.Duration {
	return t.delay
}

func (t *Task) AvgDelay() time.Duration {
	return t.delay / time.Duration(t.numCalls)
}

//----------------------------------------------------------------------------

const (
	loadMeasureLengthMS = 5000
	loadSlotLengthMS    = 200
	loadNumSlots        = loadMeasureLengthMS / loadSlotLengthMS
)

type Dispatcher struct {
	readyList          *Task
	termList           [loadNumSlots]time.Duration
	currSlot, lastSlot uint8
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{}
}

func (d *Dispatcher) AddTask(t *Task) {
	currentTime := Now()

	if t.execTime.Before(currentTime) {
		t.execTime = currentTime.Add(t.interval)
	}
	if t.isHalting {
		t.isHalting = false
	} else {
		d.insert(t)
	}
}

func (d *Dispatcher) Print() {
	ptr := d.readyList
	for i := 0; ptr != nil; i++ {
		println("[", i, "]")
		println("  numCalls  :", ptr.numCalls)
		println("  execTime  :", ptr.execTime.String())
		println("  interval  :", ptr.interval.String())
		println("  isHalting :", ptr.isHalting)
		println("  term      :", ptr.AvgTerm().String())
		println("  delay     :", ptr.AvgDelay().String())
		ptr = ptr.next
	}
}

// Ueber diese Methode wird das gesamte Dispatching gesteuert. Sie sollte
// einer Endlos-Schleife des Hauptprogrammes ohne weitere Funktionen oder
// Methoden aufgerufen werden. Wenn die Applikation konsequent auf die
// Verwendung von Tasks umgestellt wurde, dann sieht die main()-Funktion
// im Wesentlichen wie folgt aus:
//
//	for {
//	    tinylib.Disp.Tick()
//	}
func (d *Dispatcher) Tick() {
	currentTime := Now()

	for {
		task := d.pop(currentTime)
		if task == nil {
			return
		}
		if task.isHalting {
			task.isHalting = false
			continue
		}
		task.Start(currentTime)
		d.currSlot = uint8((NowMS() % loadMeasureLengthMS) / loadSlotLengthMS)
		if d.currSlot != d.lastSlot {
			d.lastSlot = (d.lastSlot + 1) % loadNumSlots
			for d.lastSlot != d.currSlot {
				d.termList[d.lastSlot] = time.Duration(0)
				d.lastSlot = (d.lastSlot + 1) % loadNumSlots
			}
			d.termList[d.currSlot] = task.lastTerm
		} else {
			d.termList[d.currSlot] += task.lastTerm
		}
		if task.interval > 0 {
			d.insert(task)
		}
	}
}

// Liefert die Anzahl der im Dispatcher registrierten Tasks. Da diese Methode
// nicht synchronisiert ist, kann es zu geringfuegigen Abweichungen kommen.
func (d *Dispatcher) NumTasks() int {
	n := 0
	ptr := d.readyList
	for ptr != nil {
		n++
		ptr = ptr.next
	}
	return n
}

// Liefert die Systemlast in Prozent. Fuer die Berechnung wird das Verhaeltnis
// der Task-Laufzeiten zum definierten Zeitfenster berechnet.
func (d *Dispatcher) Load() uint8 {
	sumDur := time.Duration(0)
	for _, dur := range d.termList {
		sumDur += dur
	}
	return uint8(100.0 * float64(sumDur) / float64(loadMeasureLengthMS * time.Millisecond))
}

// Diese Methode dient dazu, den naechsten zur Ausfuehrung bereiten Task zu
// ermitteln. Liefert nil, falls aktuell kein Task zur Ausfuehrung bereit
// steht.
func (d *Dispatcher) pop(t time.Time) *Task {
	var ptr *Task

	if d.readyList == nil {
		return nil
	}
	ptr = d.readyList
	if ptr.execTime.After(t) {
		return nil
	}
	d.readyList = ptr.next
	ptr.next = nil
	return ptr
}

// Stellt den Task sortiert in die Taskliste. Das Feld execTime des Tasks
// muss vorgaengig korrekt ausgefuellt worden sein, diese Methode greift
// auf dieses Feld nur lesend zu.
func (d *Dispatcher) insert(task *Task) {
	ptr := d.readyList

	if ptr == nil {
		d.readyList = task
	} else if task.execTime.Before(ptr.execTime) {
		task.next = ptr
		d.readyList = task
	} else {
		for ptr.next != nil {
			if task.execTime.Before(ptr.next.execTime) {
				break
			}
			ptr = ptr.next
		}
		task.next = ptr.next
		ptr.next = task
	}
}
