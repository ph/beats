package aws

import (
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/x-pack/beatless/provider"
)

func TestSQS(t *testing.T) {
	cfg := common.MustNewConfigFrom(map[string]interface{}{
		"name": "foobar",
	})

	t.Run("when publish is succesful", func(t *testing.T) {
		client := &arrayBackedClient{}
		s, err := NewSQS(&provider.DefaultProvider{}, cfg)
		if !assert.NoError(t, err) {
			return
		}

		c, _ := s.(*SQS)
		handler := c.createHandler(client)
		err = handler(generateSQSEvent())
		assert.NoError(t, err)
	})

	t.Run("when publish is not succesful", func(t *testing.T) {
		e := errors.New("something bad")
		client := &arrayBackedClient{err: e}

		s, err := NewSQS(&provider.DefaultProvider{}, cfg)
		if !assert.NoError(t, err) {
			return
		}

		c, _ := s.(*SQS)
		handler := c.createHandler(client)
		err = handler(generateSQSEvent())
		assert.Equal(t, e, err)
	})
}

func generateSQSEvent() events.SQSEvent {
	return events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				MessageId:     "1234",
				ReceiptHandle: "12345",
				Body:          "hello world",
			},
		},
	}
}
