package statsdclient

import (
	"testing"
	"time"
)

func TestResolvingUDPConn(t *testing.T) {
	c := ResolvingUDPConn{Addr: "127.0.0.1:15261", Window: 0 * time.Second}
	conn1 := c.conn
	c.Write([]byte{})
	conn2 := c.conn
	if conn1 == conn2 {
		t.Fatal("after one write, we failed to open another UDP socket")
	}

	c.Write([]byte{})
	conn3 := c.conn
	if conn2 == conn3 {
		t.Fatal("after two writes, we failed to open another UDP socket")
	}

	if conn1 == conn3 {
		t.Fatal("strangely, the first connection is the same as the last one")
	}
}
