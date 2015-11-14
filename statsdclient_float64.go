// Package statsdclient provides a StatsD client that simply formats output to any underlying io.Writer.
//
// this file is auto-generated. do not edit.
//
// parameters:
//  type: float64
//  suffix: Float64
//  valueconv: b = strconv.AppendFloat(b, val, 'f', 10, 64)
//
package statsdclient

import "strconv"

// SampledIncFloat64 emits a 'sampled' float64 counter with the given rate.
func (c *Client) SampledIncFloat64(stat string, val float64, rate float32) error {
	b := c.getbuf(stat)
	b = strconv.AppendFloat(b, val, 'f', 10, 64)
	return c.finishbuf(b, "|c", rate)
}

// SampledGaugeFloat64 emits a 'sampled' float64 gauge with the given rate.
func (c *Client) SampledGaugeFloat64(stat string, val float64, rate float32) error {
	b := c.getbuf(stat)
	b = strconv.AppendFloat(b, val, 'f', 10, 64)
	return c.finishbuf(b, "|g", rate)
}

// SampledDeltaFloat64 emits a 'sampled' float64 delta with the given rate.
func (c *Client) SampledDeltaFloat64(stat string, val float64, rate float32) error {
	b := c.getbuf(stat)
	if val > 0 {
		b = append(b, '+')
	}
	b = strconv.AppendFloat(b, val, 'f', 10, 64)
	return c.finishbuf(b, "|g", rate)
}

// IncFloat64 emits a float64 counter.
func (c *Client) IncFloat64(stat string, val float64) error {
	return c.SampledIncFloat64(stat, val, 1.0)
}

// GaugeFloat64 emits a float64 gauge.
func (c *Client) GaugeFloat64(stat string, val float64) error {
	return c.SampledGaugeFloat64(stat, val, 1.0)
}

// DeltaFloat64 emits a float64 delta.
func (c *Client) DeltaFloat64(stat string, val float64) error {
	return c.SampledDeltaFloat64(stat, val, 1.0)
}
