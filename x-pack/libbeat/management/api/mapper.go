package api

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// {
//   _type: 'filebeat.inputs',
//   _description: 'Something here',
//   _sub_type: 'elasticsearch', // this key will define things like output type, or module type
//   hosts: ['foo', 'bar'],
//   username: 'foo',
//   password: 'string',
// }

type configType string

type genericResponse struct {
	Data    json.RawMessage
	Error   error
	Success bool
}

func (g *genericResponse) UnmarshalJSON(b []byte) error {
	resp := struct {
		Data    json.RawMessage `json:"data"`
		Success bool            `json:"success"`
		Err     string          `json:"error,omitempty"`
	}{}

	if err := json.Unmarshal(b, &resp); err != nil {
		return err
	}

	*g = genericResponse{
		Data:    resp.Data,
		Error:   errors.New(resp.Err),
		Success: resp.Success,
	}

	return nil
}

type configurationResponse struct {
	Type        string
	Description string
	SubType     string
	Config      map[string]interface{}
}

func (c *configurationResponse) UnmarshalJSON(b []byte) error {
	var header := struct{
		Type        string `json:"_type"`
		Description string `json:"_description"`
		SubType     string `json:"_sub_type"`
	}

}
