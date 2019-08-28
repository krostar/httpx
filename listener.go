package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/pkg/errors"
)

// NewListener returns a new listener. The built listener can be
// customized (keepalive duration, tls config, ...) via options.
// Default (no options) is valid and probably enough for most servers.
func NewListener(ctx context.Context, address string, opts ...ListenerOption) (net.Listener, error) {
	var o = listenerOptions{
		network:   "tcp4",
		keepAlive: 1 * time.Minute,
	}

	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, errors.Wrap(err, "unable to apply option")
		}
	}

	var lc net.ListenConfig
	l, err := lc.Listen(ctx, o.network, address)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to listen on %s", address)
	}

	l = tcpListenerKeepAlive{
		TCPListener: l.(*net.TCPListener),
		period:      o.keepAlive, // if period is 0, keepalive is disabled
	}

	if o.tlsConfig != nil {
		l = tls.NewListener(l, o.tlsConfig)
	}

	return l, err
}
