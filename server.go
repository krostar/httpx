package httpx

import (
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
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       8 * time.Second,
		// this timeout, with TLS, start counting down at the very
		// beginning of the request, right after the accept() call,
		// and not when start writing (which is the behavior without
		// tls). Settings for this value should consider it.
		WriteTimeout:   50 * time.Second,
		IdleTimeout:    1 * time.Minute,
		MaxHeaderBytes: 10 << 10, // 10 ko
		TLSConfig:      o.tlsConfig,
	}
}
