// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/x-pack/beatless/core"
	"github.com/elastic/beats/x-pack/beatless/provider"
	"github.com/elastic/beats/x-pack/beatless/provider/aws/transformer"
)

// SQS receives events from the web service and forward them to elasticsearch.
type SQS struct {
	log *logp.Logger
}

// NewSQS creates a new function to receives events from a SQS queue.
func NewSQS(provider provider.Provider, config *common.Config) (provider.Function, error) {
	return &SQS{log: logp.NewLogger("sqs")}, nil
}

// Run starts the lambda function and wait for web triggers.
func (s *SQS) Run(_ context.Context, client core.Client) error {
	lambda.Start(func(request events.SQSEvent) error {
		// defensive checks
		if len(request.Records) == 0 {
			s.log.Error("no log events received from sqs")
			return errors.New("no event received")
		}

		s.log.Debugf("received %d events", len(request.Records))

		events := transformer.SQS(request)
		if err := client.PublishAll(events); err != nil {
			s.log.Errorf("could not publish events to the pipeline, error: %s")
			return err
		}
		client.Wait()
		return nil
	})

	return nil
}

// Name return the name of the lambda function.
func (s *SQS) Name() string {
	return "sqs"
}
