// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package core

import (
	"context"

	"github.com/joeshaw/multierror"

	"github.com/elastic/beats/libbeat/logp"
)

// Runner is the interface that the coordinator will follow to manage a function goroutine.
type Runner interface {
	Run(context.Context) error
}

// Coordinator takes care of managing the function goroutine, it receives the list of functions that
// need to be executed and manage the goroutine.  If an error happen and its not handled by the
// function, we assume its a fatal error and we will
// stop all the other goroutine and beatless will terminate.
type Coordinator struct {
	log     *logp.Logger
	runners []Runner
}

// NewCoordinator create a new coordinator objects receiving the clientFactory and the runner.
func NewCoordinator(log *logp.Logger,
	runners ...Runner,
) *Coordinator {
	if log == nil {
		log = logp.NewLogger("coordinator")
	}
	return &Coordinator{log: log, runners: runners}
}

// Run starts each functions into an independent goroutine and wait until all the goroutine are
// stopped to exit.
func (r *Coordinator) Run(ctx context.Context) error {
	r.log.Debug("coordinator started")
	defer r.log.Debug("coordinator stopped")

	// When an errors happen in a function and its not handled by the running function, we log an error
	// and we trigger a shutdown of all the others goroutine and start will return and beatless
	// will stop.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make(chan error)
	defer close(results)

	r.log.Debugf("starting %d functions", len(r.runners))
	for _, rfn := range r.runners {
		go r.startFunc(ctx, cancel, rfn, results)
	}

	// Wait for goroutine to complete and aggregate any errors from the goroutine and
	// raise them back to the main program.
	var aggErr multierror.Errors
	for i := 0; i < len(r.runners); i++ {
		select {
		case err := <-results:
			if err != nil {
				aggErr = append(aggErr, err)
			}
		}
	}
	return aggErr.Err()
}

func (r *Coordinator) startFunc(
	ctx context.Context,
	cancel context.CancelFunc,
	rfn Runner,
	output chan<- error,
) {
	r.log.Info("starting function: %s", rfn)
	defer r.log.Infof("function stopped: %s", rfn)

	err := rfn.Run(ctx)
	defer func() { output <- err }()

	if err != nil {
		r.log.Error(
			"nonrecoverable error when executing the function: '%s', error: '%s', terminating all the functions",
			rfn,
			err,
		)
		cancel()
	}
}
