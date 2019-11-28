// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by x-pack/dev-tools/cmd/buildspec/buildspec.go - DO NOT EDIT.

package program

import "github.com/elastic/beats/x-pack/agent/pkg/packer"

var Supported []Spec

func init() {
	// Packed Files
	// spec/auditbeat.yml
	// spec/filebeat.yml
	// spec/heartbeat.yml
	// spec/journalbeat.yml
	// spec/metricbeat.yml
	unpacked := packer.MustUnpack("eJzUV9Fzo7YTfv/9GXn9dVoMZ/fozD0YcgiIgxucSEJvSHIAW4BrsAl0+r93AGODY18vTWZu+uCZREi73+5++6305022WbJf/B2Pcrr085/LWNz8dkNjIyePafCA4drHbkSQIemxk5PbNPD6a4k9ogBW8yDNLUA21HQFE6rsoZcRwc4zl4XkG2pJEBdLczrRo2lg6dqW4PXEMt0xA093lj4NLNMW80hbtrYf9rMoO6w7ggK44kAt55FWcdMWHhpVeiQFFBnFPEiDeh8DhuTfpgFT3JIgI2fFYT2BGcH3B7/T3DLd0kPjigJDIlBVfOymWLEFkeEndpu2PpufVnjYTeeRJi2xJg7nW5vx08QyxiFH7obrn9NZlG3635uYb9PAvv16N1toK6po4w4nSew9fXyNk4PPEwsYIwKE5MUw5NM6n6PwiMnUQg6C4R593MtHg7nk6EWwqotDClgMdx4aic5/8wOkpLI0zH8bs0SrNKBIXdd26hyfn+/w9zkwSxyJxTCk+P6A1d3PIy2nsttwo7ZDYpFxAMtjvWSYEeRIVLGrbm220NqcAlt4+OEYA1WgxEwoHWsIcrHEzVmhJzwl6NPE+hqGTBISQaPqToe/InmQpyMWjKWNpXtxu38cUgQrBowVWQRJ58cui8CWYeZhR/KRUxFklJ4cJLNpjZOXBLkbVrLM0vnWQ2TrLZq/MyrzigOj8nW20YMvX25+anvrORLLC63lxiIjqA3Vk0freaTFPnoRbfk3x1LOFlru4WPrdLT8x9Jz0w5p4mQedofl76j0uvy5h7WCYGvSp7Ola4LHcMd1tcE3j7SMyqx3rqYUlAkaSzPFkQgQuwYfcDe077emDyaCJg97Dpyibi1PmV6xA9fsNg04coszLLXvPQHqypdhLQlrKjtbgq0zOy97UqptqUZqRYGqUKAm80hr14ohrvrbTHHGTH7Zk6puT1GfSfrUxQoJqQlPOb4qTdck7ux8Te3GRydNbe5O/8M1mx5bo5crKejiP+21BQGi6tlqYylqGXH2zFx3/IlZrOYX8DatT3qcnC3qXNUt8nR+9sTT1xLU2PHRuODYfUvr97h/XNtz7Ba8JwW+DMdMcfcsfrrjcijoKg0WwKgW2BUU36c2zmkvhl2H411tr08TFtdyd99JwHe1fbj0t5dG6gLAkCXDker31z54pFKk7nzs7ln0eXdXnNfAmtS86+T7I7QGK3zDQfjMYpgQHBbfqT1lnRscj/Y0FjWu/d2t+vB7u+//V0fswc7bxmwoPGyfxqxpi0Nt7gbjsb/vCs+PuU2a3NUclk76qVWk5rTSfD9dQYb+38X1Pm9mx9zZP4Tvq3S3TXxxifGxOmLxOKRgMO3+oAosaWxkFybeIZP2M4+NjCNYvVaPIXPbyuVi+ZgG3BQFeThU/Oz8bNFOEL/DNH2bQvpnsXxDJYfxfZxSnuXtnQwaxvMxivkvL0rxMt9G7AKBHhGUWCxWB6KsKHqYWGAkuGlvPLktxOmmf2Fkj67c9C8Q5nC7zwl2Sx/d/xdeIkesb5HMj3oheLJaLM+l87wxDy+CFqfTjrlvvC5ov97XpbfBwIremFy9vpoNrlsDHr2vcc4w/pgXxl//+zsAAP//k9VT0g==")

	for f, v := range unpacked {
		s, err := NewSpecFromBytes(v)
		if err != nil {
			panic("Cannot read spec from " + f)
		}
		Supported = append(Supported, s)
	}
}
