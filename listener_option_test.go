package httpx

import (
	"crypto/tls"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func Test_ListenWithNetwork(t *testing.T) {
	var o listenerOptions
	err := ListenWithNetwork("net")(&o)
	assert.NilError(t, err)
	assert.Equal(t, "net", o.network)
}

func Test_ListenWithIntermediateTLSConfig(t *testing.T) {
	var o listenerOptions
	assert.NilError(t, ListenWithIntermediateTLSConfig("./testdata/cert.crt", "./testdata/cert.key")(&o))
	assert.Check(t, o.tlsConfig != nil, "tls config should not be nil")
	assert.Check(t, len(o.tlsConfig.Certificates) != 0, "tls config should contain certificates")
	assert.Check(t, o.tlsConfig.MinVersion == tls.VersionTLS12, "tls config min version should be 1.2")
	assert.ErrorContains(t, ListenWithIntermediateTLSConfig("./dont/exists", "./testdata/cert.key")(&o), "unable to load tls key pair")
}

func Test_ListenWithTLSConfig(t *testing.T) {
	var o listenerOptions
	assert.NilError(t, ListenWithTLSConfig(&tls.Config{ServerName: "meee", MinVersion: tls.VersionTLS12})(&o))
	assert.Check(t, o.tlsConfig != nil, "tls config should not be nil")
	assert.Check(t, o.tlsConfig.ServerName == "meee")
}

func Test_ListenWithKeepAlive(t *testing.T) {
	var o listenerOptions
	err := ListenWithKeepAlive(3 * time.Millisecond)(&o)
	assert.NilError(t, err)
	assert.Equal(t, 3*time.Millisecond, o.keepAlive)
}

func Test_ListenWithoutKeepAlive(t *testing.T) {
	var o listenerOptions
	o.keepAlive = 10 * time.Second
	err := ListenWithoutKeepAlive()(&o)
	assert.NilError(t, err)
	assert.Equal(t, time.Duration(-1), o.keepAlive)
}
