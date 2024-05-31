// Copyright (c) 2018 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package mock

import (
	"testing"

	"github.com/aristanetworks/clock"
	"go.uber.org/mock/gomock"
)

// TestMocks validates that the mock clock is in sync with the current
// Clock, Timer and Ticker interfaces
func TestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	var _ clock.Clock = NewMockClock(ctrl)
	var _ clock.Ticker = NewMockTicker(ctrl)
	var _ clock.Timer = NewMockTimer(ctrl)
}
