// Copyright (c) 2014 Ben Johnson

package clock

import (
	"sync"
	"testing"
	"time"
)

// Ensure that the clock's After channel sends at the correct time.
func TestClock_After(t *testing.T) {
	clk := New()
	st := clk.Now()
	<-clk.After(100 * time.Millisecond)
	elapsed := clk.Now().Sub(st)
	if elapsed < 10*time.Millisecond {
		t.Fatalf("Expected more time to elapse for 100ms sleep than %s", elapsed)
	}
}

// Ensure that the clock's AfterFunc executes at the correct time.
func TestClock_AfterFunc(t *testing.T) {
	clk := New()
	st := clk.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	clk.AfterFunc(100*time.Millisecond, func() {
		wg.Done()
	})
	wg.Wait()
	elapsed := clk.Now().Sub(st)
	if elapsed < 10*time.Millisecond {
		t.Fatalf("Expected more time to elapse for 100ms sleep than %s", elapsed)
	}
}

// Ensure that the clock's time matches the standary library.
func TestClock_Now(t *testing.T) {
	a := time.Now().Round(time.Minute)
	b := New().Now().Round(time.Minute)
	if !a.Equal(b) {
		t.Errorf("not equal: %s != %s", a, b)
	}
}

// Ensure that the clock sleeps for the appropriate amount of time.
func TestClock_Sleep(t *testing.T) {
	clk := New()
	st := clk.Now()
	clk.Sleep(100 * time.Millisecond)
	elapsed := clk.Now().Sub(st)
	if elapsed < 10*time.Millisecond {
		t.Fatalf("Expected more time to elapse for 100ms sleep than %s", elapsed)
	}
}

// Ensure that the clock ticks correctly.
func TestClock_Tick(t *testing.T) {
	clk := New()
	st := clk.Now()
	c := clk.Tick(50 * time.Millisecond)
	<-c
	<-c
	elapsed := clk.Now().Sub(st)
	if elapsed < 10*time.Millisecond {
		t.Fatalf("Expected more time to elapse for two 50ms ticks than %s", elapsed)
	}
}

// Ensure that the clock's ticker ticks correctly.
func TestClock_Ticker(t *testing.T) {
	clk := New()
	st := clk.Now()
	ticker := clk.Ticker(50 * time.Millisecond)
	<-ticker.C()
	<-ticker.C()
	elapsed := clk.Now().Sub(st)
	if elapsed < 10*time.Millisecond {
		t.Fatalf("Expected more time to elapse for two 50ms ticks than %s", elapsed)
	}
}

// Stop a ticker
func TestClock_Ticker_Stp(t *testing.T) {
	ticker := New().Ticker(20 * time.Millisecond)
	<-ticker.C()
	ticker.Stop()
}

// Ensure that the clock's timer waits correctly.
func TestClock_Timer(t *testing.T) {
	clk := New()
	st := clk.Now()

	timer := clk.Timer(100 * time.Millisecond)
	<-timer.C()

	elapsed := clk.Now().Sub(st)
	if elapsed < 10*time.Millisecond {
		t.Fatalf("Expected more time to elapse for a 100ms timer than %s", elapsed)
	}

	if timer.Stop() {
		t.Fatal("timer still running")
	}
}

// Ensure that the clock's timer can be stopped.
func TestClock_Timer_Stop(t *testing.T) {
	timer := New().Timer(20 * time.Millisecond)
	if !timer.Stop() {
		t.Fatal("timer not running")
	}
	if timer.Stop() {
		t.Fatal("timer wasn't cancelled")
	}
}

// Ensure that the clock's timer can be reset.
func TestClock_Timer_Reset(t *testing.T) {
	timer := New().Timer(10 * time.Millisecond)
	if !timer.Reset(20 * time.Millisecond) {
		t.Fatal("timer not running")
	}
	<-timer.C()
}
