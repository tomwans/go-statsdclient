package statsdclient

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"time"
)

// keep a pool of buffers to write entire metric lines out, which will
// then be flushed to our underlying writer.
var bufPool = &sync.Pool{New: func() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 128))
}}

// Client represents a connection to statsd. You can call several
// measurement methods on the client, and they will be formatted and
// written to the underlying Writer.
type Client struct {
	// prefix we should write for every measurement. it should not end with '.'
	Prefix string
	// underlying writer we will write all of our measurements to. it
	// will most likely be a UDP socket for the case of statsd, but can
	// be anything at all.
	Writer io.Writer
}

func (c *Client) write(stat, value, suffix string, rate float32) error {
	buf := bufPool.Get().(*bytes.Buffer)

	if c.Prefix != "" {
		buf.WriteString(c.Prefix)
		buf.WriteString(".")
	}

	buf.WriteString(stat)
	buf.WriteString(":")
	buf.WriteString(value)
	buf.WriteString(suffix)
	if rate != 1.0 {
		buf.WriteString("|@")
		buf.WriteString(strconv.FormatFloat(float64(rate), 'f', 4, 32))
	}
	buf.WriteString("\n")
	// todo: io.Copy might be a bad idea here.
	_, err := io.Copy(c.Writer, buf)
	buf.Reset()
	bufPool.Put(buf)
	return err
}

// SampledInc will emit a counter metric + sampling rate.
func (c *Client) SampledInc(stat string, value int64, rate float32) {
	c.write(stat, strconv.FormatInt(value, 10), "|c", rate)
}

// SampledGauge will emit a gauge metric + sampling rate.
func (c *Client) SampledGauge(stat string, value int64, rate float32) {
	c.write(stat, strconv.FormatInt(value, 10), "|g", rate)
}

// SampledDelta will emit a gauge delta metric + sampling rate.
func (c *Client) SampledDelta(stat string, value int64, rate float32) {
	prefix := ""
	if value > 0 {
		prefix = "+"
	}
	c.write(stat, prefix+strconv.FormatInt(value, 10), "|g", rate)
}

// SampledTiming will emit a timing metric + sampling rate.
func (c *Client) SampledTiming(stat string, duration time.Duration, rate float32) {
	ms := duration.Nanoseconds() / 1000000
	c.write(stat, strconv.FormatInt(ms, 10), "|ms", rate)
}

// Inc will emit a counter with a sample rate of 1.0
func (c *Client) Inc(stat string, value int64) {
	c.SampledInc(stat, value, 1.0)
}

// Gauge will emit a gauge with a sample rate of 1.0
func (c *Client) Gauge(stat string, value int64) {
	c.SampledGauge(stat, value, 1.0)
}

// Delta will emit a gauge delta with a sample rate of 1.0
func (c *Client) Delta(stat string, value int64) {
	c.SampledDelta(stat, value, 1.0)
}

// Timing will emit a timing metric with a sample rate of 1.0
func (c *Client) Timing(stat string, duration time.Duration) {
	c.SampledTiming(stat, duration, 1.0)
}
