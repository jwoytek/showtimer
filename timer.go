package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"
)

const (
	CountUp   = iota
	CountDown = iota
)

type Timer struct {
	Name            string
	Key             string
	index           int
	timerType       int
	totalSecs       float64
	initialDuration time.Duration
	set             bool
	running         bool
	timerStop       chan bool
}

func NewTimer(timerType int, name string, key string, idx int, initialDuration time.Duration) (t *Timer, err error) {
	if timerType < CountUp || timerType > CountDown {
		err = errors.New("Timer type not one of 'CountUp' or 'CountDown'")
		return nil, err
	}
	log.Printf("Creating new timer '%s' with initial duration of %fs", name, initialDuration.Seconds())
	t = new(Timer)
	t.timerType = timerType
	t.Name = name
	t.index = idx
	t.Key = key
	t.initialDuration = initialDuration
	if timerType == CountDown {
		t.totalSecs = initialDuration.Seconds()
		t.set = true
	} else {
		t.totalSecs = 0
	}

	return t, nil
}

func (t *Timer) update(duration time.Duration) {
	//log.Printf("%s updated from %f to %f", t.Name, t.totalSecs, duration.Seconds())
	t.set = true
	t.totalSecs = math.Round(duration.Seconds())
	//e.hours = int(e.totalSecs/(60*60)) % 24
	//e.minutes = int(e.totalSecs/60) % 60
	//e.seconds = int(e.totalSecs) % 60
}

func (t *Timer) Start() {
	if t.timerStop == nil {
		t.timerStop = make(chan bool)
	}

	go func() {
		log.Printf("%s timer started", t.Name)
		t.running = true
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		start := time.Now()
		end := start.Add(time.Duration(t.initialDuration))
		for {
			select {
			case <-t.timerStop:
				log.Printf("%s timer stopped", t.Name)
				return
			case tick := <-ticker.C:
				switch t.timerType {
				case CountUp:
					t.update(tick.Sub(start))
					//log.Printf("%s elapsed: %s", t.Name, t.HMS())
				case CountDown:
					t.update(end.Sub(tick))
					//log.Printf("%s remaining: %s", t.Name, t.HMS())
				}
			}
		}
	}()
}

func (t *Timer) Stop() {
	if t.running {
		t.timerStop <- true
		t.running = false
	}
}

func (t *Timer) Reset() {
	t.Stop()
	if t.timerType == CountUp {
		t.set = false
	}
	t.totalSecs = t.initialDuration.Seconds()
}

func (t Timer) HMS() string {
	if !t.set {
		return "--:--:--"
	}
	//log.Printf("totalSecs = %f", t.totalSecs)
	secs := t.totalSecs
	if t.totalSecs < 0 {
		secs = math.Abs(t.totalSecs)
	}
	hours := int(secs/(60*60)) % 24
	minutes := int(secs/60) % 60
	seconds := int(secs) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func (t Timer) HMSIndicator() string {
	indicator := " "
	if t.timerType == CountDown {
		if t.totalSecs < 0 {
			indicator = "+"
		} else {
			indicator = "-"
		}
	} else {
		if t.totalSecs < 0 {
			indicator = "-"
		}
	}
	return indicator
}

func (t Timer) HMSWithIndicator() string {
	return fmt.Sprintf("%s %s", t.HMSIndicator(), t.HMS())
}

func (t Timer) Seconds() int {
	return int(t.totalSecs)
}

func (t Timer) Over() bool {
	return t.totalSecs < 0
}

func (t Timer) Type() int {
	return t.timerType
}

func (t Timer) Running() bool {
	return t.running
}
