package httpx

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerWithModernTLSConfig(t *testing.T) {
	var (
		o   serverOptions
		cfg tls.Config
	)

	tlsSetModernConfig(&cfg)

	ServerWithModernTLSConfig()(&o)
	assert.Equal(t, &cfg, o.tlsConfig)
}

func TestServerWithTLSConfig(t *testing.T) {
	var (
		o   serverOptions
		cfg = &tls.Config{ServerName: "hello"}
	)

	ServerWithTLSConfig(cfg)(&o)
	assert.Equal(t, cfg, o.tlsConfig)
}
