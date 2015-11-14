// Package statsdclient provides a StatsD client that simply formats output to any underlying io.Writer.
//
// this file is auto-generated. do not edit.
//
// parameters:
//  type: int
//  suffix:
//  valueconv: b = strconv.AppendInt(b, int64(val), 10)
//
package statsdclient

import "strconv"

// SampledInc emits a 'sampled' int counter with the given rate.
func (c *Client) SampledInc(stat string, val int, rate float32) error {
	b := c.getbuf(stat)
	b = strconv.AppendInt(b, int64(val), 10)
	return c.finishbuf(b, "|c", rate)
}

// SampledGauge emits a 'sampled' int gauge with the given rate.
func (c *Client) SampledGauge(stat string, val int, rate float32) error {
	b := c.getbuf(stat)
	b = strconv.AppendInt(b, int64(val), 10)
	return c.finishbuf(b, "|g", rate)
}

// SampledDelta emits a 'sampled' int delta with the given rate.
func (c *Client) SampledDelta(stat string, val int, rate float32) error {
	b := c.getbuf(stat)
	if val > 0 {
		b = append(b, '+')
	}
	b = strconv.AppendInt(b, int64(val), 10)
	return c.finishbuf(b, "|g", rate)
}

// Inc emits a int counter.
func (c *Client) Inc(stat string, val int) error {
	return c.SampledInc(stat, val, 1.0)
}

// Gauge emits a int gauge.
func (c *Client) Gauge(stat string, val int) error {
	return c.SampledGauge(stat, val, 1.0)
}

// Delta emits a int delta.
func (c *Client) Delta(stat string, val int) error {
	return c.SampledDelta(stat, val, 1.0)
}
