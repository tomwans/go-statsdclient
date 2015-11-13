package statsdclient

import (
	"bytes"
	"testing"
	"time"
)

type Experiment struct {
	Stat     string
	Val      int
	Rate     float32
	Expected string
}

var incs = []Experiment{
	{Stat: "testa", Val: 0, Rate: 1.0, Expected: "testa:0|c\n"},
	{Stat: "testb", Val: 0, Rate: 0.1, Expected: "testb:0|c|@0.1000\n"},
	{Stat: "testc", Val: 420, Rate: 1.0, Expected: "testc:420|c\n"},
	{Stat: "testd", Val: -420, Rate: 1.0, Expected: "testd:-420|c\n"},
	{Stat: "teste", Val: 420, Rate: 0.01, Expected: "teste:420|c|@0.0100\n"},
	{Stat: "testf", Val: -420, Rate: 0.01, Expected: "testf:-420|c|@0.0100\n"},
}

func TestInc(t *testing.T) {
	buf := &bytes.Buffer{}
	c := Client{Prefix: "", Writer: buf}
	for _, experiment := range incs {
		buf.Reset()
		c.SampledInc(experiment.Stat, experiment.Val, experiment.Rate)
		got := string(buf.Bytes())
		if got != experiment.Expected {
			t.Fatalf("Inc w/ args: %s %d %f: got: %s, but expected: %s", experiment.Stat, experiment.Val, experiment.Rate, got, experiment.Expected)
		}
	}
}

var gauges = []Experiment{
	{Stat: "test", Val: 0, Rate: 1.0, Expected: "test:0|g\n"},
	{Stat: "test", Val: 0, Rate: 0.1, Expected: "test:0|g|@0.1000\n"},
	{Stat: "test", Val: 420, Rate: 1.0, Expected: "test:420|g\n"},
	{Stat: "test", Val: -420, Rate: 1.0, Expected: "test:-420|g\n"},
	{Stat: "test", Val: 420, Rate: 0.01, Expected: "test:420|g|@0.0100\n"},
	{Stat: "test", Val: -420, Rate: 0.01, Expected: "test:-420|g|@0.0100\n"},
}

func TestGauges(t *testing.T) {
	buf := &bytes.Buffer{}
	c := Client{Prefix: "", Writer: buf}
	for _, experiment := range gauges {
		buf.Reset()
		c.SampledGauge(experiment.Stat, experiment.Val, experiment.Rate)
		got := string(buf.Bytes())
		if got != experiment.Expected {
			t.Fatalf("Gauges w/ args: %s %d %f: got: %s, but expected: %s", experiment.Stat, experiment.Val, experiment.Rate, got, experiment.Expected)
		}
	}
}

var deltas = []Experiment{
	{Stat: "test", Val: 0, Rate: 1.0, Expected: "test:0|g\n"},
	{Stat: "test", Val: 0, Rate: 0.1, Expected: "test:0|g|@0.1000\n"},
	{Stat: "test", Val: 420, Rate: 1.0, Expected: "test:+420|g\n"},
	{Stat: "test", Val: -420, Rate: 1.0, Expected: "test:-420|g\n"},
	{Stat: "test", Val: 420, Rate: 0.01, Expected: "test:+420|g|@0.0100\n"},
	{Stat: "test", Val: -420, Rate: 0.01, Expected: "test:-420|g|@0.0100\n"},
}

func TestDeltas(t *testing.T) {
	buf := &bytes.Buffer{}
	c := Client{Prefix: "", Writer: buf}
	for _, experiment := range deltas {
		buf.Reset()
		c.SampledDelta(experiment.Stat, experiment.Val, experiment.Rate)
		got := string(buf.Bytes())
		if got != experiment.Expected {
			t.Fatalf("Delta w/ args: %s %d %f: got: %s, but expected: %s", experiment.Stat, experiment.Val, experiment.Rate, got, experiment.Expected)
		}
	}
}

type TimingExperiment struct {
	Stat     string
	Val      time.Duration
	Rate     float32
	Expected string
}

var timings = []TimingExperiment{
	{Stat: "testa", Val: time.Second * 2, Rate: 1.0, Expected: "testa:2000|ms\n"},
	{Stat: "testb", Val: time.Second * 3, Rate: 1.0, Expected: "testb:3000|ms\n"},
	{Stat: "testc", Val: time.Second * 4, Rate: 0.73, Expected: "testc:4000|ms|@0.7300\n"},
	{Stat: "testd", Val: time.Second * 59, Rate: 0.73, Expected: "testd:59000|ms|@0.7300\n"},
}

func TestTimings(t *testing.T) {
	buf := &bytes.Buffer{}
	c := Client{Prefix: "", Writer: buf}
	for _, experiment := range timings {
		buf.Reset()
		c.SampledTiming(experiment.Stat, experiment.Val, experiment.Rate)
		got := string(buf.Bytes())
		if got != experiment.Expected {
			t.Fatalf("Inc w/ args: %s %d %f: got: %s, but expected: %s", experiment.Stat, experiment.Val, experiment.Rate, got, experiment.Expected)
		}
	}
}
