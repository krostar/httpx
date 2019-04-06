package httpx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	srv := NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}))

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go func() {
		defer srv.Shutdown(context.Background()) // nolint: errcheck
		resp, err := http.Get("http://" + l.Addr().String())
		require.NoError(t, err)
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
	}()

	assert.Equal(t, http.ErrServerClosed, srv.Serve(l))
}

func TestNewServer_with_tls(t *testing.T) {
	srv := NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusCreated)
	}), ServerWithModernTLSConfig())
	srv.Addr = ":7777"

	go func() {
		defer srv.Shutdown(context.Background()) // nolint: errcheck

		rootCAPEM, err := ioutil.ReadFile("./testdata/ca.crt")
		require.NoError(t, err)
		rootCAs := x509.NewCertPool()
		require.True(t, rootCAs.AppendCertsFromPEM(rootCAPEM))

		var client http.Client
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{RootCAs: rootCAs, ServerName: "go-test"}}
		resp, err := client.Get("https://:7777")
		require.NoError(t, err)
		require.NotNil(t, resp.TLS)
		assert.Equal(t, resp.StatusCode, http.StatusCreated)
	}()

	assert.Equal(t, http.ErrServerClosed, srv.ListenAndServeTLS("./testdata/cert.crt", "./testdata/cert.key"))
}
