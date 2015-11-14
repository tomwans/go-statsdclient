// Package statsdclient provides a StatsD client that simply formats output to any underlying io.Writer.
//
// this file is auto-generated. do not edit.
//
// parameters:
//  type: int64
//  suffix: Int64
//  valueconv: b = strconv.AppendInt(b, val, 10)
//
package statsdclient

import "strconv"

// SampledIncInt64 emits a 'sampled' int64 counter with the given rate.
func (c *Client) SampledIncInt64(stat string, val int64, rate float32) error {
	b := c.getbuf(stat)
	b = strconv.AppendInt(b, val, 10)
	return c.finishbuf(b, "|c", rate)
}

// SampledGaugeInt64 emits a 'sampled' int64 gauge with the given rate.
func (c *Client) SampledGaugeInt64(stat string, val int64, rate float32) error {
	b := c.getbuf(stat)
	b = strconv.AppendInt(b, val, 10)
	return c.finishbuf(b, "|g", rate)
}

// SampledDeltaInt64 emits a 'sampled' int64 delta with the given rate.
func (c *Client) SampledDeltaInt64(stat string, val int64, rate float32) error {
	b := c.getbuf(stat)
	if val > 0 {
		b = append(b, '+')
	}
	b = strconv.AppendInt(b, val, 10)
	return c.finishbuf(b, "|g", rate)
}

// IncInt64 emits a int64 counter.
func (c *Client) IncInt64(stat string, val int64) error {
	return c.SampledIncInt64(stat, val, 1.0)
}

// GaugeInt64 emits a int64 gauge.
func (c *Client) GaugeInt64(stat string, val int64) error {
	return c.SampledGaugeInt64(stat, val, 1.0)
}

// DeltaInt64 emits a int64 delta.
func (c *Client) DeltaInt64(stat string, val int64) error {
	return c.SampledDeltaInt64(stat, val, 1.0)
}
