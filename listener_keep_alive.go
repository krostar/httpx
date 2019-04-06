package httpx

import (
	"net"
	"time"

	"github.com/pkg/errors"
)

type tcpListenerKeepAlive struct {
	*net.TCPListener
	period time.Duration
}

// Accept is a copy of how net/http is doing the same thing
// but with a possible configuration of the keep alive period.
func (l tcpListenerKeepAlive) Accept() (net.Conn, error) {
	conn, err := l.AcceptTCP()
	if err != nil {
		return nil, err
	}

	if l.period > 0 {
		if err := conn.SetKeepAlive(true); err != nil {
			return nil, errors.Wrap(err, "unable to set keepalive")
		}
		if err := conn.SetKeepAlivePeriod(l.period); err != nil {
			return nil, errors.Wrap(err, "unable to set keepalive period")
		}
	}

	return conn, nil
}
