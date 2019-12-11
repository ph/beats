package scheduler

import "time"

// Scheduler simple interface that encapsulate the scheduling logic, this is useful if you want to
// test asynchronous code in a synchronous way.
type Scheduler interface {
	WaitTick() <-chan time.Time
	Stop()
}

// Stepper is a scheduler where each Tick is manully triggered, this is useful in scenario
// when you want to test the behavior of asynchronous code in a synchronous way.
type Stepper struct {
	C chan time.Time
}

// Next trigger the WaitTick unblock manually.
func (s *Stepper) Next() {
	s.C <- time.Now()
}

// WaitTick returns a channel to watch for ticks.
func (s *Stepper) WaitTick() <-chan time.Time {
	return s.C
}

// Stop is stopping the scheduler, in the case of the Stepper scheduler nothing is done.
func (s *Stepper) Stop() {}

// NewStepper returns a new Stepper scheduler where the tick is manually controlled.
func NewStepper() *Stepper {
	return &Stepper{
		C: make(chan time.Time),
	}
}

// Periodic wraps a time.Timer as the scheduler.
type Periodic struct {
	Ticker *time.Ticker
}

// NewPeriodic returns a Periodic scheduler that will unblock the WaitTick based on a duration.
func NewPeriodic(d time.Duration) *Periodic {
	return &Periodic{Ticker: time.NewTicker(d)}
}

// WaitTick wait on the duration to be experied to unblock the channel.
func (p *Periodic) WaitTick() <-chan time.Time {
	return p.Ticker.C
}

// Stop stops the internal Ticker.
// Note this will not close the internal channel is up to the developer to unblock the goroutine
// using another mechanism.
func (p *Periodic) Stop() {
	p.Ticker.Stop()
}
