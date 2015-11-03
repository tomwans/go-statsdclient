package statsdclient

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"time"
)

var bufPool = &sync.Pool{New: func() interface{} {
	return bytes.NewBuffer(make([]byte, 0, 128))
}}

type Client struct {
	Writer io.Writer
	Prefix string
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
	_, err := io.Copy(c.Writer, buf)
	buf.Reset()
	bufPool.Put(buf)
	return err
}

func (c *Client) SampledInc(stat string, value int64, rate float32) {
	c.write(stat, strconv.FormatInt(value, 10), "|c", rate)
}

func (c *Client) SampledGauge(stat string, value int64, rate float32) {
	c.write(stat, strconv.FormatInt(value, 10), "|g", rate)
}

func (c *Client) SampledDelta(stat string, value int64, rate float32) {
	prefix := ""
	if value > 0 {
		prefix = "+"
	}
	c.write(stat, prefix+strconv.FormatInt(value, 10), "|g", rate)
}

func (c *Client) SampledTiming(stat string, duration time.Duration, rate float32) {
	ms := duration.Nanoseconds() / 1000000
	c.write(stat, strconv.FormatInt(ms, 10), "|ms", rate)
}

func (c *Client) Inc(stat string, value int64) {
	c.SampledInc(stat, value, 1.0)
}

func (c *Client) Gauge(stat string, value int64) {
	c.SampledGauge(stat, value, 1.0)
}

func (c *Client) Delta(stat string, value int64) {
	c.SampledDelta(stat, value, 1.0)
}

func (c *Client) Timing(stat string, duration time.Duration) {
	c.SampledTiming(stat, duration, 1.0)
}
