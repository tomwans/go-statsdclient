package statsdclient

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

type FmtTest func(*Client) (error, string, string)

var TimingTests = []FmtTest{
	func(c *Client) (error, string, string) {
		test := "BasicTiming"
		expect := "BasicTiming:2000|ms"
		return c.Timing(test, 2*time.Second), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "NegativeTiming"
		expect := "NegativeTiming:-2000|ms"
		return c.Timing(test, -2*time.Second), test, expect
	},
}

func TestTimingFormatting(t *testing.T) {
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	for _, test := range TimingTests {
		buf.Reset()
		err, name, expected := test(c)
		if err != nil {
			t.Errorf("%s test: error: %s", name, err)
		}
		got := string(buf.Bytes())
		if got != expected+"\n" {
			t.Errorf("%s test:\ngot:    '%s'\nexpect: '%s'", name, strings.TrimSpace(got), expected)
		}
	}
}

func TestPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "test_without_trailing_dot", Writer: buf}
	c.Inc("ok", 1)
	got := string(buf.Bytes())
	expected := "test_without_trailing_dot.ok:1|c\n"
	if got != expected {
		t.Fatalf("got != expected: %s, %s", got, expected)
	}
}

func TestPrefixWithDot(t *testing.T) {
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "test_with_dot.", Writer: buf}
	c.Inc("ok", 1)
	got := string(buf.Bytes())
	expected := "test_with_dot..ok:1|c\n"
	if got != expected {
		t.Fatalf("got != expected: %s, %s", got, expected)
	}
}

func BenchmarkParallel(b *testing.B) {
	b.StopTimer()
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := c.Inc("test.stat", 420)
			if err != nil {
				b.Fatalf("error: %s", err)
			}
		}
	})
}

func BenchmarkIntFormatting(b *testing.B) {
	b.StopTimer()
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := c.Inc("test.stat", 420)
		if err != nil {
			b.Fatalf("error: %s", err)
		}
	}
}

func BenchmarkIntFormattingSampled(b *testing.B) {
	b.StopTimer()
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := c.SampledInc("test.stat", 420, 0.73)
		if err != nil {
			b.Fatalf("error: %s", err)
		}
	}
}
