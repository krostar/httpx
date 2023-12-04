package httpx

import (
	"crypto/tls"
	"fmt"
	"time"
)

type listenerOptions struct {
	network   string
	keepAlive time.Duration
	tlsConfig *tls.Config
}

// ListenerOption defines options applier for the listener.
type ListenerOption func(*listenerOptions) error

// ListenWithNetwork sets the network option for net.Listen().
func ListenWithNetwork(network string) ListenerOption {
	return func(o *listenerOptions) error {
		o.network = network
		return nil
	}
}

// ListenWithIntermediateTLSConfig sets the tls configuration for tls.NewListener.
// "Intermediate" is defined based on this website: https://wiki.mozilla.org/Security/Server_Side_TLS.
func ListenWithIntermediateTLSConfig(certFile, keyFile string) ListenerOption {
	return func(o *listenerOptions) error {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("unable to load tls key pair: %w", err)
		}

		return ListenWithTLSConfig(&tls.Config{
			Certificates: []tls.Certificate{cert},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384},
		})(o)
	}
}

// ListenWithTLSConfig sets the tls configuration for tls.NewListener.
func ListenWithTLSConfig(cfg *tls.Config) ListenerOption {
	return func(o *listenerOptions) error {
		o.tlsConfig = cfg
		return nil
	}
}

// ListenWithKeepAlive sets keepalive period.
func ListenWithKeepAlive(keepAlive time.Duration) ListenerOption {
	return func(o *listenerOptions) error {
		o.keepAlive = keepAlive
		return nil
	}
}

// ListenWithoutKeepAlive disables the keepalive on the listener.
func ListenWithoutKeepAlive() ListenerOption {
	return func(o *listenerOptions) error {
		o.keepAlive = -1
		return nil
	}
}
