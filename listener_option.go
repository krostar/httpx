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

// ListenWithModernTLSConfig sets the tls configuration for tls.NewListener.
func ListenWithModernTLSConfig(certFile, keyFile string) ListenerOption {
	return func(o *listenerOptions) error {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("unable to load tls key pair: %w", err)
		}

		o.tlsConfig = new(tls.Config)
		tlsSetModernConfig(o.tlsConfig)
		o.tlsConfig.Certificates = []tls.Certificate{cert}
		return err
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
func ListenWithKeepAlive(keepPeriod time.Duration) ListenerOption {
	return func(o *listenerOptions) error {
		o.keepAlive = keepPeriod
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
