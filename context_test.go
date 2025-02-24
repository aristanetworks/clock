package clock

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

// delay is the extra delay use to ensure that enough time has passed
const delay = 10 * time.Millisecond

// Ensure that WithDeadline is cancelled when deadline exceeded.
func TestMock_WithDeadline(t *testing.T) {
	c := New()
	cause := errors.New("example cause")
	ctx, _ := ContextWithDeadlineCause(context.Background(), c, c.Now().Add(time.Second), cause)
	time.Sleep(time.Second + delay)
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Error("invalid type of error returned when deadline exceeded")
		}
		if context.Cause(ctx) != cause {
			t.Error("cause does not match expected cause")
		}
	default:
		t.Error("context is not cancelled when deadline exceeded")
	}
}

// Ensure that WithDeadline does nothing when the deadline is later than the current deadline.
func TestMock_WithDeadlineLaterThanCurrent(t *testing.T) {
	c := New()
	cause := errors.New("example cause")
	ctx, _ := ContextWithDeadlineCause(context.Background(), c, c.Now().Add(time.Second), cause)
	ctx, _ = ContextWithDeadlineCause(ctx, c, c.Now().Add(10*time.Second), cause)
	time.Sleep(time.Second + delay)
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Error("invalid type of error returned when deadline exceeded")
		}
		if context.Cause(ctx) != cause {
			t.Error("cause does not match expected cause")
		}
	default:
		t.Error("context is not cancelled when deadline exceeded")
	}
}

// Ensure that WithDeadline cancel closes Done channel with context.Canceled error.
func TestMock_WithDeadlineCancel(t *testing.T) {
	c := New()
	ctx, cancel := ContextWithDeadline(context.Background(), c, c.Now().Add(time.Second))
	cancel()
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Error("invalid type of error returned after cancellation")
		}
	case <-time.After(time.Second + delay):
		t.Error("context is not cancelled after cancel was called")
	}
}

// Ensure that WithDeadline closes child contexts after it was closed.
func TestMock_WithDeadlineCancelledWithParent(t *testing.T) {
	c := New()
	parent, cancel := context.WithCancel(context.Background())
	ctx, _ := ContextWithDeadline(parent, c, c.Now().Add(time.Second))
	cancel()
	runtime.Gosched()
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.Canceled) {
			t.Errorf("invalid type of error returned after cancellation: %v", ctx.Err())
		}
	default:
		t.Error("context is not cancelled when parent context is cancelled")
	}
}

// Ensure that WithDeadline cancelled immediately when deadline has already passed.
func TestMock_WithDeadlineImmediate(t *testing.T) {
	c := New()
	ctx, _ := ContextWithDeadline(context.Background(), c, c.Now().Add(-time.Second))
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Error("invalid type of error returned when deadline has already passed")
		}
	default:
		t.Error("context is not cancelled when deadline has already passed")
	}
}

// Ensure that WithTimeout is cancelled when deadline exceeded.
func TestMock_WithTimeout(t *testing.T) {
	c := New()
	ctx, _ := ContextWithTimeout(context.Background(), c, time.Second)
	time.Sleep(time.Second + delay)
	select {
	case <-ctx.Done():
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			t.Error("invalid type of error returned when time is over")
		}
	default:
		t.Error("context is not cancelled when time is over")
	}
}
