clock [![Coverage Status](https://coveralls.io/repos/aristanetworks/clock/badge.png?branch=master)](https://coveralls.io/r/aristanetworks/clock?branch=master) [![GoDoc](https://godoc.org/github.com/aristanetworks/clock?status.png)](https://godoc.org/github.com/aristanetworks/clock) ![Project status](http://img.shields.io/status/experimental.png?color=red)
=====
NOTE:  This README has not yet been updated to reflect the fixes/changes.  Please refer to the source code in the meantime.

Clock is a small library for mocking time in Go. It provides an interface
around the standard library's [`time`][time] package so that the application
can use the realtime clock while tests can use the mock clock.

[time]: http://golang.org/pkg/time/


## Usage

### Realtime Clock

Your application can maintain a `Clock` variable that will allow realtime and
mock clocks to be interchangable. For example, if you had an `Application` type:

```go
import "github.com/aristanetworks/clock"

type Application struct {
	Clock clock.Clock
}
```

You could initialize it to use the realtime clock like this:

```go
var app Application
app.Clock = clock.New()
...
```

Then all timers and time-related functionality should be performed from the
`Clock` variable.


### Mocking time

In your tests, you will want to use a `Mock` clock:

```go
import (
	"testing"

	"github.com/aristanetworks/clock/mock"
)

func TestApplication_DoSomething(t *testing.T) {
	mock := mock.NewMockClock(ctrl)
	app := Application{Clock: mock}
	...
}
```

Now that you've initialized your application to use the mock clock, you can
use the standard gomock methods to mock any call for any method.


### Examples

#### fake Sleep

```go
mock := clock.NewMockClock(ctrl)

mock.EXPECT().Sleep(gomock.Any()).AnyTimes()
```

#### Return a fake time

```go
mock := clock.NewMockClock(ctrl)

now := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
mock.EXPECT().Now().Return(now).AnyTimes()
```
