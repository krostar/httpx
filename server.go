// Package httpx ...
package httpx

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// NewServer returns a configured server http. Default (no
// given options) is valid and probably enough for most servers.
func NewServer(handler http.Handler, opts ...ServerOption) *http.Server {
	var o serverOptions

	for _, opt := range opts {
		opt(&o)
	}

	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 2 * time.Second, // nolint: gomnd
		ReadTimeout:       8 * time.Second, // nolint: gomnd
		// this timeout, with TLS, start counting down at the very
		// beginning of the request, right after the accept() call,
		// and not when start writing (which is the behavior without
		// tls). Settings for this value should consider it.
		WriteTimeout:   50 * time.Second, // nolint: gomnd
		IdleTimeout:    1 * time.Minute,  // nolint: gomnd
		MaxHeaderBytes: 10 << 10,         // nolint: gomnd
		TLSConfig:      o.tlsConfig,
	}
}

// Serve serves the server through the provided listener.
// On context cancellation, the server tries to gracefully
// shut down for as long as shutdownTimeout. Once this timeout
// is reached, the server is stopped, any way.
func Serve(ctx context.Context, s *http.Server, l net.Listener, shutdownTimeout time.Duration) error {
	cerr := make(chan error)
	go func() {
		if err := s.Serve(l); err != http.ErrServerClosed {
			cerr <- fmt.Errorf("unable to serve listener %q: %w", l.Addr().String(), err)
			return
		}
		cerr <- nil
	}()

	for {
		select {
		case err := <-cerr:
			return err
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
			defer cancel()
			if err := s.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
				return fmt.Errorf("unable to shut server down: %w", err)
			}
		}
	}
}

// ListenAndServe is a shortcut for NewListener and Serve.
func ListenAndServe(ctx context.Context, s *http.Server) error {
	l, err := NewListener(ctx, s.Addr)
	if err != nil {
		return err
	}
	return Serve(ctx, s, l, 15*time.Second) // nolint: gomnd
}
