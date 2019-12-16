// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package application

import (
	"fmt"
	"log"

	"github.com/pkg/errors"

	"github.com/elastic/beats/x-pack/agent/pkg/config"
	"github.com/elastic/beats/x-pack/agent/pkg/fleetapi"
)

type handlerPolicyChange struct {
	log     *log.Logger
	emitter emitterFunc
}

func (h *handlerPolicyChange) Handle(a action) error {
	h.log.Debugf("HandlerPolicyChange: action '%+v' received", a)
	action, ok := a.(*fleetapi.ActionPolicyChange)
	if !ok {
		return fmt.Errorf("invalid type, expected ActionPolicyChange and received %T", a)
	}

	c, err := config.NewConfigFrom(action.Policy)
	if err != nil {
		return errors.Wrap(err, "could not parse the configuration from the policy")
	}

	h.log.Debug("HandlerPolicyChange: emit configuration")

	return h.emitter(c)
}
