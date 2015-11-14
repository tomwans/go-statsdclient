package statsdclient

import (
	"bytes"
	"math"
	"strings"
	"testing"
)

var Float64Tests = []FmtTest{
	func(c *Client) (error, string, string) {
		test := "IncFloat64.Max"
		expect := "IncFloat64.Max:179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.0000000000|c"
		return c.IncFloat64(test, math.MaxFloat64), test, expect
	},
	func(c *Client) (error, string, string) {
		// yeah, we basically don't support this one. I don't know how to
		// format it if statsd can't accept any kind of exponent notation.
		test := "IncFloat64.SmallestNonzeroFloat64"
		expect := "IncFloat64.SmallestNonzeroFloat64:0.0000000000|c"
		return c.IncFloat64(test, math.SmallestNonzeroFloat64), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "IncFloat64.Pi"
		expect := "IncFloat64.Pi:3.1415926536|c"
		return c.IncFloat64(test, math.Pi), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "GaugeFloat64.Pi"
		expect := "GaugeFloat64.Pi:3.1415926536|g"
		return c.GaugeFloat64(test, math.Pi), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "DeltaFloat64.Pi"
		expect := "DeltaFloat64.Pi:+3.1415926536|g"
		return c.DeltaFloat64(test, math.Pi), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "DeltaFloat64.Pi"
		expect := "DeltaFloat64.Pi:-3.1415926536|g"
		return c.DeltaFloat64(test, -1*math.Pi), test, expect
	},
	func(c *Client) (error, string, string) {
		test := "DeltaFloat64.Zero"
		expect := "DeltaFloat64.Zero:0.0000000000|g"
		return c.DeltaFloat64(test, 0), test, expect
	},
}

func TestFloat64Formatting(t *testing.T) {
	buf := &bytes.Buffer{}
	c := &Client{Prefix: "", Writer: buf}
	for _, test := range Float64Tests {
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
