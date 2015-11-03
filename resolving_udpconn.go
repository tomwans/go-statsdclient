package statsdclient

import (
	"net"
	"time"
)

type ResolvingUDPConn struct {
	Addr    string
	Window  time.Duration
	expires time.Time
	conn    *net.UDPConn
}

func (s *ResolvingUDPConn) possiblyResolve() error {
	// do we need to re-resolve?
	if time.Now().After(s.expires) || s.conn == nil {
		// re-resolve.
		raddr, err := net.ResolveUDPAddr("udp", s.Addr)
		if err != nil {
			return err
		}
		conn, err := net.DialUDP("udp", nil, raddr)
		if err != nil {
			return err
		}
		s.conn = conn
		s.expires = time.Now().Add(s.Window)
	}
	return nil
}

func (s *ResolvingUDPConn) Write(b []byte) (int, error) {
	err := s.possiblyResolve()
	if err != nil {
		return 0, err
	}
	return s.conn.Write(b)
}
