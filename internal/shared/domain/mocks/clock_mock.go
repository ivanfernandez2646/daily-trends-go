package mocks

import (
	"time"

	"github.com/jonboulle/clockwork"
)

type ClockMock struct {
	clock clockwork.Clock
}

func NewClockMock() ClockMock {
	return ClockMock{clock: clockwork.NewFakeClock()}
}

func (cm ClockMock) Now() time.Time {
	return cm.clock.Now()
}
