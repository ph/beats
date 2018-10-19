// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package aws

import (
	"errors"

	"github.com/elastic/beats/libbeat/logp"
)

var (
	errNeverRun        = errors.New("executor was never executed")
	errCannotAdd       = errors.New("cannot add to an already executed executor")
	errAlreadyExecuted = errors.New("executor already executed")
)

type executor struct {
	operations []doer
	undos      []undoer
	completed  bool
	log        *logp.Logger
}

type doer interface {
	Execute() error
}

type undoer interface {
	Rollback() error
}

func newExecutor(log *logp.Logger) *executor {
	if log == nil {
		log = logp.NewLogger("")
	}

	log = log.Named("executor")
	return &executor{log: log}
}

func (e *executor) Execute() (err error) {
	e.log.Debugf("Executing %d operations", len(e.operations))
	if e.IsCompleted() {
		return errAlreadyExecuted
	}
	for _, operation := range e.operations {
		err = operation.Execute()
		if err != nil {
			break
		}
		v, ok := operation.(undoer)
		if ok {
			e.undos = append(e.undos, v)
		}
	}
	if err == nil {
		e.log.Debug("All operations have been successful")
	}
	e.markCompleted()
	return err
}

func (e *executor) Rollback() (err error) {
	e.log.Debugf("Rolling back previous execution, %d operations", len(e.undos))
	if !e.IsCompleted() {
		return errNeverRun
	}
	for i := len(e.undos) - 1; i >= 0; i-- {
		operation := e.undos[i]
		err = operation.Rollback()
		if err != nil {
			break
		}
	}

	if err == nil {
		e.log.Debug("Rollback successful")
	} else {
		e.log.Debug("Rollback incomplete")
	}
	return err
}

func (e *executor) Add(operation ...doer) error {
	if e.IsCompleted() {
		return errCannotAdd
	}
	e.operations = append(e.operations, operation...)
	return nil
}

func (e *executor) markCompleted() {
	e.completed = true
}

func (e *executor) IsCompleted() bool {
	return e.completed
}
