package clock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ContextWithTimeout(parent context.Context, clock Clock, timeout time.Duration) (context.Context, context.CancelFunc) {
	return ContextWithTimeoutCause(parent, clock, timeout, nil)
}

func ContextWithTimeoutCause(parent context.Context, clock Clock, timeout time.Duration, cause error) (context.Context, context.CancelFunc) {
	return ContextWithDeadlineCause(parent, clock, clock.Now().Add(timeout), cause)
}

func ContextWithDeadline(parent context.Context, clock Clock, deadline time.Time) (context.Context, context.CancelFunc) {
	return ContextWithDeadlineCause(parent, clock, deadline, nil)
}

func ContextWithDeadlineCause(parent context.Context, clock Clock, deadline time.Time, cause error) (context.Context, context.CancelFunc) {
	// using WithCancelCause to facilitate adding a Cause
	// delegating the actual context logic to a cancel context, using the background context here
	// since using the parent thread can cause this context to cancel before the timerCtx has been
	// populated with an error/cause
	wrapped, cancelFunc := context.WithCancelCause(context.Background())
	ctx := &timerCtx{
		clock:   clock,
		Context: wrapped,
		parent:  parent,
		cancelFunc: func() {
			cancelFunc(cause)
		},
		deadline: deadline,
	}
	propagateCancel(parent, ctx)
	dur := deadline.Sub(clock.Now())
	if dur <= 0 {
		ctx.cancel(context.DeadlineExceeded) // deadline has already passed
		return ctx, func() {}
	}
	if ctx.Err() == nil {
		ctx.Lock()
		ctx.timer = clock.AfterFunc(dur, func() {
			ctx.cancel(context.DeadlineExceeded)
		})
		ctx.Unlock()
	}
	return ctx, func() { ctx.cancel(context.Canceled) }
}

// propagateCancel arranges for child to be canceled when parent is.
func propagateCancel(parent context.Context, child *timerCtx) {
	if parent.Done() == nil {
		return // parent is never canceled
	}
	go func() {
		select {
		case <-parent.Done():
			child.cancel(parent.Err())
		case <-child.Done():
		}
	}()
}

type timerCtx struct {
	context.Context // wrapped cancelCtx

	clock      Clock
	parent     context.Context
	cancelFunc context.CancelFunc
	deadline   time.Time

	sync.Mutex // lock to protect the following fields
	err        error
	timer      Timer
}

func (c *timerCtx) cancel(err error) {
	c.Lock()
	defer c.Unlock()
	if c.err != nil {
		return // already canceled
	}
	c.err = err
	c.cancelFunc()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) { return c.deadline, true }

func (c *timerCtx) Err() error {
	c.Lock()
	defer c.Unlock()
	return c.err
}

func (c *timerCtx) Value(key any) any {
	parentValue := c.parent.Value(key)
	if parentValue != nil {
		return parentValue
	}
	return c.Context.Value(key)
}

func (c *timerCtx) String() string {
	return fmt.Sprintf("clock.WithDeadline(%s [%s])", c.deadline, c.deadline.Sub(c.clock.Now()))
}
