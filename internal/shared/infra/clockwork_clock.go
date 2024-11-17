package infra

import (
	"time"

	"github.com/jonboulle/clockwork"
)

type ClockworkClock struct {
	clock clockwork.Clock
}

func NewClockworkClock() ClockworkClock {
	return ClockworkClock{
		clock: clockwork.NewRealClock(),
	}
}

func (c ClockworkClock) Now() time.Time {
	return c.clock.Now()
}
