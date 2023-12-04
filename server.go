package httpx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"go.uber.org/multierr"
)

// NewServer returns a configured http server with robust defaults.
func NewServer(handler http.HandlerFunc, opts ...ServerOption) *http.Server {
	var o serverOptions
	for _, opt := range opts {
		opt(&o)
	}

	return &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 15 * time.Second,
		IdleTimeout:       3 * time.Minute,
		MaxHeaderBytes:    5 << 10, // 5ko
		ErrorLog:          o.errorLogger,
	}
}

// Serve serves the server through the provided listener.
// On context cancellation, the server tries to gracefully
// shut down for as long as shutdownTimeout. Once this timeout
// is reached, the server is stopped, any way.
func Serve(ctx context.Context, s *http.Server, l net.Listener, shutdownTimeout time.Duration) error {
	cerr := make(chan error)
	go func() {
		if err := s.Serve(l); !errors.Is(err, http.ErrServerClosed) {
			cerr <- fmt.Errorf("unable to serve listener %s: %w", l.Addr().String(), err)
			return
		}
		cerr <- nil
	}()

	select {
	case err := <-cerr:
		return err
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) { //nolint:contextcheck // we don't want to provide the function's context to give some time to the server to gracefully shut down
			return fmt.Errorf("unable to shut server down: %w", multierr.Combine(err, <-cerr))
		}
	}

	return <-cerr
}

// ListenAndServe is a shortcut for NewListener and Serve.
func ListenAndServe(ctx context.Context, address string, handler http.HandlerFunc) error {
	l, err := NewListener(ctx, address)
	if err != nil {
		return err
	}
	return Serve(ctx, NewServer(handler), l, time.Minute)
}
