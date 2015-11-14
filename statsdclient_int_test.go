package statsdclient

import (
	"bytes"
	"strings"
	"testing"
)

var IntTests = []FmtTest{
	func(c *Client) (error, string, string) {
		test := "Inc.420"
		expect := "Inc.420:420|c"
		return c.Inc(test, 420), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "Inc.Negative420"
		expect := "Inc.Negative420:-420|c"
		return c.Inc(test, -420), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "Inc.Zero"
		expect := "Inc.Zero:0|c"
		return c.Inc(test, 0), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "SampledInc.Zero.73"
		expect := "SampledInc.Zero.73:0|c|@0.7300"
		return c.SampledInc(test, 0, 0.73), test, expect
	},
}

func TestIntFormatting(t *testing.T) {
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	for _, test := range IntTests {
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
