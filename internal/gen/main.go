// Generates Counter, Gauge, Delta code for a given datatype
package main

import (
	"flag"
	"io"
	"os"
	"text/template"
)

type CodeDesc struct {
	TypeStr string // e.g. float32
	Suffix  string // e.g. Float32, or could possibly be blank.
	Val     string // e.g., strconv.FormatInt(int64(val), 10)
}

var CodeTmpl = template.Must(template.New("code").Parse(`// Package statsdclient provides a StatsD client that simply formats output to any underlying io.Writer.
//
// this file is auto-generated. do not edit.
//
// parameters:
//  type: {{ .TypeStr }}
//  suffix: {{ .Suffix }}
//  valueconv: {{ .ValAsStr }}
//
package statsdclient

import "strconv"

// SampledInc{{ .Suffix }} emits a 'sampled' {{ .TypeStr }} counter with the given rate.
func (c *Client) SampledInc{{ .Suffix }}(stat string, val {{ .TypeStr }}, rate float32) error {
  b := c.getbuf(stat)
  {{ .Val }}
  return c.finishbuf(b, "|c", rate)
}

// SampledGauge{{ .Suffix }} emits a 'sampled' {{ .TypeStr }} gauge with the given rate.
func (c *Client) SampledGauge{{ .Suffix }}(stat string, val {{ .TypeStr }}, rate float32) error {
  b := c.getbuf(stat)
  {{ .Val }}
  return c.finishbuf(b, "|g", rate)
}

// SampledDelta{{ .Suffix }} emits a 'sampled' {{ .TypeStr }} delta with the given rate.
func (c *Client) SampledDelta{{ .Suffix }}(stat string, val {{ .TypeStr }}, rate float32) error {
  b := c.getbuf(stat)
  if val > 0 {
    b = append(b, '+')
  }
  {{ .Val }}
  return c.finishbuf(b, "|g", rate)
}

// Inc{{ .Suffix }} emits a {{ .TypeStr }} counter.
func (c *Client) Inc{{ .Suffix }}(stat string, val {{ .TypeStr }}) error {
  return c.SampledInc{{ .Suffix }}(stat, val, 1.0)
}

// Gauge{{ .Suffix }} emits a {{ .TypeStr }} gauge.
func (c *Client) Gauge{{ .Suffix }}(stat string, val {{ .TypeStr }}) error {
  return c.SampledGauge{{ .Suffix }}(stat, val, 1.0)
}

// Delta{{ .Suffix }} emits a {{ .TypeStr }} delta.
func (c *Client) Delta{{ .Suffix }}(stat string, val {{ .TypeStr }}) error {
  return c.SampledDelta{{ .Suffix }}(stat, val, 1.0)
}`))

var types = map[string]CodeDesc{
	"int": {
		TypeStr: "int", Suffix: "", Val: "in64(val)",
	},
	"int64": {
		TypeStr: "int64", Suffix: "Int64", Val: "val",
	},
	"float64": {
		TypeStr: "float64", Suffix: "Float64", Val: "val",
	},
}

func main() {
	var (
		fType     = flag.String("type", "int", "Type to generate code for")
		fFileDest = flag.String("file-dest", "", "File to write output to")
	)
	flag.Parse()

	c := types[*fType]

	var codeOut io.Writer
	if *fFileDest == "" {
		codeOut = os.Stdout
	} else {
		f, err := os.Create(*fFileDest)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		codeOut = f
	}

	err := CodeTmpl.Execute(codeOut, c)
	if err != nil {
		panic(err)
	}
}
