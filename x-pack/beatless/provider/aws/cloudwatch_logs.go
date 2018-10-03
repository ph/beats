// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package aws

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/goformation/cloudformation"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/x-pack/beatless/core"
	"github.com/elastic/beats/x-pack/beatless/provider"
	"github.com/elastic/beats/x-pack/beatless/provider/aws/transformer"
)

const handlerName = "beatless"

// CloudwatchLogsConfig is the configuration for the cloudwatchlogs event type.
type CloudwatchLogsConfig struct {
	Triggers     []*CloudwatchLogsTriggerConfig `config:"triggers"`
	Description  string                         `config:"description"`
	Name         string                         `config:"name" validate:"nonzero,required"`
	LambdaConfig *lambdaConfig                  `config:",inline"`
}

// CloudwatchLogsTriggerConfig is the configuration for the specific triggers for cloudwatch.
type CloudwatchLogsTriggerConfig struct {
	LogGroupName  string `config:"log_group_name" validate:"nonzero,required"`
	FilterPattern string `config:"filter_pattern"`
}

// Validate validates the configuration.
func (cfg *CloudwatchLogsConfig) Validate() error {
	if len(cfg.Triggers) == 0 {
		return errors.New("you need to specify at least one trigger")
	}
	return nil
}

// CloudwatchLogs receives CloudwatchLogs events from a lambda function and forward the logs to
// an Elasticsearch cluster.
type CloudwatchLogs struct {
	log    *logp.Logger
	config *CloudwatchLogsConfig
}

// NewCloudwatchLogs create a new function to listen to cloudwatch logs events.
func NewCloudwatchLogs(provider provider.Provider, cfg *common.Config) (provider.Function, error) {
	config := &CloudwatchLogsConfig{
		LambdaConfig: DefaultLambdaConfig,
	}
	if err := cfg.Unpack(config); err != nil {
		return nil, err
	}
	return &CloudwatchLogs{log: logp.NewLogger("cloudwatch_logs"), config: config}, nil
}

// Run start the AWS lambda handles and will transform any events received to the pipeline.
func (c *CloudwatchLogs) Run(_ context.Context, client core.Client) error {
	lambda.Start(c.createHandler(client))
	return nil
}

func (c *CloudwatchLogs) createHandler(client core.Client) func(request events.CloudwatchLogsEvent) error {
	return func(request events.CloudwatchLogsEvent) error {
		parsedEvent, err := request.AWSLogs.Parse()
		if err != nil {
			c.log.Errorf("could not parse events from cloudwatch logs, error: %s", err)
			return err
		}

		c.log.Debugf(
			"received %d events (logStream: %s, owner: %s, logGroup: %s, messageType: %s)",
			len(parsedEvent.LogEvents),
			parsedEvent.LogStream,
			parsedEvent.Owner,
			parsedEvent.LogGroup,
			parsedEvent.MessageType,
		)

		events := transformer.CloudwatchLogs(parsedEvent)
		if err := client.PublishAll(events); err != nil {
			c.log.Errorf("could not publish events to the pipeline, error: %s")
			return err
		}
		client.Wait()
		return nil
	}
}

// Name returns the name of the function.
func (c CloudwatchLogs) Name() string {
	return "cloudwatch_logs"
}

// AWSLogsSubscriptionFilter overrides the type from goformation to allow to pass an empty string.
// The API support an empty string, but requires one, the original type does not permit that.
type AWSLogsSubscriptionFilter struct {
	DestinationArn string `json:"DestinationArn,omitempty"`
	FilterPattern  string `json:"FilterPattern"`
	LogGroupName   string `json:"LogGroupName,omitempty"`
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r AWSLogsSubscriptionFilter) MarshalJSON() ([]byte, error) {
	type Properties AWSLogsSubscriptionFilter
	return json.Marshal(&struct {
		Type           string
		Properties     Properties
		DeletionPolicy cloudformation.DeletionPolicy `json:"DeletionPolicy,omitempty"`
	}{
		Type:       r.AWSCloudFormationType(),
		Properties: (Properties)(r),
	})
}

// AWSCloudFormationType return the AWS type.
func (r *AWSLogsSubscriptionFilter) AWSCloudFormationType() string {
	return "AWS::Logs::SubscriptionFilter"
}

// Template returns the cloudformation template for configuring the service with the specified triggers.
func (c *CloudwatchLogs) Template() *cloudformation.Template {
	prefix := func(suffix string) string {
		return "btl" + c.config.Name + suffix
	}

	template := cloudformation.NewTemplate()
	for idx, trigger := range c.config.Triggers {
		template.Resources[prefix("Permission"+strconv.Itoa(idx))] = &cloudformation.AWSLambdaPermission{
			Action:       "lambda:InvokeFunction",
			FunctionName: cloudformation.GetAtt(prefix(""), "Arn"),
			Principal: cloudformation.Join("", []string{
				"logs.",
				cloudformation.Ref("AWS::Region"), // Use the configuration region.
				".",
				cloudformation.Ref("AWS::URLSuffix"), // awsamazon.com or .com.ch
			}),
			SourceArn: cloudformation.Join(
				"",
				[]string{
					"arn:",
					cloudformation.Ref("AWS::Partition"),
					":logs:",
					cloudformation.Ref("AWS::Region"),
					":",
					cloudformation.Ref("AWS::AccountId"),
					":log-group:",
					trigger.LogGroupName,
					":*",
				},
			),
		}

		normalize := func(c string) string {
			return strings.Replace(c, "/", "", -1)
		}

		template.Resources[prefix("SubscriptionFilter"+normalize(trigger.LogGroupName))] = &AWSLogsSubscriptionFilter{
			DestinationArn: cloudformation.GetAtt(prefix(""), "Arn"),
			FilterPattern:  trigger.FilterPattern,
			LogGroupName:   trigger.LogGroupName,
		}
	}
	return template
}

// LambdaConfig returns the configuration to use when creating the lambda.
func (c *CloudwatchLogs) LambdaConfig() *lambdaConfig {
	return c.config.LambdaConfig
}
