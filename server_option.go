package httpx

import (
	"crypto/tls"
)

type serverOptions struct {
	tlsConfig *tls.Config
}

// ServerOption defines options applier for the server.
type ServerOption func(*serverOptions)

// ServerWithModernTLSConfig parses two files as x509 certificates and load
// a modern (based on https://wiki.mozilla.org/Security/Server_Side_TLS
// and https://blog.cloudflare.com/exposing-go-on-the-internet/) configuration.
// This configuration should be compatible at least with Firefox 27, Chrome 30,
// IE 11 on Windows 7, Edge, Opera 17, Safari 9, Android 5.0, and Java 8.
func ServerWithModernTLSConfig() ServerOption {
	return func(o *serverOptions) {
		o.tlsConfig = new(tls.Config)
		tlsSetModernConfig(o.tlsConfig)
	}
}

// ServerWithTLSConfig sets the tls configuration for the server.
func ServerWithTLSConfig(cfg *tls.Config) ServerOption {
	return func(o *serverOptions) {
		o.tlsConfig = cfg
	}
}
