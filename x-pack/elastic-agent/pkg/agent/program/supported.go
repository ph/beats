// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by x-pack/dev-tools/cmd/buildspec/buildspec.go - DO NOT EDIT.

package program

import (
	"strings"

	"github.com/elastic/beats/v7/x-pack/elastic-agent/pkg/packer"
)

var Supported []Spec
var SupportedMap map[string]Spec

func init() {
	// Packed Files
	// spec/endpoint.yml
	// spec/filebeat.yml
	// spec/metricbeat.yml
	unpacked := packer.MustUnpack("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")
	SupportedMap = make(map[string]Spec)

	for f, v := range unpacked {
		s, err := NewSpecFromBytes(v)
		if err != nil {
			panic("Cannot read spec from " + f)
		}
		Supported = append(Supported, s)
		SupportedMap[strings.ToLower(s.Cmd)] = s
	}
}
