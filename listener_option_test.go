package httpx

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListenWithNetwork(t *testing.T) {
	var o listenerOptions
	err := ListenWithNetwork("net")(&o)
	require.NoError(t, err)
	assert.Equal(t, "net", o.network)
}

func TestListenWithModernTLSConfig(t *testing.T) {
	var o listenerOptions

	assert.NoError(t, ListenWithModernTLSConfig("./testdata/cert.crt", "./testdata/cert.key")(&o))
	assert.Error(t, ListenWithModernTLSConfig("./dont/exists", "./testdata/cert.key")(&o))
}

func TestListenWithTLSConfig(t *testing.T) {
	var o listenerOptions
	err := ListenWithTLSConfig(&tls.Config{ServerName: "meee"})(&o)
	require.NoError(t, err)
	assert.Equal(t, "meee", o.tlsConfig.ServerName)
}

func TestListenWithKeepAlive(t *testing.T) {
	var o listenerOptions
	err := ListenWithKeepAlive(3 * time.Millisecond)(&o)
	require.NoError(t, err)
	assert.Equal(t, 3*time.Millisecond, o.keepAlive)
}

func TestListenWithoutKeepAlive(t *testing.T) {
	var o listenerOptions
	o.keepAlive = 10 * time.Second
	err := ListenWithoutKeepAlive()(&o)
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), o.keepAlive)
}
