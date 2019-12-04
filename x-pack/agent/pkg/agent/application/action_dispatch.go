// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package application

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/elastic/beats/x-pack/agent/pkg/core/logger"
)

type action interface{}

type actionHandler interface {
	Handle(a action) error
}

type actionHandlers map[string]actionHandler

type actionDispatcher struct {
	log      *logger.Logger
	handlers actionHandlers
	def      actionHandler
}

func newActionDispatcher(log *logger.Logger, def actionHandler) (*actionDispatcher, error) {
	var err error
	if log == nil {
		log, err = logger.New()
		if err != nil {
			return nil, err
		}
	}

	if def == nil {
		return nil, errors.New("missing default handler")
	}

	return &actionDispatcher{
		log:      log,
		handlers: make(actionHandlers),
		def:      def,
	}, nil
}

func (ad *actionDispatcher) Register(a action, handler actionHandler) {
	t := reflect.TypeOf(a).String()
	ad.handlers[t] = handler
}

func (ad *actionDispatcher) Dispatch(actions ...action) error {
	ad.log.Debugf("Dispatch %d actions", len(actions))
	for _, action := range actions {
		if err := ad.dispatchAction(action); err != nil {
			ad.log.Debugf("Failed to dispatch action %+v error %+v", action, err)
			return err
		}
		ad.log.Debugf("Succesfully dispatch action %+v", action)
	}
	return nil
}

func (ad *actionDispatcher) dispatchAction(a action) error {
	t := reflect.TypeOf(a).String()
	handler, found := ad.handlers[t]
	if !found {
		return ad.def.Handle(a)
	}

	return handler.Handle(a)
}
