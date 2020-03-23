package httpx

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
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
			return nil, fmt.Errorf("unable to apply option: %w", err)
		}
	}

	lc := net.ListenConfig{
		KeepAlive: o.keepAlive,
	}

	l, err := lc.Listen(ctx, o.network, address)
	if err != nil {
		return nil, fmt.Errorf("unable to listen on %s: %w", address, err)
	}

	if o.tlsConfig != nil {
		l = tls.NewListener(l, o.tlsConfig)
	}

	return l, err
}
