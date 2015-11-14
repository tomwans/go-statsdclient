// Package statsdclient provides a client that can format StatsD
// metrics to a an io.Writer.
//
//go:generate go run internal/gen/main.go -type=int -file-dest=./statsdclient_int.go
//go:generate go run internal/gen/main.go -type=int64 -file-dest=./statsdclient_int64.go
//go:generate go run internal/gen/main.go -type=float64 -file-dest=./statsdclient_float64.go
//go:generate gofmt -s -w .
package statsdclient

import (
	"io"
	"strconv"
	"sync"
	"time"
)

// Client represents a connection to StatsD. One can view this
// essentially as a StatsD formatter with convenient methods for
// measuring counters, gauges, and so on, for different types (int,
// in64, float64). The formatter just writes a StatsD formatted string
// to an underlying io.Writer, where in most cases the Writer will be
// a connection to the StatsD server.
//
// A Client can be used across several goroutines.
type Client struct {
	// Prefix for every stat. It should not end with '.'
	Prefix string
	// Writer we will write all of our measurements to. It will most
	// likely be a UDP socket for the case of statsd, but can be
	// anything at all.
	Writer io.Writer

	bufMu   sync.Mutex
	tempBuf [512]byte
}

// getbuf returns a buffer that is prepped with the 'stat' and prefix
// already filled out. The caller is expected to append bytes directly
// to the buffer, and when complete, will pass the buf to finishbuf.
func (c *Client) getbuf(stat string) []byte {
	c.bufMu.Lock()
	buf := c.tempBuf[0:0]
	if c.Prefix != "" {
		buf = append(buf, c.Prefix...)
		buf = append(buf, '.')
	}
	buf = append(buf, stat...)
	buf = append(buf, ':')
	return buf
}

// finishbuf will take a buffer and finish writing the suffix, rate,
// etc. It will also flush the buffer's contents to the underlying
// writer.
func (c *Client) finishbuf(buf []byte, suffix string, rate float32) error {
	buf = append(buf, suffix...)
	if rate != 1.0 {
		buf = append(buf, "|@"...)
		buf = strconv.AppendFloat(buf, float64(rate), 'f', 4, 32)
	}
	buf = append(buf, '\n')
	_, err := c.Writer.Write(buf)
	c.bufMu.Unlock()
	return err
}

// SampledTiming will emit a timing metric + sampling rate.
func (c *Client) SampledTiming(stat string, duration time.Duration, rate float32) error {
	buf := c.getbuf(stat)
	ms := duration.Nanoseconds() / 1000000
	buf = strconv.AppendInt(buf, ms, 10)
	return c.finishbuf(buf, "|ms", rate)
}

// Timing will emit a timing metric with a sample rate of 1.0
func (c *Client) Timing(stat string, duration time.Duration) error {
	return c.SampledTiming(stat, duration, 1.0)
}

// TimingSince will return a timer since a given time.
func (c *Client) TimingSince(stat string, since time.Time) error {
	return c.Timing(stat, time.Now().Sub(since))
}
