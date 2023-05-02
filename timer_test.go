package main

import (
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	tmrUp, err := NewTimer(CountUp, "Test Up Timer", "test_up", 0, time.Duration(0))
	if err != nil || tmrUp == nil {
		t.Errorf("failed to create up timer: %s", err)
	}
	if tmrUp.totalSecs != 0 {
		t.Errorf("up timer totalSecs want 0 got %f", tmrUp.totalSecs)
	}

	d := time.Duration(1 * time.Minute)
	tmrDown, err := NewTimer(CountDown, "Test Down Timer", "test_down", 1, d)
	if err != nil || tmrDown == nil {
		t.Errorf("failed to create down timer: %s", err)
	}
	if tmrDown.totalSecs != d.Seconds() {
		t.Errorf("down timer totalSecs want %f got %f", d.Seconds(), tmrDown.totalSecs)
	}
}

func TestTimerUpdate(t *testing.T) {
	tmr, err := NewTimer(CountUp, "Test Timer", "test", 0, time.Duration(0))
	if err != nil || tmr == nil {
		t.Errorf("failed to create timer: %s", err)
	}

	d, _ := time.ParseDuration("1s")
	cur := tmr.totalSecs
	tmr.update(d)
	if tmr.totalSecs != cur+d.Seconds() {
		t.Errorf("update wanted %f, got %f", cur+d.Seconds(), tmr.totalSecs)
	}
}

/*
func TestTimerStart(t *testing.T) {
	tmr, err := NewTimer(CountUp, "Test Timer", "test", 0, time.Duration(0))
	if err != nil || tmr == nil {
		t.Errorf("failed to create timer: %s", err)
	}

	tmr.Start()
	time.Sleep(100 * time.Millisecond)
	if !tmr.running {
		t.Errorf("failed to start timer")
	}
}

func TestTimerStop(t *testing.T) {
	tmr, err := NewTimer(CountUp, "Test Timer", "test", 0, time.Duration(0))
	if err != nil || tmr == nil {
		t.Errorf("failed to create timer: %s", err)
	}

	tmr.Start()
	time.Sleep(100 * time.Millisecond)
	tmr.Stop()
	if tmr.running {
		t.Errorf("failed to stop timer")
	}
}
*/

func TestTimerHMS(t *testing.T) {
	tmr, err := NewTimer(CountUp, "Test Timer", "test", 0, time.Duration(0))
	if err != nil || tmr == nil {
		t.Errorf("failed to create timer: %s", err)
	}
	tmr.totalSecs = 60.0
	tmr.set = true
	hms := tmr.HMS()
	if hms != "00:01:00" {
		t.Errorf("HMS want 00:01:00, got %s", hms)
	}
}

func TestTimerHMSIndicator(t *testing.T) {
	tmr, err := NewTimer(CountDown, "Test Timer", "test", 0, time.Duration(1*time.Minute))
	if err != nil || tmr == nil {
		t.Errorf("failed to create timer: %s", err)
	}
	tmr.set = true
	hmsIndicator := tmr.HMSIndicator()
	if hmsIndicator != "-" {
		t.Errorf("want indicator -, got %s", hmsIndicator)
	}

	tmr.totalSecs = -60.0
	hmsIndicator = tmr.HMSIndicator()
	if hmsIndicator != "+" {
		t.Errorf("want indicator +, got %s", hmsIndicator)
	}

	tmr.timerType = CountUp
	hmsIndicator = tmr.HMSIndicator()
	if hmsIndicator != "-" {
		t.Errorf("want indicator -, got %s", hmsIndicator)
	}
}

func TestTimerOver(t *testing.T) {
	tmr, err := NewTimer(CountDown, "Test Timer", "test", 0, time.Duration(-1*time.Minute))
	if err != nil || tmr == nil {
		t.Errorf("failed to create timer: %s", err)
	}

	if !tmr.Over() {
		t.Error("timer over not set when it should be")
	}
}
