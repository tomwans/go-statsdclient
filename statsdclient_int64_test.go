package statsdclient

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

var Int64Tests = []FmtTest{
	func(c *Client) (error, string, string) {
		test := "IncInt64.Max"
		expect := "IncInt64.Max:9223372036854775807|c"
		return c.IncInt64(test, math.MaxInt64), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "IncInt64.Min"
		expect := "IncInt64.Min:-9223372036854775808|c"
		return c.IncInt64(test, math.MinInt64), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "IncInt64.Zero"
		expect := "IncInt64.Zero:0|c"
		return c.IncInt64(test, 0), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "SampledIncInt64.Zero.73"
		expect := "SampledIncInt64.Zero.73:0|c|@0.7300"
		return c.SampledIncInt64(test, 0, 0.73), test, expect
	},
}

func TestInt64Formatting(t *testing.T) {
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	for _, test := range Int64Tests {
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
