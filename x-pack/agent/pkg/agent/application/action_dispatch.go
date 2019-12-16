// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package application

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"reflect"

	"github.com/pkg/errors"

	"github.com/elastic/beats/x-pack/agent/pkg/core/logger"
)

func init() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
}

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

func (ad *actionDispatcher) Register(a action, handler actionHandler) error {
	k := ad.key(a)
	_, ok := ad.handlers[k]
	if ok {
		return fmt.Errorf("action with type %T is already registered", a)
	}
	ad.handlers[k] = handler
	return nil
}

func (ad *actionDispatcher) MustRegister(a action, handler actionHandler) {
	err := ad.Register(a, handler)
	if err != nil {
		panic("could not register action, error: " + err.Error())
	}
}

func (ad *actionDispatcher) key(a action) string {
	return reflect.TypeOf(a).String()
}

func (ad *actionDispatcher) Dispatch(actions ...action) error {
	ad.log.Debugf("Dispatch %d actions", len(actions))
	for _, action := range actions {
		if err := ad.dispatchAction(action); err != nil {
			ad.log.Debugf("Failed to dispatch action '%+v', error: %+v", action, err)
			// fmt.Println(err)
			return err
		}
		ad.log.Debugf("Succesfully dispatched action: '%+v'", action)
	}
	return nil
}

func (ad *actionDispatcher) dispatchAction(a action) error {
	handler, found := ad.handlers[(ad.key(a))]
	if !found {
		return ad.def.Handle(a)
	}

	return handler.Handle(a)
}
